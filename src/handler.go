package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"studyAndRepeat/src/database"
	"studyAndRepeat/src/handlers"
)

func main() {

	pref := tele.Settings{Token: os.Getenv("TG_TOKEN"), Verbose: false, Synchronous: true}

	db := database.New()
	defer db.Close()

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", handlers.Start(db))
	b.Handle("/add", handlers.AddCard(db))
	b.Handle("/delete", handlers.DeleteCard(db))
	b.Handle("/dictionary", handlers.GetDictionary(db))

	showAnswer := tele.Btn{Text: "Show answer", Unique: "answer"}
	b.Handle("/train", handlers.Train(db, showAnswer))
	answerAgain := tele.Btn{Text: "Again, in < 10min", Unique: "again"}
	answerHard := tele.Btn{Text: "Hard, in 1 day", Unique: "hard"}
	answerGood := tele.Btn{Text: "Good, in next days", Unique: "good"}
	answers := []tele.Btn{answerAgain, answerHard, answerGood}
	b.Handle(&showAnswer, handlers.ShowAnswer(db, answers))

	b.Handle(&answerAgain, handlers.SaveAnswer(db, 0, showAnswer))
	b.Handle(&answerHard, handlers.SaveAnswer(db, 24, showAnswer))
	b.Handle(&answerGood, handlers.SaveAnswer(db, -1, showAnswer))
	b.Handle("/ping", func(c tele.Context) error { return c.Send("pong 🎾") })
	b.Handle(tele.OnText, handlers.SetDefinition(db))

	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var u tele.Update
		if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
			log.Print(req.Body)
			b.ProcessUpdate(u)
			return events.APIGatewayProxyResponse{Body: "processed", StatusCode: 200, IsBase64Encoded: false}, nil
		}

		log.Printf("can't process %s", req.Body)
		return events.APIGatewayProxyResponse{Body: "error", StatusCode: 400, IsBase64Encoded: false}, err
	})
}
