package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"proxy/config"
	"proxy/proxy"
	"proxy/requests"
	"proxy/saver"
)

func main() {

	conf := config.NewConfig()

	mongo := &saver.Saver{}

	mongo.MongoConnect()

	p := proxy.Proxy{
		Saver: mongo,
	}

	handlers := &requests.Handlers{
		Saver: mongo,
		Proxy: p,
	}

	handlers.Proxy = p

	server := http.Server{
		Handler:      &p,
		Addr:         ":8080",
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	router := mux.NewRouter()
	router.HandleFunc("/requests", handlers.GetRequests).Methods(http.MethodGet)
	router.HandleFunc("/requests/{id}", handlers.GetRequestByID).Methods(http.MethodGet)
	router.HandleFunc("/repeat/{id}", handlers.RepeatRequest).Methods(http.MethodGet)
	router.HandleFunc("/dirsearch/{id}", handlers.DirSearch).Methods(http.MethodGet)

	repeatServer := http.Server{
		Handler:      router,
		Addr:         ":8081",
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if conf.HTTPS {
		fmt.Println("Start serving TLS")
		go repeatServer.ListenAndServe()
		if err := server.ListenAndServeTLS("server.pem", "server.key"); err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println("Start serving HTTP")
		go repeatServer.ListenAndServe()
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}

}
