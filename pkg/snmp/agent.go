package snmp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hexqueller/SNMP-proxy/internal/config"
)

type Data struct {
	Target string
	Data   map[string]string
}

var (
	collectedData   map[string]Data
	collectedDataMu sync.Mutex
)

func init() {
	collectedData = make(map[string]Data)
}

func PollAgent(agent config.AgentConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Polling agent: %s\n", agent.Target)

	version, err := getSNMPVersion(agent.Version)
	if err != nil {
		log.Fatalf("Invalid SNMP version: %v", err)
	}

	params := &gosnmp.GoSNMP{
		Target:    agent.Target,
		Port:      agent.Port,
		Community: agent.Community,
		Version:   version,
		Timeout:   time.Duration(5) * time.Second,
		Retries:   3,
	}

	err = params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("Close err: %v", err)
		}
	}(params.Conn)

	log.Println("Connected successfully to agent:", agent.Target)

	result := make(map[string]string)
	err = params.Walk(agent.OID, func(pdu gosnmp.SnmpPDU) error {
		switch pdu.Type {
		case gosnmp.OctetString:
			result[pdu.Name] = string(pdu.Value.([]byte))
		default:
			result[pdu.Name] = fmt.Sprintf("%v", pdu.Value)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Walk() err: %v", err)
	}

	collectedDataMu.Lock()
	collectedData[agent.Name] = Data{Target: agent.Target, Data: result}
	collectedDataMu.Unlock()

	log.Println("Completed SNMP walk for agent:", agent.Target)
}

func SaveDataToFiles() {
	collectedDataMu.Lock()
	defer collectedDataMu.Unlock()

	for name, data := range collectedData {
		filename := fmt.Sprintf("%s_snmp_data.json", name)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Failed to close file: %v", err)
			}
		}(file)

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			log.Fatalf("Failed to write data to file: %v", err)
		}
		log.Println("Data saved to file:", filename)
	}
}

func getSNMPVersion(version string) (gosnmp.SnmpVersion, error) {
	switch version {
	case "1":
		return gosnmp.Version1, nil
	case "2c":
		return gosnmp.Version2c, nil
	case "3":
		return gosnmp.Version3, nil
	default:
		return 0, fmt.Errorf("unsupported SNMP version: %s", version)
	}
}
