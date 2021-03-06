package orderup

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

// Command struct.
type Cmd struct {
	Name string   // Command name
	Args []string // List of command arguments
}

type Orderup struct {
	db       *bolt.DB
	password string
}

func NewOrderup(dbFile, password string) (*Orderup, error) {
	db, err := initDb(dbFile)
	if err != nil {
		return nil, err
	}

	return &Orderup{
		db:       db,
		password: password,
	}, nil
}

// Serve web API.
func (o *Orderup) MakeAPI(apiVersion string, mux *mux.Router) {
	switch apiVersion {
	case V1:
		for _, route := range o.getAPIv1().Routes {
			mux.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Methods...)
		}

	default:
		panic("Unknown API version.")
	}
}

// Serve Slack API.
func (o *Orderup) MakeRequestHandler(mux *mux.Router) {
	mux.HandleFunc("/orderup", o.requestHandler)
}

// Open an initialize database.
func initDb(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(QUEUES))
		return err
	})

	return db, err
}

// Handle requests to orderup bot.
func (o *Orderup) requestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get command text from the request and split arguments.
	cmd := o.parseCmd(r.PostForm["text"][0])

	// Execute command

	response, inChannel, cmdErr := o.execCmd(cmd)

	if cmdErr != nil {
		switch cmdErr.ErrType {
		case ARG_ERR:
			response = o.errorMessage(cmdErr.Error())
		default:
			response = cmdErr.Error()
		}
	}

	if inChannel {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"response_type":"in_channel","text":"%s"}`, response))
	} else {
		io.WriteString(w, response)
	}
}

// Parse command from the request string.
func (o *Orderup) parseCmd(cmd string) *Cmd {
	if cmdLst := strings.Split(cmd, " "); len(cmdLst) == 1 {
		return &Cmd{
			Name: cmdLst[0],
		}
	} else {
		return &Cmd{
			Name: cmdLst[0],
			Args: cmdLst[1:],
		}
	}
}

// Execute command.
func (o *Orderup) execCmd(cmd *Cmd) (string, bool, *OrderupError) {
	switch cmd.Name {
	case CREATE_Q_CMD:
		return o.createQueueCmd(cmd)
	case DELETE_Q_CMD:
		return o.deleteQueueCmd(cmd)
	case CREATE_ORDER_CMD:
		return o.createOrderCmd(cmd)
	case FINISH_ORDER_CMD:
		return o.finishOrderCmd(cmd)
	case LIST_CMD:
		return o.listCmd(cmd)
	case HISTORY_CMD:
		return o.historyCmd(cmd)
	default:
		return o.helpCmd(cmd)
	}
}

// Safely close db and shutdown.
func (o *Orderup) Shutdown() {
	o.db.Close()
	log.Print("Bye!")
}
