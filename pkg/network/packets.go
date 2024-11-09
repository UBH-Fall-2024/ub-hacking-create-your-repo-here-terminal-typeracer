package network

// Requests are sent from the client to the server
type Request uint8

const (
	// When the client connects to the server initially
	//
	// "<name>"
	Connect Request = iota

	// When the client wants to connect to a lobby
	JoinLobby

	// (While in game) client's progress through the given text
	//
	// "<0-100%>" (neglect the %, it's just gonna be a number)
	Progress
)

// Events are sent from the server to the client to tell the client to do things
type Event uint8

const (
	// Sent to all clients in a lobby when a new person joins the lobby, and
	// also sent to the client that joined for every client already inside the
	// lobby
	//
	// "<id>,<name>"
	JoinedLobby = iota
	//Sent by the server to ask clients to send their progress
	//
	// ""
	ProgressPls
	// Sent to all clients when any client sends their progress. This will allow
	// each client to update the ui to reflect who is winning
	//
	// "<id>,<prog>"
	ProgUpdate
	// Will tell clients its starting and give them all their text
	//
	// "<text>"
	GameStart
	// Countdown before game starts, probably 3 of these.
	//
	//	"<num>"
	GameStarting
)

type Message struct {
	// Can either be Request or Event. It doesn't matter which, it will be
	// parsed accordingly clientside/serverside
	Header uint8
	Data   string
}
