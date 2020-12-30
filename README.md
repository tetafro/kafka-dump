# Kakfa Dump

Read kafka topic from timestamp, filter and save messages.

- Uses [kafka-go](https://github.com/segmentio/kafka-go) package.
- Only works with Kafka >= v0.10.0.
- Only works with JSON-encoded payload.
- May take some time on start, when consuming from the exact timestamp, because
it collects and resets offsets for all topic partitions.

## Install

```sh
go get github.com/tetafro/kafka-dump
```

# Run

Get sample config and set values
```sh
curl -o config.yaml https://raw.githubusercontent.com/tetafro/kafka-dump/master/config.example.yaml
```

Run (no flags, only config values)
```sh
kafka-dump
```

Output
```
INFO[15:45:39] Starting...
INFO[15:45:39] Saving messages to mongodb://localhost:27017
INFO[15:45:49] Read messages from 2020-12-30 14:22:00 to 2020-12-30 14:53:01 (total 140346, saved 1428)
INFO[15:45:59] Read messages from 2020-12-30 14:22:00 to 2020-12-30 14:54:00 (total 334520, saved 3425)
INFO[15:46:09] Read messages from 2020-12-30 14:22:01 to 2020-12-30 14:54:01 (total 525725, saved 5463)
```
