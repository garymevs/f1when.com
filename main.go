package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

// F1 race data
var raceData F1Data
var dataPullTime time.Time

const dataRefreshPeriod time.Duration = time.Hour * 24

var devMode = false
var apiKey = ""

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Immutable: true,
		Views:     engine,
	})

	// check for DEV env variable
	devMode, _ = strconv.ParseBool(os.Getenv("DEV"))

	envFile, _ := godotenv.Read(".env")
	apiKey = envFile["APIKEY"]
	if apiKey == "" {
		log.Fatal("No json source api key supplied. Please specify APIKEY= in .env")
	}

	// Download the json data
	err := refreshRaceData()
	if err != nil {
		log.Fatal(err)
	}

	// Serve static resources
	app.Static("/static", "./static")

	// Serve pages
	app.Get("/", webRoot)
	// app.Get("/round/:round", webRound) !Removed in favour of new API (might be readded)

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

func refreshRaceData() error {
	log.Println("Requesting data from api...")
	// Create http client and make GET request
	httpClient := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("GET", "https://api.formula1.com/v1/event-tracker", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Apikey", apiKey)
	req.Header.Set("Locale", "en")

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Read the data and unmarshal it to the structs in f1data.go
	byteValue, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// TODO: add some actual error handling
	// At the moment because the F1Data struct has everything as optional we will never get an error here

	// Unmarshal read bytes into JSON object
	raceData, err = UnmarshalF1Data(byteValue)
	if err != nil {
		return err
	}

	// Sort sessions into date time order
	// 2023-11-16T20:30:00
	sort.Slice(raceData.SeasonContext.Timetables, func(i, j int) bool {
		iDate, _ := time.Parse(time.RFC3339, *raceData.SeasonContext.Timetables[i].StartTime+"Z")
		log.Println(iDate)
		jDate, _ := time.Parse(time.RFC3339, *raceData.SeasonContext.Timetables[j].StartTime+"Z")
		log.Println(jDate)
		return iDate.Before(jDate)
	})

	// Add UTC time to the sessions
	for index, timetable := range raceData.SeasonContext.Timetables {
		startTimeWithGMTOffset, _ := time.Parse(time.RFC3339, *timetable.StartTime+*timetable.GmtOffset)
		raceData.SeasonContext.Timetables[index].StartTimeUTC = startTimeWithGMTOffset.UTC().String()
		endTimeWithGMTOffset, _ := time.Parse(time.RFC3339, *timetable.EndTime+*timetable.GmtOffset)
		raceData.SeasonContext.Timetables[index].EndTimeUTC = endTimeWithGMTOffset.UTC().String()
	}

	log.Println("Data pull complete")
	dataPullTime = time.Now()

	return nil
}

type Page struct {
	F1Data *F1Data
}

// Function for handling requests to /
func webRoot(c *fiber.Ctx) error {
	log.Printf("%s requested from %s", c.BaseURL(), c.IP())
	if time.Now().After(dataPullTime.Add(dataRefreshPeriod)) {
		log.Println(time.Now(), "is more than 24 hours after", dataPullTime, "reloading data")
		err := refreshRaceData()
		if err != nil {
			c.SendString("error loading next race: " + err.Error())
			return errors.New("error loading next race: " + err.Error())
		}
	}

	// Render the page template using the race retrieved in the previous function
	return c.Render("page-template", &Page{&raceData})
}
