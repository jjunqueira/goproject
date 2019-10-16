# goproject
A Go project structure generator for my personal use

## Goals
* Consistent folder structure
* Go modules
* Templates for various app types: CLI, REST, Kafka consumer, Kafka producer

## Commands

### New
The new command creates bootstraps a new Go project based on a template. 
All new projects will create a project using Go modules.

example: `goproject new default`

### Tidy

Like the Go modules tidy command if you change you configuration this will execute the changes in the configuration.

example `goproject tidy`

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
