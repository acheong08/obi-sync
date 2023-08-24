# Obsidian Sync

> [!WARNING]
> This is an experimental proof of concept. It was written hastily without knowledge of the real internal mechanisms. Expect a thousand bugs and inefficiencies

> [!NOTE]
> If you have the time and energy, feel free to help out with PRs or suggestions.

## Setup
- `git clone https://github.com/acheong08/obsidian-sync`
- `cd obsidian-sync`
- `export HOST=<YOUR DOMAIN NAME>`
- `go run cmd/obsidian-sync/main.go`
- Use nginx or cloudflare to proxy & handle TLS/SSL

Signup is currently manual. You edit the database yourself. I will add a tool *soon™*

## Usage
As ObsidianMD is written with Electron, you can unzip the resource pack and modify it to suite your needs

- Download [Obsidian](https://github.com/obsidianmd/obsidian-releases/releases/)
- `tar -xvzf obsidian-1.3.7.tar.gz`
- `cd obsidian-1.3.7/resources`
- `npx asar extract obsidian.asar obsidian`
- `sed -i 's|api.obsidian.md|<YOUR DOMAIN NAME>|g' obsidian/starter.js` (Remember to replace https with http if running on localhost)
- `sed -i 's|api.obsidian.md|<YOUR DOMAIN NAME>|g' obsidian/app.js`
- `npx asar pack obsidian obsidian.asar`
- Run the binary
