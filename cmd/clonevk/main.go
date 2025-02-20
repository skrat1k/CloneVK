package main

import (
	"CloneVK/internal/storage"
	"context"
	"fmt"
	"log"
)

// Потом сделать подгрузку из файла окружения
const (
	UsernameDB = "postgres"
	PasswordDB = "admin"
	HostDB     = "localhost"
	PortDB     = "5432"
	NameDB     = "clonevk"
)

func main() {
	conn, err := storage.CreatePostgresConnection(storage.ConnectionInfo{
		Username: UsernameDB,
		Password: PasswordDB,
		Host:     HostDB,
		Port:     PortDB,
		DBName:   NameDB,
	})
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close(context.Background())

	var test string
	err = conn.QueryRow(context.Background(), "SELECT testvalue FROM Test WHERE testid=4").Scan(&test)
	if err != nil {
		fmt.Println("Ehe pizda", err)
	}

	fmt.Println(test)
}
