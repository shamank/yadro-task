package usecases

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Incoming Event IDs
const (
	EInClientArrived = iota + 1
	EInClientTookTable
	EInClientWaiting
	EInClientLeft
)

// Outgoing Event IDs
const (
	EOutClientLeft = iota + 11
	EOutClientTookTable
	EOutError
)

type Event struct {
	ID        int
	TriggerAt time.Time
	Body      []string
}

// Create new Event
func NewEvent(ID int, triggerAt time.Time, body []string) *Event {
	return &Event{
		ID:        ID,
		TriggerAt: triggerAt,
		Body:      body,
	}
}

// Create an event from a string
func NewEventFromString(row string) (*Event, error) {
	data := strings.Split(row, " ")
	if len(data) < 3 || len(data) > 4 {
		return nil, errors.New(row)
	}
	triggerAt, err := time.Parse(timeFormat, data[0])
	if err != nil {
		return nil, errors.New(row)
	}

	eventID, err := strconv.Atoi(data[1])
	if err != nil {
		return nil, errors.New(row)
	}

	body := data[2:]

	if eventID == EInClientTookTable && len(body) != 2 {
		return nil, errors.New(row)
	}

	if match, _ := regexp.MatchString("^[a-z0-9_-]*$", body[0]); !match {
		return nil, errors.New(row)
	}

	return NewEvent(eventID, triggerAt, body), nil
}

// NewEventsFromScanner create events from a bufio.Scanner
func NewEventsFromScanner(sc *bufio.Scanner) ([]*Event, error) {

	events := make([]*Event, 0)

	for sc.Scan() {
		ev, err := NewEventFromString(sc.Text())
		if err != nil {
			return nil, err
		}
		events = append(events, ev)
	}
	return events, nil
}

// GenerateOutString Generate string for output
func (e *Event) GenerateOutString(ID int, msg string) string {

	return fmt.Sprintf("%s %d %s",
		e.TriggerAt.Format(timeFormat),
		ID,
		msg)
}

// ToString Convert Event to string
func (e *Event) ToString() string {
	return e.GenerateOutString(e.ID, strings.Join(e.Body, " "))
}

// Error Generate error message
func (e *Event) Error(err error) string {
	return e.GenerateOutString(EOutError, err.Error())
}
