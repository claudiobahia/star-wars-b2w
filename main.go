package main

// consuming star wars api planets to mongo db, create rest api with those datas

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Planet struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Climate string             `json:"climate,omitempty" bson:"climate,omitempty"`
	Terrain string             `json:"terrain,omitempty" bson:"terrain,omitempty"`
	Films   []string           `json:"films,omitempty" bson:"films,omitempty"`
}

func getPlanet(pNum string) {
	//total of 61 planets
	response, _ := http.Get("https://swapi.co/api/planets/" + pNum)
	var planet Planet
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(string(data)), &planet)
	fmt.Println(planet)
}

func CreatePlanetEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var planet Planet
	json.NewDecoder(request.Body).Decode(&planet)
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, planet)
	json.NewEncoder(response).Encode(result)
}

func main() {
	// fmt.Println("Starting")
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// router := mux.NewRouter()
	// router.HandleFunc("/planet", CreatePlanetEndpoint).Methods("POST")
	// http.ListenAndServe(":12345", router)
	x := "2"
	getPlanet(x)
}
