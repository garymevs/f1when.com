// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    f1Data, err := UnmarshalF1Data(bytes)
//    bytes, err = f1Data.Marshal()

package main

import "encoding/json"

func UnmarshalF1Data(data []byte) (F1Data, error) {
	var r F1Data
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *F1Data) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type F1Data struct {
	RaceHubID              *string            `json:"raceHubId,omitempty"`
	Locale                 *string            `json:"locale,omitempty"`
	CreatedAt              *string            `json:"createdAt,omitempty"`
	UpdatedAt              *string            `json:"updatedAt,omitempty"`
	FomRaceID              *string            `json:"fomRaceId,omitempty"`
	BrandColourHexadecimal *string            `json:"brandColourHexadecimal,omitempty"`
	Headline               *string            `json:"headline,omitempty"`
	CuratedSection         *CuratedSection    `json:"curatedSection,omitempty"`
	CircuitSmallImage      *CircuitSmallImage `json:"circuitSmallImage,omitempty"`
	Links                  []Link             `json:"links,omitempty"`
	SeasonContext          *SeasonContext     `json:"seasonContext,omitempty"`
	RaceResults            []RaceResult       `json:"raceResults,omitempty"`
	Race                   *Race              `json:"race,omitempty"`
	SeasonYearImage        *string            `json:"seasonYearImage,omitempty"`
	SessionLinkSets        *SessionLinkSets   `json:"sessionLinkSets,omitempty"`
	MeetingContext         *MeetingContext    `json:"meetingContext,omitempty"`
}

type CircuitSmallImage struct {
	Title             *string `json:"title,omitempty"`
	Path              *string `json:"path,omitempty"`
	URL               *string `json:"url,omitempty"`
	PublicID          *string `json:"public_id,omitempty"`
	RawTransformation *string `json:"raw_transformation,omitempty"`
}

type CuratedSection struct {
	ContentType *string `json:"contentType,omitempty"`
	Title       *string `json:"title,omitempty"`
	Items       []Item  `json:"items,omitempty"`
}

type Item struct {
	ID              *string    `json:"id,omitempty"`
	UpdatedAt       *string    `json:"updatedAt,omitempty"`
	RealUpdatedAt   *string    `json:"realUpdatedAt,omitempty"`
	Locale          *string    `json:"locale,omitempty"`
	Title           *string    `json:"title,omitempty"`
	Slug            *string    `json:"slug,omitempty"`
	ArticleType     *string    `json:"articleType,omitempty"`
	MetaDescription *string    `json:"metaDescription,omitempty"`
	Thumbnail       *Thumbnail `json:"thumbnail,omitempty"`
	IsProtected     *bool      `json:"isProtected,omitempty"`
	MediaIcon       *string    `json:"mediaIcon,omitempty"`
}

type Thumbnail struct {
	Caption *string `json:"caption,omitempty"`
	Image   *Image  `json:"image,omitempty"`
}

type Image struct {
	Title             *string     `json:"title,omitempty"`
	Path              *string     `json:"path,omitempty"`
	URL               *string     `json:"url,omitempty"`
	PublicID          *string     `json:"public_id,omitempty"`
	RawTransformation *string     `json:"raw_transformation,omitempty"`
	Width             *int64      `json:"width,omitempty"`
	Height            *int64      `json:"height,omitempty"`
	Renditions        *Renditions `json:"renditions,omitempty"`
}

type Renditions struct {
	The2Col       *string `json:"2col,omitempty"`
	The2ColRetina *string `json:"2col-retina,omitempty"`
	The3Col       *string `json:"3col,omitempty"`
	The3ColRetina *string `json:"3col-retina,omitempty"`
	The4Col       *string `json:"4col,omitempty"`
	The4ColRetina *string `json:"4col-retina,omitempty"`
	The6Col       *string `json:"6col,omitempty"`
	The6ColRetina *string `json:"6col-retina,omitempty"`
	The9Col       *string `json:"9col,omitempty"`
	The9ColRetina *string `json:"9col-retina,omitempty"`
}

type Link struct {
	Text *string `json:"text,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type MeetingContext struct {
	Season      *string     `json:"season,omitempty"`
	MeetingKey  *string     `json:"meetingKey,omitempty"`
	IsTestEvent *bool       `json:"isTestEvent,omitempty"`
	State       *string     `json:"state,omitempty"`
	SeasonState *string     `json:"seasonState,omitempty"`
	Timetables  []Timetable `json:"timetables,omitempty"`
}

type Timetable struct {
	State        *string `json:"state,omitempty"`
	Session      *string `json:"session,omitempty"`
	GmtOffset    *string `json:"gmtOffset,omitempty"`
	Description  *string `json:"description,omitempty"`
	EndTime      *string `json:"endTime,omitempty"`
	EndTimeUTC   string
	StartTime    *string `json:"startTime,omitempty"`
	StartTimeUTC string
}

type Race struct {
	MeetingCountryName  *string `json:"meetingCountryName,omitempty"`
	MeetingStartDate    *string `json:"meetingStartDate,omitempty"`
	MeetingOfficialName *string `json:"meetingOfficialName,omitempty"`
	MeetingEndDate      *string `json:"meetingEndDate,omitempty"`
}

type RaceResult struct {
	DriverTLA        *string `json:"driverTLA,omitempty"`
	DriverFirstName  *string `json:"driverFirstName,omitempty"`
	DriverLastName   *string `json:"driverLastName,omitempty"`
	TeamName         *string `json:"teamName,omitempty"`
	PositionNumber   *string `json:"positionNumber,omitempty"`
	RaceTime         *string `json:"raceTime,omitempty"`
	TeamColourCode   *string `json:"teamColourCode,omitempty"`
	GapToLeader      *string `json:"gapToLeader,omitempty"`
	DriverImage      *string `json:"driverImage,omitempty"`
	DriverNameFormat *string `json:"driverNameFormat,omitempty"`
}

type SeasonContext struct {
	ID                           *string     `json:"id,omitempty"`
	ContentType                  *string     `json:"contentType,omitempty"`
	CreatedAt                    *string     `json:"createdAt,omitempty"`
	UpdatedAt                    *string     `json:"updatedAt,omitempty"`
	Locale                       *string     `json:"locale,omitempty"`
	SeasonYear                   *string     `json:"seasonYear,omitempty"`
	CurrentOrNextMeetingKey      *string     `json:"currentOrNextMeetingKey,omitempty"`
	State                        *string     `json:"state,omitempty"`
	EventState                   *string     `json:"eventState,omitempty"`
	LiveEventID                  *string     `json:"liveEventId,omitempty"`
	LiveTimingsSource            *string     `json:"liveTimingsSource,omitempty"`
	LiveBlog                     *LiveBlog   `json:"liveBlog,omitempty"`
	SeasonState                  *string     `json:"seasonState,omitempty"`
	RaceListingOverride          *int64      `json:"raceListingOverride,omitempty"`
	DriverAndTeamListingOverride *int64      `json:"driverAndTeamListingOverride,omitempty"`
	Timetables                   []Timetable `json:"timetables,omitempty"`
	ReplayBaseURL                *string     `json:"replayBaseUrl,omitempty"`
	SeasonContextUIState         *int64      `json:"seasonContextUIState,omitempty"`
}

type LiveBlog struct {
	ContentType     *string `json:"contentType,omitempty"`
	Title           *string `json:"title,omitempty"`
	ScribbleEventID *string `json:"scribbleEventId,omitempty"`
}

type SessionLinkSets struct {
	ReplayLinks []ReplayLink `json:"replayLinks,omitempty"`
}

type ReplayLink struct {
	Text     *string `json:"text,omitempty"`
	URL      *string `json:"url,omitempty"`
	LinkType *string `json:"linkType,omitempty"`
	Session  *string `json:"session,omitempty"`
}
