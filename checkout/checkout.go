package main

import (
	"checkout/queue"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

type Order struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	ProductId string `json:"product_id"`
}

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func displayCheckout(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	response, err := http.Get(productsUrl + "/product/" + vars["id"])
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))

	var product Product
	json.Unmarshal(data, &product)

	t := template.Must(template.ParseFiles("templates/checkout.html"))
	t.Execute(w, product)
}

func finish(w http.ResponseWriter, r *http.Request)  {
	var order Order
	order.Name = r.FormValue("name")
	order.Email = r.FormValue("email")
	order.Phone = r.FormValue("phone")
	order.ProductId = r.FormValue("product_id")

	data, _ := json.Marshal(order)
	fmt.Println(string(data))

	connection := queue.Connect()
	queue.Notify(data, "checkout_ex", "", connection)

	w.Write([]byte("Processou!"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/finish", finish)
	r.HandleFunc("/{id}", displayCheckout)
	http.ListenAndServe(":8082", r)
}
