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

if ! command -v go &> /dev/null; then
    error "Go is not installed or not in PATH. Please install Go first."
    exit 1
fi

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$PROJECT_ROOT/build"
BIN_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/harbrixd.service"

info "Building harbrix release binaries..."

go mod tidy

mkdir -p "$BUILD_DIR"
go build -trimpath -ldflags="-s -w" -o "$BUILD_DIR/harbrix"  "$PROJECT_ROOT/cmd/harbrix"
go build -trimpath -ldflags="-s -w" -o "$BUILD_DIR/harbrixd" "$PROJECT_ROOT/cmd/harbrixd"

info "Installing binaries to $BIN_DIR..."

sudo install -m 755 "$BUILD_DIR/harbrix"  "$BIN_DIR/harbrix"
sudo install -m 755 "$BUILD_DIR/harbrixd" "$BIN_DIR/harbrixd"

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