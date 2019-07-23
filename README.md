<img src="https://raw.githubusercontent.com/battlebas/go-surprise/master/surprise_logo.png" width="256">

# go-surprise
The goal of this project is the provide a backend application that will help
manage the people participating in the Dutch tradition surprise (sur-preeze).

Don't have a hat to draw names from? Don't want to accidently lose the name of the person your suppose to give a gift to? Then no worries, let go-surprise do it for you!

## What is Surprise (sur-preeze)?
The Dutch tradition Surprise is nearly identical to the tradition "Secret Santa" but with the gifts wrapped in original and funny ways and usually accompanied with a poem. It is traditional celebrated on the evening before Sinterklaas's birthday on December 5th, that is why the Go Gopher is dressed the way he is!

# Development Instructions

Run the unit tests with `go test`

To locally test the application
1. `docker build -t surprise .`
2. `docker run -p 8080:8081 -e "PORT=8081" surprise`
3. Open [Postman](https://www.getpostman.com/)
4. Send a POST request to localhost:8080 and copy/paste the contents of "people-test.json" into the body.
5. The application will return a JSON that randomly matched each person with another