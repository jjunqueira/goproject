# WIP DO NOT USE - goproject
A Go project structure generator for my personal use

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

You can optionally set the config directory: `goproject init --path /my/config/path`

### New
The new command creates bootstraps a new Go project based on a template. 
All new projects will create a project using Go modules.

example: `goproject new empty`

## Project Templates

## Empty

The empty template project is the simplest project type and only does the following:
* Makes the project directory
* Initializes a git repository (using the default remote repository defined in the users configuration or the one defined in the parameters)
* Initializes Go modules project
* Creates a simple main file

## CLI

The cli template project creates the standard files and folders needed for a basic cli tool. It is opinionated on how the project should be structured and tested. It attempts to adhere to the standard Go project layout (https://github.com/golang-standards/project-layout) and the tiny main pattern.
* Makes the project directory
* Initializes a git repository (using the default remote repository defined in the users configuration or the one defined in the parameters)
* Initializes Go modules project
* Creates a simple main file
* Creates the internal package
* Creates an app package
* Creates an config package
* Creates a commands packages

## Custom
You can create custom templates utilizing the Custom project template type. Custom templates are configured via the goproject config. An example of setting up custom template names and paths is included in the default configuration.

Usage: `goproject new custom example`
