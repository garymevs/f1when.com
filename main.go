package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

// F1 season data
var raceTable RaceTable
var dataPullTime time.Time

const dataRefreshPeriod time.Duration = time.Hour * 24

var devMode = false

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Immutable: true,
		Views:     engine,
	})

	// check for DEV env variable
	devMode, _ = strconv.ParseBool(os.Getenv("DEV"))

	// Download the json data from ergast.com
	refreshRaceTable()

	// Serve static resources
	app.Static("/static", "./static")

	// Serve pages
	app.Get("/", webRoot)
	app.Get("/round/:round", webRound)

	// Start the web server
	if devMode {
		fmt.Println("Starting dev server on :8081")
		//http.ListenAndServe(":8081", nil)
		log.Fatal(app.Listen(":8081"))
	} else {
		fmt.Println("Starting server on :8080")
		//http.ListenAndServe(":8080", nil)
		log.Fatal(app.Listen(":8080"))
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
	baseData, err := loadSpecificSeasonData("current")
	if err != nil {
		return baseData, err
	}
	_, seasonEnded, _ := getNextRace(&baseData.MRData.RaceTable)
	if seasonEnded {
		log.Println("Current season is marked as ended. Pulling season for this year...")
		currentSeason, err := strconv.Atoi(baseData.MRData.RaceTable.Season)
		if err != nil {
			return baseData, err
		}
		if time.Now().Year() > currentSeason {
			return loadSpecificSeasonData(strconv.Itoa(time.Now().Year()))
		}
	}
	return baseData, err
}

// Load season using the provided string
func loadSpecificSeasonData(season string) (BaseData, error) {
	if devMode {
		log.Println("Running as dev, loading data.json")
		jsonFile, err := os.Open("data.json")
		if err != nil {
			return BaseData{}, errors.New("running as dev mode but unable to load data.json")
		}
		defer jsonFile.Close()
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return BaseData{}, errors.New("error reading dev mode data.json")
		}
		baseData := BaseData{}
		json.Unmarshal(byteValue, &baseData)
		return baseData, nil
	}
	log.Println("Requesting data from ergast...")
	baseData := BaseData{}
	// Data from here: https://ergast.com/api/f1/current.json
	// Create http client and make GET request
	httpClient := &http.Client{Timeout: time.Second * 30}
	res, err := httpClient.Get("https://ergast.com/api/f1/" + season + ".json")
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

// Iterate over the Races in the RaceTable and find the first race with a date after today
// Returns Race, seasonEnded, err
func getNextRace(raceTable *RaceTable) (*Race, bool, error) {
	for raceIdx, race := range raceTable.Races {
		parsedDate, err := time.Parse("2006-01-02 15:04:05Z", race.Date+" "+race.Time)
		if err != nil {
			log.Println(err.Error())
			break
		}

		// Return the race if the time now is before the race + 2 hours
		// Example: If the race starts at 3PM it will still be returned until 5PM
		if time.Now().Before(parsedDate.Add(time.Hour * 2)) {
			return &race, false, nil
		}

		// If we reach the last element and the previous block isn't matched then we have probably finished the season
		if raceIdx >= len(raceTable.Races)-1 {
			return &Race{}, true, nil
		}
	}
	return nil, false, errors.New("couldn't find next race")
}

// Function for handling requests to /
func webRoot(c *fiber.Ctx) error {
	log.Printf("%s requested from %s", c.BaseURL(), c.IP())
	// If we haven't pulled the data from ergast in 24 hours then pull it again (Might increse this in the future)
	if time.Now().After(dataPullTime.Add(dataRefreshPeriod)) {
		log.Println(time.Now(), "is more than 24 hours after", dataPullTime, "reloading data")
		refreshRaceTable()
	}
	race, _, err := getNextRace(&raceTable)
	if err != nil {
		c.SendString("error loading next race: " + err.Error())
		return errors.New("error loading next race: " + err.Error())
	}

	// Render the page template using the race retrieved in the previous function
	return c.Render("page-template", &Page{&raceTable.Races, race})
}

func webRound(c *fiber.Ctx) error {
	log.Printf("%s requested from %s", c.BaseURL(), c.IP())
	// If we haven't pulled the data from ergast in 24 hours then pull it again (Might increse this in the future)
	if time.Now().After(dataPullTime.Add(dataRefreshPeriod)) {
		log.Println(time.Now(), "is more than 24 hours after", dataPullTime, "reloading data")
		refreshRaceTable()
	}
	//race, _, err := getNextRace(&raceTable)
	round, err := strconv.Atoi(c.Params("round"))
	race := &raceTable.Races[round-1]
	if err != nil {
		c.SendString("error loading round: " + err.Error())
		return errors.New("error loading round: " + err.Error())
	}

	// Render the page template using the race retrieved in the previous function
	return c.Render("page-template", &Page{&raceTable.Races, race})
}
