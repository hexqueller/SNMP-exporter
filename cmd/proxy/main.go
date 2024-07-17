package main

import (
	"flag"
	"github.com/hexqueller/SNMP-proxy/internal/config"
	"github.com/hexqueller/SNMP-proxy/pkg/snmp"
	"log"
	"sync"
)

type Config struct {
	V3Username                 string
	V3AuthenticationPassphrase string
	V3PrivacyPassphrase        string
	Community                  string
}

func main() {
	// Определение флагов
	configPath := flag.String("c", "./configs/default.yaml", "Path to loadConfig file")
	flag.Parse()

	// Чтение конфигурационного файла
	loadConfig, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load loadConfig: %v", err)
	}

	var wg sync.WaitGroup
	for _, agent := range loadConfig.Agents {
		wg.Add(1)
		go snmp.PollAgent(agent, &wg)
	}

	wg.Wait()
	snmp.SaveDataToFiles()
	log.Println("Finished waiting")
}
