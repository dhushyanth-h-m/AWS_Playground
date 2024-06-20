package main

import (
	"fmt"
	"net/http"

	"lambda-func/app"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"` //`json:"username"` is a struct tag to convert the raw json to the struct field

}

// take in a payload and do something with it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username is empty")
	}

	return fmt.Sprintf("Successfully called by - %s\n", event.Username), nil
}
func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUser(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
