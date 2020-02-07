package main

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"net/http"
	"order/db"
	"order/queue"
	"os"
	"time"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"phone"`
	CreatedAt time.Time `json:"created_at,string"`
}

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.StartConsuming(connection, in)

	for payload := range in {
		createOrder(payload)
		fmt.Println(string(payload))
	}
}

func createOrder(payload []byte) {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()
	order.Uuid = uuid.String()
	order.Status = "pendente"
	order.CreatedAt = time.Now()
	saveOrder(order)
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()

	err := connection.Set(order.Uuid, string(json),0).Err()
	if err != nil {
		panic(err.Error())
	}

}

func getProductById(id string) Product {
	response, err := http.Get(productsUrl + "/product/" + id)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	json.Unmarshal(data, &product)
	return product
}
