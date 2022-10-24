package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Main HTML page template
var pageTemplate *template.Template

// F1 season data
var raceTable RaceTable
var dataPullTime time.Time

const dataRefreshPeriod time.Duration = time.Hour * 24

func main() {
	// check for DEV env variable
	devMode, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		devMode = false
	}

	// Load the page template
	pageTemplate, err = loadPageTemplate("template/page-template.html")
	if err != nil {
		log.Fatalf("error loading page template: %s", err)
	}

	// Download the json data from ergast.com
	refreshRaceTable()

	// Serve static resources
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Serve the main page
	http.HandleFunc("/json", jsonRoot)
	http.HandleFunc("/", webRoot)
	// Dev mode pages
	if devMode {
		http.HandleFunc("/reload-template", templateReload)
	}
	// Start the web server
	if devMode {
		fmt.Println("Starting dev server on :8081")
		http.ListenAndServe(":8081", nil)
	} else {
		fmt.Println("Starting server on :8080")
		http.ListenAndServe(":8080", nil)
	}
}

func refreshRaceTable() {
	seasonData, err := loadSeasonData()
	if err != nil {
		log.Fatalf("error loading season data: %s", err)
	}
	// Update dataPullTime
	dataPullTime = time.Now()

	// Setup data objects
	raceTable = seasonData.MRData.RaceTable
}

// Download the JSON from ergast.com
func loadSeasonData() (BaseData, error) {
	log.Println("Requesting data from ergast...")
	baseData := BaseData{}
	// Data from here: https://ergast.com/api/f1/current.json
	// Create http client and make GET request
	httpClient := &http.Client{}
	res, err := httpClient.Get("https://ergast.com/api/f1/current.json")
	if err != nil {
		return baseData, err
	}
	defer res.Body.Close()

	// Read the data and unmarshal it to the structs in f1data.go
	byteValue, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return baseData, err
	}

	// Unmarshal read bytes into JSON object
	err = json.Unmarshal(byteValue, &baseData)
	if err != nil {
		return baseData, err
	}

	log.Println("Data pull complete")
	return baseData, nil
}

func loadPageTemplate(pageFilePath string) (*template.Template, error) {
	log.Println("Loading page template...")
	template, err := template.ParseFiles(pageFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template file: %s", pageFilePath)
	}
	log.Println("Page template load complete")
	return template, nil
}

// Iterate over the Races in the RaceTable and find the first race with a date after today
func getNextRace(raceTable *RaceTable) (*Race, error) {
	for _, race := range raceTable.Races {
		parsedDate, err := time.Parse("2006-01-02 15:04:05Z", race.Date+" "+race.Time)
		if err != nil {
			log.Println(err.Error())
			break
		}

		// Return the race if the time now is before the race + 2 hours
		// Example: If the race starts at 3PM it will still be returned until 5PM
		if time.Now().Before(parsedDate.Add(time.Hour * 2)) {
			return &race, nil
		}
	}
	return nil, errors.New("couldn't find next race")
}

// Function for handling requests to /
func webRoot(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request from: %s", req.RemoteAddr)
	// If we haven't pulled the data from ergast in 24 hours then pull it again (Might increse this in the future)
	if time.Now().After(dataPullTime.Add(dataRefreshPeriod)) {
		log.Println(time.Now(), "is more than 24 hours after", dataPullTime, "reloading data")
		refreshRaceTable()
	}
	race, err := getNextRace(&raceTable)
	if err != nil {
		fmt.Fprintf(w, "error loading next race: %s", err.Error())
	}

	// Render the page template using the race retrieved in the previous function
	pageTemplate.Execute(w, &Page{race})
}

func jsonRoot(w http.ResponseWriter, req *http.Request) {
	if time.Now().After(dataPullTime.Add(dataRefreshPeriod)) {
		log.Println(time.Now(), "is more than 24 hours after", dataPullTime, "reloading data")
		refreshRaceTable()
	}
	race, err := getNextRace(&raceTable)
	if err != nil {
		fmt.Fprintf(w, "error loading next race: %s", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(race)
}

// Reload the page template then return the main page
func templateReload(w http.ResponseWriter, req *http.Request) {
	refreshRaceTable()
	webRoot(w, req)
}
