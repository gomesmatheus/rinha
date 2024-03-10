package service

import (
	"errors"
	"fmt"
	"time"

	dao "github.com/gomesmatheus/rinha/src/database"
	"github.com/gomesmatheus/rinha/src/models"
)

func ExecuteTransaction(transaction models.Transaction, customerId int) (models.TransactionResponse, error) {
    tResponse, err := dao.GetCustomersBalanceAndLimit(customerId)
    if err != nil {
        return tResponse, err
    }

    if transaction.Type == "c" {
        tResponse.Balance+=transaction.Value
    } else if transaction.Type == "d" {
        if transaction.Value > (tResponse.Balance + tResponse.Limit) {
            fmt.Println("Invalid transaction")
            return tResponse, errors.New("Invalid transaction")
        }
        tResponse.Balance-=transaction.Value
    }

    err = dao.UpdateCustomersBalance(customerId, tResponse.Balance)
    if err != nil {
        fmt.Println("Error updating customers balance credit operation", err)
        return tResponse, err
    }

    err = dao.RegisterTransaction(customerId, transaction)
    if err != nil {
        fmt.Println("Error registering transaction", err)
        return tResponse, err
    }
    return tResponse, nil
}

func GetClientsBankStatement(customerId int) (statement models.Statement, err error) {
    var transactions []models.Transaction

    // essa chamada poderia ser async
    transactions, err = dao.GetCustomersLastTransactions(customerId)
    if err != nil {
        fmt.Println("Error retrieving clients last 10 transactions", err)
        return
    }

    // essa aqui também
    tResponse, err := dao.GetCustomersBalanceAndLimit(customerId)
    if err != nil {
        fmt.Println("Error retrieving clients balance and limit", err)
        return
    }

    // aqui um wg de espera das duas chamadas async
    statement.Balance.Total = tResponse.Balance
    statement.Balance.Limit = tResponse.Limit
    statement.LastTransactions = transactions
    statement.Balance.StatementDate = time.Now()
    return
}

