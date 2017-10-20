package main
import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
func main(){
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	collection := session.DB("eWallet").C("customer")
	err = collection.Insert(bson.M{"name": "TestB","citizenID": "2222222222222"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert data success.")
}