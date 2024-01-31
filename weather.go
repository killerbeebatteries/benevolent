package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jlaffaye/ftp"
)

// Product was generated 2024-01-30 19:32:16 using Zek: https://github.com/miku/zek
type Product struct {
	XMLName                   xml.Name `xml:"product"`
	Text                      string   `xml:",chardata"`
	Version                   string   `xml:"version,attr"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	Amoc                      struct {
		Text   string `xml:",chardata"`
		Source struct {
			Text       string `xml:",chardata"`
			Sender     string `xml:"sender"`     // Australian Government Bur...
			Region     string `xml:"region"`     // New South Wales
			Office     string `xml:"office"`     // NSWRO
			Copyright  string `xml:"copyright"`  // http://www.bom.gov.au/oth...
			Disclaimer string `xml:"disclaimer"` // http://www.bom.gov.au/oth...
		} `xml:"source"`
		Identifier     string `xml:"identifier"`     // IDN11060
		IssueTimeUtc   string `xml:"issue-time-utc"` // 2024-01-30T04:50:00Z
		IssueTimeLocal struct {
			Text string `xml:",chardata"` // 2024-01-30T15:50:00+11:00...
			Tz   string `xml:"tz,attr"`
		} `xml:"issue-time-local"`
		SentTime             string `xml:"sent-time"`   // 2024-01-30T04:50:06Z
		ExpiryTime           string `xml:"expiry-time"` // 2024-01-31T04:50:00Z
		ValidityBgnTimeLocal struct {
			Text string `xml:",chardata"` // 2024-01-30T17:00:00+11:00...
			Tz   string `xml:"tz,attr"`
		} `xml:"validity-bgn-time-local"`
		ValidityEndTimeLocal struct {
			Text string `xml:",chardata"` // 2024-02-06T23:59:59+11:00...
			Tz   string `xml:"tz,attr"`
		} `xml:"validity-end-time-local"`
		NextRoutineIssueTimeUtc   string `xml:"next-routine-issue-time-utc"` // 2024-01-30T17:25:00Z
		NextRoutineIssueTimeLocal struct {
			Text string `xml:",chardata"` // 2024-01-31T04:25:00+11:00...
			Tz   string `xml:"tz,attr"`
		} `xml:"next-routine-issue-time-local"`
		Status      string `xml:"status"`       // O
		Service     string `xml:"service"`      // WSP
		SubService  string `xml:"sub-service"`  // FPR
		ProductType string `xml:"product-type"` // F
		Phase       string `xml:"phase"`        // NEW
	} `xml:"amoc"`
	Forecast struct {
		Text string `xml:",chardata"`
		Area []struct {
			Text           string `xml:",chardata"`
			Aac            string `xml:"aac,attr"`
			Description    string `xml:"description,attr"`
			Type           string `xml:"type,attr"`
			ParentAac      string `xml:"parent-aac,attr"`
			ForecastPeriod []struct {
				Chardata       string `xml:",chardata"`
				Index          string `xml:"index,attr"`
				StartTimeLocal string `xml:"start-time-local,attr"`
				EndTimeLocal   string `xml:"end-time-local,attr"`
				StartTimeUtc   string `xml:"start-time-utc,attr"`
				EndTimeUtc     string `xml:"end-time-utc,attr"`
				Element        []struct {
					Text  string `xml:",chardata"` // 17, 11, 0 to 4 mm, 22, 26...
					Type  string `xml:"type,attr"`
					Units string `xml:"units,attr"`
				} `xml:"element"`
				Text []struct {
					Text string `xml:",chardata"` // Possible shower., 40%!,(MISSING) Sh...
					Type string `xml:"type,attr"`
				} `xml:"text"`
			} `xml:"forecast-period"`
		} `xml:"area"`
	} `xml:"forecast"`
}

func parseXML(xmlData string) (*Product, error) {
	var product Product
	err := xml.NewDecoder(strings.NewReader(xmlData)).Decode(&product)
	return &product, err
}

func getWeather() (string, error) {
	// FTP server details
	ftpServer := os.Getenv("FTP_SERVER")
	ftpUser := os.Getenv("FTP_USER")
	ftpPassword := os.Getenv("FTP_PASSWORD")
	ftpFilePath := os.Getenv("FTP_FILE_PATH") + "/" + os.Getenv("FTP_FILE_NAME")

	fmt.Println("Retrieving data from " + ftpServer + ftpFilePath + "...")

	// Connect to FTP server
	conn, err := ftp.Dial(fmt.Sprintf("%s:%d", ftpServer, 21))
	if err != nil {
		return "", err
	}
	defer conn.Quit()

	// Login to the FTP server
	err = conn.Login(ftpUser, ftpPassword)
	if err != nil {
		return "", err
	}

	// Open the file on the FTP server
	r, err := conn.Retr(ftpFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// Read the content of the file
	xmlBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(xmlBytes), nil
}

func handleWeather(location string) ([]string, error) {
	var result []string
	xmlData, err := getWeather()

	if err != nil {
		return result, err
	}
	product, err := parseXML(xmlData)
	if err != nil {
		return result, err
	}

	for _, area := range product.Forecast.Area {
		var forecast_area string

		if area.Type == "location" && area.Description == location {
      forecast_area = fmt.Sprintf("Area: %s (%s)\n", area.Description, area.Type)
			result = append(result, forecast_area)
			var forecast_period, temp_min, temp_max, precipitation_range, precis, chance_of_rain string

      for _, period := range area.ForecastPeriod[0:2] {
				forecast_period = fmt.Sprintf("Period: %s to %s\n", period.StartTimeLocal, period.EndTimeLocal)
				result = append(result, forecast_period)

				for _, elem := range period.Element {
					if elem.Type == "air_temperature_minimum" {
						temp_min = fmt.Sprintf("Minimum temperature of: %s %s\n", elem.Text, elem.Units)
						result = append(result, temp_min)
					}
					if elem.Type == "air_temperature_maximum" {
						temp_max = fmt.Sprintf("Maximum temperature of: %s %s\n", elem.Text, elem.Units)
						result = append(result, temp_max)
					}
					if elem.Type == "precipitation_range" {
						precipitation_range = fmt.Sprintf("How much will it rain? %s %s\n", elem.Text, elem.Units)
						result = append(result, precipitation_range)
					}
				}

				for _, text := range period.Text {
					if text.Type == "precis" {
						precis = fmt.Sprintf("The forecast: %s\n", text.Text)
						result = append(result, precis)
					}
					if text.Type == "probability_of_precipitation" {
						chance_of_rain = fmt.Sprintf("Chance of rain: %s\n", text.Text)
						result = append(result, chance_of_rain)
					}
				}
			}
		}
	}

	if len(result) == 0 {
		result = append(result, "Sorry, I have no weather information for that location.")
	}

	fmt.Println(result)
	return result, nil
}
