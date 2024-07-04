package snmp

import (
	"log"
	"net"
	"strconv"

	"github.com/gosnmp/gosnmp"
	"github.com/hexqueller/SNMP-proxy/internal/config"
)

var collectedData map[string]Data

func init() {
	collectedData = make(map[string]Data)
}

func StartSNMPAgent(agent config.AgentConfig, dataChannel <-chan Data) {
	go func() {
		for data := range dataChannel {
			collectedData[agent.Target] = data
		}
	}()

	listener, err := net.Listen("udp", ":"+strconv.Itoa(int(agent.DataPort)))
	if err != nil {
		log.Fatalf("Failed to start SNMP agent: %v", err)
	}
	defer listener.Close()

	log.Printf("SNMP agent started on port %d\n", agent.DataPort)

	for {
		buffer := make([]byte, 4096)
		n, addr, err := listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			continue
		}

		go handleRequest(buffer[:n], addr, listener)
	}
}

func handleRequest(request []byte, addr net.Addr, listener net.Listener) {
	packet := &gosnmp.SnmpPacket{}
	err := packet.Unmarshal(request)
	if err != nil {
		log.Printf("Failed to unmarshal SNMP packet: %v", err)
		return
	}

	response := &gosnmp.SnmpPacket{
		Version:   packet.Version,
		Community: packet.Community,
		PDUType:   gosnmp.GetResponse,
		Variables: []gosnmp.SnmpPDU{},
	}

	target := addr.String()
	if data, exists := collectedData[target]; exists {
		for oid, value := range data.Data {
			response.Variables = append(response.Variables, gosnmp.SnmpPDU{
				Name:  oid,
				Type:  gosnmp.OctetString,
				Value: value,
			})
		}
	}

	responseBytes, err := response.Marshal()
	if err != nil {
		log.Printf("Failed to marshal SNMP response: %v", err)
		return
	}

	_, err = listener.WriteTo(responseBytes, addr)
	if err != nil {
		log.Printf("Failed to write SNMP response: %v", err)
	}
}
