# Kakfa Dump

Read kafka topic from timestamp, filter and save messages to a text file or
to mongodb.

- Uses [kafka-go](https://github.com/segmentio/kafka-go) package.
- Only works with Kafka >= v0.10.0.
- Only works with JSON-encoded payload.
- May take some time on start, when consuming from the exact timestamp, because
it collects and resets offsets for all topic partitions.

## Install

```sh
go get github.com/tetafro/kafka-dump
```

Get sample config and set values
```sh
curl -o config.yaml https://raw.githubusercontent.com/tetafro/kafka-dump/master/config.example.yaml
```

## Run

Run (no flags, only config values)
```sh
kafka-dump
```

Output
```
INFO[15:45:39] Starting...
INFO[15:45:39] Saving messages to messages.txt
INFO[15:45:49] Read messages from 2020-12-30 14:22:00 to 2020-12-30 14:53:01 (total 140346, saved 1428)
INFO[15:45:59] Read messages from 2020-12-30 14:22:00 to 2020-12-30 14:54:00 (total 334520, saved 3425)
INFO[15:46:09] Read messages from 2020-12-30 14:22:01 to 2020-12-30 14:54:01 (total 525725, saved 5463)
```

## Using mongodb

Run mongodb in docker, publish port to localhost
```sh
docker run -d \
    --publish 127.0.0.1:27017:27017 \
    --name mongo \
    mongo:4.2
```

Setup mongodb parameters in config
```yaml
mongo:
  addr: mongodb://localhost:27017
  database: kafka
  collection: events
```

Run
```
$ kafka-dump
INFO[15:45:39] Starting...
INFO[15:45:39] Saving messages to mongodb://localhost:27017
```

Check
```
$ mongo mongodb://localhost:27017/kafka
> db.events.count()
166714
```
