package models

import "time"


type Transaction struct {
    Value int `json:"valor"`
    Type string `json:"tipo"`
    Description string `json:"descricao"`
    Date time.Time `json:"realizada_em"`
}

type TransactionResponse struct {
    Limit int `json:"limit"`
    Balance int `json:"saldo"`
}

type Statement struct {
    Balance struct {
        Total int `json:"total"`
        StatementDate time.Time `json:"data_extrato"`
        Limit int `json:"limite"`
    } `json:"saldo"` 
    LastTransactions [] Transaction `json:"ultimas_transacoes"`
}
