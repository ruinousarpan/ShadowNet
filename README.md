# ShadowNet

A decentralized, open-access, cross-platform VPN network that allows anyone to run a node and connect to the network.

## Features

- Free and open access
- No login required
- Decentralized node discovery
- Automatic node selection based on performance
- Built on WireGuard for security and performance
- GitHub integration test

## Quick Start

### Running a Node

1. Clone this repository
2. Run `./nodes/run_node.sh` to start a public node
3. Your node will be automatically added to the network

### Connecting as a Client

1. Install the ShadowNet client:
   ```bash
   go install github.com/yourusername/shadownet/client@latest
   ```
2. Run the client:
   ```bash
   cd client
   go build
   sudo ./client
   ```

## Project Structure

```
shadownet/
├── client/          # Client application
├── nodes/           # Node configuration and scripts
├── bootstrap/       # Bootstrap node information
└── docs/           # Documentation
```

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## License

MIT License - See LICENSE file for details 