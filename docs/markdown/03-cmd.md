# Command-Line Interface

## Overview

The Authonomy CLI is a tool designed for managing the Authonomy service. Built using the Cobra library, it provides a user-friendly command-line interface to start, configure, and manage the Authonomy service. The service is built over an SSI (Self-Sovereign Identity) service, providing advanced identity management solutions.

## Commands

### Start

Starts the Authonomy service.

**Usage:** `authonomy start [flags]`

**Flags:**

- `--reset, -r`: Resets the service.

### Configuration

The service uses Viper for configuration management. Configuration values can be set in a file named `config` or through environment variables.

**Configurable Properties:**

- `service.port`: The port on which the service runs. Default is `8081`.
- `service.badger_path`: The path to the database.
- `service.db_encryption_key`: The encryption key for the database.
- `service.ssi_service_url`: The URL for the SSI service.

## Examples

Starting the service:

```shell
authonomy start
