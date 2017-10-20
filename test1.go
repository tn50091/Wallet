package main
import (
	"fmt"
	"regexp"
	//"os"
	"strings"
	//"log"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)
//English Character only + “,” + “-“ + “.” + “ “
var IsLetter = regexp.MustCompile(`^[a-zA-Z ,.-]+$`).MatchString

type Person struct {
	Name string
	ZID string

}

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

}