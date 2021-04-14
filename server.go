package main

import "net/http"

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
    
}

func newSpeciesHandlers() *speciesHandlers{

    return &speciesHandlers{
        store: map[string]Species{
            
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
