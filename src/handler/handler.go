package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gomesmatheus/rinha/src/models"
	"github.com/gomesmatheus/rinha/src/service"
)


func HandleRoute(w http.ResponseWriter, r *http.Request) {
    pathSegments := strings.Split(r.URL.Path, "/")[1:]
    method := r.Method

    customerId, err := strconv.Atoi(pathSegments[1])
    if err != nil {
        fmt.Println("Error parsing id path param")
        w.WriteHeader(400)
        w.Write([]byte("400 bad request"))
        return
    }

    if method == "POST" && pathSegments[2] == "transacoes" {
        transactionRoute(w, r, customerId)
        return
    }

    if method == "GET" && pathSegments[2] == "extrato" {
        statementRoute(w, r, customerId)
        return
    }

    w.WriteHeader(404)
    w.Write([]byte("Essa não é uma rota da rinha"))
}

func transactionRoute(w http.ResponseWriter, r *http.Request, customerId int) {
    body, err := io.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        fmt.Println("Error reading body", err)
        w.WriteHeader(400)
        w.Write([]byte("400 bad request"))
        return
    }

    var transaction models.Transaction
    json.Unmarshal(body, &transaction)
    if isPayloadInvalid(transaction) {
        fmt.Println("Invalid payload")
        w.WriteHeader(400)
        w.Write([]byte("400 bad request"))
        return
    }

    tResponse, err := service.ExecuteTransaction(transaction, customerId)
    if err != nil {
        w.WriteHeader(422)
        return
    }

    response, _ := json.Marshal(tResponse)
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
    return
}

func statementRoute(w http.ResponseWriter, r *http.Request, customerId int) {
    bankStatement, err := service.GetClientsBankStatement(customerId)
    if err != nil {
        fmt.Println("Internal server error", err)
        w.WriteHeader(500)
        w.Write([]byte("500 deu ruim"))
        return
    }

    response, _ := json.Marshal(bankStatement)
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
    return
}

func isPayloadInvalid(t models.Transaction) bool {
    return (t.Type != "c" && t.Type != "d") || t.Description == "" || len(t.Description) > 10 || t.Value == 0
}
