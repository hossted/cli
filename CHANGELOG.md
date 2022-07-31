# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [0.2.4] - 2022-07-31
### Added
- Changed registerUser api endpoint to new instances api
 
## [0.2.3] - 2022-06-29
### Added
- [After set domain, change MOTD accordingly](https://github.com/hossted/cli/issues/41)
  - After setting the domain, the new domain name does not match those in MOTD file. It is matching now after this release.

- [Support /opt/hossted in addition to /opt/linnovate](https://github.com/hossted/cli/issues/42)
  - Support both /opt/hossted (Preferred) and /opt/linnovate as the default path of the cli program.


## [0.2.2] - 2022-06-07
### Added

- [Make set domain generally available](https://github.com/hossted/cli/issues/39)
  - SetDomain was an app-specific command, changed to available for all

- [Changes to default / help text](https://github.com/hossted/cli/issues/37)
  - Change default help text as custom text to customize the ordering of important commands or subcommands


## [0.2.1] - 2022-05-18
### Added

- [Remove httpopen](https://github.com/hossted/cli/issues/34)
  - Removed deprecated command httpopen (changed to set auth false)
  - Removed old and development related file

- [Change register fail message when a user already registered the machine](https://github.com/hossted/cli/issues/27)


## [0.1.9] - 2022-05-17
### Added

- [Remove gitbucket hardcoding](https://github.com/hossted/cli/issues/33)
  - Change set auth false from app-specific command to general available command. the command should be available to apps besides gitbucket.
  - NOT included the set auth true command (#32)
  - NOT included the general available for set domain command, still only available to specific app (prometheus, airflow, wordpress), unless specified otherwise.


## [0.1.8] - 2022-05-14
### Added

- [New text when user runs "hossted" with no command](https://github.com/hossted/cli/issues/28)
  - Changed greeting message

- [Change register fail message when a user already registered the machine](https://github.com/hossted/cli/issues/27)
  - Update error message if the email already registered

- [implement hossted set auth appName false](https://github.com/hossted/cli/issues/15)
  - Change httpopen to set auth false
  - only gitbucket is supported now
  - deprecate old httpopen command


### Bug fixed
- Bug fix: Error prompt for No available app commands for unregistered user.



## [0.1.7] - 2022-05-04
### Added

- [Allow unregistered access](https://github.com/hossted/cli/issues/20)
  - Currently user need to run **hossted register** before any commands, will release this restriction.

- [use default app when appName is missing](https://github.com/hossted/cli/issues/25)
  - If only one app in the vm, use that as a default app
  - If it's under the app directory, use it as default app


## [0.1.6] - 2022-04-19

### Added

- [Lock/Unlock master key](https://github.com/hossted/cli/issues/17)
  - Added command - hossted set remote-support true|false

- [Multi environment support](https://github.com/hossted/cli/issues/21)
  - Added environment support in build flag (LDFLAGS/DEVFLAGS) in Makefile
  - Prod/Dev should be pointing to corresponding environment

### Changed
- [Set Domain](https://github.com/hossted/cli/issues/7)
  - Change set URL to set domain


## [0.1.5] - 2022-04-06

### Added
(TBA)

### Changed
(TBA)


A
## Planned/Backlog
- [Multi environment support](https://github.com/hossted/cli/issues/21)
  - Study to add a new flag **--continuous**, and disable the master key when user press Ctrl-C again.
