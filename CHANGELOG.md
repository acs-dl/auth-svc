# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.3] - 2023-03-28

### Fixed

- Repeating generated refresh tokens, by adding created at time in nanoseconds 

## [1.0.2] - 2023-03-27

### Added

- Worker to delete expired tokens from database

### Fixed

- Add expires at date in structure for creating refresh token
- Split date from structure in function where converting claims map to claims structure
- Runners blocked API, so changed runners starting logic

## [1.0.1] - 2023-03-22

### Added

- Rabbit-mq channels for communicating with other modules

### Changed

- Receiving modules with their permissions from `orchestrator`
- Update database for new module - permissions structure

## [1.0.0] - 2023-03-15

### Added

- Database.
- API handlers.


[1.0.0]: https://gitlab.com/distributed_lab/acs/auth
[1.0.1]: https://gitlab.com/distributed_lab/acs/auth/-/compare/main...feature/dynamic_roles
