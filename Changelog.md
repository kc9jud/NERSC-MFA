# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
###
- Add debug mode for, well, debugging

## [1.0.3] - 2019-01-28
### Added

### Changed
- Made order of arguments in usage message match description (again)

## [1.0.2] - 2019-01-28
### Added
- Added -v flag to display version number and exit, and updated usage message with same
- Cleaned up Usage message, added documentation of new options

### Changed
- Fixed bug in password quoting which caused failure for passwords with spaces
- Fixed bug with password reading which caused passwords with "\" to fail
- Added successful completion message when obtaining PuTTY keys

## [1.0.1] - 2019-01-27
### Added
- Included target username in password prompt
- support for PuTTY key format (-p)

### Changed


## [<1.0.0>] - 2018-10-31
### Added
- Extracts public key from private key, in addition to providing cert.  Some openssh distros can use certs without also having the public key present.  Go figure.
- "-a" flag to add keys to intelligently add keys ssh-agent including telling ssh-agent the expiration time.
