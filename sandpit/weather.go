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

func main() {
  // Replace the following line with code to fetch the XML data from the FTP server
  xmlData, err := getWeather()
  if err != nil {
     fmt.Println(err)
  }

  // Parse XML data
  product, err := parseXML(xmlData)
  if err != nil {
    log.Fatal(err)
  }

  // Access parsed data
  fmt.Printf("Product Version: %s\n", product.Version)
  fmt.Printf("Product Xsi: %s\n", product.Xsi)
  fmt.Printf("Product NoNamespaceSchemaLocation: %s\n", product.NoNamespaceSchemaLocation)

  fmt.Printf("AMOC Identifier: %s\n", product.Amoc.Identifier)
  fmt.Printf("AMOC Issue Time (UTC): %s\n", product.Amoc.IssueTimeUtc)
  fmt.Printf("AMOC Issue Time (Local): %s %s\n", product.Amoc.IssueTimeLocal.Text, product.Amoc.IssueTimeLocal.Tz)
  fmt.Printf("AMOC Sent Time: %s\n", product.Amoc.SentTime)
  fmt.Printf("AMOC Expiry Time: %s\n", product.Amoc.ExpiryTime)
  fmt.Printf("AMOC Validity Begin Time (Local): %s %s\n", product.Amoc.ValidityBgnTimeLocal.Text, product.Amoc.ValidityBgnTimeLocal.Tz)
  fmt.Printf("AMOC Validity End Time (Local): %s %s\n", product.Amoc.ValidityEndTimeLocal.Text, product.Amoc.ValidityEndTimeLocal.Tz)
  fmt.Printf("AMOC Next Routine Issue Time (UTC): %s\n", product.Amoc.NextRoutineIssueTimeUtc)
  fmt.Printf("AMOC Next Routine Issue Time (Local): %s %s\n", product.Amoc.NextRoutineIssueTimeLocal.Text, product.Amoc.NextRoutineIssueTimeLocal.Tz)
  fmt.Printf("AMOC Status: %s\n", product.Amoc.Status)
  fmt.Printf("AMOC Service: %s\n", product.Amoc.Service)
  fmt.Printf("AMOC SubService: %s\n", product.Amoc.SubService)
  fmt.Printf("AMOC Product Type: %s\n", product.Amoc.ProductType)
  fmt.Printf("AMOC Phase: %s\n", product.Amoc.Phase)

  for _, area := range product.Forecast.Area {
    fmt.Printf("\nArea: %s (%s)\n", area.Description, area.Type)
    for _, period := range area.ForecastPeriod {
        fmt.Printf("  Period: %s to %s\n", period.StartTimeLocal, period.EndTimeLocal)
        for _, elem := range period.Element {
            fmt.Printf("    Element Type: %s, Value: %s %s\n", elem.Type, elem.Text, elem.Units)
        }
        for _, text := range period.Text {
            fmt.Printf("    Text Type: %s, Value: %s\n", text.Type, text.Text)
        }
    }
  }
}

