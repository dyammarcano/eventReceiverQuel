# Hacks and tools, (event hub reveiver)

## Description

This is a simple tool to receive messages from an event hub. It is intended to be used for testing purposes.


## Template Config

```yaml
eventHub:
  topic:
  accountName:
  accountKey:
  consumerGroup:
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
