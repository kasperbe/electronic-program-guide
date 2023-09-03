package epg

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const input = `
{
 
 "monday": [ 
 { 
 "title": "Nyhederne", 
 "state": "begin", 
 "time": 21600 
 }, 
 { 
 "title": "Nyhederne", 
 "state": "end", 
 "time": 36000 
 }, 
 { 
 "title": "Dybvaaaaad", 
 "state": "begin", 
 "time": 36000 
 }, 
 { 
 "title": "Dybvaaaaad", 
 "state": "end", 
 "time": 38100 
 } 
 ], 
 "tuesday": [], 
 "wednesday": [ 
 { 
 "title": "Nyhederne", 
 "state": "begin", 
 "time": 21600 
 }, 
 { 
 "title": "Nyhederne", 
 "state": "end", 
 "time": 43200 
 }, 
 { 
 "title": "Fodbold", 
 "state": "begin", 
 "time": 50400 
 }, 
 { 
 "title": "Fodbold", 
 "state": "end", 
 "time": 55800 
 }, 
 { 
 "title": "Nyhederne", 
 "state": "begin", 
 "time": 75600 
 }, 
 { 
 "title": "Nyhederne", 
 "state": "end", 
 "time": 77400 
 } 
 ], 
 "thursday": [ 
 { 
 "title": "ESL", 
 "state": "begin", 
 "time": 43200 
 }, 
 { 
 "title": "ESL",
 "state": "end", 
 "time": 46800 
 }, 
 {
 "title": "ESLPro", 
 "state": "begin", 
 "time": 82800 
 } 
 ], 
 "friday": [ 
 { 
 "title": "ESLPro", 
 "state": "end", 
 "time": 3600 
 } 
 ], 
 "saturday": [ 
 { 
 "title": "Comedy", 
 "state": "begin", 
 "time": 52200 
 }, 
 { 
 "title": "Comedy", 
 "state": "end", 
 "time": 59400 
 }, 
 { 
 "title": "Nybyggerne", 
 "state": "begin", 
 "time": 81600 
 } 
 ], 
 "sunday": [ 
 { 
 "title": "Nybyggerne", 
 "state": "end", 
 "time": 5400 
 }, 
 { 
 "title": "Dybvvvvvad", 
 "state": "begin", 
 "time": 41400 
 }, 
 { 
 "title": "Dybvvvvvad", 
 "state": "end", 
 "time": 43200 
 } 
 ] 

}
`

func TestTranslator(t *testing.T) {
	var in TranslateInput
	err := json.Unmarshal([]byte(input), &in)
	if err != nil {
		log.Fatalf("unmarshal: %v", err)
	}

	pg := New(in)

	if err != nil {
		log.Fatalf("translate: %v", err)
	}

	expected := [][]*Program{
		{
			{Title: "Nyhederne", Times: []Slot{{Start: 21600, End: 36000}}, Day: 0},
			{Title: "Dybvaaaaad", Times: []Slot{{Start: 36000, End: 38100}}, Day: 0},
		},
		nil,
		{
			{Title: "Nyhederne", Times: []Slot{{Start: 21600, End: 43200}, {Start: 75600, End: 77400}}, Day: 2},
			{Title: "Fodbold", Times: []Slot{{Start: 50400, End: 55800}}, Day: 2},
		},
		{
			{Title: "ESL", Times: []Slot{{Start: 43200, End: 46800}}, Day: 3},
			{Title: "ESLPro", Times: []Slot{{Start: 82800, End: 3600}}, Day: 3},
		},
		nil,
		{
			{Title: "Comedy", Times: []Slot{{Start: 52200, End: 59400}}, Day: 5},
			{Title: "Nybyggerne", Times: []Slot{{Start: 81600, End: 5400}}, Day: 5},
		},
		{
			{Title: "Dybvvvvvad", Times: []Slot{{Start: 41400, End: 43200}}, Day: 6},
		},
	}

	assert.Equal(t, expected, pg.program)
}

func TestTranslateToString(t *testing.T) {
	var in TranslateInput
	err := json.Unmarshal([]byte(input), &in)
	if err != nil {
		log.Fatalf("unmarshal: %v", err)
	}

	pg := New(in)

	if err != nil {
		log.Fatalf("translate: %v", err)
	}

	expected := `Monday: Nyhederne 7 - 11 / Dybvaaaaad 11 - 11:35
Tuesday: Nothing aired today
Wednesday: Nyhederne 7 - 13,22 - 22:30 / Fodbold 15 - 16:30
Thursday: ESL 13 - 14 / ESLPro 0 - 2
Friday: Nothing aired today
Saturday: Comedy 15:30 - 17:30 / Nybyggerne 23:40 - 2:30
Sunday: Dybvvvvvad 12:30 - 13`

	assert.Equal(t, expected, pg.ToString())
}
