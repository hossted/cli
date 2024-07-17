# hossted cli
The hossted cli - built to interact with a hossted container
<br/>

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
# Table of Contents
- [Command](#command)
- [Pre-requisite](#pre-requisite)
    - [uuid](#uuid)
    - [software.txt](#softwaretxt)
    - [sudo access](#sudo-access)
- [Installation](#installation)
    - [Binary](#binary)
    - [Source](#source)
    - [Manual](#manual)
        - [Prod](#prod)
        - [Dev](#dev)
- [Version](#version)
- [Usage](#usage)

<!-- markdown-toc end -->


# Command
| Command   | Command Example                            | Descriptions                                                               |
|-----------|--------------------------------------------|----------------------------------------------------------------------------|
| register  | hossted register                           | Register email and organization                                            |
| set       |                                            | Change application settings                                                |
|           | hossted set list                           | List all the commands of the available applications                        |
|           | hossted set auth \<AppName\> true          | Set authorization of the provided application                              |
|           | hossted set domain \<AppName\> example.com | Set the domain of the provided application                                 |
|           | hossted set ssl \<AppName\> sign           | (TBC)                                                                      |
|           | hossted set remote-support true            | To enable or disable remote ssh access with our maintanece and support key |
| logs      | hossted logs \<AppName\>                   | View Application logs                                                      |
| ps        | hossted ps \<AppName\>                     | docker compose ps of the application                                       |
| version   | hossted version                            | Get the version of the hossted CLI program                                 |
|           |                                            |                                                                            |
| ip        | -                                          | (TBC) Get external and internal ip addresses                               |
| dashboard | -                                          | (TBC) Open browser with dashboard                                          |


<br/><br/>

# Pre-requisite
## uuid
  uuid being saved under  `/opt/hossted/run/uuid.txt` or `/opt/linnovate/run/uuid.txt`
  ```
  <uuid>
  ```

## software.txt
  `software.txt` being saved under `/opt/hossted/run/software.txt` or `/opt/linnovate/run/software.txt`, and it should be in this format to get the available applications
  ```
  Linnovate-<CloudProvider>-<Application>
  ```

  __Example__
  ```
  Linnovate-AWS-gitbucket
  ```

## sudo access
- user should have **sudo access** for most of the change setting commands, docker commands, etc.. to work.

<br/><br/>

# Installation
## Binary
Generally it is not a good idea to download the binary file directly from anywhere on the web. But if you do not have Go environment setup, you can download the compiled file here.

| Operating System    | Branch | Binary                                                           |
|---------------------|--------|------------------------------------------------------------------|
| Linux Prod (64-bit) | Main   | [Here](https://github.com/hossted/cli/raw/main/bin/linux/hossted) |
| Linux Dev (64-bit)  | Dev    | [Here](https://github.com/hossted/cli/raw/dev/bin/dev/hossted)   |

<br/>

## Source
Or you can just install it with `go install` from the source
```
git clone https://github.com/hossted/cli.git
cd cli
go install .
```

<br/>

## Manual

#### Prod
```bash
wget https://github.com/hossted/cli/raw/main/bin/linux/hossted
chmod 755 hossted
sudo cp ./hossted /usr/local/bin
```

<br/>

#### Dev
```bash
wget https://github.com/hossted/cli/raw/dev/bin/dev/hossted
mv hossted hossted-dev
chmod 755 hossted-dev
sudo cp ./hossted-dev /usr/local/bin
```

# Version
Get the version with following command
```
hossted version
```

__Result__
```
# or more recent
hossted version v0.1.5.
Built on 2022-04-10 (24b8619)
```

<br/><br/>

# Usage
0. Generate help on the commands with `-h` or `--help`<br/>

   ```
   hossted -h
   hossted register -h
   ```
   <br/>

1. Register users
   `
   hossted register
   `
   <br/>

2. Config file is saved under `~/.hossted/config.yaml`
   ```yaml
   email: abc@hossted.com
   userToken:
   sessionToken: eyJhbGdaOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjIzLCJpYXQiOjE2NDY4NDIxOTAsImV4cCI6MTY0NjkyODU5MH0.JMUCLFMHLznZ7Dc0uNFhFFS0J-LqoB_mAehnMFFwgfs
   uuidPath: /opt/linnovate/run/uuid.txt
   applications:
       - appName: prometheus
         appPath: /opt/prometheus
       - appName: demoapp
         appPath: /opt/demoapp
   ```
t2
