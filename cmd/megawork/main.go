package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	m "github.com/rcy/megaworking/internal/model"

	"github.com/charmbracelet/huh"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", os.Getenv("SQLITE_DB"))
	if err != nil {
		panic(err)
	}

	q := m.New(db)

	err = prepForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	prep, err := q.CreatePreparation(context.TODO(), prepParams)
	if err != nil {
		panic(err)
	}

	err = cycleForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	cycle, err := q.CreateCycle(context.TODO(), cycleParams)
	if err != nil {
		panic(err)
	}
	_, _ = cycle, prep
}

var prepParams = m.CreatePreparationParams{}

var prepForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("What am I trying to accomplish?").
			Value(&prepParams.Accomplish),
		huh.NewInput().
			Title("Why is this important and valuable?").
			Value(&prepParams.Important),
		huh.NewInput().
			Title("How will I know when this is complete?").
			Value(&prepParams.Complete),
		huh.NewInput().
			Title("Any risks / hazards? Potential distractions, procrastination, etc.").
			Value(&prepParams.Distractions),
		huh.NewInput().
			Title("Is this concrete / measurable or subjective / ambiguous?").
			Value(&prepParams.Measurable),
		huh.NewInput().
			Title("Anything else noteworthy?").
			Value(&prepParams.Noteworthy),
	),
)

var cycleParams = m.CreateCycleParams{}

var cycleForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("What am I trying to accomplish this cycle?").
			Value(&cycleParams.Accomplish),
		huh.NewInput().
			Title("How will I get started?").
			Value(&cycleParams.Started),
		huh.NewInput().
			Title("Any hazards present?").
			Value(&cycleParams.Hazards),
		huh.NewSelect[int64]().
			Title("Energy").
			Options(
				huh.NewOption("High", int64(1)),
				huh.NewOption("Medium", int64(0)),
				huh.NewOption("Low", int64(-1)),
			).
			Value(&cycleParams.Energy),
		huh.NewSelect[int64]().
			Title("Morale").
			Options(
				huh.NewOption("High", int64(1)),
				huh.NewOption("Medium", int64(0)),
				huh.NewOption("Low", int64(-1)),
			).
			Value(&cycleParams.Morale),
	),
)
