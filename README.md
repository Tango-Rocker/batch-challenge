# Batch-Challenge

## Overview
Simple batch processing system designed to read, validate, and process monthly payment records from CSV files.

## Prerequisites

- Docker
- Docker Compose

## installation
To install the application, follow these steps:

- Clone the repository to your local machine.
- Navigate to the root directory of the project.
- Run the installation script:
- chmod +x install.sh
- ./install.sh

#### The script will interactively guide you through the setup process including:

- Prompting for the directory path for storing app data.
- Building the application binary.
- Verifying the presence of Docker and Docker Compose, with an option to install if they're missing.
- Setting up and starting the Docker containers for the app and the database.

 
## Technical Aspects



The system comprises multiple Go packages, each with a specific role:

- **Main Package**: Orchestrates the application's flow, loading configuration, and setting up dependencies.
- **App Package**: Defines the core application structure, managing the file reading process and coordinating goroutines for concurrent processing.
- **CSV Package**: Handles the parsing and validation of CSV files according to a defined schema, utilizing various validators for different data types.

## Configuration and Schema

- The application configuration is loaded at runtime, specifying paths for the source CSV and schema definition.
- The schema is defined in JSON format, detailing the required columns, their data types, and any additional validation formats or rules.

## Processing Flow

1. The main function initializes the application by loading the configuration and setting up CSV parsing dependencies.
2. The `Application` struct, containing configuration and parser, manages the reading from the source CSV file.
3. A channel is used to pass valid records from the CSV reader to the processor, which can then construct the desired output.
4. The application employs synchronization primitives to ensure that all go-routines complete their execution before the application exits.

## Validation and Transformation

- Each CSV record is validated against the schema definition to ensure the correct format and data types.
- Custom validators are used for different column types, with the ability to transform data (e.g., dates) into a consistent format.
