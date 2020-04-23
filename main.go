package main

import "os"
import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

// Sensor is the output structure for ESP8266
type Sensor struct {
    Temperature string `json:"temperature"`
    Humidity string `json:"humidity"`
}

func get() {
    url := "http://192.168.1.29/data"

    res, err := http.Get(url)

    if err != nil {
        panic(err.Error())
    }

    body, err := ioutil.ReadAll(res.Body)

    if err != nil {
        panic(err.Error())
    }

	var data Sensor
	json.Unmarshal(body, &data)
	fmt.Println(data.Temperature)
	fmt.Println(data.Humidity)
    os.Exit(0)
}

func main() {
    get()
}
