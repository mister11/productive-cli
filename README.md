# Productive CLI

![Build status](https://github.com/mister11/productive-cli/workflows/Build/badge.svg) [![License MIT](https://img.shields.io/badge/License-MIT-brightgreen)](https://github.com/mister11/productive-cli/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/mister11/productive-cli)](https://github.com/mister11/productive-cli) [![codecov](https://codecov.io/gh/mister11/productive-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/mister11/productive-cli)

CLI enables you easier tracking on [Productive](https://productive.io) from your terminal.

**DISCLAIMER: I do not take any responsibility for incorrectly tracked projects. This is still experimental and visual confirmation in the UI that everything is correct is a must**.

## Installation

Download ZIP file from [releases page](https://github.com/mister11/productive-cli/releases) and unzip it. Optionally, you can add executable to the path to make it available from everywhere.

## Usage

You can look up usage by providing a `-h` flag to any CLI command (e.g. `productive-cli -h`).

### Initialization

Before the first use, you have to run `productive-cli init` so that CLI can setup your user information. Command will ask you to provide your personal access token for Productive which can be obtained in Settings -> Security.

### Tracking

Tracking has 2 available commands:
* food
* project

Each command has an optional `-d` flag to specify a particular date in YYYY-MM-DD format (e.g. `productive-cli -d 2020-04-01 track food`). In case flag is not provided, command will used today's date.

When asked for time, you can provide number of minutes (e.g. 120) or standard HH:mm format (e.g 8:00; leading zero is optional).

#### Food

`food` command will track 30 minutes on food budget.
Alongside generic `-d` flag, you have `-w` flag to track food for the whole week.

#### Project

`project` command enables you to track time for any project you are working on. There are no additional flags - everything is entered interactively after the command is ran.

When command is ran, you'll first see all the projects you've already tracked in the past. If you want to track a new project, press `Ctrl+C`. 

At the moment, there's no nice way to delete already tracked projects. You can edit `~/.productive/config` file where they are saved.
