package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var db map[int]TodoObject

// m := make(db)

type TodoObject struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	t := TodoObject{}
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	}
	t.Id = len(db) + 1
	db[t.Id] = t
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(dbMapToSlice())
}

func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(dbMapToSlice())
}

func MarkDone(w http.ResponseWriter, r *http.Request) {
	var resp string
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("error while conversion ::", err)
	}
	if t, ok := db[id]; ok {
		t.Status = "Done"
		db[id] = t

		resp = "Updated successfully!"
	} else {
		resp = "No TODO found"
	}
	json.NewEncoder(w).Encode(resp)

}

func dbMapToSlice() []TodoObject {
	resp := []TodoObject{}
	for _, v := range db {
		resp = append(resp, v)

	}
	return resp
}

func main() {
	fmt.Println("Hello Gopher!")
	db = make(map[int]TodoObject)

	r := mux.NewRouter()
	r.HandleFunc("/todo/new", AddTodo)
	r.HandleFunc("/todo/all", GetAllTodo)
	r.HandleFunc("/todo/id/{id}/done", MarkDone)
	r.Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

	r.Use(mux.CORSMethodMiddleware(r))

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
