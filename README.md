# WIP DO NOT USE - goproject
A Go project generator

## Getting Started

```bash
# Clone the repo
git clone https://github.com/jjunqueira/goproject.git

# Build the binary
cd goproject
make
cp target/bin/goproject-darwin18 /usr/local/bin/

# Initialize configuration (downloads templates, sets configuration parameters)
goproject init

# Create a new test project to make sure everything is working correctly
goproject new basic testproj
cd testproj
make
```

## Goals
* Consistent folder structure
* Go modules
* Templates for various app types: CLI, REST, Kafka consumer, Kafka producer

## Commands

### Init
The init command initializes goproject to run for the first time. By default it will do the following:
1. Create a directory in the users home folder .config/goproject
2. Create a default config .config/goproject/config.toml
3. Create a templates directory .config/goproject/templates
4. Download default templates to .config/goproject/templates

example: `goproject init`

### New
The new command creates bootstraps a new Go project based on a template. 
All new projects will create a project using Go modules.

example: `goproject new basic testproj`

### Update
The update command executes a git pull in the goproject templates directory to make sure templates are up to date

## Project Templates

## Basic

The basic template project is the simplest project type and only does the following:
* Makes the project directory
* Initializes a git repository (using the default remote repository defined in the users configuration or the one defined in the parameters)
* Initializes Go modules project
* Creates a simple main file
* Creates a simple Makefile
