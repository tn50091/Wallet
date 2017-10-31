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
	"container/list"
	"strconv"
	"strings"
)

// Set url and port for Postman
var addr = flag.String("addr", "127.0.0.1:8082", "Rest Service")

// English Character only + “,” + “-“ + “.” + “ “
var IsChar = regexp.MustCompile(`^[a-zA-Z ,.-]+$`).MatchString

// Create Wallet Account
type RqBodyCreate struct {
	Rqbodycre struct {
		Fullname	string		`json:"fullname"`
		Citizenid	string		`json:"citizenid"`
	} `json:"rqbodycre"`
}
type RsBodyCreate struct {
	Rsbodycre struct{
		Walletid		string		`json:"walletid"`
		Opendatetime	time.Time	`json:"opendatetime"`
	}`json:"rsbodycre"`
	Rserrcre		[]Error		`json:"rserrcre"`
}

type Error struct {
	Listerror		list.List	`json:"listerror"`
	Errorcode		string		`json:"errorcode"`
	Errordesc		string		`json:"errordesc"`
}

type Seq struct {
	Widseq			string
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func main(){
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	http.HandleFunc("/", homepage)
	http.HandleFunc("/ewallets", createEwallet(session))
	http.ListenAndServe(*addr, nil)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("eWallet").C("account")

	index := mgo.Index{
		Key:        []string{"citizen_id"},
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

func genSeqWID(s *mgo.Session) string {
	session := s.Copy()
	defer session.Close()

	c := session.DB("eWallet").C("cidseq")

	var seq []Seq
	err := c.Find(bson.M{"wid_seq": bson.M{"$exists": true}}).All(&seq)
	if err != nil {
		return "Error"
	}
	log.Println("Test1: ", seq)

	var wids Seq
	var widseqs string
	var widseqi int
	log.Println("Len: ", len(seq))
	for i:=0; i< len(seq); i++ {
		wids = seq[i]
		widseqs = wids.Widseq
		log.Println("Test2: ", wids)
		log.Println("Test3: ", widseqs)
	}
	if widseqs == "" {
		widseqi=1

		c = session.DB("eWallet").C("cidseq")
		err = c.Insert(bson.M{"wid_seq": "0000000000"})
		if err != nil {
			return "Error"
		}
	} else {
		widseqi, _ = strconv.Atoi(widseqs)
		widseqi = widseqi + 1
	}

	var widseq string
	widseqs = strconv.Itoa(widseqi)

	widseq = leftPad(widseqs, "0", 10)

	return widseq
}

func leftPad(s string, padStr string, pLen int) string {
	pLen = pLen-len(s)
	return strings.Repeat(padStr, pLen) + s
}

func genWalletID(widseq string) string {
	var widseqs string
	var widseqi int
	widseqs = "1" + widseq

	seq := widseqs
	widseqi = 1
	var x string
	var y int
	for i:=1; i< len(seq); i++ {
		x = seq[i:i+1]
		y,_ = strconv.Atoi(x)
		widseqi = widseqi + (y * (i + 1))
	}

	widseqi = widseqi % 10
	widseqs = widseqs + (strconv.Itoa(widseqi))

	return widseqs
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to E-Wallet")
}

func createEwallet(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var rqbodycreate RqBodyCreate
		err := json.NewDecoder(r.Body).Decode(&rqbodycreate)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			log.Println("Incorrect body: ", rqbodycreate)
			return
		}
		log.Println("--------------------------------------------------")
		log.Println("Create Wallet Account")

		var fullname string
		var citizenid string
		fullname = rqbodycreate.Rqbodycre.Fullname
		citizenid = rqbodycreate.Rqbodycre.Citizenid
		log.Println("Full Name: ", fullname)
		log.Println("Citizen ID: ", citizenid)

		if IsChar(fullname) != true {
			ErrorWithJSON(w, "Invalid Full Name", http.StatusBadRequest)
			log.Println("Invalid Full Name: ", fullname)
			return
		}

		if checkCitizenID(citizenid) != false {
			ErrorWithJSON(w, "Invalid Citizen ID", http.StatusBadRequest)
			log.Println("Invalid Citizen ID: ", citizenid)
			return
		}

		wid := genSeqWID(session)
		if wid == "Error" {
			ErrorWithJSON(w, "Invalid Wallet ID", http.StatusInternalServerError)
			log.Println("Invalid Wallet ID: ", wid)
			return
		}
		walletid := genWalletID(wid)
		log.Println("Wallet ID: ", walletid)

		opendt := time.Now()
		log.Println("Open DateTime: ", opendt)

		ensureIndex(session)
		c := session.DB("eWallet").C("account")
		err = c.Insert(bson.M{"wallet_id": walletid,"citizen_id": citizenid,"full_name": fullname,"open_datetime": opendt,"ledger_balance": "0.00"})
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Citizen ID already exists", http.StatusBadRequest)
				log.Println("Citizen ID already exists: ", err)
				return
			}

			ErrorWithJSON(w, "Failed insert E-Wallet", http.StatusInternalServerError)
			log.Println("Failed insert E-Wallet: ", err)
			return
		}

		c = session.DB("eWallet").C("cidseq")
		err = c.Update(nil,bson.M{"wid_seq": wid})
		if err != nil {
			ErrorWithJSON(w, "Failed insert E-Wallet(seq)", http.StatusInternalServerError)
			log.Println("Failed insert E-Wallet(seq): ", err)
			return
		}

		var rsbodycreate RsBodyCreate
		rsbodycreate.Rsbodycre.Walletid = walletid
		rsbodycreate.Rsbodycre.Opendatetime = opendt

		if err := json.NewEncoder(w).Encode(rsbodycreate); err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			log.Println("Incorrect body: ", rsbodycreate)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		log.Println("--------------------------------------------------")
	}
}

