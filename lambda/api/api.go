package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse the request body
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		// return an error
		return events.APIGatewayProxyResponse{
			Body:       "Username or password is empty",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// check if the user already exists

	userExists, err := api.dbStore.DoesUserExist(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error checking if user exists",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	user, err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error hashing password",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error hashing password: %v", err)

	}

	// insert the user into the database
	err = api.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error inserting user",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user: %v", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "User registered successfully",
		StatusCode: http.StatusOK,
	}, nil

}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error getting user: %v", err)
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid username or password",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Login successful",
		StatusCode: http.StatusOK,
	}, nil
}
