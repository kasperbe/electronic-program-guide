package epg

import (
	"fmt"
	"strings"
	"time"
)

type Slot struct {
	Start int
	End   int
}

type Program struct {
	Title string
	State string
	Times []Slot
	Day   int
}

type InputProgram struct {
	Title string
	State string
	Time  int
}

type TranslateInput struct {
	Monday    []InputProgram
	Tuesday   []InputProgram
	Wednesday []InputProgram
	Thursday  []InputProgram
	Friday    []InputProgram
	Saturday  []InputProgram
	Sunday    []InputProgram
}

type EPG struct {
	program [][]*Program
}

// New translates some input to a domain representation of a Program.
// This is done in O(n) time and space on average, by storing a reference to each entry in a map
// And inserting this reference into an ordered slice, so we don't have to scan our slice for each iteration.
//
// We assume correct data here, we could easily add validation as a part of this algorithm as well, but since the rules for correct data
// is not clearly laid out as a part of the program, I will refrain from doing so, however this is why I chose to split up New and ToString.
func New(input TranslateInput) *EPG {
	program := map[string]*Program{}
	out := make([][]*Program, 7)

	// Order the input by day.
	days := [][]InputProgram{
		input.Monday,
		input.Tuesday,
		input.Wednesday,
		input.Thursday,
		input.Friday,
		input.Saturday,
		input.Sunday,
	}

	for currentDay, slots := range days {
		for _, slot := range slots {
			if slot.State == "begin" {
				if existing, ok := program[slot.Title]; ok && existing.Day == currentDay {
					existing.Times = append(existing.Times, Slot{Start: slot.Time})
				} else {
					x := &Program{
						Title: slot.Title,
						Times: []Slot{
							{Start: slot.Time},
						},

						Day: currentDay,
					}

					program[slot.Title] = x
					out[currentDay] = append(out[currentDay], x)
				}
			} else {
				x := program[slot.Title]
				x.Times[len(x.Times)-1].End = slot.Time
			}
		}

	}

	return &EPG{
		program: out,
	}
}

func (epg *EPG) ToString() string {
	out := ""

	for day, program := range epg.program {
		if len(program) == 0 {
			out += "Nothing aired today\n"
			continue
		}

		out += formatDayIndexToString(day)

		for i, slot := range program {
			times := []string{}
			for _, t := range slot.Times {
				start := time.Unix(int64(t.Start), 0)
				end := time.Unix(int64(t.End), 0)

				times = append(times, fmt.Sprintf("%s - %s", formatTime(start), formatTime(end)))
			}

			out += fmt.Sprintf("%s %s", slot.Title, strings.Join(times, ","))
			if i != len(program)-1 {
				out += " / "
			}
		}

		out += "\n"
	}

	return strings.Trim(out, "\n ")
}

func formatTime(t time.Time) string {
	if t.Minute() == 0 {
		return fmt.Sprintf("%d", t.Hour())
	}

	return fmt.Sprintf("%d:%d", t.Hour(), t.Minute())
}

func formatDayIndexToString(day int) string {
	switch day {
	case 0:
		return "Monday: "
	case 1:
		return "Tuesday: "
	case 2:
		return "Wednesday: "
	case 3:
		return "Thursday: "
	case 4:
		return "Friday: "
	case 5:
		return "Satuday: "
	case 6:
		return "Sunday: "
	default:
		// In reality we will never get here, but we will handle the case gracefully nonetheless.
		fmt.Printf("Indexing unknown day, got %d\n", day)
		return "Unknown day: "
	}
}
