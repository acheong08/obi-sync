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

- Obsidian publish

## To do

- Fix bugs
- Publish

## Setup

<details>

<summary>

### Manual building
</summary>

- `git clone https://github.com/acheong08/obsidian-sync`
- `cd obsidian-sync`
- `export HOST=<YOUR DOMAIN NAME>` - Not necessary when running on localhost
- `go run cmd/obsidian-sync/main.go`

</details>

<details>

<summary>

### Docker
</summary>

https://github.com/acheong08/rev-obsidian-sync/pkgs/container/rev-obsidian-sync

</details>

HTTPS is not required.

When you're done, configure the [plugin](#sync-override-plugin)

<details>
<summary>
	
### Nginx configuration
</summary>

```nginx
map $http_upgrade $connection_upgrade {
        default upgrade;
        '' close;
}
server {
	listen 80;
	listen [::]:80;
	location / {
		proxy_http_version 1.1;
            	proxy_set_header Upgrade $http_upgrade;
            	proxy_set_header Connection $connection_upgrade;
           	proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:3000/;
	}
	server_name <your domain name>; # e.g. api.obsidian.md
}
# This is for obsidian publish (Optional)
server {
	listen 80;
	listen [::]:80;
	location / {
		proxy_http_version 1.1;
            	proxy_set_header Upgrade $http_upgrade;
            	proxy_set_header Connection $connection_upgrade;
           	proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:3000/published/;
	}
	server_name <another domain name>; # e.g. publish.obsidian.md
}
```

You can use `certbot` or cloudflare to handle HTTPS although it is not mandatory.

</details>

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
  - e.g. `https://obsidian.yourdomain.com`
  - For development: `http://127.0.0.1:3000`
