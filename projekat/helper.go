package main

import (
	"encoding/json"
	cs "github.com/jovana112/Go-Projekat/projekat/configstore"
	"io"
	"net/http"
)

func decodeBodyForConfig(r io.Reader) (*cs.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *cs.Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeBodyForConfigs(r io.Reader) ([]*cs.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt []*cs.Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeBodyForGroup(r io.Reader) (*cs.Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *cs.Group
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
