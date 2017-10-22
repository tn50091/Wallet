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
	collection := session.DB("walletacc").C("account")
	err = collection.Insert(bson.M{"wallet_id": "001"},bson.M{"citizen_id": "4705104043181"},
		bson.M{"full_name": "Somchai Suksiri"},bson.M{"open_datetime": "20/10/2017"},bson.M{"ledger_balance": "0.00"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert data success.")
}