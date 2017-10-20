package main
import (
	"fmt"
	"regexp"
	//"os"
	"strings"
	//"log"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"time"
)

type Kwallet struct {
	Name 		string `json:"Name",bson:"Name"`
	CitizenId          string `json:"CitizenId",bson:"_CitizenId,omitempty"`
	WID  int    `json:"WID,bson:"WID"`
	OpenDateTime        time.Time `json:"OpenDateTime",bson:"OpenDateTime"`
	Balance     string `json:"Balance",bson:"Balance"`
}
//English Character only + “,” + “-“ + “.” + “ “
var IsLetter = regexp.MustCompile(`^[a-zA-Z ,.-]+$`).MatchString

func main() {
	var Name string = "Mr.Test Jung"
	var ZID string = "1234567890123"
	var zwseq string = "19"
	fmt.Println(Name)
	fmt.Println(ZID)
	fmt.Println(IsLetter(Name)) // true

	if IsLetter(Name) != true {
		return
	}
	//Upper case
	nupper := strings.ToUpper(Name)
	fmt.Println(nupper)

	length := len(ZID)
	fmt.Println(length)
	//Citizen ID lenght 13
	if length != 13 {
		return
	}
	//ttttt
	fmt.Println(len(zwseq))
	padz := 10-len(zwseq)+1
	fmt.Println(padz)
	for i := padz; i <= 10; i++ {
		fpadz := 10-i+1
		fmt.Println(fpadz)
		//substring := zwseq[fpadz:fpadz]
		//fmt.Println(substring)
		//chkdigit += i*substring
	}
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/Kwallet")

	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("Kwallet").C("insrt")

	doc := Kwallet{
		Name:           Name,
		CitizenId:       ZID,
		WID: 1,
		OpenDateTime:    time.Now(),
		Balance:   "0.00",

	}

	err = c.Insert(doc)
	if err != nil {
		panic(err)
	}

}