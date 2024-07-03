package snmp

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/hexqueller/SNMP-proxy/internal/config"
	"log"
	"time"
)

func PollAgent(agent config.AgentConfig) {
	params := &gosnmp.GoSNMP{
		Target:    agent.Target,
		Port:      agent.Port,
		Community: agent.Community,
		Version:   gosnmp.SnmpVersion(agent.Version),
		Timeout:   time.Duration(2) * time.Second,
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	result, err := params.Get([]string{agent.OID})
	if err != nil {
		log.Fatalf("Get() err: %v", err)
	}

	for _, variable := range result.Variables {
		switch variable.Type {
		case gosnmp.OctetString:
			fmt.Printf("Agent %s: %s: %s\n", agent.Target, variable.Name, string(variable.Value.([]byte)))
		default:
			fmt.Printf("Agent %s: %s: %v\n", agent.Target, variable.Name, variable.Value)
		}
	}
}
