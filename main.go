package main

// consuming star wars api planets to mongo db, create rest api with those datas

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dao = Data{}

var client *mongo.Client

type Data struct { //used to save in db
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Climate string             `json:"climate,omitempty" bson:"climate,omitempty"`
	Terrain string             `json:"terrain,omitempty" bson:"terrain,omitempty"`
	Films   int                `json:"films,omitempty" bson:"films,omitempty"`
}
type Planet struct { //my planet
	Name    string
	Climate string
	Terrain string
	Films   []string
}

// ------------------------------------------------ test start
type PlanetsAPI struct { //returned by swapi
	Count    int
	Next     string
	previous string
	Results  []Planet
}

func getPlanetsFromAPI() {
	response, _ := http.Get("https://swapi.co/api/planets")
	var planets PlanetsAPI
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &planets)
	fmt.Println(planets.Results)
}

// ------------------------------------------------ tests end

func getPlanetFromAPIbyID(pNum string) {
	// returns planet by id
	response, _ := http.Get("https://swapi.co/api/planets/" + pNum)
	var planet Planet
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &planet)
}

func getPlanetFromAPIbyName(name string) PlanetsAPI {
	// returns the planet by name
	response, _ := http.Get("https://swapi.co/api/planets?search=" + name)
	var planet PlanetsAPI
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &planet)
	return planet
}

func POSTplanet(response http.ResponseWriter, request *http.Request) {
	// on POST method, this will be called.
	// it gets the name of the planet that will be saved at mongo
	// search it on the swapi and check if exists
	// than save to collection
	// also create db if it doesn't exists (mongodb f)
	response.Header().Add("content-type", "application/json")
	var planet Planet
	json.NewDecoder(request.Body).Decode(&planet)
	if planet.Name != "" {
		planetAPI := getPlanetFromAPIbyName(planet.Name)
		if planetAPI.Count != 1 {
			fmt.Println("Specify the planet name in a better way.")
		} else {
			var data Data
			data.Name = planetAPI.Results[0].Name
			data.Climate = planetAPI.Results[0].Climate
			data.Terrain = planetAPI.Results[0].Terrain
			if len(planetAPI.Results[0].Films) == 0 {
				data.Films = -1
			} else {
				data.Films = len(planetAPI.Results[0].Films)
			}
			collection := client.Database("starwars").Collection("planets")
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			result, _ := collection.InsertOne(ctx, data)
			json.NewEncoder(response).Encode(result)
		}
	} // else get by id
}

func GETallPlanetsFromMongoDB(response http.ResponseWriter, request *http.Request) {
	// return all database into a json.
	response.Header().Add("content-type", "application/json")
	var planets []Data
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.Find(ctx, bson.M{})
	for result.Next(ctx) {
		var planet Data
		result.Decode(&planet)
		planets = append(planets, planet)
	}
	json.NewEncoder(response).Encode(planets)

	if planets == nil {
		fmt.Println("Empty Database")
	}

}

func DELETEplanetFromMongoDB(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection.DeleteOne(ctx, filter)
}

func GETPlanetByIdFromMongoDB(response http.ResponseWriter, request *http.Request) {

}

func main() {
	fmt.Println("Starting...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/planet", POSTplanet).Methods("POST")
	router.HandleFunc("/planet/{id}", DELETEplanetFromMongoDB).Methods("DELETE")
	router.HandleFunc("/planet/{id}", GETPlanetByIdFromMongoDB).Methods("GET")
	router.HandleFunc("/planet", GETallPlanetsFromMongoDB).Methods("GET")
	fmt.Println("Started")
	http.ListenAndServe(":12345", router)
}
