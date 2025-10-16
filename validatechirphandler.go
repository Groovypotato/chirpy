package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func validatechirpHandler (w http.ResponseWriter, r *http.Request) {
	type returnError struct {
		Error string `json:"error"`
	}
	if r.Method != "POST"{
		errorMsg := returnError {
			Error: "method not allowed",
		}
		dat, err := json.Marshal(errorMsg)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(405)
			fmt.Fprintf(w, "error: method nto allowed\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(405)
		w.Write(dat)
		return
	}
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValid struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorMsg := returnError {
			Error: "Something went wrong",
		}
		dat, errs := json.Marshal(errorMsg)
		if errs != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			fmt.Fprintf(w, "error: incorrect format\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	if len(params.Body) > 140 {
		errorMsg := returnError {
			Error: "Chirp is too long",
		}
		dat, errs := json.Marshal(errorMsg)
		if errs != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			fmt.Fprintf(w, "error: Chirp is too long\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}else {
		validMsg := returnValid {
			Valid: true,
		}
		dat, err := json.Marshal(validMsg)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			fmt.Fprintf(w, "error: something went wrong\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}
}