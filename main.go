package main

import (
	"github.com/Leandroreign/crud/storage"
)

func main() {
	driver := storage.Postgres
	storage.New(driver)

}
