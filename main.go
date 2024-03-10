package main

import (
	"log"
	"net/http"
	dao "github.com/gomesmatheus/rinha/src/database"
	"github.com/gomesmatheus/rinha/src/handler"
)

// docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres
func main() {
    err := dao.InitDb()
    if err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
    defer dao.CloseDb()

    http.HandleFunc("/clientes/", handler.HandleRoute)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

