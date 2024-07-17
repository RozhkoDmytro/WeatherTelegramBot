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

type ApiMongoDB struct {
	Client  *mongo.Client
	Logger  *slog.Logger
	baseURL string
	context *context.Context
	db      *mongo.Database
}

type Subscribers struct {
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

func NewApiMongoDB(url string, l *slog.Logger) (*ApiMongoDB, error) {
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
		db:      client.Database(defaultDBName),
		Logger:  l,
	}, nil
}

func (apiMongoDB *ApiMongoDB) CloseApiMongoDB() error {
	if err := apiMongoDB.Client.Disconnect(*apiMongoDB.context); err != nil {
		return err
	}
	return nil
}

func (apiMongoDB *ApiMongoDB) Subscribe(chatId int, lat, lon float64, t time.Time) error {
	s := Subscribers{ChatId: chatId, Location: Location{Longitude: lon, Latitude: lat}, Hour: t.Hour()}

	// firstly delete all previos subscribe for this chatId
	apiMongoDB.Unsubscribe(chatId)

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

func (apiMongoDB *ApiMongoDB) Unsubscribe(chatId int) error {
	collection := apiMongoDB.Client.Database(defaultDBName).Collection(defualtTableSubscribe)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	filter := bson.D{{Key: "chatid", Value: chatId}}

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	apiMongoDB.Logger.Info("Deleted object(s)", "Count", res.DeletedCount)

	return nil
}

func (apiMongoDB *ApiMongoDB) GetSubsribersByTime(h int) ([]primitive.M, error) {
	collection := apiMongoDB.Client.Database(defaultDBName).Collection(defualtTableSubscribe)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	filter := bson.D{{Key: "hour", Value: h}}
	cursor, err := collection.Find(ctx, filter)
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

	apiMongoDB.Logger.Info("Results of ChatId", "Subscribers", results)

	return results, nil
}

func (apiMongoDB *ApiMongoDB) GetAllSubsribers() ([]int, error) {
	collection := apiMongoDB.Client.Database(defaultDBName).Collection(defualtTableSubscribe)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defualtTimeOut)
	defer cancel()

	// Define the filter for the document you want to delete
	cursor, err := collection.Find(ctx, bson.D{})
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

		if value, ok := document["chatid"]; ok {
			results = append(results, int(value.(int32)))
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	apiMongoDB.Logger.Info("Results of ChatId", "Subscribers", results)

	return results, nil
}

func SendReport() {
}
