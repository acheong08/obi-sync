# Rev Obsidian Sync

Reverse engineered obsidian sync server (NOT OFFICIAL)

> [!WARNING]
> The main branch is the development branch. For stable usage, use the latest release.

> [!NOTE]
> If you have the time and energy, feel free to help out with PRs or suggestions.

## Features

- End to end encryption
- Live sync (across devices)
- File history/recovery/snapshots
- Works natively on IOS/Android/Linux/MacOS/Windows... (via the plugin)
- Vault sharing

### Experimental

These features are not in the latest release but in the main branch. They might not be fully tested and are probably unstable.

- <N/A>

## To do

- Fix bugs
- Publish

## Setup

- `git clone https://github.com/acheong08/obsidian-sync`
- `cd obsidian-sync`
- `export HOST=<YOUR DOMAIN NAME>`
- `go run cmd/obsidian-sync/main.go`
- Use nginx or cloudflare to proxy & handle TLS/SSL

### Nginx configuration

```nginx
map $http_upgrade $connection_upgrade {
        default upgrade;
        '' close;
}
server {
	listen 80 default_server;
	listen [::]:80 default_server;
	location / {
		proxy_http_version 1.1;
            	proxy_set_header Upgrade $http_upgrade;
            	proxy_set_header Connection $connection_upgrade;
           	proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:3000/;
	}
	server_name _;
}
```

HTTPS _should_ be required. I use `certbot` or Cloudflare

## Adding a new user

`go run cmd/signup/main.go`

## Sync override plugin

Tested on

- IOS
- Linux (Flatpak)

### Usage

> While we have no qualms with reverse engineering as a playground for experimentation, Obsidian Sync is a service we intend to keep first-party only for the foreseeable future. - https://github.com/obsidianmd/obsidian-releases/pull/2353

This plugin will not be part of the official community plugins list.

- Install https://github.com/acheong08/rev-obsidian-sync-plugin
- Go to settings
- Set API endpoint

Known bugs:

- ~~Cannot restart plugin (for whatever reason you might want to do that...) - Restart the app if you want to reload this particular plugin~~

Report all bugs in this repository.
