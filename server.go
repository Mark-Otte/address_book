package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type Entry struct {
	ForeName    string `json:"forename"`
	SurName     string `json:"surname"`
	PhoneNumber int    `json:"phonenumber"`
}

type EntryHandlers struct {
	sync.Mutex
	store []Entry
}

func (h *EntryHandlers) ManageEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.post(w, r)
		return
	case "DELETE":
		h.delete(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *EntryHandlers) GetEntriesForeName(w http.ResponseWriter, r *http.Request) {
	entries := make([]Entry, len(h.store))

	h.Lock()
	sort.SliceStable(h.store, func(i, j int) bool {
		return h.store[i].ForeName < h.store[j].ForeName
	})

	i := 0
	for _, entry := range h.store {
		entries[i] = entry
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *EntryHandlers) GetEntriesSurName(w http.ResponseWriter, r *http.Request) {
	entries := make([]Entry, len(h.store))

	h.Lock()

	sort.SliceStable(h.store, func(i, j int) bool {
		return h.store[i].SurName < h.store[j].SurName
	})

	i := 0
	for _, entry := range h.store {
		entries[i] = entry
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *EntryHandlers) SearchEntries(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}
	type search struct {
		Searchword string `json:"searchterm"`
	}
	var searchterm search
	err = json.Unmarshal(bodyBytes, &searchterm)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	entries := make([]Entry, len(h.store))
	entriesToReturn := make([]Entry, 0)
	h.Lock()
	i := 0
	for _, entry := range h.store {
		entries[i] = entry
		if strings.Contains(entry.ForeName, searchterm.Searchword) {
			entriesToReturn = append(entriesToReturn, entry)
		}
		if strings.Contains(entry.SurName, searchterm.Searchword) {
			entriesToReturn = append(entriesToReturn, entry)
		}
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(entriesToReturn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *EntryHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var entry Entry
	err = json.Unmarshal(bodyBytes, &entry)
	fmt.Print(entry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	h.Lock()
	h.store = append(h.store, entry)
	defer h.Unlock()
}

func (h *EntryHandlers) delete(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}
	type search struct {
		Forename string `json:"forename"`
		Surname  string `json:"surname"`
	}
	var searchterm search
	err = json.Unmarshal(bodyBytes, &searchterm)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	entries := make([]Entry, len(h.store))
	entriesToKeep := make([]Entry, 0)
	h.Lock()
	i := 0
	for _, entry := range h.store {
		entries[i] = entry
		if entry.ForeName == searchterm.Forename {
			if entry.SurName == searchterm.Surname {
				//This is the slice we do not want
			} else {
				entriesToKeep = append(entriesToKeep, entry)
			}
		} else {
			entriesToKeep = append(entriesToKeep, entry)
		}
		i++
	}
	h.store = nil
	h.store = append(h.store, entriesToKeep...)
	h.Unlock()

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func newentryHandlers() *EntryHandlers {
	return &EntryHandlers{
		store: []Entry{},
	}
}

func main() {
	entryHandlers := newentryHandlers()
	http.HandleFunc("/entry", entryHandlers.ManageEntry)
	http.HandleFunc("/entriesfn", entryHandlers.GetEntriesForeName)
	http.HandleFunc("/entriessn", entryHandlers.GetEntriesSurName)
	http.HandleFunc("/search", entryHandlers.SearchEntries)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
