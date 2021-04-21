package main

import "net/http"
import "encoding/json"

type Species struct{
    ID string `json:"id"`
    GenericName string `json:"genericname"`
    SpecificName string `json:"specificname"`
    WorkersLength float64 `json:"workerslength"`
    QueenLength float64 `json:"queenlength"`
}

type speciesHandlers struct{

    store map[string]Species

}
func (h *speciesHandlers) get(w http.ResponseWriter, r *http.Request){
    
    speciesList := make([]Species, len(h.store))
    
    i := 0
    for _, species := range h.store {
        speciesList[i] = species
        i++
    }

    jsonBytes, err := json.Marshal(speciesList)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }
    
    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
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
    http.HandleFunc("/species", speciesHandlers.get)
    err := http.ListenAndServe(":8080", nil)
    if err != nil{
        panic(err)
    }
}
