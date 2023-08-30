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

- N/A

## To do

- Fix bugs
- Publish

## Setup

[Quickstart with Docker](https://github.com/acheong08/rev-obsidian-sync/wiki/Docker-Compose)

### Environment variables
#### Required:
- `DOMAIN_NAME` - The domain name or IP address of your server. Include port if not on 80 or 433. The default is `localhost:3000`
#### Optional
- `ADDR_HTTP` - Server listener address. The default is `127.0.0.1:3000`
- `SIGNUP_KEY` - Signup API is at `/user/signup`. This optionally restricts users who can sign up.
- `DATA_DIR` - Where data is saved. Default `.`

### Building & Running

- `git clone https://github.com/acheong08/obsidian-sync`
- `cd obsidian-sync`
- `go run cmd/obsidian-sync/main.go`

Optional:
- Configure [nginx](https://github.com/acheong08/rev-obsidian-sync/wiki/Nginx-Configuration)
- HTTPS is recommended.

When you're done, install and configure the [plugin](https://github.com/acheong08/rev-obsidian-sync-plugin)

## Adding a new user

`go run cmd/signup/main.go`

Alternatively:
```bash
curl --request POST \
  --url https://yourdomain.com/user/signup \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "example@example.com",
	"password": "example_password",
	"name": "Example User",
	"signup_key": "<SIGNUP_KEY>"
}'
```
You can set the signup key via the `SIGNUP_KEY` environment variable. If it has not been set, you can exclude it from the request.
