package mongodb

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BaseURL               = "mongodb://localhost:27017"
	defaultDBName         = "Telegram"
	defualtTableSubscribe = "Subscribe"
	defualtTimeOut        = 5
)

type ApiMongoDB struct {
	Client  *mongo.Client
	Logger  *slog.Logger
	baseURL string
	context *context.Context
	db      *mongo.Database
}

type Subscribe struct {
	ChatId   int           `json:"chatId"`
	Location Location      `json:"location"`
	Hour     time.Duration `json:"hour"`
}

// Location represents a point on the map.
type Location struct {
	// Longitude as defined by sender
	Longitude float64 `json:"longitude"`
	// Latitude as defined by sender
	Latitude float64 `json:"latitude"`
}

func NewApiMongoDB(url string, nameDataBase string, l *slog.Logger) (*ApiMongoDB, error) {
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

	return &ApiMongoDB{
		baseURL: url,
		Client:  client,
		context: &ctx,
		db:      client.Database(nameDataBase),
		Logger:  l,
	}, nil
}

func (apiMongoDB *ApiMongoDB) CloseApiMongoDB() error {
	if err := apiMongoDB.Client.Disconnect(*apiMongoDB.context); err != nil {
		return err
	}
	return nil
}

func (apiMongoDB *ApiMongoDB) AddSubscribe(chatId int, lon, lat float64, t *time.Time) error {
	s := Subscribe{ChatId: chatId, Location: Location{Longitude: lon, Latitude: lat}, Hour: time.Hour}

	collection := apiMongoDB.Client.Database(defaultDBName).Collection(defualtTableSubscribe)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	res, err := collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}

	apiMongoDB.Logger.Info("New object", "ID", res.InsertedID)

	return nil
}

func (apiMongoDB *ApiMongoDB) DeleteSubscribe(chatId int) error {
	collection := apiMongoDB.Client.Database(defaultDBName).Collection(defualtTableSubscribe)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	filter := bson.D{{Key: "chatId", Value: chatId}}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	apiMongoDB.Logger.Info("Deleted object(s)", "Count", res.DeletedCount)

	return nil
}

func (apiMongoDB *ApiMongoDB) GetSubsribersByTime(t time.Duration) ([]int, error) {
	collection := apiMongoDB.Client.Database(defaultDBName).Collection(defualtTableSubscribe)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	filter := bson.D{{Key: "time", Value: t}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []int
	for cursor.Next(ctx) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}

		if value, ok := document["ChatId"]; ok {
			results = append(results, value.(int))
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	apiMongoDB.Logger.Info("Results of ChatId", "Subscribers", results)

	return results, nil
}
