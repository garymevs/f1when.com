package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Function for handling requests to /
func webRoot(w http.ResponseWriter, req *http.Request) {
    fmt.Println("Request from:", req.RemoteAddr)
    // If we haven't pulled the data from ergast in 24 hours then pull it again (Might increse this in the future)
    if time.Now().After(dataPullTime.Add(time.Hour*24)) {
        fmt.Println(time.Now(), "is more than 24 hours after", dataPullTime, "reloading data")
        loadSeasonData()
    }
    // Get the next race (duh!)
    race, err := getNextRace(&raceTable)
    if err != nil {
        fmt.Fprintf(w, err.Error())
    }
    // Render the page template using the race retrieved in the previous function
    pageTemplate.Execute(w, &Page{race})
}

// Reload the page template then return the main page 
func templateReload(w http.ResponseWriter, req *http.Request) {
    loadPageTemplate()
    webRoot(w, req)
}

// True to enable dev features
var devMode = false

// Store the data in a publicly accessible area
var baseData BaseData
var raceTable RaceTable
var pageTemplate *template.Template

// Store when we last pulled the data
var dataPullTime time.Time

func main() {
    // Load the page template
    loadPageTemplate()

    // Download the json data from ergast.com
    loadSeasonData()

    // Setup data objects 
    f1Data := baseData.MRData
    raceTable = f1Data.RaceTable

    // Serve static resources
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    // Serve the main page 
    http.HandleFunc("/", webRoot)
    // Dev mode pages
    if devMode {
        http.HandleFunc("/reload-template", templateReload)
    }
    // Start the web server
    if devMode {
        http.ListenAndServe(":8081", nil)
        fmt.Println("Dev mode server started on :8081")
    } else {
        http.ListenAndServe(":8080", nil)
        fmt.Println("Server started on :8080")
    }
}

func loadPageTemplate() {
    fmt.Println("Loading page template...")
    // Load template file
    template, err := template.ParseFiles("template/page-template.html")
    if err != nil {
        fmt.Println(err.Error())
    }
    pageTemplate = template
    fmt.Println("Page template load complete")
}

// Download the JSON from ergast.com
func loadSeasonData() {
    fmt.Println("Requesting data from ergast...")
    // Data from here: https://ergast.com/api/f1/current.json
    // Create http client and make GET request
    httpClient := &http.Client{}
    res, err := httpClient.Get("https://ergast.com/api/f1/current.json")
    if err != nil {
        fmt.Println("Error pulling data from ergast")
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer res.Body.Close()

    // Read the data and unmarshal it to the structs in f1data.go
    byteValue, _ := ioutil.ReadAll(res.Body)
    json.Unmarshal(byteValue, &baseData)

    // Update dataPullTime
    dataPullTime = time.Now()
    fmt.Println("Data pull complete at:", dataPullTime)
}

// Iterate over the Races in the RaceTable and find the first race with a date after today
func getNextRace(raceTable *RaceTable) (*Race, error) {
    for _, race := range raceTable.Races {
        parsedDate, err := time.Parse("2006-01-02 15:04:05Z", race.Date + " " + race.Time)
        if err != nil {
            fmt.Println(err.Error())
            break
        }

        if time.Now().Before(parsedDate.Add(time.Hour * 2)) {
            return &race, nil
        }
    }
    return nil, errors.New("Couldn't find next race")
}
