package main
import (
	"fmt"
	"regexp"
	//"os"
	"strings"
)


func main() {
	var IsLetter = regexp.MustCompile(`^[a-zA-Z ,.-]+$`).MatchString
	var Name string = "Mr.Test Jung"
	var ZID string = "1234567890123"
	var zwseq string = "19"
	fmt.Println(Name)
	fmt.Println(ZID)
	fmt.Println(IsLetter(Name)) // true

	if IsLetter(Name) != true {
		return
	}
	nupper := strings.ToUpper(Name)
	fmt.Println(nupper)

	length := len(ZID)
	fmt.Println(length)
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