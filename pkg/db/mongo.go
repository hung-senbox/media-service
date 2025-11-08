package db

import (
	"context"
	"fmt"
	"log"
	"media-service/pkg/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var TopicCollection *mongo.Collection
var PDFCollection *mongo.Collection
var TopicResourceCollection *mongo.Collection
var VideoUploaderCollection *mongo.Collection
var MediaAssetCollection *mongo.Collection

func ConnectMongoDB() {
	d := config.AppConfig.Database.Mongo

	var uri string
	if d.User != "" && d.Password != "" {
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/?directConnection=true",
			d.User, d.Password, d.Host, d.Port,
		)
	} else {
		uri = fmt.Sprintf(
			"mongodb://%s:%s/?directConnection=true",
			d.Host, d.Port,
		)
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := MongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	TopicCollection = MongoClient.Database(d.Name).Collection("topics")
	PDFCollection = MongoClient.Database(d.Name).Collection("pdf_resources")
	TopicResourceCollection = MongoClient.Database(d.Name).Collection("topic_resources")
	VideoUploaderCollection = MongoClient.Database(d.Name).Collection("video_uploaders")
	MediaAssetCollection = MongoClient.Database(d.Name).Collection("media_assets")
	log.Println("Connected to MongoDB and loaded 'topics', 'pdf_resources', 'topic_resources', 'video_uploaders', 'media_assets' collections")
}
