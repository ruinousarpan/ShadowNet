package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Node struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	PublicKey string `json:"public_key"`
	Country   string `json:"country"`
	Ping      int    `json:"ping"`
}

func main() {
	// Create config directory if it doesn't exist
	configDir := filepath.Join(os.Getenv("HOME"), ".shadownet")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	// Load nodes from bootstrap.json
	nodes, err := loadNodes("../bootstrap/bootstrap.json")
	if err != nil {
		log.Fatalf("Failed to load nodes: %v", err)
	}

	// Select the fastest node
	fastestNode := selectFastestNode(nodes)
	if fastestNode == nil {
		log.Fatalf("No nodes available")
	}

	fmt.Printf("Connecting to the fastest node: %s (%s) with ping %dms\n", fastestNode.IP, fastestNode.Country, fastestNode.Ping)

	// Generate WireGuard configuration
	err = generateWireGuardConfig(fastestNode, configDir)
	if err != nil {
		log.Fatalf("Failed to generate WireGuard config: %v", err)
	}

	// Establish WireGuard connection
	err = connectWireGuard(configDir)
	if err != nil {
		log.Fatalf("Failed to connect to WireGuard: %v", err)
	}

	fmt.Println("Connected successfully!")
	fmt.Println("Press Ctrl+C to disconnect")

	// Keep the program running
	select {}
}

func loadNodes(filePath string) ([]Node, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func selectFastestNode(nodes []Node) *Node {
	var fastestNode *Node
	for _, node := range nodes {
		if fastestNode == nil || node.Ping < fastestNode.Ping {
			fastestNode = &node
		}
	}
	return fastestNode
}

func generateWireGuardConfig(node *Node, configDir string) error {
	// Generate private key
	privateKey, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Generate public key
	echoCmd := exec.Command("echo", string(privateKey))
	publicKeyCmd := exec.Command("wg", "pubkey")
	stdin, err := echoCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %v", err)
	}
	publicKeyCmd.Stdin = stdin
	publicKey, err := publicKeyCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate public key: %v", err)
	}

	// Create WireGuard configuration
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = 10.0.0.2/24
DNS = 1.1.1.1

[Peer]
PublicKey = %s
Endpoint = %s:%d
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
`, string(privateKey), string(publicKey), node.IP, node.Port)

	configPath := filepath.Join(configDir, "wg0.conf")
	if err := ioutil.WriteFile(configPath, []byte(config), 0600); err != nil {
		return fmt.Errorf("failed to write WireGuard config: %v", err)
	}

	return nil
}

func connectWireGuard(configDir string) error {
	// Check if WireGuard is installed
	if _, err := exec.LookPath("wg-quick"); err != nil {
		return fmt.Errorf("WireGuard is not installed: %v", err)
	}

	// Start WireGuard
	cmd := exec.Command("sudo", "wg-quick", "up", filepath.Join(configDir, "wg0.conf"))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start WireGuard: %v", err)
	}

	return nil
}
