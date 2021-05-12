package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Species struct {
	ID            string  `json:"id"`
	GenericName   string  `json:"genericname"`
	SpecificName  string  `json:"specificname"`
	WorkersLength float64 `json:"workerslength"`
	QueenLength   float64 `json:"queenlength"`
}

type speciesHandlers struct {
	sync.Mutex
	store map[string]Species
}

func (h *speciesHandlers) get(w http.ResponseWriter, r *http.Request) {

	speciesList := make([]Species, len(h.store))

	h.Lock()
	i := 0
	for _, species := range h.store {
		speciesList[i] = species
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(speciesList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *speciesHandlers) getSpecies(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Lock()
	species, ok := h.store[parts[2]]
	h.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(species)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *speciesHandlers) post(w http.ResponseWriter, r *http.Request) {
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
	var species Species
	err = json.Unmarshal(bodyBytes, &species)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	species.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.store[species.ID] = species
	defer h.Unlock()
}

func (h *speciesHandlers) species(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func newSpeciesHandlers() *speciesHandlers {

	return &speciesHandlers{
		store: map[string]Species{},
	}
}

type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")

	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorized"))
		return
	}

	w.Write([]byte("<html><h1>You've gained authorized access to the admin portal</h1></html>"))
}

func main() {
	admin := newAdminPortal()
	speciesHandlers := newSpeciesHandlers()
	http.HandleFunc("/species", speciesHandlers.species)
	http.HandleFunc("/species/", speciesHandlers.getSpecies)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

/* json input data for testing

curl localhost:8080/species -X POST -d '{"genericname": "Lasius", "specificname": "niger", "workerslength": 4, "queenlength" : 9}' -H "Content-Type: application/json"

*/
