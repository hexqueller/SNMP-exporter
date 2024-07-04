package snmp

import (
	"log"
	"net"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hexqueller/SNMP-proxy/internal/config"
)

type Data struct {
	Timestamp time.Time
	Data      map[string]interface{}
}

func PollAgent(agent config.AgentConfig, dataChannel chan<- Data) {
	log.Printf("Polling agent: %s\n", agent.Target)

	var version gosnmp.SnmpVersion
	switch agent.Version {
	case "1":
		version = gosnmp.Version1
	case "2c":
		version = gosnmp.Version2c
	case "3":
		version = gosnmp.Version3
	default:
		log.Fatalf("Unknown SNMP version: %s", agent.Version)
	}

	params := &gosnmp.GoSNMP{
		Target:    agent.Target,
		Port:      agent.Port,
		Community: agent.Community,
		Version:   version,
		Timeout:   time.Duration(10) * time.Second,
		Retries:   3,
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("Close() err: %v", err)
		}
	}(params.Conn)

	for {
		log.Println("Starting SNMP walk")
		data := make(map[string]interface{})
		err = params.Walk(agent.OID, func(pdu gosnmp.SnmpPDU) error {
			log.Printf("Received PDU: %v", pdu)
			switch pdu.Type {
			case gosnmp.OctetString:
				data[pdu.Name] = string(pdu.Value.([]byte))
			default:
				data[pdu.Name] = pdu.Value
			}
			return nil
		})

		if err != nil {
			log.Printf("Walk() err: %v", err)
		} else {
			dataChannel <- Data{
				Timestamp: time.Now(),
				Data:      data,
			}
		}

		time.Sleep(time.Duration(agent.Interval) * time.Second)
	}
}
