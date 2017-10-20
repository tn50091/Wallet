package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Name struct {
	FullName string `json:"FullName"`
	CitizenID    int    `json:"CitizenID"`
	//ID    int    `json:"id"`
	//Title string `json:"title"`
	//Url   string `json:"url"`
}
func main() {
	//สร้าง struct ว่างๆเพื่อเอาไว้รับค่า
	todos := make(map[string]Name, 0)
	//อ่านไฟล์
	content, err := ioutil.ReadFile("name1.json")
	if err == nil {
		//unmarshall ข้อมูลที่ได้เก็บลงไปใน todos
		json.Unmarshal(content, &todos)
		for _, todo := range todos {
			fmt.Println("Name: ", todo)
		}
	}
}