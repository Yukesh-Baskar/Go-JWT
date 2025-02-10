package main

import (
	"fmt"
	"go-jwt/app"
)

func main() {
	if err := app.StartApplication(); err != nil {
		fmt.Println(err)
	}
}
