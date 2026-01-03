# harbrix

A lightweight user-isolated service supervisor for Linux.

harbrix allows each user to define and manage their own background services, running them securely under their own user privileges. A single central daemon (harbrixd) runs as root and supervises services for all users. The `harbrix` CLI lets users control only their own services.

Services are configured in simple TOML files located in `~/.local/share/harbrix/services/`. Enabled services are automatically started at system boot.

## Why harbrix

- No need for root privileges to manage personal services.
- Services run isolated under the owning user's UID/GID.
- Simpler and lighter than full system supervisors for personal/long-running processes.
- Supports restart policies and optional logging.

## Features

- Per-user service isolation
- TOML-based service definitions
- Restart policies: `never`, `on-failure`, `always`
- Optional stdout/stderr logging
- Auto-start on boot for enabled services
- Commands: `list`, `start`, `stop`, `restart`, `status`, `log`, `enable`, `disable`, `new`, `edit`, `delete`, `reload`

## Dev/Build Requirements

- go >= 1.25.5
- Linux (?)

## Installation


### Download and Install (recommended)

Run the following command in your shell:

```bash
curl -fsSL https://raw.githubusercontent.com/mhs003/harbrix/main/net-install.sh | sudo bash
```
This command downloads the latest release binaries for your CPU architecture and installs them on your system.

### Build and Install

```bash
./install.sh
```

This builds the binaries, installs them to `/usr/local/bin`, sets up the systemd service for `harbrixd`, enables and starts it.

---

### Manual build and install *(test mode)*

```bash
sudo make install   # build and Install to /usr/local/bin
```

The daemon must run as root:

```bash
sudo harbrixd &
```

For production, use the provided `install.sh` or `net-install.sh` to manage the systemd service.

## Uninstallation

```bash
sudo ./uninstall.sh
```

Removes binaries and the systemd service. User data in `~/.local/share/harbrix` is preserved.

## Usage

### Create a service

```bash
harbrix new myservice
```

### Edit the generated service:

```bash
harbrix edit myservice
```

service files are simple:

```toml
name = "myservice"
description = "Example service" # optional
author = "yourusername"         # optional

[service]
command = "your-command-here"
workdir = "/optional/path"      # optional; default=~/.local/share/harbrix
log = true                      # optional; default=false

[restart]
policy="always" # always, on-failure, never; default=never
delay="3s"      # optional; default=0s
limit=5         # optional
maxfailed=3     # optional; default=5
```

### Common commands

```bash
harbrix list                      # List your services
harbrix start myservice
harbrix stop myservice
harbrix restart myservice
harbrix status myservice
harbrix log myservice             # View logs
harbrix log -f myservice          # Follow logs
harbrix enable myservice          # Auto-start on boot
harbrix enable myservice --now    # Enable and start now
harbrix disable myservice
harbrix delete myservice          # Only if stopped and disabled
harbrix reload                    # Reload service definitions
```

## Logs

Daemon logs:

```bash
sudo journalctl -u harbrixd.service -f
```

Service logs (if `log = true`):

```bash
harbrix log myservice
```

## FAQ

<details>
    <summary>With so many supervisor tools out there, why build another one?</summary>
    
    No deep reason—just laziness.
    
    I'm too forgetful to keep track of where systemd wants its unit files, and I got tired of sudo-ing just to tweak a personal service.
    
    PM2 is nice, but dragging in Node.js for this kind of stuffs? Hard pass.
    
    So I made harbrix to scratch my own itch.
</details>

<details>
    <summary>Can multiple users run services at the same time?</summary>
    
    Yep. One daemon rules them all (as root), but each user's services run under their own UID/GID with full isolation.
</details>

<details>
    <summary>Is this thing safe for production?</summary>
    
    I'm running it in production myself, so it's good enough for me.
    
    That said, it hasn't undergone extensive third-party testing, so you can use it at your own discretion if you want to.
</details>

## TODO

- [X] Service restart limitation
- [X] Restart delay time
- [ ] Environment variables in service files
- [ ] Reimplement CLI in better approach
- [ ] Better CLI Outputs

## License

MIT (see [LICENSE](https://github.com/mhs003/harbrix/blob/main/LICENSE) file)
