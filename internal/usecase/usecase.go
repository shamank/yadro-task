package usecases

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

const timeFormat = "15:04"

// StartWork function handles the work process in the club
func StartWork(sc *bufio.Scanner, ch chan string) {
	// close channel after function finish
	defer close(ch)

	club, err := NewClubFromScanner(sc)
	if err != nil {
		ch <- err.Error()
		return
	}

	events, err := NewEventsFromScanner(sc)
	if err != nil {
		ch <- err.Error()
		return
	}

	ch <- club.OpenAt.Format(timeFormat)

	// looping through each event
	// and the result of processing the event (if any) are sent to the channel
	for _, event := range events {
		ch <- event.ToString()
		result := processEvent(event, club)
		if result != "" {
			ch <- result
		}

	}

	// creating an exit event from the club for all remaining customers
	for _, client := range club.GetAllClients() {
		client := client

		ev := Event{
			ID:        EOutClientLeft,
			TriggerAt: club.CloseAt,
			Body:      []string{client},
		}

		ch <- ev.ToString()

		_, _, err := club.ClientLeft(ev.TriggerAt, ev.Body[0])
		if err != nil {
			ch <- ev.Error(err)
		}
	}

	ch <- club.CloseAt.Format(timeFormat)

	// calculation of profit and work time for each table
	for i := 1; i <= len(club.Tables); i++ {
		table := club.Tables[i]

		h := int64(table.WasBusyFor.Hours())
		m := int64(table.WasBusyFor.Minutes()) % 60

		ch <- fmt.Sprintf("%d %d %02d:%02d", i, table.Profit, h, m)
	}

	return
}

// processEvent function handles how to process each event based on its ID
func processEvent(event *Event, club *Club) string {
	switch event.ID {
	case EInClientArrived:
		err := club.ClientArrived(event.TriggerAt, event.Body[0])
		if err != nil {
			return event.Error(err)
		}

	case EInClientTookTable:
		clientName := event.Body[0]
		computerNumber, err := strconv.Atoi(event.Body[1])
		if err != nil {
			return event.Error(err)
		}
		err = club.ClientTookTable(event.TriggerAt, clientName, computerNumber)
		if err != nil {
			return event.Error(err)
		}

	case EInClientWaiting:
		err := club.ClientToQueue(event.TriggerAt, event.Body[0])
		if err != nil {
			if errors.Is(err, errClientLeft) {
				return event.GenerateOutString(EOutClientLeft, event.Body[0])
			} else {
				return event.Error(err)
			}
		}

	case EInClientLeft:
		client, table, err := club.ClientLeft(event.TriggerAt, event.Body[0])

		if err != nil {
			if errors.Is(err, errClientTookTableFromQ) {

				ev := Event{
					ID:        EOutClientTookTable,
					TriggerAt: event.TriggerAt,
					Body:      []string{client, strconv.Itoa(table)},
				}
				return ev.ToString()
			} else {
				return event.Error(err)
			}

		}
	}
	return ""
}
