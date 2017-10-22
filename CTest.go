package main
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type rqCreate struct {
	WalletID		string	`json:"walletid"`
	CitizenID		string	`json:"citizenid"`
	FullName		string	`json:"fullname"`
	OpenDateTime	string	`json:"opendatetime"`
	LedgerBal		string	`json:"ledgerbal"`
}

func main(){
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	collection := session.DB("eWallet").C("customer")
	err = collection.Insert(bson.M{"name": "TestB","citizenID": "2222222222222"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert data success.")
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("eWallet").C("account")

	index := mgo.Index{
		Key:        []string{"walletid"},
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
