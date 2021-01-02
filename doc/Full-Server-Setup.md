# Guide on Setting up a Full Server

This document will tell you exactly how to set up and operate a full server. This document was written with non-developers in mind, so people new to linux or command line operations should be able to follow along without any trouble. 

You can of course run the GoDBLedger software on any operating system that you can get it to build on, but for the purposes of this document, the instructions apply to running a Service Node on a remote Ubuntu 18.04 server. If that isn’t what you want to do, syntax and server set up will of course differ according to whatever OS you choose to run your server.

## Table of Contents

- [Overview of GoDBLedger Nodes](#Overview)
- [New User Guide](#new-user-guide)
    - Step 1 [Server](#step-1-get-a-server)
    - Step 2 [Server Prep](#step-2-prepare-your-server)
    - Step 3 [Download Binaries](#step-3-download-binaries)
    - Step 4 [Service File](#step-4-loki-launcher-as-a-system-service)

---

## Overview

For now, all you need to know is that:

-   A server running GoDBLedger will allow various software, app and scripts to communicate financial transactions to it. 

-   The method of communication is a GRPC endpoint

-   GoDBLedger will take double entry bookeeping transactions and save them into a database of your choosing. SQLite and MySQL are built in automatically and can be linked immediately.

-   Viewing the recorded transactions is primarily intended to be done via the database. GoDBLedger by itself only contains minimal features to view the data. It does however come with a `reporter` executable that will execute some prebuilt SQL commands so you can view transactions as expected (Trial Balance and General Ledgers directly from SQL!!)

It is also worth noting that from a usability point of view, the "Front End" apps that would be the primary method for humans to physically input their financial data are yet to be developed. A few basic scripts and apps have been created to allow for initial usage by the creators but these are very basic. Right now it is not possible to run a full financial system that can keep up with all the functions that a business would require. However this is very much intended for the future! But right now its usability is limited to how well you can spin up a script that can communicate with a GRPC Endpoint. If this is not you right now I apologise, working as hard as possible to get this usable for business owners and accountants and not just developers.

---

## New User Guide

This section of this guide is for new users to servers and the CLI interface. 


### Step 1 - Get a Server

Righto! Let’s get started. The eventual end goal of GoDBLedger is for a digital consultant to spin up a Virtual Private Server (VPS) and be able to hit the ground running. 

Right now its likely to be run on a spare computer you have at home. This is very much the best option right now as the endpoints are not secured. Eventually this will be protected with signed transactions and authentication but while we are alpha you definitely want to restrict access to the server to those on your local network

When selecting your operating system, choose Ubuntu 18.04 64 bit or Ubuntu 20.04 64 bit if you want to follow this guide. If you feel more confident or wish to run your server on another distribution or operating system, the commands in this guide will still apply.

---

### Step 2 - Prepare your Server

To access your server, you will need a SSH client for your operating system. Because we’re on Windows today, we’ll download PuTTY, Mac users can also use PuTTY. If you’re a Linux user, you probably don’t want us telling you where to get a SSH client from.

To connect to our server we will need to paste the IP address into the SSH client’s “Host Name (or IP address)” input box and click the “Open” button. The Port number can usually just be left as `22`.

A terminal window will now appear prompting for your log-in details, username(root) and password, which were provided by your provider or was set up when creating the home server. When entering your password, nothing will visually appear in the terminal. This is normal. Hit enter when it’s typed or pasted, and you should be logged in to your VPS.

#### 2.1 Create a non-root User

Best practice when running a server is to not run your software as the root user.  Although
it is possible to do everything as root, it is strongly recommended that you create a non-root user
on our VPS by running the following command:

```
adduser <username>
```

Replacing `<username>` with a name you will log-in with. For this user-guide we will use `godbledger` as our username.

```
adduser godbledger
```

The terminal will prompt you for a new password for our newly created user. Use a secure password that is different password from the root password.

Once the password has been set, the terminal will prompt for a few details about the individual running the user. You can just hit enter through each of the inputs as the details are not important for the purposes of running the server.

Once that’s done, run the following two commands to give our new account admin privileges and to change to such account.

```
usermod -aG sudo godbledger
sudo su - godbledger
```

Before we proceed further, it is advised to close your terminal and reopen PuTTY to set up a saved session with our godbledger user. Your SSH client will have a load and save session function. For PuTTY we will need to type in our IP address again, on the same screen type godbledger under “Saved Session”. Click on “Data” under the drop-down menu “Connection”, and type in godbledger (or your username defined before) into the input box “Auto-login username”. Go back to your session screen, where we entered the IP address, and click “Save”. You can load this session whenever you want to check on your Service Node.


#### 2.2 Hot Tips for using the Console

Consoles don't work like the rest of your computer. Here are some basic tips for navigating your way around the command line!

- Don't try copying something by using the usual Ctrl + C hotkey! If you want to copy something, do so by highlighting text and then right clicking it. Pasting works by right clicking a blank area in the console.
  
- If you want to kill a process or stop something from running, press Ctrl + C. (This is why you shouldn't try copying something with this hotkey.)

- You can always check the directory you are in by typing `pwd` and list its contents by typing `ls`.
    
- You can always return to your home directory by typing `cd`.

- You can move into a given directory by typing `cd <name>` or move back up one level by typing `cd ..`.

- PuTTY allows you to easily duplicate or restart a session by right clicking the top of the window. Handy if you’re trying to do a few things at once.

Once we have logged in correctly to the server for the first time, the server may be configured to prompt for a new password for the root account. The terminal will require you to enter the new password twice before we can start running commands.  If you aren't prompted for a new `root` password but want to change it anyway, type `sudo passwd`.  Choose something very secure!

#### 2.3 Server Preparation Continued

We should update our package lists, the below command downloads the package lists from the repositories and "updates" them to get information on the newest versions of packages and their dependencies. It will do this for all repositories and PPAs.

```
sudo apt update && sudo apt upgrade
sudo apt dist-upgrade && sudo apt autoremove
```

You will be prompted to authorise the use of disk space, type `y` and enter to authorise.

If you are prompted during the upgrade that a new version of any file is available then click the up and down arrows until you are hovering over `install the package maintainer’s version` and click enter.

Alright, good to go. Our server is now set up, up to date, and is not running as root. On to the fun part!

---

### Step 3 - Release Binaries

#### 3.1 - Download GoDBLedger Binaries

We will download the binaries by running the following command:

```
wget https://github.com/darcys22/godbledger/releases/download/v0.3.0/godbledger-linux-x64-v0.4.0.tar.gz
```

#### 3.2 - Unzip the binaries to the home directory

To unzip run the following command:

```
tar -xvzf godbledger-linux-x64-v0.4.0.tar.gz 
```

#### 3.3 - Create a Symbolic Link to the folder
```
cd ~
ln -s godbledger-linux-x64-v0.4.0 godbledger
```
#### 3.4 - Copy the binaries to /usr/local

Copy the contents of the archive to the usr/local/godbledger directory where we will run GoDBLedger as a service.
```
sudo cp -a ~/godbledger/. /usr/local/godbledger
```

    Note: You will need to do this step each time you pull latest code and build a new version of the GoDBLedger binary.

#### 4.1 - Creating the Service File

To create our godbledger.service file run the following command:

```
sudo nano /etc/systemd/system/godbledger.service
```

Next copy the text below and paste it into your new file.

> To paste in putty you can right mouse click the terminal screen.

```
[Unit]
Description=godbledger
After=network-online.target

[Service]
Type=simple
User=godbledger
ExecStart=/usr/local/godbledger/godbledger 
RestartSec=30s

[Install]
WantedBy=multi-user.target
```

If you chose a username other than `godbledger` then change godbledger in the `User=` line to the alternative username.

Once completed, save the file and quit nano: `CTRL+X -> Y -> ENTER`.

#### 4.2 (Optional) - Initialise GoDBLedger and choose database 

If you simply want to run SQLite3 as the backend you can run 
```
~/godbledger/godbledger init
```

If you want to run MySQL Backend follow steps in this guide here to setup MySQL Server

Then run the init with the mysql parameter, and set the connection string
```
~/godbledger/godbledger init -m
```

#### 4.3 - Enabling the Service File

Reload systemd to reflect the changes.
```
sudo systemctl daemon-reload
```

Start the service and check to make sure it’s running correctly.

```
sudo systemctl start godbledger.service
sudo systemctl status godbledger.service
```

Enable godbledger.service so that it starts automatically upon boot:

```
sudo systemctl enable godbledger.service
```

Then reboot your system to check if the service file is working correctly.

```
sudo reboot
```

#### 4.3 (Optional) - Testing with a transaction from the command line

Log back into your server and run the following command:

```
~/godbledger/ledger_cli single
```

To view this transaction run
```
~/godbledger/reporter trialbalance
```

Well done! You're server is now setup. 

---

### Default Directories for GoDBLedger Files

|         Name              |               Directory              |
|:-------------------------:|:------------------------------------:|
|    Godbledger Binaries    |     `/usr/local/godbledger/`         |
|       Config Files        |         `/home/<user>/.ledger`       |
|      SQLite Database      |  `/home/<user>/.ledger/ledgerdata/`  |
  
---

If you can improve this guide, please submit a pull request. 

