package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Planet struct {
	Name    string
	Climate string
	Terrain string
}

func getPlanet(pNum int) {
	//total of 61 planets
	response, err := http.Get("https://swapi.co/api/planets/" + string(pNum))
	if err != nil {
		panic(err)
	} else {
		var planet Planet
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal([]byte(string(data)), &planet)
		fmt.Println(planet)
	}
}

func main() {

}
