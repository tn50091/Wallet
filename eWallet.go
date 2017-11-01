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
	"strings"
	"strconv"
)

// Set url and port for Postman
var addr = flag.String("addr", "127.0.0.1:8082", "Rest Service")

//English Character only + “,” + “-“ + “.” + “ “
var IsChar = regexp.MustCompile(`^[a-zA-Z ,.-]+$`).MatchString

type Ewallet struct {
	Fullname		string	`json:"fullname"`
	Citizenid		string	`json:"citizenid"`
}
type WallSeq struct {
	Typ 	string `bson:"typ"`
	Seq		int	`bson:"seq"`
}
type RsBody struct {
	Walletid		string		`json:"walletid"`
	Opendatetime	time.Time	`json:"opendatetime"`
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

//func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	w.WriteHeader(code)
//	w.Write(json)
//}

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
		fname := strings.ToUpper(fullname)

		zwallid := genWalletID(s)
		fmt.Print(zwallid)

		c := session.DB("eWallet").C("account")

		err = c.Insert(bson.M{"wallet_id": zwallid ,"citizen_id": citizen,"full_name": fname,"open_datetime": time.Now(),"ledger_balance": "0.00"})
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "E-Wallet already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert E-Wallet: ", err)
			return
		}

		var rsbody RsBody
		rsbody.Walletid = zwallid
		rsbody.Opendatetime = time.Now()

		if err := json.NewEncoder(w).Encode(rsbody); err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			log.Println("Incorrect body: ", rsbody)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func checkCitizenID(cid string) bool {
	err := false
	if len(cid) != 13 {
		err = true
	}
	return err
}

func genWalletID(s *mgo.Session) string {
	//fmt.Println("a")
	session := s.Copy()
	defer session.Close()
	result :=  WallSeq {}
	c := session.DB("eWallet").C("WallSeq")
	err := c.Find(bson.M{"typ":"SEQ"}).One(&result)
	//fmt.Println(result.Seq)
	if err != nil {
		//fmt.Println("b")
		fmt.Println(err)
	}

	//fmt.Println("Results All: ", result)
	result.Seq ++

	if result.Seq ==1 {
		//fmt.Println("c")
		err = c.Insert(bson.M{"typ":"SEQ","seq": 1})
		if err != nil {
			//error
			//fmt.Println("d")
			fmt.Println(err)
		}
	} else {
		//fmt.Println("e")
		err = c.Update(bson.M{"typ":"SEQ"}, bson.M{"$set": bson.M{"seq": result.Seq}})
		if err != nil {
			//fmt.Println("f")
			fmt.Println(err)
		}
	}

	za := result.Seq
	zb := format(za)

	//fmt.Println("write zb"+strconv.Itoa(zb))
	//return result.Seq
	zc := fmt.Sprintf("%010d",za)
	return "1"+zc+strconv.Itoa(zb)
}

func format(zwseq int)int {

	zstr := strconv.Itoa(zwseq)
	padz := len(zstr)
	//fmt.Println(padz)
	//fmt.Println(zstr)
	zsum :=0
	for i := padz; i >= 1; i-- {
		//fmt.Println("for loop")
		//fmt.Println(zsum)
		//fmt.Println(i)
		substring := string(zstr[i-1])
		//fmt.Println(substring)
		zintcon,_ := strconv.Atoi(substring)
		//fmt.Println(zintcon)
		if i-1 == 1 {zsum += zintcon*11}
		if i-1 == 2 {zsum += zintcon*10}
		if i-1 == 3 {zsum += zintcon*9}
		if i-1 == 4 {zsum += zintcon*8}
		if i-1 == 5 {zsum += zintcon*7}
		if i-1 == 6 {zsum += zintcon*6}
		if i-1 == 7 {zsum += zintcon*5}
		if i-1 == 8 {zsum += zintcon*4}
		if i-1 == 9 {zsum += zintcon*3}
		if i-1 == 10 {zsum += zintcon*2}

		//fmt.Println(zsum)
	}
	//fmt.Println(zsum)
	return zsum%10
}
