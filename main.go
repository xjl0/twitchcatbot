package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"twitchcatbot/pkg/twitch"
)

func main() {
	fmt.Println("Чтение config.yaml")
	if err := initConfig(); err != nil {
		log.Fatalf("Error init config file: %s", err.Error())
	}

	bot := twitch.NewTwitchBot(viper.GetString("twitch.host"), viper.GetString("twitch.nick"), viper.GetString("twitch.oauth"), viper.GetString("twitch.vhost"))
	err := bot.Connect()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	err = bot.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	defer bot.Disconnect()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		bot.Disconnect()
		os.Exit(1)
	}()

	go bot.ReadMessages()

	time.Sleep(time.Second * 3)

	var word string
	var reply int

	log.Println("Введите команду")
	_, err = fmt.Scanf("%s\n", &word)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Введите промежуток повтора в сек.")
	fmt.Scanf("%d\n", &reply)
	if err != nil {
		log.Fatal(err)
	}

	twitch.Handler(bot, word, reply)
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
