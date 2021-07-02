package main

import (
	"fmt"
	"log"

	"github.com/move-ton/ton-client-go/domain"

	goton "github.com/move-ton/ton-client-go"
)

func TestTon() string {
	ton, err := goton.NewTon(domain.BaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer ton.Client.Destroy()

	value, err := ton.Client.Version()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Version bindings is: ", value.Version)
	return value.Version
}
