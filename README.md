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
|           | hossted set ssl \<AppName\> sign           | -                                                                          |
|           | hossted set remote-support true            | To enable or disable remote ssh access with our maintanece and support key |
| httpopen  | hossted httpopen \<AppName\>               | httpopen appname                                                           |
|           |                                            |                                                                            |
| list      |                                            | List hossted apps on node                                                  |
| support   |                                            | Open support tickets                                                       |
| ip        |                                            | Get external and internal ip addresses                                     |
| dashboard |                                            | Open browser with dashboard                                                |
|           |                                            |                                                                            |
| logs      |                                            | Read docker-compose logs                                                   |
| ps        |                                            | docker-compose ps                                                          |
| support   |                                            | Open support tickers                                                       |
| htopen    |                                            | Remove httpauth from CLI                                                   |
| url       |                                            | Set front-end URL                                                          |
| ssl       | signed                                     | Change to custom signed SSL                                                |

### Binary
Generally it is not a good idea to download the binary file directly from anywhere on the web. But if you do not have Go environment setup, you can download the compiled file here.

| Operating System | Binary                    |
|------------------|---------------------------|
| Linux (64-bit)   | [Here](bin/linux/hossted) |
| Dev (64-bit)     | [Here](bin/osx/hossted)   |


### Source
Or you can just install it with `go install` from the source

```
git clone https://github.com/hossted/cli.git
cd cli
go install .
```
# Pre-requiste
  uuid being saved under `/opt/linnovate/run/uuid.txt`
  ```
  <uuid>
  ```


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
