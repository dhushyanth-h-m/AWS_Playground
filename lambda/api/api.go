package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		// return an error
		return fmt.Errorf("required fields are empty")
	}

	// check if the user already exists

	userExists, err := api.dbStore.DoesUserExist(event.Username)

	if err != nil {
		return fmt.Errorf("error checking if user exists %w", err)
	}

	if userExists {
		return fmt.Errorf("user already exists")
	}

	// insert the user into the database
	err = api.dbStore.InsertUser(event)

	if err != nil {
		return fmt.Errorf("error inserting user %w", err)
	}

	return nil
}
