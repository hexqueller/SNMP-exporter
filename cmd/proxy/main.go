package main

import (
	"flag"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type AgentConfig struct {
	Target    string `yaml:"target"`
	Port      uint16 `yaml:"port"`
	Community string `yaml:"community"`
	Version   int    `yaml:"version"`
	OID       string `yaml:"oid"`
}

type Config struct {
	Agents []AgentConfig `yaml:"agents"`
}

func main() {
	// Определение флагов
	timeout := flag.Int("t", 10, "Timeout in seconds")
	configPath := flag.String("c", "./configs/default.yaml", "Path to config file")
	flag.Parse()

	// Чтение конфигурационного файла
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	for _, agent := range config.Agents {
		go pollAgent(agent)
	}

	// Ожидание, чтобы все агенты успели ответить
	time.Sleep(time.Duration(*timeout) * time.Second)
}

func pollAgent(agent AgentConfig) {
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
