package main

import (
	"fmt"

	"github.com/BattleBas/go-surprise/pkg/email"
	"github.com/BattleBas/go-surprise/pkg/matching"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(givers matching.Group) (matching.Matches, error) {

	for _, g := range givers.People {
		if !email.IsValid(g.Email) {
			return matching.Matches{}, fmt.Errorf("Email isn't valid")
		}
	}

	recievers := givers
	recievers.People = make([]matching.Person, len(givers.People))
	copy(recievers.People, givers.People)

	matches, err := matching.CreateMatches(givers, recievers)
	if err != nil {
		return matching.Matches{}, err
	}

	err = email.SendMasterList(matches)
	if err != nil {
		return matching.Matches{}, err
	}

	return matches, nil
}

func main() {
	lambda.Start(handler)
}
