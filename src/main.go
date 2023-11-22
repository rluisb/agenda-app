package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/rluisb/agenda-app/src/api"
	"github.com/rluisb/agenda-app/src/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	listenAddrs := flag.String("listen", ":8080", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	contactStore := db.NewMongoContactStore(client)
	contactHandler := api.NewContactHandler(contactStore)

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/api/v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				contactHandler.HandleGetContact(w, r)
				return
			}
			contactHandler.HandleListContacts(w, r)
		case http.MethodPost:
			contactHandler.HandlePostContact(w, r)
		case http.MethodDelete:
			contactHandler.HandleDeleteContact(w, r)
		case http.MethodPut:
			contactHandler.HandleUpdateContact(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	log.Println("Starting server on port 8080")
	http.ListenAndServe(*listenAddrs, mainMux)
}
