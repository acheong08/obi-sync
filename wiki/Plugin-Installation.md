## Plugin Installation

In order for your new self-hosted Publish and Sync server to function properly, you must install a plugin to your copy of Obsidian.md

<br /><br />

### Install with WGET

Navigate to your vault's `.obsidian` folder:

```shell
cd /path/to/vault/.obsidian
```

Create a new folder for the plugin and enter the folder:

```shell
mkdir -p plugins/custom-sync-plugin && cd plugins/custom-sync-plugin
```

Use `wget` to download the plugin files:

```shell
wget https://github.com/acheong08/rev-obsidian-sync-plugin/raw/master/main.js https://github.com/acheong08/rev-obsidian-sync-plugin/raw/master/manifest.json
```

<br /><br />

### Install Manually

To manually get the latest copy of the unofficial Obi-Sync Plugin, [download here](https://github.com/acheong08/rev-obsidian-sync-plugin/releases).

- Navigate to the folder where your vault is on your local machine, and enter the `.obsidian\plugins\` folder.
- Create a new folder: `custom-sync-plugin`
- Inside the `custom-sync-plugin` folder, install the plugin files:
  - 📄 `main.js`
  - 📄 `manifest.json`

<br />

> [!NOTE]
> Alternatively, you can use https://github.com/TfTHacker/obsidian42-brat which can be found in the official community plugins list.

<br /><br />

### Enable Plugin

Once the plugin is installed, activate it by launching Obsidian.md.

- Open `Obsidian Settings` ![uJ5MSWk](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/f5695ae4-0730-496c-b182-3bf4836ba571)
- On left menu, click `Community Plugins`
- Scroll down to `Custom Native Sync` and enable ![f8iiGTI](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e38ac70d-60ea-4cf7-939a-ab56d5302f11)
- To the left of the enable button for `Custom Native Sync`, press the options ![j4JYNxN](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/1de68f02-17be-4759-b21b-d344579477a4) button.
- Configure the `Obsidian Sync URL` in the Plugin settings after setting up your sync server
  - By default, it is `https://api.domain.com`
- Go to the `Core Plugins` section of Obsidian
- Locate the core plugin `Sync` and enable ![f8iiGTI](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e38ac70d-60ea-4cf7-939a-ab56d5302f11)

![HzDfnB4](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/44363513-699b-49c9-84fe-94faf1785981)

![ioHy4jQ](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/61d154df-430d-4785-983c-633418607d12)

<br /><br />

### Configure Sync

Before doing these steps, ensure that your obi-sync self-hosted server is running.

Open your Obsidian.md Settings ![uJ5MSWk](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/f5695ae4-0730-496c-b182-3bf4836ba571) and select `About`.

You need to sign in and connect Obsidian.md to your self-hosted server by clicking `Log In`:

![fs9PioG](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/941f8b5a-485f-4f6b-bdd0-ad0e6ab092a5)

The login dialog will appear:

![aVD8YRs](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/5d065b57-708c-4bbe-8373-8d11354b09f5)

Fill in your email address and password that you used to register your account with the `cURL` command in the section [Creating New User](Server-Installation#creating-new-user).

```powershell
curl --request POST `
  --url https://api.domain.com/user/signup `
  --header 'Content-Type: application/json' `
  --data '{
	"email": "example@example.com",      <----- The email you will use
	"password": "example_password",	     <----- The password you will use
	"name": "Example User",
	"signup_key": "<SIGNUP_KEY>"
}'
```

<br />

> [!NOTE]
> After signing in successfully, the `Log In` button will change to `Manage settings`. Keep in mind though that clicking the `Manage Settings` button will take you to Obsidian's official website. It has nothing to do with your own self-hosted server.

<br />

Next, on the left side under `Core Plugins`, select `Sync`.

On the right side, click the `Choose` button:

![eV6KKFy](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/4a58bc85-2fb1-4bf6-882e-596a681f9385)

A new dialog will then appear. Select `Create New Vault`

![97sJzUP](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/dca34861-a62d-4714-819e-acc71b579671)

Then fill in the information for your new vault. The `Encryption Password` can be anything you want. Do **not** lose this password. You cannot unlock / decrypt your vault if you can't remember it.

![tj19O8m](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e204489c-5a03-46bd-b5ba-94402280477e)

Your new vault will be created, along with a few options you can select, and a `Connect` button.

![RDFF70c](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/c264782c-6f92-4f7c-97b5-bfb0b0c857cd)

Click the `Connect` button to start Syncing between your local vault and your self-hosted sync server. If you get a warning message labeled `Confirm Merge Vault`, click `Continue`.

If you provided an `Encryption Password` a few steps back, you will now be asked to enter that password, and then click `Unlock Vault`.

You will then get one more confirmation that your vault is now synced.

![6znaxbj](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/2dbe1e86-3b91-4cb4-a5da-fc354d13e40d)

From this point on, you can adjust any of the Sync settings you wish to modify, which includes syncing things like images, videos, plugins, settings, and more. Whatever options you enable, will be included in the Sync job.

Ensure `Sync Status` is set to `running`. You can enable / disable the sync's state by clicking the button to the right.

![UkqH5w1](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/0c786122-3ac6-45e7-843a-2aa4cc8312da)

You can confirm that your files are successfully syncing by clicking the `View` button to the right of `Sync Activity`.

![iyeEBrs](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/bc1c591b-dc75-4ae2-a3f3-ac121ffe9643)

Once pressing the `View` button, a new dialog will appear and show you the status of your Sync job.

```console
2023-09-09 01:38 - Uploading file Test Folder/Subfolder/My File 1.md
2023-09-09 01:38 - Upload complete Test Folder/Subfolder/My File 1.md
2023-09-09 01:38 - Uploading file Test Folder/Subfolder/My File 2.md
2023-09-09 01:38 - Upload complete Test Folder/Subfolder/My File 2.md
2023-09-09 01:40 - Fully synced
```

<br /><br />

### Configure Publish

These instructions help guide you through setting up `Obsidian Publish`. This feature will allow you to view your Vault .md notes on your webserver.

Open your Obsidian.md settings, and under `Options` on the left, select `Core Plugins`. In the middle of the screen, enable ![f8iiGTI](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e38ac70d-60ea-4cf7-939a-ab56d5302f11) `Publish`.

![NvmH8hV](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/794d9820-1f45-4f9c-afbc-243fd3537260)

Back on the main Obsidian interface, select the `Publish Changes` ![Publish Changes Button](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e9cd7054-0a41-472f-accb-d9fa0426436d) button in your Obsidian side / toolbar.

![iYTDAi4](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/e7f24105-432b-4436-a8a1-71862c873640)

In the new dialog, enter a `Site ID` in the textfield provided. This will become the `slug` name for your new vault.

![nofm4EB](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/4254380c-331e-4a80-8299-3204fb098cc2)

Once you have a name, click `Create`.

You will then see a list of all your vault's associated files. Select the checkbox to the left of each folder you wish to upload to your self-hosted Publish server.

![FiUou2Z](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/936fbb7f-4867-4b44-9e80-02b7221f3c62)

Once you've selected the desired folders, click `Publish` in the bottom right.

Yet another dialog will appear which confirms your uploaded vault documents.

![sTxX4Ax](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/1a2501d8-338d-455e-8829-0ef2b9439494)

<br />

> [!NOTE]
> Because we are hosting our own server, the link provided at the bottom of the dialog **will not work** since it goes to Obsidian's official publish server.

<br />

Open your web browser and go to the URL for your self-hosted publish server. In our example, we would view our `Testing Folder/Tags.md` page by entering the following URL:

```
https://publish.domain.com/myvault/Testing Folder/Tags.md
```

- `myvault`: Name given to example vault
- `Testing Folder`: Folder name the note resides in
- `Tags.md`: Note filename

To see an overview of files you have uploaded to your publish server, go to your publish subdomain, and add your vault name at the end. Nothing else needs added.

```
https://publish.domain.com/myvault
```

In our example vault, visiting the URL above displays `JSON` and includes a list of all files uploaded from the test vault:

![uBIwPoS](https://github.com/Aetherinox/obi-sync-docs/assets/118329232/c73e23cd-3a3e-46ca-8bf1-62bb948930d2)

From this point, you can upload whatever files you wish to have published and play around with the settings.

<br />

<sub>[🔝 Top](#plugin-installation)</sub>

<br />

---

<br />
