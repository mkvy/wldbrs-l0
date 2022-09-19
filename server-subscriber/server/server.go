package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
	"github.com/mkvy/wldbrs-l0/server-subscriber/store"
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	Srv     *http.Server
	storage store.StoreService
	Addr    string
}

func InitServer(store store.StoreService, addr string) *Server {
	server := Server{
		storage: store,
		Addr:    addr,
	}
	return &server
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/orders/{o_id}", s.ordersHandler)
	s.Srv = &http.Server{Addr: s.Addr, Handler: router}
	log.Println("Server is starting")
	err := s.Srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	log.Println("Server stops")
	return s.Srv.Close()
}

func (s *Server) ordersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["o_id"]
	od := s.storage.GetFromCacheByUID(id)
	if od.OrderUid == "" {
		w.WriteHeader(http.StatusBadRequest)
		parsedTemplate, _ := template.ParseFiles("./server/templates/notFound.html")
		err := parsedTemplate.Execute(w, struct{ Id string }{Id: id})
		if err != nil {
			log.Printf("Error occurred while executing the template : ", id)
			return
		}
		return
	}

	dataItem := model.DataItem{
		ID:        id,
		OrderData: od,
	}
	parsedTemplate, _ := template.ParseFiles("./server/templates/index.html")
	fmt.Println("DATA ITEM", dataItem)
	err := parsedTemplate.Execute(w, dataItem)
	if err != nil {
		log.Printf("Error occurred while executing the template : ", dataItem)
		return
	}
}
