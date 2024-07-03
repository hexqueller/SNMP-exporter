package main

import (
	"flag"
	"log"
	"time"

	"github.com/hexqueller/SNMP-proxy/internal/config"
	"github.com/hexqueller/SNMP-proxy/pkg/snmp"
)

func main() {
	// Определение флагов
	timeout := flag.Int("t", 10, "Timeout in seconds")
	configPath := flag.String("c", "./configs/default.yaml", "Path to config file")
	flag.Parse()

	// Чтение конфигурационного файла
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	for _, agent := range cfg.Agents {
		go snmp.PollAgent(agent)
	}

	// Ожидание, чтобы все агенты успели ответить
	log.Printf("Waiting for %d seconds\n", *timeout)
	time.Sleep(time.Duration(*timeout) * time.Second)
	log.Println("Finished waiting")
}
