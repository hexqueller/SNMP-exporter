package main

import (
	"flag"
	"github.com/hexqueller/SNMP-proxy/internal/config"
	"github.com/hexqueller/SNMP-proxy/pkg/snmp"
	"log"
	"sync"
)

func main() {
	// Определение флагов
	configPath := flag.String("c", "./configs/default.yaml", "Path to config file")
	flag.Parse()

	// Чтение конфигурационного файла
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	log.Println("Configuration loaded successfully")

	var wg sync.WaitGroup
	for _, agent := range cfg.Agents {
		dataChannel := make(chan snmp.Data)
		wg.Add(1)
		go func(agent config.AgentConfig) {
			defer wg.Done()
			snmp.StartSNMPAgent(agent, dataChannel)
		}(agent)

		wg.Add(1)
		go func(agent config.AgentConfig) {
			defer wg.Done()
			snmp.PollAgent(agent, dataChannel)
		}(agent)
	}

	wg.Wait()
}
