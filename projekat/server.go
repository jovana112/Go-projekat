package main

import (
	"errors"
	"github.com/gorilla/mux"
	cs "github.com/jovana112/Go-Projekat/projekat/configstore"
	"mime"
	"net/http"
	"net/url"
)

type Service struct {
	store *cs.ConfigStore
}

func (ts *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {

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

	rt, err := decodeBodyForConfig(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf, err := ts.store.CreateConfig(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, conf)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	configs, ok := ts.store.GetConfig(id, version)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, configs)
}

func (ts *Service) deleteConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	configs, ok := ts.store.DeleteConfig(id, version)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, configs)
}

func (ts *Service) getAllConfigHandler(w http.ResponseWriter, _ *http.Request) {
	allTasks, err := ts.store.GetAllConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}
func (ts *Service) updateConfigWithNewVersionHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

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

	rt, err := decodeBodyForConfig(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf, err := ts.store.UpdateConfigWithNewVersion(rt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, conf)
}

//GROUP

func (ts *Service) createGroupHandler(w http.ResponseWriter, req *http.Request) {

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

	rt, err := decodeBodyForGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := ts.store.CreateGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, group)
}
func (ts *Service) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	groups, ok := ts.store.GetGroup(id, version)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, groups)
}
func (ts *Service) getAllGroupHandler(w http.ResponseWriter, _ *http.Request) {
	allTasks, err := ts.store.GetAllGroup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}

func (ts *Service) deleteGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	groups, ok := ts.store.DeleteGroup(id, version)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, groups)
}
func (ts *Service) updateGroupWithNewVersionHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

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

	rt, err := decodeBodyForGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := ts.store.UpdateGroupWithNewVersion(rt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, group)
}

func (ts *Service) extendConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
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
	version := mux.Vars(req)["version"]
	group, ok := ts.store.GetGroup(id, version)
	if ok != nil {
		err := errors.New("Group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	rt, err := decodeBodyForConfigs(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, config := range rt {
		group.Configs = append(group.Configs, config.Entries)
	}

	updatedGroup, err := ts.store.UpdateGroup(group, id)
	if err != nil {
		return
	}
	renderJSON(w, updatedGroup)
}

func (ts *Service) getConfigsByLabelsHandler(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["version"]
	id := mux.Vars(req)["id"]

	req.ParseForm()
	params := url.Values.Encode(req.Form)
	labels, err := ts.store.GetConfigsByLabels(id, ver, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, labels)
}
