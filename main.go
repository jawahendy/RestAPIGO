package main

import (
    "encoding/json"
    "fmt"
    "log"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"
)

type Customer struct {
    Id      string    `json:"Id"`
    Name string `json:"Name"`
    Alamat  string `json:"Alamat"`
    Gender  string `json:"Gender"`
}

type Response struct {
    Status  string    `json:"Status"`
    Message string    `json:"Message"`
}

var Customers []Customer
var response []Response

func Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Home!")
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(Customers)
}

func GetCustomerSpecific(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    for _, Customer := range Customers {
        if Customer.Id == key {
            json.NewEncoder(w).Encode(Customer)
        }
    }
}


func NewCustomer(w http.ResponseWriter, r *http.Request) {   
    reqBody, _ := ioutil.ReadAll(r.Body)
    var customer Customer 
    json.Unmarshal(reqBody, &customer)
    Customers = append(Customers, customer)

    json.NewEncoder(w).Encode(customer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {   
    vars := mux.Vars(r)
    id := vars["id"]
    for index, customer := range Customers {
        if customer.Id == id {
            Customers = append(Customers[:index], Customers[index+1:]...)
            var NewCustomer Customer
            json.NewDecoder(r.Body).Decode(&NewCustomer)
            NewCustomer.Id = id
            Customers = append(Customers, NewCustomer)
            json.NewEncoder(w).Encode(NewCustomer)
            return
        }
    }
}

func RemoveCustomer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var response Response
    id := vars["id"]

    for index, customer := range Customers {
        if customer.Id == id {
            Customers = append(Customers[:index], Customers[index+1:]...)
        }
    }

    response.Status = "200"
    response.Message = "Success Delete"
    json.NewEncoder(w).Encode(response)

}

func handleRequests() {
    CusRouter := mux.NewRouter().StrictSlash(true)
    CusRouter.HandleFunc("/", Home)
    CusRouter.HandleFunc("/customers", GetCustomer)
    CusRouter.HandleFunc("/customer", NewCustomer).Methods("POST")
    CusRouter.HandleFunc("/customer/{id}", RemoveCustomer).Methods("DELETE")
    CusRouter.HandleFunc("/Ecustomer/{id}", UpdateCustomer).Methods("PUT")
    CusRouter.HandleFunc("/customers/{id}", GetCustomerSpecific)
    log.Fatal(http.ListenAndServe(":9000", CusRouter))
}

func main() {
    Customers = []Customer{
        Customer{Id: "1", Name: "Hendy", Alamat: "Jl.SituAja", Gender: "male"},
        Customer{Id: "2", Name: "Siska", Alamat: "Jl.SanaAja", Gender: "Female"},
    }
    handleRequests()
}