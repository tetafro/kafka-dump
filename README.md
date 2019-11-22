# Kakfa Dump

Read kafka topic from timestamp, filter and save messages.

- Uses [kafka-go](https://github.com/segmentio/kafka-go) package.
- Only works with Kafka >= v0.10.0.
- Only works with JSON-encoded payload.
- May take some time on start, when consuming from timestamp (not just oldest
  or newest), since it collects and resets offsets for all topic partitions.

Install
```sh
go get github.com/tetafro/kafka-dump
```

Get sample config and set values
```sh
curl -o config.yaml https://raw.githubusercontent.com/tetafro/kafka-dump/master/config.example.yaml
```

Run (no flags, only config values)
```sh
kafka-dump
```

Example output (shows statistics for total reads and number of saved messages)
```sh
INFO[21:23:11] Starting...
INFO[21:23:11] Start consumer
INFO[21:23:11] Start storage
INFO[21:23:20] Read all messages until 2019-11-19 01:07:00 (total 104171, saved 0)
INFO[21:23:30] Read all messages until 2019-11-19 01:08:00 (total 335857, saved 0)
INFO[21:23:40] Read all messages until 2019-11-19 01:09:00 (total 561615, saved 0)
```
