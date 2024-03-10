package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/gomesmatheus/rinha/src/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

const (
    createTransactions = `
        CREATE TABLE IF NOT EXISTS transactions (
            customer_id INTEGER NOT NULL,
            value INTEGER NOT NULL,
            type VARCHAR(1) NOT NULL,
            description VARCHAR(10) NOT NULL,
            time TIMESTAMP NOT NULL,

            CONSTRAINT fk_customer FOREIGN KEY(customer_id) REFERENCES customers(id)
        );
    `
    createCustomers = `
        CREATE TABLE IF NOT EXISTS customers (
            id INTEGER NOT NULL PRIMARY KEY,
            "limit" INTEGER NOT NULL,
            balance INTEGER NOT NULL
        );
    `

    createIndex = "CREATE INDEX IF NOT EXISTS customer_id_index ON transactions(customer_id);"
)

func InitDb() (err error) {
    url := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
    config, err := pgxpool.ParseConfig(url)
    if err != nil {
        fmt.Println("Error parsing config", err)
    }
    db, err = pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        fmt.Println("Error creating connection", err)
    }
    
    if _, err := db.Exec(context.Background(), createCustomers); err != nil {
        fmt.Println("Error creating table customers", err)
    }

    if _, err := db.Exec(context.Background(), createTransactions); err != nil {
        fmt.Println("Error creating table transactions", err)
    }
 
    if _, err := db.Exec(context.Background(), createIndex); err != nil {
        fmt.Println("Error creating index on table transactions", err)
    }

    return 
}

func CloseDb() {
    db.Close()
}

func GetCustomersLastTransactions(customerId int) (transactions []models.Transaction, err error) {
    rows, err := db.Query(context.Background(), "SELECT value, type, description, time FROM transactions WHERE customer_id = $1 ORDER BY time DESC LIMIT 10", customerId)
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
    row := db.QueryRow(context.Background(), `SELECT balance, "limit" FROM customers WHERE id = $1`, customerId)
    if err = row.Scan(&tResponse.Balance, &tResponse.Limit); err != nil {
        fmt.Println("Error scanning balance and limit", err)
    }
    return
}

func UpdateCustomersBalance(customerId int, balance int) error {
    _, err := db.Exec(context.Background(), "UPDATE customers SET balance = $1 WHERE id = $2", balance, customerId)
    return err
}

func RegisterTransaction(customerId int, t models.Transaction) error {
    _, err := db.Exec(context.Background(), "INSERT INTO transactions (customer_id, value, type, description, time) VALUES ($1, $2, $3, $4, $5)", customerId, t.Value, t.Type, t.Description, time.Now())
    return err
}

