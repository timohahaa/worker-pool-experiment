package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/timohahaa/croc"
)

func main() {
	cc := croc.New()
	err := cc.Get("https://api.ipify.org?format=json").End()
	if err != nil {
		log.Fatal(err)
	}
	body := cc.RawRespBody()
	s := struct {
		Str string `json:"ip"`
	}{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Str)
}
