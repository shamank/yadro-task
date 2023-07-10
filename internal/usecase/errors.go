package usecases

import "errors"

var (
	// if a client comes, but he is already in the club
	errNotPass = errors.New("YouShallNotPass")

	// if the client came during the club's non-working hours
	errNotOpen = errors.New("NotOpenYet")

	// if a customer sits down at a table that is already occupied
	errPlaceBusy = errors.New("PlaceIsBusy")

	// if the user is not in the club
	errUnknownClient = errors.New("ClientUnknown")

	// if the user is waiting, although there are free tables
	errCanWait = errors.New("ICanWaitNoLonger!")

	// if the client leaves. required for outgoing event
	errClientLeft = errors.New("ClientLeft!")

	// if a client from the queue sits on a free table. required for outgoing event
	errClientTookTableFromQ = errors.New("ClientFromTheQueueTookATable")
)
