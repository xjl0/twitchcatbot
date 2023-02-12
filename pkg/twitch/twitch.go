package twitch

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type TwitchBot interface {
	Connect() error
	Authenticate() error
	SendMessage(message string) error
	Disconnect() error
}

type twitchBot struct {
	conn    net.Conn
	server  string
	nick    string
	pass    string
	vhost   string
	reader  *bufio.Reader
	message string
}

func NewTwitchBot(server string, nick string, pass string, vhost string) *twitchBot {
	return &twitchBot{server: server, nick: nick, pass: pass, vhost: vhost}
}

// Connect Create a connection
func (b *twitchBot) Connect() error {
	log.Println("Подключение к серверу")
	conn, err := net.Dial("tcp", b.server)
	if err != nil {
		return err
	}
	b.conn = conn
	return nil
}

// Authenticate Bot authentication in chat
func (b *twitchBot) Authenticate() error {
	log.Println("Авторизация на сервере")
	fmt.Fprintf(b.conn, "PASS %s\r\n", b.pass)
	fmt.Fprintf(b.conn, "NICK %s\r\n", b.nick)
	fmt.Fprintf(b.conn, "JOIN #%s\r\n", b.vhost)
	return nil
}

func (b *twitchBot) ReadMessages() error {
	log.Println("Запуск чтения сообщений")
	b.reader = bufio.NewReader(b.conn)
	for {
		message, err := b.reader.ReadString('\n')
		if err != nil {
			return err
		}
		b.HandleMessage(message)
	}
}

func (b *twitchBot) HandleMessage(message string) {
	// Handle ping/pong to keep connection alive
	if strings.Contains(message, "PING") {
		pong := strings.Replace(message, "PING", "PONG", 1)
		fmt.Fprintf(b.conn, "%s\r\n", pong)
	}
	log.Println(message)
}

// SendMessage Sending messages to chat
func (b *twitchBot) SendMessage(message string) error {
	_, err := fmt.Fprintf(b.conn, "PRIVMSG #%s :%s\r\n", b.vhost, message)
	return err
}

// Disconnect Chat disconnect handler
func (b *twitchBot) Disconnect() error {
	fmt.Fprintf(b.conn, "PART #%s\r\n", b.vhost)
	fmt.Fprintf(b.conn, "QUIT\r\n")
	fmt.Println("CLOSED SUCCESS")
	return b.conn.Close()
}
