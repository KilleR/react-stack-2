package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to read index: "+err.Error())
		return
	}

	_, err = fmt.Fprint(w, string(index))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Faiiled to write output")
		return
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "test handler!")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
