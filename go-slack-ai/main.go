package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	// for event := range analyticsChannel {
	// 	fmt.Println("Command Events")
	// 	fmt.Println(event.Timestamp)
	// 	fmt.Println(event.Command)
	// 	fmt.Println(event.Parameters)
	// 	fmt.Println(event.Event)
	// 	fmt.Println()

	// }
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_TOKEN")}

	go printCommandEvents(bot.CommandEvents())

	bot.Command("!chat <message>", &slacker.CommandDefinition{
		Description: "Send any question to wolfram",
		Examples:    []string{"Who is the president of Indonesia"},
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			query := r.Param("message")
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			data, _ := json.MarshalIndent(msg, "", "  ")
			rough := string(data[:])
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			answer := value.String()

			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)

			if err != nil {
				fmt.Println("there is an error ", err)
				w.Reply("Error, Check server for log")
			}

			w.Reply(res)

		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
