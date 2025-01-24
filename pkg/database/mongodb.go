package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient establishes a connection to a MongoDB instance using provided URI and auth credentials.
func NewClient(uri, username, password string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	if username != "" && password != "" {
		fmt.Println(username)
		fmt.Println(password)
		opts.SetAuth(options.Credential{
			Username: username, Password: password,
		})
	}

	// Use mongo.Connect instead of mongo.NewClient
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	// Verify connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}
