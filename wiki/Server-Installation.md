## Server Installation
Follow the instructions below to set up obsidian sync with your [Obsidian.md](https://obsidian.md/) program.

<br />

### Create DNS Records / SSL
If you are wanting to host your sync server outside of `localhost`, you can use [Cloudflare](https://cloudflare.com/) to take care of your SSL certificate needs. If you do not want to use Cloudflare but still require an SSL certificate, you can skip the instructions below and utilize [Certbot](https://certbot.eff.org/).

Access your Domain or Cloudflare control panel and create new records for two new subdomains.

The `Content` value should be the `IP Address` your self-hosted sync server will be hosted on.

![NNOxQSI](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/2cf070f5-4ee5-4733-9e6e-db2770fb4599)

<br />

Each subdomain plays the following roles:

<br />

| Subdomain | Purpose |
| --- | --- |
| `api.domain.com` | <br />Used for `Obsidian Sync`. <br /><br />This URL needs to be plugged into your `docker-compose.yml`, and set up in the settings of the [Unofficial obi-sync plugin](https://github.com/acheong08/rev-obsidian-sync-plugin) as the `Obsidian Sync URL`<br /><br /> |
| `publish.domain.com` | <br />Used for `Obsidian Publish`. <br /><br />Configure your self-hosted server, enable the `Obsidian Plugin`; then upload your vault files using `Obsidian Publish`.<br /><br />View your uploaded files at: <br />`https://publish.domain.com/vaultname/Path/To/File.md`<br /><br /> |

<br />

If you want Cloudflare to handle the SSL certificate, set each record's `Proxy Status` to ![cloudflare-enabled](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/34e80e5f-6751-4741-a063-683cf51948a9) **Proxied**

<br /><br /><br />

### Docker-Compose (Option 1)
If using this option, you must decide which setup you would like to use. 
- One File: `docker-compose.yml`
- Two Fles: `docker-compose.yml` and `.env`

<br /><br />

Using the `two file` option is slightly more secure in regards to your `SIGNUP_KEY` being exposed to logs. 

- If you are not concerned about security, use [DOCKER-COMPOSE.YML ONLY](#-docker-composeyml-only) option.
- If you want the extra security, use [DOCKER-COMPOSE.YML + .ENV](#-docker-composeyml--env) option.

<br /><br />

---

<br /><br />

##### 🟢 DOCKER-COMPOSE.YML ONLY
Create a new file named `docker-compose.yml`, and then paste the corresponding code provided below:
```yml
version: '3.8'

services:
  obsidian_sync:
    image: ghcr.io/acheong08/obi-sync:latest
    container_name: obsidian_sync
    restart: always
    ports:
      - 3000:3000
    environment:
      - DOMAIN_NAME=api.domain.com
      - ADDR_HTTP=0.0.0.0:3000
      - SIGNUP_KEY=YOUR_PRIVATE_STRING_HERE
      - DATA_DIR=/obi-sync/
      - MAX_STORAGE_GB=10
      - MAX_SITES_PER_USER=5
    volumes:
      - ./obi-sync:/obi-sync

volumes:
  obi-sync:
```

<br />

<sub>[🔹Continue With Installation](#Variable-List)</sub>

<br /><br />

---

<br /><br />

##### 🟢 DOCKER-COMPOSE.YML + .ENV
Create the following two files in the same folder, and then paste the corresponding code provided below:
- [docker-compose.yml](#-docker-composeyml)
- [.env](#-env)

<br />

#### 📄 docker-compose.yml
```yml

version: '3.8'
services:
  obsidian_sync:
    image: ghcr.io/acheong08/obi-sync:latest
    container_name: obsidian_sync
    restart: always
    ports:
      - 3000:3000
    environment:
      DOMAIN_NAME: ${DOMAIN}
      ADDR_HTTP: ${ADDR_HTTP}
      SIGNUP_KEY: ${USER_SIGNUP_KEY}
      DATA_DIR: ${DIR_DATA}
      MAX_STORAGE_GB: ${USER_MAX_STORAGE}
      MAX_SITES_PER_USER: ${USER_MAX_SITES}
    volumes:
      - ./obi-sync:/obi-sync

volumes:
  obi-sync:
```

#### 📄 .env
```env
DOMAIN='api.domain.com'
ADDR_HTTP='0.0.0.0:3000'
DIR_DATA='/obi-sync/'
USER_SIGNUP_KEY='YOUR_PRIVATE_STRING_HERE'
USER_MAX_STORAGE=10
USER_MAX_SITES=5
```

<br />

<sub>[🔹Continue With Installation](#Variable-List)</sub>

<br />

---

<br />

#### Variable List

After creating `docker-compose.yml` / `.env` file(s), edit the variables to match your domain `DOMAIN_NAME`, registration signup key `SIGNUP_KEY`, etc. 
<br /><br />
A description of each variable is provided below:

<br />

| Variable | Description | Required | Default |
| --- | --- | --- | --- | 
| `DOMAIN_NAME` | This is the URL to your API subdomain | Yes | `localhost:3000` |
| `ADDR_HTTP` | The address to run Obi-Sync on | No | `127.0.0.1:3000` |
| `SIGNUP_KEY` | Required later when creating users who will be able to access your self-hosted server | No | None |
| `DATA_DIR` | Where encrypted `vault.db` and other files will be stored | No | `./` |
| `MAX_STORAGE_GB` | The maximum storage per user in GB. | No | `10` |
| `MAX_SITES_PER_USER` | The maximum number of sites per user. | No | `5` |

<br /><br />

#### Start/Stop Docker

After you finish configuring your docker's variables, you need to `cd` to the folder where you created the file(s) and execute:
```shell
docker compose up -d
```

To shut down the container:
```shell
docker compose down
```

To confirm container is running:
```shell
docker ps
```

![9XdXESs](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/81f54b21-b4e5-4541-afb7-212cd227053c)


<br /><br /><br />

---

### Docker (Option 2)
To install obi-sync using docker `pull`, you can run the default command below:

```shell
docker pull ghcr.io/acheong08/obi-sync:latest
```
<br />

For a full list of pull commands based on your operating system, view [Package Release](https://github.com/acheong08/obi-sync/pkgs/container/obi-sync).


<br /><br />

---

<br /><br />

### Nginx Configuration
You must configure nginx to handle your new subdomains by adding a few entries to your nginx config file. If you wish to add this to your current running nginx webserver, you can paste it somewhere in `\etc\nginx\sites-enabled\default`.

<br />
<br />

> [!NOTE]
> Do not copy/paste the server blocks below without looking them over. If you have changed port `3000` in your docker container, you must change it below to that same port. Also change `domain.com` to your correct domain.

<br />
<br />

```nginx
map $http_upgrade $connection_upgrade {
        default upgrade;
        '' close;
}

#
#   Obsidian Sync
#
#   defined within the obi-sync plugin for Obsidian so that Obsidian and
#   the self-hosted server can sync files.
#

    server
    {
        listen 443 ssl http2;
        listen [::]:443 ssl http2;
        listen 80;
        listen [::]:80;

        server_name         	www.api.domain.com api.domain.com;

        location /
        {
            proxy_http_version  1.1;
            proxy_set_header    Upgrade $http_upgrade;
            proxy_set_header    Connection $connection_upgrade;
            proxy_set_header    Host $host;
            proxy_pass          http://127.0.0.1:3000/;
        }
        
    }

#
#   Obsidian Publish
#
#   used for viewing published Obsidian .md files on a webserver.
#

    server
    {
        listen 443 ssl http2;
        listen [::]:443 ssl http2;
        listen 80;
        listen [::]:80;

        server_name         	www.publish.domain.com publish.domain.com;
        
        location /
        {
            proxy_pass     	http://127.0.0.1:3000/published/;
        }
    }
```
<br />

After modifying your nginx service, **restart** it to load up the new configs.
```shell
sudo systemctl restart nginx
```
Or restart nginx via docker if you run it through a container.

<br /><br />

---

<br /><br />

### Creating New User
You must create a user that will be allowed to access your self-hosted server from within Obsidian.md.
This can be done by opening Powershell in Windows, or Terminal in Linux and executing the following:

<br />

> [!NOTE]
> `email`: This is the email address you will use in Obsidian's `About` tab to sign into your self-hosted server.
> 
> `password`: Pick any password you wish to use. This is the password you will use in the Obsidian.md program once you configure the plugin to connect to your self-hosted server.
>
> `name`: Can be anything, not super important.
>
> `signup_key`: Variable you provided in `docker-compose.yml` or `.env` file.<br />
> If you removed `signup_key` from your docker container's variables and don't want to require a signup key for registration, remove that line from the command below.

<br />

#### ⚙️ Windows `Powershell`
```powershell
curl --request POST `
  --url https://api.domain.com/user/signup `
  --header 'Content-Type: application/json' `
  --data '{
	"email": "example@example.com",
	"password": "example_password",
	"name": "Example User",
	"signup_key": "<SIGNUP_KEY>"
}'
```

#### ⚙️ Linux `Terminal`
```bash
curl --request POST \
  --url https://api.domain.com/user/signup \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "example@example.com",
	"password": "example_password",
	"name": "Example User",
	"signup_key": "<SIGNUP_KEY>"
}'
```

<br />

A **successful** registration will return the following response:
```json
{"email":"example@example.com","name":"Example User"}
```

<br />

A **failed** registration will return the following response:
> [!WARNING]
> If you receive a **failure** from registration, make sure you're not trying to sign up the same user multiple times.
```json
{"error":"not sign up!"}
```

<br />

<sub>[🔝 Top](#Index)</sub>

<br />

---

<br />

Huge thanks to @Aetherinox for writing the comprehensive documentation