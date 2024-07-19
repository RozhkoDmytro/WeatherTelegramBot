package mongodb_test

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"projecttelegrambot/pkg/mongodb"

	"github.com/stretchr/testify/assert"
)

func TestNewApiMongoDB(t *testing.T) {
	// Create logger
	logger, err := createLogger("app.log")
	if err != nil {
		panic(err)
	}
	_, err = mongodb.NewMongoDBService("http://wrongURL", logger)

	if err == nil {
		assert.Equal(t, "True", "False", "Wrong url")
	}
}

// Create logger and set fields
func createLogger(NameLog string) (*slog.Logger, error) {
	// Create logger
	file, err := os.OpenFile(NameLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	w := io.MultiWriter(os.Stderr, file)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}
