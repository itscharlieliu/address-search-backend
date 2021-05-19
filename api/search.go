package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type test struct {
	Test1 string
	Test2 int
}

func (handler *BaseHandler) Search(w http.ResponseWriter, r *http.Request) {

	fmt.Println(handler.addresses)

	bytes, err := json.Marshal(test{Test1: "hello", Test2: 2})

	if err != nil {
		log.Fatal("Unable to generate json: " + err.Error())
	}

	w.Write([]byte(bytes))
}
