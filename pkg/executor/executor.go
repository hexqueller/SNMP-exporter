package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
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

func StartAgents() {
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

func runPythonScript(filePath string, dataPort uint16) error {
	scriptPath := "./agent.py"

	cmd := exec.Command("python3", scriptPath, filePath, fmt.Sprintf("%d", dataPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
