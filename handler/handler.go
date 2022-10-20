package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
)

func main() {

	pref := tele.Settings{
		Token:       os.Getenv("TG_TOKEN"),
		Verbose:     false,
		Synchronous: true,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello! How are you doing?")
	})

	b.Handle("/hi", func(c tele.Context) error {
		return c.Send("Hi there!")
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("For now I don't know how to respond")
	})

	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var u tele.Update
		if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
			log.Printf("update text is: %s", u.Message.Text)
			b.ProcessUpdate(u)
			return events.APIGatewayProxyResponse{Body: "processed", StatusCode: 200, IsBase64Encoded: false}, nil
		}

		log.Printf("can't process %s", req.Body)
		return events.APIGatewayProxyResponse{Body: "error", StatusCode: 400, IsBase64Encoded: false}, err
	})
}
