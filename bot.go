package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

// Command struct.
type Cmd struct {
	Name string   // Command name
	Args []string // List of command arguments
}

type Orderup struct {
	db *bolt.DB
}

func NewOrderup(dbFile string) (*Orderup, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Orderup{
		db: db,
	}, nil
}

// Open an initialize database.
func initDb(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	if _, err := tx.CreateBucketIfNotExists([]byte(RESTAURANTS)); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return db, nil
}

// Handle requests to orderup bot.
func (o *Orderup) RequestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get command text from the request and split arguments.
	cmd := o.parseCmd(r.PostForm["text"][0])

	// Execute command
	response := o.execCmd(cmd)

	io.WriteString(w, response)
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
func (o *Orderup) execCmd(cmd *Cmd) string {
	switch cmd.Name {
	case "create-restaurant":
		return o.createRestaurant(cmd)
	default:
		return o.help(cmd)
	}
}

// Safely close db and shutdown.
func (o *Orderup) Shutdown() {
	o.db.Close()
	log.Print("Bye!")
}
