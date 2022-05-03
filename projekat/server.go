package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type Servis struct {
	Data map[string][]*Config
}

func (ts *Servis) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	for _, v := range rt {
		ts.Data[id] = append(ts.Data[id], v)
	}

	renderJSON(w, rt)
}

func (ts *Servis) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	configs, ok := ts.Data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, configs)
}
func (ts *Servis) deleteConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	configs, ok := ts.Data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	delete(ts.Data, id)
	renderJSON(w, configs)
}

func (ts *Servis) getAllHandler(w http.ResponseWriter, _ *http.Request) {
	renderJSON(w, ts.Data)
}

func (ts *Servis) extendConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	id := mux.Vars(req)["id"]
	configs, ok := ts.Data[id]
	if !ok {
		err := errors.New("config not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if len(configs) == 1 {
		err := errors.New("only config groups can be extended")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, conf := range rt {
		configs = append(configs, conf)
	}
	renderJSON(w, configs)
}
