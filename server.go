package main

import "net/http"
import "encoding/json"
import "sync"
import "io/ioutil"
import "time"
import "fmt"

type Species struct{
    ID string `json:"id"`
    GenericName string `json:"genericname"`
    SpecificName string `json:"specificname"`
    WorkersLength float64 `json:"workerslength"`
    QueenLength float64 `json:"queenlength"`
}

type speciesHandlers struct{
    sync.Mutex
    store map[string]Species

}
func (h *speciesHandlers) get(w http.ResponseWriter, r *http.Request){

    speciesList := make([]Species, len(h.store))

    h.Lock()
    i := 0
    for _, species := range h.store {
        speciesList[i] = species
        i++
    }
    h.Unlock()

    jsonBytes, err := json.Marshal(speciesList)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}

func (h *speciesHandlers) post(w http.ResponseWriter, r *http.Request){
    bodyBytes, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil{
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
    if err != nil{
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    species.ID = fmt.Sprintf("%d", time.Now().UnixNano())

    h.Lock()
    h.store[species.ID] = species
    defer h.Unlock()
}

func (h *speciesHandlers) species(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "GET":
        h.get(w,r)
        return
    case "POST":
        h.post(w,r)
        return
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("method not allowed"))
        return
    }
}

func newSpeciesHandlers() *speciesHandlers{

    return &speciesHandlers{
        store: map[string]Species{
            "id1": Species{
                ID: "id1",
                GenericName: "Myrmica",
                SpecificName: "rubra",
                WorkersLength: 6,
                QueenLength: 7.5,
            },
        },
    }
}

func main(){
    speciesHandlers := newSpeciesHandlers()
    http.HandleFunc("/species", speciesHandlers.species)
    err := http.ListenAndServe(":8080", nil)
    if err != nil{
        panic(err)
    }
}
