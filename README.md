# Hacks and tools, (event hub reveiver)

## Description

This is a simple tool to receive messages from an event hub. It is intended to be used for testing purposes.


## Template Config

```yaml
azure:
  eventHub:
    accountName:
    accountKey:
    topic:
  storage:
    accountKey:
    accountName:
```

## Usage

```bash
go run main.go --connection-string "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>" --partition-id 0 --consumer-group "$Default"
```

## Build

Windows

```bash
go build -o .\hacks\bin\event-hub-receiver.exe .\hacks\receiver\receiver.go
```

Linux

```bash
go build -o ./hacks/bin/event-hub-receiver ./hacks/receiver/receiver.go
```

## Docker

```bash
docker build -t event-hub-reveiver .
```
