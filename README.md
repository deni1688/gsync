# gsync

![](https://img.shields.io/badge/stage-work%20in%20progress-orange)

This tool synchronizes files between a specified local directory and remote directory. 
Currently, the implementation is limited to google drive with plans to add more providers like AWS S3, OneDrive as well ad a VPS over SSH. Once those are added it
will be possible to use the tool to synchronize files between multiple providers.

## Table of contents

* [Usage](#usage)
* [License](#license)

## Installation

### Requirements

- Git
- Go >= 1.16
- make

1. Clone the repository

```
git clone https://github.com/deni1688/gsync.git
```

2. Install dependencies

```
go get -u github.com/deni1688/gsync/...
```

3. Run gsync directly

```
go run cmd/gsync/main.go
```
4. Build and run gsync CLI binary

```
make build_cli && ./bin/cli/gsync
```

## Usage

When using gsync to synchronize files between a local directory and google drive, the following steps are required:

1. Add Google Drive API credentials to $HOME/.gsync with the name `credentials.json`.
2. Run the cli binary

You should be prompted to authorize the application via OAuth. From there you should be able to use the cli as follows:


```

Usage:
gsync [command]

Available Commands:
completion Generate the autocompletion script for the specified shell
help Help about any command
pull pull files from remote gs to local directory
push push files from local directory to remote gs
sync sync files between a remote gs and local directory

```

## License

gsync is licensed under the [MIT license](https://github.io/deni1688/gsync/blob/master/LICENSE)
