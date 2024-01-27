package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	m "github.com/rcy/megaworking/internal/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	_ "modernc.org/sqlite"
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

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

	// if !discount {
	// 	fmt.Println("What? You didn’t take the discount?!")
	// }

	// p := tea.NewProgram(initialModel())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
}

var (
	burger       string
	toppings     []string
	sauceLevel   int
	name         string
	instructions string
	discount     bool
)

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

var form = huh.NewForm(
	huh.NewGroup(
		// Ask the user for a base burger and toppings.
		huh.NewSelect[string]().
			Title("Choose your burger").
			Options(
				huh.NewOption("Charmburger Classic", "classic"),
				huh.NewOption("Chickwich", "chickwich"),
				huh.NewOption("Fishburger", "fishburger"),
				huh.NewOption("Charmpossible™ Burger", "charmpossible"),
			).
			Value(&burger), // store the chosen option in the "burger" variable

		// Let the user select multiple toppings.
		huh.NewMultiSelect[string]().
			Title("Toppings").
			Options(
				huh.NewOption("Lettuce", "lettuce").Selected(true),
				huh.NewOption("Tomatoes", "tomatoes").Selected(true),
				huh.NewOption("Jalapeños", "jalapeños"),
				huh.NewOption("Cheese", "cheese"),
				huh.NewOption("Vegan Cheese", "vegan cheese"),
				huh.NewOption("Nutella", "nutella"),
			).
			Limit(4). // there’s a 4 topping limit!
			Value(&toppings),

		// Option values in selects and multi selects can be any type you
		// want. We’ve been recording strings above, but here we’ll store
		// answers as integers. Note the generic "[int]" directive below.
		huh.NewSelect[int]().
			Title("How much Charm Sauce do you want?").
			Options(
				huh.NewOption("None", 0),
				huh.NewOption("A little", 1),
				huh.NewOption("A lot", 2),
			).
			Value(&sauceLevel),
	),

	// Gather some final details about the order.
	huh.NewGroup(
		huh.NewInput().
			Title("What's your name?").
			Value(&name).
			// Validating fields is easy. The form will mark erroneous fields
			// and display error messages accordingly.
			Validate(func(str string) error {
				if str == "Frank" {
					return errors.New("Sorry, we don’t serve customers named Frank.")
				}
				return nil
			}),

		huh.NewText().
			Title("Special Instructions").
			CharLimit(400).
			Value(&instructions),

		huh.NewConfirm().
			Title("Would you like 15% off?").
			Value(&discount),
	),
)
