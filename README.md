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

Windows

```bash
eventReceiverQuel.exe --config "path to config yaml"
```

Linux

```bash
./eventReceiverQuel --config "path to config yaml"
```
