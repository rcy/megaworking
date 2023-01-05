package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	irc "github.com/thoj/go-ircevent"
)

func getenv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s not set!", key)
	} else {
		log.Printf("%s=%s\n", key, val)
	}

	return val
}

func ircmain(nick, channel, server string) (*irc.Connection, error) {
	ircnick1 := nick
	irccon := irc.IRC(ircnick1, "github.com/rcy/mxwork")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Privmsgf("nickserv", "identify %s", getenv("IRC_NICKSERV_PASSWORD"))
		irccon.Join(channel)
	})

	err := irccon.Connect(server)

	return irccon, err
}

const MINUTE = 60
const FOCUS_LENGTH = 30 * MINUTE
const BREAK_LENGTH = 10 * MINUTE

const START_BREAK_AT = 0
const END_BREAK_AT = BREAK_LENGTH
const START_FOCUS_AT = END_BREAK_AT
const END_FOCUS_AT = START_FOCUS_AT + FOCUS_LENGTH
const TOTAL_CYCLE_SECONDS = END_FOCUS_AT

var users string

func main() {
	channel := getenv("IRC_CHANNEL")
	nick := getenv("IRC_NICK")
	conn, err := ircmain(nick, channel, getenv("IRC_SERVER"))
	if err != nil {
		log.Fatal(err)
	}

	conn.AddCallback("353", func(e *irc.Event) {
		users = e.Arguments[len(e.Arguments)-1]
	})

	initialized := false
	conn.AddCallback("366", func(e *irc.Event) {
		if !initialized {
			go loop(conn, channel)

			initialized = true
		}
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick != nick {
			//conn.Privmsgf(channel, "hello %s", e.Nick)
			conn.SendRawf("NAMES %s", channel)
		}
	})

	conn.AddCallback("PART", func(e *irc.Event) {
		if e.Nick != nick {
			//conn.Privmsgf(channel, "goodbye %s", e.Nick)
			conn.SendRawf("NAMES %s", channel)
		}
	})

	conn.Loop()
}

func loop(conn *irc.Connection, channel string) {
	// HACK: wait until chanserv gives ops
	time.Sleep(10 * time.Second)

	state := "init"

	keepTime(func(newState string, secondsRemaining int) {
		minutesRemaining := secondsRemaining / 60

		if state != newState {
			if newState == "break" {
				conn.Privmsgf("chanserv", fmt.Sprintf("unquiet %s *!*@*", channel))
				time.Sleep(5 * time.Second)
				conn.SendRawf("TOPIC %s :BREAK UNTIL %v", channel, time.Now().Add(time.Second*time.Duration(secondsRemaining)).Format("15:04 MST"))
				conn.Privmsgf(channel, "BREAK for %d minutes %s\n", minutesRemaining, users)
			} else if newState == "focus" {
				conn.Privmsg("chanserv", fmt.Sprintf("quiet %s *!*@*", channel))
				time.Sleep(5 * time.Second)
				conn.SendRawf("TOPIC %s :FOCUS UNTIL %v", channel, time.Now().Add(time.Second*time.Duration(secondsRemaining)).Format("15:04 MST"))
				conn.Privmsgf(channel, "FOCUS for %d minutes %s\n", minutesRemaining, users)
			}

			state = newState
		}
	})
}

func origin() time.Time {
	return time.Date(2022, time.June, 0, 1, 0, 0, 0, time.UTC) // top of odd hours is a beginning of 10 minute break
}

func keepTime(callback func(string, int)) {
	c := time.Tick(1 * time.Second)

	for range c {
		since := time.Since(origin())
		cs := int(since.Seconds()) % TOTAL_CYCLE_SECONDS

		var state string
		var remaining int

		if cs < END_BREAK_AT {
			state = "break"
			remaining = END_BREAK_AT - cs
		} else {
			state = "focus"
			remaining = END_FOCUS_AT - cs
		}

		callback(state, remaining)
	}
}
