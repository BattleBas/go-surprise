[![Build Status](https://travis-ci.com/BattleBas/go-surprise.svg?branch=master)](https://travis-ci.com/BattleBas/go-surprise)

<img src="https://raw.githubusercontent.com/BattleBas/go-surprise/master/assets/surprise_logo.png" width="256">

# go-surprise
The goal of this project is the provide a backend application that will help
manage the people participating in the Dutch tradition surprise (sur-preeze).

Don't have a hat to draw names from? Don't want to accidently lose the name of the person your suppose to give a gift to? Then no worries, let go-surprise do it for you!

## What is Surprise (sur-preeze)?
The Dutch tradition Surprise is nearly identical to the tradition "Secret Santa" but with the gifts wrapped in original and funny ways and usually accompanied with a poem. It is traditional celebrated on the evening before Sinterklaas's birthday on December 5th, that is why the Go Gopher is dressed the way he is!

# Build requirements

To build `go-surprise` you need:

* Docker (18.03 or above)
* Go (1.12 or above)
* PostGres (11.4 or above)

# Development Instructions

Run the unit tests with `go test`

To locally test the application
1. `docker build -t surprise .`
2. `docker run --rm -p 8080:8081 --env-file local-env.list surprise` (Note: update "\<password\>" with the one you set when setting up PostGres)
3. Open [Postman](https://www.getpostman.com/)
4. You can now send requests to the server