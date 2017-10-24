package main
import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"log"
	"fmt"
	"time"
	"flag"
	"regexp"
)

var addr = flag.String("addr", "127.0.0.1:8082", "Rest Service")

//English Character only + “,” + “-“ + “.” + “ “
var IsChar = regexp.MustCompile(`^[a-zA-Z ,.-]+$`).MatchString

type Ewallet struct {
	Fullname		string	`json:"fullname"`
	Citizenid		string	`json:"citizenid"`
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func main(){
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	http.HandleFunc("/", homepage)
	http.HandleFunc("/ewallets", createEwallet(session))
	http.ListenAndServe(*addr, nil)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("eWallet").C("account")

	index := mgo.Index{
		Key:        []string{"wallet_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func checkCitizenID(cid string) bool {
	error := false
	if len(cid) != 13 {
		error = true
	}
	return error
}

func findCitizenID(s *mgo.Session, cid string) bool {
	error := false

	session := s.Copy()
	defer session.Close()

	return error
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to E-Wallet")
}

func createEwallet(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var ewallet Ewallet
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ewallet)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		fullname := ewallet.Fullname
		if IsChar(fullname) != true {
			ErrorWithJSON(w, "Invalid Full Name", http.StatusBadRequest)
			return
		}

		citizen := ewallet.Citizenid
		if checkCitizenID(citizen) != false {
			ErrorWithJSON(w, "Invalid Citizen ID", http.StatusBadRequest)
			return
		}

		c := session.DB("eWallet").C("account")

		err = c.Insert(bson.M{"wallet_id": "00003","citizen_id": citizen,"full_name": fullname,"open_datetime": time.Now(),"ledger_balance": "0.00"})
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "E-Wallet already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert E-Wallet: ", err)
			return
		}

		/*
		var resultTask Task
		resultTask.ID = task.ID
		resultTask.Title = task.Title
		resultTask.Done = task.Done

		if err := json.NewEncoder(w).Encode(resultTask); err != nil {
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
		*/
		//w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Location", r.URL.Path+"/"+book.ISBN)
		w.WriteHeader(http.StatusCreated)
	}
}

