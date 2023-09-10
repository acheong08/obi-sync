# Reversed Obsidian Sync (Unofficial)

![GitHub release (with filter)](https://img.shields.io/github/v/release/acheong08/obi-sync?label=Sync%20Server) ![GitHub release (with filter)](https://img.shields.io/github/v/release/acheong08/rev-obsidian-sync-plugin?label=Plugin&color=%23f51159)

This is an unofficial Obsidian Sync library which allows you to host your own server for [syncing](https://obsidian.md/sync) and [publishing](https://obsidian.md/publish) obsidian notes. Behind the scenes, it re-routes Sync and Publish tasks from the official Obsidian servers to your own self-hosted server.

This system consists of two parts:
1. Self-hosted Server | [Repo](https://github.com/acheong08/obi-sync)
2. Obsidian Plugin | [Repo](https://github.com/acheong08/rev-obsidian-sync-plugin)

You will first need to configure the self-hosted server which will be responsible for storing all of your Obsidian.md vault plugins and files. Then you will need to install the associated plugin in your Obsidian.md software which will be responsible for telling Obsidian to use your servers instead of the official Obsidian servers.

<br /><br />

Huge thanks to @Aetherinox for writing the comprehensive documentation