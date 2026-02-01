// This file is used by Atlas to load GORM models

package main

import (
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.Product{},
		&models.ProductImages{},
		&models.Cart{},
		&models.CartItems{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		panic(err)
	}
	io.WriteString(os.Stdout, stmts)
}
