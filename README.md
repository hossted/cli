# cli
The hossted cli - built to interact with a hossted container



#  Command
| Command   | Command Example                            | Descriptions                                                               |
|-----------|--------------------------------------------|----------------------------------------------------------------------------|
| register  | hossted register                           | Register email and organization                                            |
| set       |                                            | Change application settings                                                |
|           | hossted set list                           | List all the commands of the available applications                        |
|           | hossted set auth \<AppName\> true          | Set authorization of the provided application                              |
|           | hossted set domain \<AppName\> example.com | Set the domain of the provided application                                 |
|           | hossted set ssl \<AppName\> sign           | (TBC)                                                                      |
|           | hossted set remote-support true            | To enable or disable remote ssh access with our maintanece and support key |
| httpopen  | hossted httpopen \<AppName\>               | httpopen appname                                                           |
| logs      | hossted logs \<AppName\>                   | View Applicatin logs                                                       |
| ps        | hossted ps \<AppName\>                     | docker-compose ps of the application                                       |
|           |                                            |                                                                            |
| ip        | -                                          | (TBC) Get external and internal ip addresses                               |
| dashboard | -                                          | (TBC) Open browser with dashboard                                          |

# Pre-requiste
## uuid
  uuid being saved under `/opt/linnovate/run/uuid.txt`
  ```
  <uuid>
  ```

## sudo access
- user should have **sudo access** to for most of the change setting, docker commands, etc..

# Installation
## Binary
Generally it is not a good idea to download the binary file directly from anywhere on the web. But if you do not have Go environment setup, you can download the compiled file here.

| Operating System | Branch | Binary                                                           |
|------------------|--------|------------------------------------------------------------------|
| Linux (64-bit)   | Main   | [Here](https://github.com/hossted/cli/raw/dev/bin/linux/hossted) |
| Dev (64-bit)     | Dev    | [Here](https://github.com/hossted/cli/raw/dev/bin/dev/hossted)   |


## Source
Or you can just install it with `go install` from the source
```
git clone https://github.com/hossted/cli.git
cd cli
go install .
```


## Manual
blah blah




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
   organization: hossted
   userToken:
   sessionToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjcxLCJpYXQiOjE2NDY1NTE5MTgsImV4cCI6MTY0NjYzODMxOH0.jgweC-by2l7ksJ9NZUtjgIqvpu27ls7NZEsZgKrmkGA
   uuidPath: /opt/linnovate/run/uuid.txt
   ```
