# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased/dev] - 2022-04-19

### Added

- [Lock/Unlock master key](https://github.com/hossted/cli/issues/17)
  - Added command - hossted set remote-support true|false

- [Multi environment support](https://github.com/hossted/cli/issues/21)
  - Added environment support in build flag (LDFLAGS/DEVFLAGS) in Makefile
  - Prod/Dev should be pointing to corresponding environment

### Changed
- [Set Domain](https://github.com/hossted/cli/issues/7)
  - Change set URL to set domain

### Planned
- [Allow unregistered access](https://github.com/hossted/cli/issues/20)
  - Currently user need to run **hossted register** before any commands, will release this restriction.

- [Multi environment support](https://github.com/hossted/cli/issues/21)
  - Study to add a new flag **--continuous**, and disable the master key when user press Ctrl-C again.


## [0.1.5] - 2022-04-06

### Added
(TBA)

### Changed
(TBA)
