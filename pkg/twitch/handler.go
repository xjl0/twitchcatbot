package twitch

import (
	"fmt"
	"log"
	"time"
)

func Handler(bot *twitchBot, word string, reply int) {
	for {
		err := bot.SendMessage(word)
		if err != nil {
			fmt.Println(err)
		}
		log.Println(word, time.Now())
		time.Sleep(time.Second * time.Duration(reply))
	}
}
