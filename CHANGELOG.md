# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog].

This project adheres to [Semantic Versioning].

## [Unreleased]

## [1.3.2] - 2022-04-11

### Changed

- Greatly improved git branch inference.

## [1.3.1] - 2022-04-07

### Fixed

- Fixed issues with git branch inference.
- Fixed erroneous detection of the git commit from environment variables in a
  GitHub Actions workflow when merging via pull request.

## [1.3.0] - 2022-03-24

### Added

- Added support for detecting the current git branch name and git commit from
  environment variables specific to the CI provider (if any).

## [1.2.0] - 2022-03-03

### Added

- Added support for triggering the run of all currently enabled test flows for
  an app. Optionally, the set of test flows to run can be controlled by
  supplying a rule name.

## 1.1.0 - _(skipped)_

## [1.0.0] - 2021-11-02

Initial public release.

[Unreleased]:   https://github.com/waldoapp/waldo-go-lib/compare/v1.3.2...HEAD
[1.3.2]:        https://github.com/waldoapp/waldo-go-lib/compare/v1.3.1...v1.3.2
[1.3.1]:        https://github.com/waldoapp/waldo-go-lib/compare/v1.3.0...v1.3.1
[1.3.0]:        https://github.com/waldoapp/waldo-go-lib/compare/v1.2.0...v1.3.0
[1.2.0]:        https://github.com/waldoapp/waldo-go-lib/compare/v1.0.0...v1.2.0
[1.0.0]:        https://github.com/waldoapp/waldo-go-lib/compare/7a87b12...v1.0.0

[Keep a Changelog]:     https://keepachangelog.com
[Semantic Versioning]:  https://semver.org
