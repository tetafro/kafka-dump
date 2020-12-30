package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Storage describes a storage for messages.
type Storage interface {
	Save(Message) error
	Close()
}

// FileSystemStorage is a storage that saves messages to a file.
type FileSystemStorage struct {
	file *os.File
}

// NewFileSystemStorage creates new filesystem storage.
func NewFileSystemStorage(file string) (*FileSystemStorage, error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o600) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("open file %s: %v", file, err)
	}
	return &FileSystemStorage{file: f}, nil
}

// Save formats a message as an indented json and saves it to the file.
func (s *FileSystemStorage) Save(msg Message) error {
	data, err := json.MarshalIndent(msg.data, "", "    ")
	if err != nil {
		return fmt.Errorf("marshall message: %v", err)
	}
	_, err = s.file.Write(append(data, []byte("\n")...))
	if err != nil {
		return fmt.Errorf("write data to file: %v", err)
	}
	return nil
}

// Close properly closes the file.
func (s *FileSystemStorage) Close() {
	s.file.Close() // nolint: errcheck,gosec
}

// MongoStorage is a storage that saves messages to mongodb.
type MongoStorage struct {
	client     *mongo.Client
	database   string
	collection string
}

// NewMongoStorage creates new mongodb storage.
func NewMongoStorage(conf MongoConf) (*MongoStorage, error) {
	if conf.Addr == "" {
		return nil, fmt.Errorf("mongo address is empty")
	}
	if conf.Database == "" {
		return nil, fmt.Errorf("mongo database is empty")
	}
	if conf.Collection == "" {
		return nil, fmt.Errorf("mongo collection is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	m, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Addr))
	if err != nil {
		return nil, fmt.Errorf("create connection: %v", err)
	}

	if err := m.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping primary node: %v", err)
	}

	s := &MongoStorage{
		client:     m,
		database:   conf.Database,
		collection: conf.Collection,
	}
	return s, nil
}

// Save saves a message to mongodb.
func (s *MongoStorage) Save(msg Message) error {
	collection := s.client.Database(s.database).Collection(s.collection)
	_, err := collection.InsertOne(context.Background(), msg.data)
	if err != nil {
		return fmt.Errorf("write data to file: %v", err)
	}
	return nil
}

// Close properly closes mongodb connection.
func (s *MongoStorage) Close() {
	s.client.Disconnect(context.Background()) // nolint: errcheck,gosec
}
