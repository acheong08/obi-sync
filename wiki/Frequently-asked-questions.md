## FAQ

### Where Are Files Stored?
This depends on how you've configured obsidian sync server to run.
If you installed this project using [Install with Docker-compose](Server-Installation.md#docker-compose-option-1), and specified the environment variable `DATA_DIR`, then the files should exist in the same folder you specified the variable to use. By default, vault files are placed in the folder `obi-sync`.

If you installed this project using [Install with Docker](Server-Installation.md#docker-option-2) or did not specify a variable for `DATA_DIR`, then the files are typically stored in `/var/lib/docker/*`

Some users may be running [Portainer](https://www.portainer.io/), which allows you to view your docker containers and volumes within a graphical user interface. Login to the portainer web admin panel and under `Volumes`, find out which volume is assigned to `obsidian_sync`, and copy the value `Mount path`. Open your File Explorer and go to that location to view your vault files.

<br />

---

<br />

### What Files Are Created On Server?
This project will store your vault data in the following files:
```
📁 obi-sync
   📄 publish.db
   📄 secret.gob
   📄 vault.db
   📄 vaultfiles.db
```

<br />

---

<br />

### Error: User Not Signed Up
This error typically occurs when you're creating a user with `cURL` and shows the following:
```json
{"error":"not sign up!"}
```
<br />

If you receive this error when creating your first API user, ensure the user did not previously exist.

<br />

---

<br />

### Network Error Occured. Network Unavailable
If you are attempting to use the `Obsidian Publish` feature and receive the error:

![publish-network-error](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/6094db18-a523-40d2-891c-f59d2e868556)

Ensure you have provided the correct environment variable `DOMAIN_NAME=api.domain.com`.
If you provided the wrong domain name and need to change it, you must do the following:
- Edit `docker-compose.yml`
- Change `DOMAIN_NAME` variable to correct domain.
- Locate the project file `publish.db` and DELETE it completely.
- Restart the docker container and Obsidian.md program

Once back in Obsidian.md program, click `Publish Changes` ![Publish Changes Button](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e9cd7054-0a41-472f-accb-d9fa0426436d) button again.

You can also attempt to locate the root cause of the issue by pressing `CTRL + SHIFT + I` inside of Obsidian.md. Then on the right side, click `Network` and then re-open the `Publish Changes` interface again.

![Publish Network inspect](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/8f462c26-503b-4cc5-a2f5-3272ee7c2ff6)

Anything listed in `red` is an error and needs to be corrected. Ensure that it is trying to fetch the correct domain from the `Request URL`. Anything listed in `white` and with a `Status Code` of `200` is properly connecting.

<br />

---

<br />

### Error: No Matching Manifest linux/amd64
Make sure you are pulling the correct docker image.

You can also visit the [Package Release](https://github.com/acheong08/obi-sync/pkgs/container/obi-sync) page, select `OS / Arch` tab, and copy the pull command listed below `linux/amd64`

<br />

---

<br />

### Docker-compose.yml vs .Env File
In the section above titled [Install with Docker-compose](Server-Installation#docker-compose-option-1), there are two ways to install obi-sync using `docker-compose`.
1. Single `docker-compose.yml` file
   - See [DOCKER-COMPOSE.YML ONLY](Server-Installation#-docker-composeyml-only)
3. `docker-compose.yml` and `.env` file
   - See [DOCKER-COMPOSE.YML + .ENV](Server-Installation#-docker-composeyml--env)

#### 🟢 Single `docker-compose.yml` File
This method involves creating a single `docker-compose.yml` file which will hold all of your settings for this project.

#### 🟢 `docker-compose.yml` and `.env` File
This method requires you to create two files which both exist in the same folder. The benefit of this method is that your `SIGNUP_KEY` environment variable will not be leaked in your docker logs and is slightly better if you are worried about security.

You can decide to use either one of the two options.

If using `Portainer` web manager for Docker, you can access the environment variables clicking `Containers` on left side menu, then click `obsidian-sync` in your list of containers. Then scroll down to the `ENV` section.

![PZkqLFM](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/50ed9c37-d17c-4b1a-a67a-8f2d404484e1)

<br />

---

<br />

### Can't Find `.obsidian` Folder or `.Env` File
The `.obsidian` plugin folder and docker `.env` file may be hidden. You can configure your operating system's File Explorer show hidden files or use a terminal.

`Linux`: By default, this OS hides all folders and files that start with a period. To view these files in your File Explorer, press the key combination `CTRL + H`. You may be asked for your account or sudo password.

`iOS`: You may need to plug the device into a computer, then to access your `.obsidian` folder, go to `Documents on iPhone`, then go to `Obsidian/<your vault name>/.obsidian`.

<br />

---

<br />

### Privacy & Security
The plugin and self-hosted server are dead simple. No data is collected. You can even run your sync server on the local network without access to the internet. All vault data is stored behind an encrypted password and inside a database file.

<br />

Huge thanks to @Aetherinox for writing the comprehensive documentation

---

<sub>[🔝 Top](#faq)</sub>