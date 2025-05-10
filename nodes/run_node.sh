#!/bin/bash

# Exit on error
set -e

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root"
    exit 1
fi

# Install required packages
if command -v apt-get &> /dev/null; then
    apt-get update
    apt-get install -y wireguard qrencode
elif command -v yum &> /dev/null; then
    yum install -y wireguard-tools qrencode
else
    echo "Unsupported package manager"
    exit 1
fi

# Generate WireGuard keys
WG_PRIVATE_KEY=$(wg genkey)
WG_PUBLIC_KEY=$(echo "$WG_PRIVATE_KEY" | wg pubkey)

# Create WireGuard configuration
mkdir -p /etc/wireguard
cat > /etc/wireguard/wg0.conf << EOF
[Interface]
PrivateKey = $WG_PRIVATE_KEY
Address = 10.0.0.1/24
ListenPort = 51820
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
EOF

# Enable IP forwarding
echo "net.ipv4.ip_forward = 1" > /etc/sysctl.d/99-wireguard.conf
sysctl -p /etc/sysctl.d/99-wireguard.conf

# Start WireGuard
systemctl enable wg-quick@wg0
systemctl start wg-quick@wg0

# Get public IP
PUBLIC_IP=$(curl -s ifconfig.me)

# Create node info file
mkdir -p ../bootstrap
cat > ../bootstrap/node_info.json << EOF
{
    "ip": "$PUBLIC_IP",
    "port": 51820,
    "public_key": "$WG_PUBLIC_KEY",
    "country": "$(curl -s ipinfo.io/country)",
    "ping": 0
}
EOF

echo "Node setup complete!"
echo "Public Key: $WG_PUBLIC_KEY"
echo "Node information saved to bootstrap/node_info.json"