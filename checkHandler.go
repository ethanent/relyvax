package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethanent/calronavaxverif"
	"gopkg.in/square/go-jose.v2"
	"io"
	"io/ioutil"
	"net/http"
)

var caKeySet *jose.JSONWebKeySet

func checkHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rdr := io.LimitReader(r.Body, 10000)

	d, err := ioutil.ReadAll(rdr)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	summ, err := calronavaxverif.GenerateSummary(string(d), caKeySet)

	if err != nil {
		http.Error(w, err.Error(), 401)
	}

	rd, err := json.Marshal(summ)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(rd)
}
