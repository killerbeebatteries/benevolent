package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
  "io/ioutil"
  "os"

  "github.com/jlaffaye/ftp"
)

// Product represents the root structure of the XML data
type Product struct {
	XMLName xml.Name `xml:"product"`
	AMOC    AMOC     `xml:"amoc"`
	Forecast Forecast `xml:"forecast"`
}

// AMOC represents the AMOC section of the XML data
type AMOC struct {
	Source                 Source `xml:"source"`
	Identifier             string `xml:"identifier"`
	IssueTimeUTC           string `xml:"issue-time-utc"`
	IssueTimeLocal         string `xml:"issue-time-local"`
	SentTime               string `xml:"sent-time"`
	ExpiryTime             string `xml:"expiry-time"`
	ValidityBgnTimeLocal   string `xml:"validity-bgn-time-local"`
	ValidityEndTimeLocal   string `xml:"validity-end-time-local"`
	NextRoutineIssueTimeUTC string `xml:"next-routine-issue-time-utc"`
	NextRoutineIssueTimeLocal string `xml:"next-routine-issue-time-local"`
	Status                 string `xml:"status"`
	Service                string `xml:"service"`
	SubService             string `xml:"sub-service"`
	ProductType            string `xml:"product-type"`
	Phase                  string `xml:"phase"`
}

// Source represents the source section of the XML data
type Source struct {
	Sender     string `xml:"sender"`
	Region     string `xml:"region"`
	Office     string `xml:"office"`
	Copyright  string `xml:"copyright"`
	Disclaimer string `xml:"disclaimer"`
}

// Forecast represents the forecast section of the XML data
type Forecast struct {
	Areas []Area `xml:"area"`
}

// Area represents a forecast area
type Area struct {
	AAC         string          `xml:"aac,attr"`
	Description string          `xml:"description,attr"`
	Type        string          `xml:"type,attr"`
	ParentAAC   string          `xml:"parent-aac,attr,omitempty"`
	Periods     []ForecastPeriod `xml:"forecast-period"`
}

// ForecastPeriod represents details of a forecast period
type ForecastPeriod struct {
	Index                 int     `xml:"index,attr"`
	StartTimeLocal        string  `xml:"start-time-local,attr"`
	EndTimeLocal          string  `xml:"end-time-local,attr"`
	StartTimeUTC          string  `xml:"start-time-utc,attr"`
	EndTimeUTC            string  `xml:"end-time-utc,attr"`
	IconCode              int     `xml:"element[type='forecast_icon_code']"`
	Precis                string  `xml:"text[type='precis']"`
	ProbabilityOfPrecip   string  `xml:"text[type='probability_of_precipitation']"`
	PrecipitationRange    string  `xml:"element[type='precipitation_range'],omitempty"`
	TemperatureMinimum    float64 `xml:"element[type='air_temperature_minimum'],attr,omitempty"`
	TemperatureMaximum    float64 `xml:"element[type='air_temperature_maximum'],attr,omitempty"`
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

func main() {
	// Replace the following line with code to fetch the XML data from the FTP server
  xmlData, err := getWeather()
  if err != nil {
     fmt.Println(err)
  }

	// Parse XML data
	parsedData, err := parseXML(xmlData)
	if err != nil {
		log.Fatal(err)
	}

	// Access parsed data
	fmt.Printf("Identifier: %s\n", parsedData.AMOC.Identifier)
	fmt.Printf("Issue Time (UTC): %s\n", parsedData.AMOC.IssueTimeUTC)
	fmt.Printf("Validity Begin Time (Local): %s\n", parsedData.AMOC.ValidityBgnTimeLocal)

	for _, area := range parsedData.Forecast.Areas {
		fmt.Printf("\nArea: %s (%s)\n", area.Description, area.Type)
		for _, period := range area.Periods {
			fmt.Printf("  Period: %s to %s\n", period.StartTimeLocal, period.EndTimeLocal)
			fmt.Printf("  Precis: %s\n", period.Precis)
			fmt.Printf("  Probability of Precipitation: %s\n", period.ProbabilityOfPrecip)
			fmt.Printf("  Temperature Range: %.2f°C to %.2f°C\n", period.TemperatureMinimum, period.TemperatureMaximum)
			if period.PrecipitationRange != "" {
				fmt.Printf("  Precipitation Range: %s\n", period.PrecipitationRange)
			}
		}
	}
}

