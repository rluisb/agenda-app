package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/rluisb/agenda-app/src/api"
	"github.com/rluisb/agenda-app/src/db"
	"github.com/rluisb/agenda-app/src/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	listenAddrs := flag.String("listen", ":8080", "The listen address of the API server")
	flag.Parse()

	args := types.Args{
		Conn: db.DBURI,
	}
	if conn := os.Getenv("DB_CONN"); conn != "" {
		log.Printf("Using DB_CONN env variable: %s", conn)
		args.Conn = conn
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(args.Conn))
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
		case http.MethodPatch:
			contactHandler.HandlePatchContact(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mainMux.HandleFunc("/api/v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mainMux.HandleFunc("/api/v1/readyz", func(w http.ResponseWriter, r *http.Request) {
		err := client.Ping(context.Background(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	log.Println("Starting server on port 8080")
	http.ListenAndServe(*listenAddrs, mainMux)
}
