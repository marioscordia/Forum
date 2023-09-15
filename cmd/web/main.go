package main

import (
	"log"
	"net/http"
	"newforum/config"
	"newforum/internal/handler"
	"newforum/internal/service"
	"newforum/internal/store"
	"newforum/internal/temp"
	"os"

	_ "github.com/mattn/go-sqlite3"
)


func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	config, err := config.NewConfig()
	if err != nil {
		errLog.Fatal(err)
	}


	tmpCache, err := temp.NewTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	db, err := store.InitializeDB(config)
	if err != nil{
		errLog.Fatal(err)
	}
	defer db.Close()

	store := store.NewStore(db)
	service := service.NewService(store)
	handler := handler.NewHandler(infoLog, errLog, tmpCache, service)

	infoLog.Printf("Starting the server on port%s\n", config.Port)

	srv := &http.Server{
		Addr: config.Port,
		ErrorLog: errLog,
		Handler: handler.Routes(),
	}
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
