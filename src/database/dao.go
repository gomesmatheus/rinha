package dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gomesmatheus/rinha/src/models"
)


var db *sql.DB

const (
    createTransactions = `
        CREATE TABLE IF NOT EXISTS transactions (
            customer_id INTEGER NOT NULL,
            value INTEGER NOT NULL,
            type VARCHAR(1) NOT NULL,
            description VARCHAR(10) NOT NULL,
            time DATETIME NOT NULL
        );
    `
    createCustomers = `
        CREATE TABLE IF NOT EXISTS customers (
            id INTEGER NOT NULL PRIMARY KEY,
            'limit' INTEGER NOT NULL,
            balance INTEGER NOT NULL
        );
    `
)

func InitDb() (err error) {
    db, err = sql.Open("sqlite3", "database.db")
    if err != nil {
        fmt.Println("Error connecting to db", err)
    }

    if _, err := db.Exec(createTransactions); err != nil {
        fmt.Println("Error creating table transactions", err)
    }

    if _, err := db.Exec(createCustomers); err != nil {
        fmt.Println("Error creating table customers", err)
    }
    return 
}

func CloseDb() {
    db.Close()
}

func GetCustomersLastTransactions(customerId int) (transactions []models.Transaction, err error) {
    rows, err := db.Query("SELECT value, type, description, time FROM transactions WHERE customer_id = ? ORDER BY time DESC LIMIT 10", customerId)
    defer rows.Close()
    if err != nil {
        fmt.Println("Error retrieving data for customerId", customerId)
        fmt.Println(err)
        return
    }
    
    for rows.Next() {
        var t models.Transaction
        if err = rows.Scan(&t.Value, &t.Type, &t.Description, &t.Date); err != nil {
            fmt.Println("Error scanning")
            fmt.Println(err)
            return
        }
        transactions = append(transactions, t)
    }
    return
}

func GetCustomersBalanceAndLimit(customerId int) (tResponse models.TransactionResponse, err error) {
    row := db.QueryRow("SELECT balance, `limit` FROM customers WHERE id = ?", customerId) 
    if err = row.Scan(&tResponse.Balance, &tResponse.Limit); err != nil {
        fmt.Println("Error scanning balance and limit", err)
    }
    return
}

func UpdateCustomersBalance(customerId int, balance int) error {
    _, err := db.Exec("UPDATE customers SET balance = ? WHERE id = ?", balance, customerId)
    return err
}

func RegisterTransaction(customerId int, t models.Transaction) error {
    _, err := db.Exec("INSERT INTO transactions (customer_id, value, type, description, time) VALUES (?, ?, ?, ?, ?)", customerId, t.Value, t.Type, t.Description, time.Now())
    return err
}

