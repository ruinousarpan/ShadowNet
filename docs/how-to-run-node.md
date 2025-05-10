# How to Run a ShadowNet Node

This guide will help you set up and run a ShadowNet node on your Linux server.

## Prerequisites

- A Linux server with a public IP address
- Root access
- WireGuard support in the kernel
- Basic networking knowledge

## Installation Steps

1. Clone the ShadowNet repository:
   ```bash
   git clone https://github.com/yourusername/shadownet.git
   cd shadownet
   ```

2. Make the node script executable:
   ```bash
   chmod +x nodes/run_node.sh
   ```

3. Run the node setup script as root:
   ```bash
   sudo ./nodes/run_node.sh
   ```

The script will:
- Install WireGuard if not already installed
- Generate WireGuard keys
- Configure the WireGuard interface
- Enable IP forwarding
- Start the WireGuard service
- Create a node information file

## Verifying the Node

After installation, you can verify that your node is running:

1. Check WireGuard status:
   ```bash
   sudo wg show
   ```

2. Check the node information:
   ```bash
   cat bootstrap/node_info.json
   ```

## Security Considerations

- Keep your WireGuard private key secure
- Regularly update your system and WireGuard
- Monitor your server's resources and network usage
- Consider implementing rate limiting if needed

## Troubleshooting

If you encounter issues:

1. Check WireGuard logs:
   ```bash
   journalctl -u wg-quick@wg0
   ```

2. Verify IP forwarding is enabled:
   ```bash
   sysctl net.ipv4.ip_forward
   ```

3. Check firewall rules:
   ```bash
   iptables -L
   ```

## Contributing

If you find any issues or have suggestions for improvement, please open an issue or submit a pull request on GitHub. 