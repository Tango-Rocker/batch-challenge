# `README.md` for Batch Processing Application

This document outlines the `batch-challenge` project, a Go application designed to parse CSV data, process transactions, and generate summaries.

## Overview

The application leverages concurrency to process large CSV files efficiently. It uses a modular design with distinct components handling parsing, writing, relaying streams, and summarizing data. The project is configured to run in a Docker environment, with PostgreSQL as the database backend.

## Requirements

- Go programming language
- Docker and Docker Compose
- PostgreSQL database
- SMTP server credentials for email notifications

## Configuration

Configuration is managed through environment variables. These include paths to the data and schema files (`DATA_PATH`, `SCHEMA_PATH`), database connection details (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`), and email server settings (`MAIL_SERVER_HOST`, `MAIL_SERVER_PORT`, `MAIL_ACCOUNT`, `MAIL_PASSWORD`).

## Components

- `Parser`: Reads and validates CSV files according to a defined schema.
- `Writer`: Buffers and writes the parsed transactions to the database.
- `StreamRelayService`: Manages the propagation of data to the writer and summarizer.
- `SummaryService`: Aggregates transaction data and sends a summary report via email.
- `EmailService`: Configures and sends emails using provided SMTP settings.

## Running the Application

To run the application:

1. Configure your environment variables as needed.
2. Use Docker Compose to build and start the services defined in `docker-compose.yml`.
3. The application will begin processing data as per the CSV file specified in the `DATA_PATH`.

## Deployment

- A Dockerfile is included for building the application image.
- Use the provided `docker-compose.yml` to deploy the application along with its database.

# Technical Architecture

## Component Descriptions

### Application (`Application`)
The central orchestrator that initializes all services and triggers the processing flow.

### Parser (`csv.Parser`)
Responsible for consuming CSV files. It validates and transforms the data based on a schema definition and sends the data to the `StreamRelayService`.

### StreamRelayService (`business.StreamRelayService`)
Acts as a conduit, taking in parsed data from the `Parser` and distributing it to both the `Writer` and `SummaryService` through a subscription system

### Writer (`business.Writer`)
Buffers the data and periodically flushes this buffer to the `Database`. It listens for data from the `StreamRelayService`.

### SummaryService (`business.SummaryService`)
Aggregates data for summary and sends an email report. It receives data from the `StreamRelayService` and utilizes the `EmailService` to dispatch emails.


### Database (`db.Repository`)
Persistently stores transaction records. It is accessed by the `Writer` to insert data records.

### EmailService (`business.EmailService`)
Configures and sends out emails. Used by the `SummaryService` to send out summary reports to users.

## Data Flow

1. **Initialization**: The `Application` starts and initializes all components, setting up their interconnections.
2. **Parsing**: The `Parser` reads the CSV file, validates, and transforms the data.
3. **Relaying**: The `StreamRelayService` receives parsed data and relays it to both the `Writer` and `SummaryService`.
4. **Writing**: The `Writer` buffers transactions and writes them to the `Database` at set intervals or buffer sizes.
5. **Summarizing**: Concurrently, the `SummaryService` aggregates data to create a summary report.
6. **Emailing**: Once the summary is ready, the `SummaryService` sends it via the `EmailService`.

![img.png](img.png)

