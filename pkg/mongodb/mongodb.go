package mongodb

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BaseURL               = "mongodb://localhost:27017"
	defaultDBName         = "Telegram"
	defualtTableSubscribe = "Subscribe"
	defualtTimeOut        = 5
)

type MongoDBService struct {
	client     *mongo.Client
	logger     *slog.Logger
	baseURL    string
	context    *context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

type Subscribe struct {
	ChatId   int      `json:"chatId"`
	Location Location `json:"location"`
	Hour     int      `json:"hour"`
}

// Location represents a point on the map.
type Location struct {
	// Longitude as defined by sender
	Longitude float64 `json:"longitude"`
	// Latitude as defined by sender
	Latitude float64 `json:"latitude"`
}

func NewMongoDBService(url string, l *slog.Logger) (*MongoDBService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	l.Info("Connected to MongoDB!")

	return &MongoDBService{
		baseURL:    url,
		client:     client,
		context:    &ctx,
		db:         client.Database(defaultDBName),
		logger:     l,
		collection: client.Database(defaultDBName).Collection(defualtTableSubscribe),
	}, nil
}

func (srv *MongoDBService) CloseApiMongoDB() error {
	if err := srv.client.Disconnect(*srv.context); err != nil {
		return err
	}
	return nil
}

func (srv *MongoDBService) Subscribe(chatId int, lat, lon float64, t time.Time) error {
	s := Subscribe{ChatId: chatId, Location: Location{Longitude: lon, Latitude: lat}, Hour: t.Hour()}

	// firstly delete all previos subscribe for this chatId
	srv.Unsubscribe(chatId)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	res, err := srv.collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}

	srv.logger.Info("New object", "ID", res.InsertedID)

	return nil
}

func (srv *MongoDBService) Unsubscribe(chatId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	filter := bson.D{{Key: "chatid", Value: chatId}}

	res, err := srv.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	srv.logger.Info("Deleted object(s)", "Count", res.DeletedCount)

	return nil
}

func (srv *MongoDBService) GetSubsribersByTime(h int) ([]primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	filter := bson.D{{Key: "hour", Value: h}}
	cursor, err := srv.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []primitive.M
	for cursor.Next(ctx) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}

		if _, ok := document["chatid"]; ok {
			results = append(results, document)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	srv.logger.Info("Results of ChatId", "Subscribers", results)

	return results, nil
}
