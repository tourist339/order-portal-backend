package main

import (
	"backend/internal/company"
	"backend/internal/model"
	"backend/internal/property"
	"backend/internal/role"
	"backend/internal/unit"
	"backend/internal/user"
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
	createTables(ms)

	propertyRepo := property.NewRepository(ms, cr)

	propertyService := property.NewService(propertyRepo)
	roleRepo := role.NewRepo(ms)
	userRepo := user.NewRepo(ms)
	userService := user.NewService(userRepo)
	tenantService := role.NewService(roleRepo, userService)
	unitRepo := unit.NewRepo(ms)
	unitService := unit.NewService(unitRepo, tenantService, propertyService)
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.TxKey{}, &model.Tx{nil})
	err = propertyService.CreateProperty(ctx, "123", "owner", "32132")
	if err != nil {
		fmt.Println("Error creating Property", err)
	}
	_, err = unitService.CreateUnit(ctx, "PR-1722218677537813", "7")
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

func createTables(ms *model.Service) {
	tables := map[string]interface{}{
		user.TABLE_NAME:     &user.User{},
		role.TABLE_NAME:     &role.Role{},
		unit.TABLE_NAME:     &unit.Unit{},
		property.TABLE_NAME: &property.Property{},
	}
	for k, v := range tables {
		_ = ms.CreateTable(k, v)
	}

}
