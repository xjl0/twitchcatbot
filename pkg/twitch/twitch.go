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
	HandleGiveVip(userWhoWasGivenVip string)
	HandleTakeVip(userWhoseVipWasTakenAway string)
	Disconnect() error
	Join(chat string) error
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

// Join chat
func (b *twitchBot) Join(chat string) error {
	fmt.Fprintf(b.conn, "JOIN #%s\r\n", chat)
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

// HandleGiveVip Handler for giving VIP status to the user
func (b *twitchBot) HandleGiveVip(userWhoWasGivenVip string) {
	b.SendMessage(fmt.Sprintf("/vip %s", userWhoWasGivenVip))
	log.Println("Gave VIP status to", userWhoWasGivenVip)
}

// HandleTakeVip Handler for taking the VIP status from the user
func (b *twitchBot) HandleTakeVip(userWhoseVipWasTakenAway string) {
	b.SendMessage(fmt.Sprintf("/unvip %s", userWhoseVipWasTakenAway))
	log.Println("Took VIP status from", userWhoseVipWasTakenAway)
}

// Disconnect Chat disconnect handler
func (b *twitchBot) Disconnect() error {
	fmt.Fprintf(b.conn, "PART #%s\r\n", b.vhost)
	fmt.Fprintf(b.conn, "QUIT\r\n")
	fmt.Println("CLOSED SUCCESS")
	return b.conn.Close()
}
