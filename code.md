## Migrar tabla de productos
```go
storageProduct := storage.NewPsqlProduct(storage.Pool())
serviceProduct := product.NewService(storageProduct)

if err := serviceProduct.Migrate(); err != nil {
    log.Fatalf("product.Migrate: %v", err)
}
```


## Migrar tabla de facturas
```go
storageInvoice := storage.NewPsqlInvoice(storage.Pool())
serviceInvoice := invoice.NewService(storageInvoice)

if err := serviceInvoice.Migrate(); err != nil {
    log.Fatalf("invoice.Migrate: %v", err)
}
```

## Migrar tabla de las lineas de las facturas
```go
storageInvoiceLines := storage.NewPsqlInvoicesLines(storage.Pool())
serviceInvoiceLines := invoiceLines.NewService(storageInvoiceLines)

if err := serviceInvoiceLines.Migrate(); err != nil {
    log.Fatalf("invoiceLines.Migrate: %v", err)
}
```