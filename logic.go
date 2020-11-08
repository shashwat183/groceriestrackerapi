package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func createGrocery(w http.ResponseWriter, r *http.Request) {
	var newGrocery grocery
	var retrievedGrocery grocery
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please provide grocery name and price")
	}
	json.Unmarshal(reqBody, &newGrocery)
	fmt.Println(newGrocery)
	result := db.Where("name = ?", newGrocery.Name).First(&retrievedGrocery)
	if result.Error != nil {
		db.Create(&newGrocery)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newGrocery)
		log.Printf("Created New grocery %v", newGrocery.Name)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(createError("Grocery record already present"))
		log.Printf("Grocery record %v already present", newGrocery.Name)
	}
}

func getGrocery(w http.ResponseWriter, r *http.Request) {
	groceryName := mux.Vars(r)["name"]

	var retrievedGrocery grocery
	result := db.Where("name = ?", groceryName).First(&retrievedGrocery)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(createError(fmt.Sprintf("Could not find grocery %v", groceryName)))
		log.Printf("Could not find grocery %v", groceryName)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(retrievedGrocery)
		log.Printf("Grocery %v retrieved", groceryName)
	}

}

func getGroceries(w http.ResponseWriter, r *http.Request) {
	var retrievedGroceries groceries
	result := db.Find(&retrievedGroceries)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(createError("Database operation encountered error"))
		// fmt.Fprint(w, "Database operation encountered error")
		log.Printf("Database operation encountered error")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(retrievedGroceries)
	log.Printf("Returned all groceries")
}

func updateGrocery(w http.ResponseWriter, r *http.Request) {
	groceryName := mux.Vars(r)["name"]
	var updatedGrocery grocery

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(createError("Please provide grocery name and price"))
		// fmt.Fprintf(w, "Please provide grocery name and price")
		log.Printf("User update input was incorrect")
	} else {
		json.Unmarshal(reqBody, &updatedGrocery)
		var retrievedGrocery grocery
		result := db.Where("name = ?", groceryName).First(&retrievedGrocery)
		if result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(createError(fmt.Sprintf("Could not find grocery %v", groceryName)))
			log.Printf("Could not find grocery %v", groceryName)
		} else {
			retrievedGrocery.Name = updatedGrocery.Name
			retrievedGrocery.Price = updatedGrocery.Price
			db.Unscoped().Where("name = ?", groceryName).Delete(grocery{})
			db.Save(&retrievedGrocery)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(retrievedGrocery)
		}
	}
}

func deleteGrocery(w http.ResponseWriter, r *http.Request) {
	groceryName := mux.Vars(r)["name"]
	var retrievedGrocery grocery

	result := db.Where("name = ?", groceryName).First(&retrievedGrocery)
	if result.Error != nil {
		//error clause
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(createError(fmt.Sprintf("Could not find grocery %v", groceryName)))
		log.Printf("Could not find grocery %v", groceryName)
	} else {
		// success clause
		db.Unscoped().Where("name = ?", groceryName).Delete(grocery{})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(createError(fmt.Sprintf("grocery %v deleted successfully", groceryName)))
		log.Printf("grocery %v deleted successfully", groceryName)
	}
}

func createError(message string) errorResponse {
	var newErrorResponse = errorResponse{
		Time:    time.Now(),
		Message: message,
	}

	return newErrorResponse
}

func createResponseMessage(message string) successMessage {
	var newSuccessMessage = successMessage{
		Message: message,
	}

	return newSuccessMessage
}
