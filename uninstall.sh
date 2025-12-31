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

if [[ $EUID -ne 0 ]]; then
   error "This script must be run as root (use sudo)."
   exit 1
fi

BIN_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/harbrixd.service"

info "Stopping harbrixd service..."
if systemctl is-active --quiet harbrixd.service; then
    systemctl stop harbrixd.service
fi

info "Disabling harbrixd service..."
if systemctl is-enabled --quiet harbrixd.service; then
    systemctl disable harbrixd.service
fi

info "Removing systemd service file..."
rm -f "$SERVICE_FILE"

info "Reloading systemd daemon..."
systemctl daemon-reload
systemctl reset-failed

info "Removing binaries..."
rm -f "$BIN_DIR/harbrix"
rm -f "$BIN_DIR/harbrixd"

echo
info "Uninstallation complete!"
echo
echo "harbrix and harbrixd have been removed from the system."
echo "User service files and logs in ~/.local/share/harbrix remain untouched."

exit 0