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

	// Fetch available nodes
	nodes, err := fetchNodes()
	if err != nil {
		log.Fatalf("Failed to fetch nodes: %v", err)
	}

	// Find the fastest node
	fastestNode := findFastestNode(nodes)
	if fastestNode == nil {
		log.Fatal("No available nodes found")
	}

	fmt.Printf("Connecting to node in %s (ping: %dms)\n", fastestNode.Country, fastestNode.Ping)

	// Generate WireGuard configuration
	if err := generateWireGuardConfig(fastestNode, configDir); err != nil {
		log.Fatalf("Failed to generate WireGuard config: %v", err)
	}

	// Start WireGuard connection
	if err := startWireGuard(configDir); err != nil {
		log.Fatalf("Failed to start WireGuard: %v", err)
	}

	fmt.Println("Connected to ShadowNet!")
	fmt.Println("Press Ctrl+C to disconnect")

	// Keep the program running
	select {}
}

func fetchNodes() ([]Node, error) {
	// In a real implementation, this would fetch from multiple sources
	// For now, we'll just read from bootstrap.json
	data, err := ioutil.ReadFile("../bootstrap/bootstrap.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read bootstrap.json: %v", err)
	}

	var nodes []Node
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("failed to parse bootstrap.json: %v", err)
	}

	return nodes, nil
}

func findFastestNode(nodes []Node) *Node {
	if len(nodes) == 0 {
		return nil
	}

	fastest := nodes[0]
	for _, node := range nodes {
		if node.Ping < fastest.Ping {
			fastest = node
		}
	}
	return &fastest
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

func startWireGuard(configDir string) error {
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
