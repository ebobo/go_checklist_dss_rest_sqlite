package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ebobo/dss_checklist/pkg/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *Server) startHTTP() error {
	m := mux.NewRouter()

	// Add CORS
	cors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		MaxAge:           31,
		Debug:            false,
	})

	// This is where you add other stuff you want to map in the mux

	// Config endpoint
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Welcome to the CHECKLIST-SERVER !")
	}).Methods("GET")

	// Get All Items
	m.HandleFunc("/api/v1/items", s.GetListItems).Methods("GET")

	// Add new Item
	m.HandleFunc("/api/v1/item", s.AddItem).Methods("POST")

	// Get Item by ID
	m.HandleFunc("/api/v1/item/{id}", s.GetItem).Methods("GET")

	// Update Item by ID can use PATCH or PUT
	m.HandleFunc("/api/v1/item/{id}", s.UpdateItem).Methods("PUT")

	// Delete Item by ID
	m.HandleFunc("/api/v1/item/{id}", s.DeleteItem).Methods("DELETE")

	httpServer := &http.Server{
		Addr:              s.httpListenAddr,
		Handler:           handlers.ProxyHeaders(cors.Handler(m)),
		ReadTimeout:       (10 * time.Second),
		ReadHeaderTimeout: (8 * time.Second),
		WriteTimeout:      (45 * time.Second),
	}

	// Set up shutdown handler
	go func() {
		<-s.ctx.Done()
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Printf("error shutting down HTTP interface '%s': %v", s.httpListenAddr, err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("starting HTTP interface '%s'", s.httpListenAddr)

		// This isn't entirely true and really represents a race condition, but
		// doing this properly is a pain in the neck.
		s.httpStarted.Done()

		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			err = errors.New("")
		}

		log.Printf("HTTP interface '%s' down %v", s.httpListenAddr, err)
		s.httpStopped.Done()
	}()

	return nil
}

func (s *Server) GetListItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := s.db.ListItems()
	if err != nil {
		log.Printf("failed to get items %v", err)
		http.Error(w, "failed to get items", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (s *Server) GetItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	itemID := mux.Vars(r)["id"]
	log.Printf("itemID: %s", itemID)
	item, err := s.db.GetItem(itemID)

	if err != nil {
		log.Printf("failed to get item %v by id %s", err, itemID)
		http.Error(w, "failed to get item", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (s *Server) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var newItem model.Item

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the item detail in order to update")
	}
	json.Unmarshal(reqBody, &newItem)

	err = s.db.AddItem(newItem)

	if err != nil {
		log.Printf("failed to add item %v", err)
		http.Error(w, "failed to add item", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func (s *Server) UpdateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	itemID := mux.Vars(r)["id"]
	item, err := s.db.GetItem(itemID)

	if err != nil {
		log.Printf("failed to get item %v by id %s", err, itemID)
		http.Error(w, "failed to update item", http.StatusBadRequest)
		return
	}

	log.Printf("update item 2")

	var updateItem model.Item

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the item detail in order to update")
	}
	json.Unmarshal(reqBody, &updateItem)

	// update item
	if updateItem.Name != "" {
		item.Name = updateItem.Name
	}
	if updateItem.Position != 0 {
		item.Position = updateItem.Position
	}
	item.Status = updateItem.Status

	if updateItem.Tag != "" {
		item.Tag = updateItem.Tag
	}

	log.Println("update item", item.Name, item.Position, item.Status, item.Tag)

	err = s.db.UpdateItem(item)

	if err != nil {
		log.Printf("failed to get item %v by id %s", err, itemID)
		http.Error(w, "failed to update item", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (s *Server) DeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	itemID := mux.Vars(r)["id"]
	err := s.db.DeleteItem(itemID)

	if err != nil {
		log.Printf("failed to delete item %v by id %s", err, itemID)
		http.Error(w, "failed to delete item", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
