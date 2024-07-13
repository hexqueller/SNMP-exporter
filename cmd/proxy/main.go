package main

import (
	"flag"
	"github.com/gosnmp/gosnmp"
	"github.com/hexqueller/SNMP-proxy/internal/config"
	"github.com/hexqueller/SNMP-proxy/pkg/snmp"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/GoSNMPServer/mibImps"
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

	conf := Config{
		V3Username:                 "myV3User",
		V3AuthenticationPassphrase: "myAuthPassphrase",
		V3PrivacyPassphrase:        "myPrivPassphrase",
		Community:                  "myCommunity",
	}
	master := GoSNMPServer.MasterAgent{
		Logger: GoSNMPServer.NewDefaultLogger(),
		SecurityConfig: GoSNMPServer.SecurityConfig{
			AuthoritativeEngineBoots: 1,
			Users: []gosnmp.UsmSecurityParameters{
				{
					UserName:                 conf.V3Username,
					AuthenticationProtocol:   gosnmp.MD5,
					PrivacyProtocol:          gosnmp.DES,
					AuthenticationPassphrase: conf.V3AuthenticationPassphrase,
					PrivacyPassphrase:        conf.V3PrivacyPassphrase,
				},
			},
		},
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				CommunityIDs: []string{conf.Community},
				OIDs:         mibImps.All(),
			},
		},
	}
	server := GoSNMPServer.NewSNMPServer(master)
	err = server.ListenUDP("udp", "127.0.0.1:1161")
	if err != nil {
		log.Fatalf("Error in listen: %+v", err)
	}
	err = server.ServeForever()
	if err != nil {
		return
	}
}
