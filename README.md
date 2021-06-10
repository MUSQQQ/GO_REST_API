# GO_REST_API
Simple REST API in Go. Made just to get to know basic API structure and http go package

# Commands
```
GET /species    returns a list of species as JSON
GET /species/{id}   returns details of a specific species as JSON
POST /species   accepts a new species to be added
POST /species   returns status 415 if content is not application/json
GET /admin  require auth
GET /species/random redirects (Status 302) to a random species
```

# JSON
```
data:
{
    "id": "someid",
    "genericname": "genusname",
    "specificname": "speciesname",
    "workerslength": "averageworkerslength",
    "queenlength": "averagequeenlenght",
}
```
