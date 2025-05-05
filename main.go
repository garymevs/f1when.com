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

const dataRefreshPeriod time.Duration = time.Hour * 1

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
		fmt.Println("Starting dev server on :8085")
		//http.ListenAndServe(":8081", nil)
		log.Fatal(app.Listen(":8085"))
	} else {
		fmt.Println("Starting server on :8080")
		//http.ListenAndServe(":8080", nil)
		log.Fatal(app.Listen(":8080"))
	}
}

// refreshRaceData fetches race data from the Formula 1 API. It can optionally fetch a specific meeting by providing a meetingKey.
func refreshRaceData() error {
	log.Println("Requesting data from api...")
	requestURL := "https://api.formula1.com/v1/event-tracker"

	// Create http client and make GET request
	httpClient := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("GET", requestURL, nil)
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

	// At the moment because the F1Data struct has everything as optional we will never get an error here

	// Unmarshal read bytes into JSON object
	raceData, err = UnmarshalF1Data(byteValue)
	if err != nil {
		return err
	}

	// Sort sessions into date time order
	// 2023-11-16T20:30:00
	if raceData.SeasonContext.Timetables == nil {
		return errors.New("no timetables found")
	}
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
		// Work around to request the next race if we are in a race completed state
		if *timetable.Session == "r" && *timetable.State == "completed" {
			log.Println("Completed race state found. Getting next race...")
			raceIDInt, err := strconv.Atoi(*raceData.FomRaceID)
			if err != nil {
				return err
			}
			raceIDInt += 1
			raceIDString := strconv.Itoa(raceIDInt)
			refreshMeetingData(raceIDString)
			return nil
		}
	}

	log.Println("Data pull complete")
	dataPullTime = time.Now()

	return nil
}

func refreshMeetingData(meetingKey string) error {
	log.Println("Requesting data from api...")
	requestURL := "https://api.formula1.com/v1/event-tracker"
	if meetingKey != "" {
		log.Println("Requesting specific meeting: " + meetingKey)
		requestURL += "/meeting/" + meetingKey
	}

	// Create http client and make GET request
	httpClient := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("GET", requestURL, nil)
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

	// At the moment because the F1Data struct has everything as optional we will never get an error here

	// Unmarshal read bytes into JSON object
	raceData, err = UnmarshalF1Data(byteValue)
	if err != nil {
		return err
	}

	// Sort sessions into date time order
	// 2023-11-16T20:30:00
	if raceData.MeetingContext.Timetables == nil {
		return errors.New("no timetables found")
	}
	sort.Slice(raceData.MeetingContext.Timetables, func(i, j int) bool {
		iDate, _ := time.Parse(time.RFC3339, *raceData.MeetingContext.Timetables[i].StartTime+"Z")
		log.Println(iDate)
		jDate, _ := time.Parse(time.RFC3339, *raceData.MeetingContext.Timetables[j].StartTime+"Z")
		log.Println(jDate)
		return iDate.Before(jDate)
	})

	// Add UTC time to the sessions
	for index, timetable := range raceData.MeetingContext.Timetables {
		startTimeWithGMTOffset, _ := time.Parse(time.RFC3339, *timetable.StartTime+*timetable.GmtOffset)
		raceData.MeetingContext.Timetables[index].StartTimeUTC = startTimeWithGMTOffset.UTC().String()
		endTimeWithGMTOffset, _ := time.Parse(time.RFC3339, *timetable.EndTime+*timetable.GmtOffset)
		raceData.MeetingContext.Timetables[index].EndTimeUTC = endTimeWithGMTOffset.UTC().String()
	}
	raceData.SeasonContext.Timetables = raceData.MeetingContext.Timetables

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
