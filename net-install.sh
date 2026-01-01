#!/bin/bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC}    $*";
}
warn() {
    echo -e "${YELLOW}[WARN]${NC}    $*";
}
error() {
    echo -e "${RED}[ERROR]${NC}   $*" >&2;
}

VERSION_FILE="https://raw.githubusercontent.com/mhs003/harbrix/refs/heads/main/version"
BIN_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/harbrixd.service"

command -v curl >/dev/null 2>&1 || {
    error "curl is required"
    exit 1
}

case "$(uname -m)" in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l|armv6l) ARCH="arm" ;;
    i386|i686) ARCH="386" ;;
    *)
        error "Unsupported architecture: $(uname -m)"
        exit 1
        ;;
esac

info "Detected architecture: $ARCH"


info "Fetching latest version..."
VERSION="$(curl -fsSL "$VERSION_FILE")"

if [[ -z "$VERSION" ]]; then
    error "Failed to fetch version"
    exit 1
fi

info "Latest version: v$VERSION"

CLI_URL="https://github.com/mhs003/harbrix/releases/download/v${VERSION}/harbrix-${ARCH}"
DAEMON_URL="https://github.com/mhs003/harbrix/releases/download/v${VERSION}/harbrixd-${ARCH}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

info "Downloading binaries..."

curl -fL "$CLI_URL" -o "$TMP_DIR/harbrix"
curl -fL "$DAEMON_URL" -o "$TMP_DIR/harbrixd"

chmod +x "$TMP_DIR/harbrix" "$TMP_DIR/harbrixd"

info "Installing binaries to $BIN_DIR..."

sudo install -m 755 "$TMP_DIR/harbrix"  "$BIN_DIR/harbrix"
sudo install -m 755 "$TMP_DIR/harbrixd" "$BIN_DIR/harbrixd"


info "Installing systemd service..."

sudo tee "$SERVICE_FILE" > /dev/null <<'EOF'
[Unit]
Description=Harbrix System Daemon
Documentation=https://github.com/mhs003/harbrix
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/local/bin/harbrixd
Restart=always
RestartSec=2
KillSignal=SIGTERM
TimeoutStopSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF


info "Reloading systemd daemon..."
sudo systemctl daemon-reload

info "Enabling and starting harbrixd service..."
sudo systemctl enable harbrixd.service
sudo systemctl start harbrixd.service

info "Checking service status..."
sleep 2
if systemctl is-active --quiet harbrixd.service; then
    info "harbrixd is now running and enabled on boot."
else
    warn "harbrixd failed to start. Check status with:"
    echo "    sudo systemctl status harbrixd.service"
    echo "    sudo journalctl -u harbrixd.service -f"
    exit 1
fi

echo
info "Installation complete!"
echo
echo "You can now use 'harbrix'. Some basic commands:"
echo "    harbrix list"
echo "    harbrix new myservice"
echo "    harbrix start myservice"
echo "    harbrix help"
echo
echo "To view daemon logs:"
echo "    sudo journalctl -u harbrixd.service -f"

exit 0