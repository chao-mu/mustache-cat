package main

import (
	"log"
	"net/http"
)

func main() {
	catsH, err := NewCatHandler("images")
	if err != nil {
		log.Fatal(err)
	}

	/*
		dataH, err := NewDataHandler("eth0")
		if err != nil {
			log.Fatal(err)
		}
	*/

	m := http.NewServeMux()

	//m.Handle("/data/", dataH.Handler())
	m.Handle("/", catsH.Handler())

	if err := http.ListenAndServe("0.0.0.0:80", m); err != nil {
		log.Fatal(err)
	}
}
