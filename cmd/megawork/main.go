package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/internal/app"
	"github.com/rcy/megaworking/internal/db"
	_ "modernc.org/sqlite"
)

func main() {
	sqldb, err := sql.Open("sqlite", os.Getenv("SQLITE_DB"))
	if err != nil {
		panic(err)
	}

	q := db.New(sqldb)

	f, _ := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	log.SetOutput(f)

	p := tea.NewProgram(app.New(q), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
