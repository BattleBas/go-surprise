package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BattleBas/go-surprise/pkg/email"
	"github.com/BattleBas/go-surprise/pkg/matching"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func failure(msg string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusInternalServerError,
		IsBase64Encoded: false,
		Body:            msg,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func handler(matches matching.Matches) (events.APIGatewayProxyResponse, error) {

	for _, p := range matches.Pairs {
		err := email.Send(p)
		if err != nil {
			log.Printf("%v", err)
			return failure(fmt.Sprintf("Failed to email matches, %s", err.Error()))
		}
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            "Successfully matched",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
