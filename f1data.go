package main

type BaseData struct {
	MRData struct {
		Total     string    `json:"total"`
		RaceTable RaceTable `json:"RaceTable"`
	}
}

type RaceTable struct {
	Season string `json:"season"`
	Races  []Race `json:"Races"`
}

type Race struct {
	Round string `json:"round"`
	// wiki URL
	URL            string  `json:"url"`
	Name           string  `json:"raceName"`
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	FirstPractice  Session `json:"FirstPractice"`
	SecondPractice Session `json:"SecondPractice"`
	ThirdPractice  Session `json:"ThirdPractice"`
	Qualifying     Session `json:"Qualifying"`
	Sprint         Session `json:"Sprint"`
	Circuit        Circuit `json:"Circuit"`
}

type Circuit struct {
	CircuitId string `json:"circuitId"`
	// wiki URL
	URL      string   `json:"url"`
	Name     string   `json:"circuitName"`
	Location Location `json:"Location"`
}

type Session struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Location struct {
	Lat      string `json:"lat"`
	Long     string `json:"long"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type Page struct {
	Races *[]Race
	Race  *Race
}
