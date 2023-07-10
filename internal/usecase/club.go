package usecases

import (
	"bufio"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Club struct {
	OpenAt  time.Time
	CloseAt time.Time

	WaitQueue []string
	Clients   map[string]int
	Tables    map[int]*Table

	CostTablePerHouse int
}

type Table struct {
	UsedBy string

	LastBusyAt time.Time
	WasBusyFor time.Duration
	Profit     int
}

// Create new Club
func NewClub(openAt, closeAt time.Time, tableCount, cost int) *Club {
	clients := make(map[string]int)
	tables := make(map[int]*Table)
	for i := 1; i <= tableCount; i++ {
		tables[i] = &Table{}
	}
	queue := make([]string, 0)
	return &Club{
		OpenAt:            openAt,
		CloseAt:           closeAt,
		Clients:           clients,
		WaitQueue:         queue,
		Tables:            tables,
		CostTablePerHouse: cost,
	}
}

// NewClubFromScanner creates a new Club object from scanner input.
func NewClubFromScanner(sc *bufio.Scanner) (*Club, error) {
	sc.Scan()
	row := sc.Text()
	tableCount, err := strconv.Atoi(row)
	if err != nil {
		return nil, errors.New(row)
	}

	sc.Scan()
	row = sc.Text()
	times := strings.Split(row, " ")
	if len(times) != 2 {
		return nil, errors.New(row)
	}

	openAt, err := time.Parse(timeFormat, times[0])
	if err != nil {
		return nil, errors.New(row)
	}

	closeAt, err := time.Parse(timeFormat, times[1])
	if err != nil {
		return nil, errors.New(row)
	}

	sc.Scan()
	row = sc.Text()
	cost, err := strconv.Atoi(row)
	if err != nil {
		return nil, errors.New(row)
	}

	return NewClub(openAt, closeAt, tableCount, cost), nil
}

// ClientArrived is called when a client arrives at the club.
func (c *Club) ClientArrived(timeAt time.Time, clientName string) error {

	if _, ok := c.Clients[clientName]; ok {
		return errNotPass
	}

	if timeAt.Before(c.OpenAt) || timeAt.After(c.CloseAt) {
		return errNotOpen
	}

	c.Clients[clientName] = 0

	return nil
}

// ClientTookTable is called when a client takes a table
func (c *Club) ClientTookTable(timeAt time.Time, clientName string, computerNumber int) error {

	// checking the table for busy
	if table := c.Tables[computerNumber]; table.UsedBy != "" {
		return errPlaceBusy
	}

	// checking the client for being in the club
	if err := c.checkClient(clientName); err != nil {
		return err
	}

	// client occupies the computer
	c.Clients[clientName] = computerNumber
	table := c.Tables[computerNumber]
	table.UsedBy = clientName
	table.LastBusyAt = timeAt

	return nil
}

// ClientToQueue adds a client to the queue if the club is full.
func (c *Club) ClientToQueue(timeAt time.Time, clientName string) error {
	if err := c.checkClient(clientName); err != nil {
		return err
	}

	if _, ok := c.checkFreeTables(); ok {
		return errCanWait
	}

	if len(c.WaitQueue) > len(c.Tables) {
		//delete(c.Clients, clientName)
		if _, _, err := c.ClientLeft(timeAt, clientName); err != nil {
			return err
		}
		return errClientLeft
	}

	c.WaitQueue = append(c.WaitQueue, clientName)

	return nil
}

// ClientLeft is called when a client leaves the club.
func (c *Club) ClientLeft(timeAt time.Time, clientName string) (string, int, error) {
	if err := c.checkClient(clientName); err != nil {
		return "", 0, err
	}

	if computerNumber := c.Clients[clientName]; computerNumber != 0 {
		table := c.Tables[computerNumber]
		table.UsedBy = ""

		tookDuration := timeAt.Sub(table.LastBusyAt)
		table.WasBusyFor += tookDuration
		table.Profit += int(math.Ceil(tookDuration.Hours())) * c.CostTablePerHouse

		if len(c.WaitQueue) > 0 && timeAt.Before(c.CloseAt) {
			clientFromQueue := c.WaitQueue[0]
			c.WaitQueue = c.WaitQueue[1:]
			if err := c.ClientTookTable(timeAt, clientFromQueue, computerNumber); err != nil {
				return "", 0, err
			}

			delete(c.Clients, clientName)
			return clientFromQueue, computerNumber, errClientTookTableFromQ
		}
	}

	delete(c.Clients, clientName)

	return "", 0, nil

}

// GetAllClients returns a sorted slice of all clients in the club.
func (c *Club) GetAllClients() []string {
	arr := make([]string, 0)

	for key := range c.Clients {
		arr = append(arr, key)
	}

	sort.Strings(arr)

	return arr
}

// checkClient checks if a client is in the club or not.
func (c *Club) checkClient(clientName string) error {
	if _, ok := c.Clients[clientName]; !ok {
		return errUnknownClient
	}
	return nil
}

// checkFreeTables checks if there are any free tables in the club or not.
func (c *Club) checkFreeTables() (int, bool) {
	for key, value := range c.Tables {
		if value.UsedBy == "" {
			return key, true
		}
	}
	return 0, false

}
