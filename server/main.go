package main

import (
	"backend/internal/company"
	"backend/internal/model"
	"backend/internal/property"
	"backend/internal/tenant"
	"backend/internal/unit"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres sslmode=disable host=localhost")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	ms := model.NewService(db)
	cr := company.NewRepository(ms)
	propertyService := property.NewRepository(ms, cr)
	tenantRepo := tenant.NewRepo(ms)
	tenantService := tenant.NewService(tenantRepo)
	unitRepo := unit.NewRepo(ms)
	unitService := unit.NewService(unitRepo, tenantService, propertyService)
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.TxKey{}, &model.Tx{nil})
	err = propertyService.CreateProperty(ctx, "123", "owner", "32132")
	if err != nil {
		fmt.Println("Error creating Property", err)
	}
	err = unitService.CreateUnit(ctx, "PR-1722060321450632", "7", tenant.Tenant{
		FirstName:   "Kartik",
		LastName:    "Kapoor",
		Email:       "kartikkapoor33",
		PhoneNumber: "4343",
	})
	if err != nil {
		fmt.Println(err)
	}
	//mux := http.NewServeMux()
	//mux.Handle("/foo", http.HandlerFunc(api.CreateWorkOrderHandler))
	//err = http.ListenAndServe(":8080", mux)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("Started Listening on port 8080")
}
