package snmp

import (
	"fmt"
	"log"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hexqueller/SNMP-proxy/internal/config"
)

func PollAgent(agent config.AgentConfig) {
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
		Version:   version,                         // Используем интерпретированную версию
		Timeout:   time.Duration(10) * time.Second, // Увеличиваем тайм-аут
		Retries:   3,
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	log.Println("Connected successfully to agent:", agent.Target)

	log.Println("Starting SNMP walk")
	err = params.Walk(agent.OID, func(pdu gosnmp.SnmpPDU) error {
		log.Printf("Received PDU: %v", pdu) // Логирование PDU для отладки
		switch pdu.Type {
		case gosnmp.OctetString:
			fmt.Printf("Agent %s: %s: %s\n", agent.Target, pdu.Name, string(pdu.Value.([]byte)))
		default:
			fmt.Printf("Agent %s: %s: %v\n", agent.Target, pdu.Name, pdu.Value)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Walk() err: %v", err)
	}

	log.Println("Completed SNMP walk for agent:", agent.Target)
}
