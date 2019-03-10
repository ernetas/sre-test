package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "regexp"
  "strings"
  "time"
)

type Message struct {
  DateOfBirth string`json:"dateOfBirth"`
}

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func handleGet(username string, w http.ResponseWriter) {
  getBirthday(username, w)
}

func handlePut(username string, r *http.Request, w http.ResponseWriter) {
  // Read request body
  b, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintln(w, "Failed to decode JSON, make sure it follows documented format.")
    log.Printf("Failed to decode JSON: ", string(b), err)
    return
  }
  // Decode JSON
  var data Message
  err = json.Unmarshal([]byte(b), &data)
  if err != nil {
    // Give HTTP status code 400 on bad JSON data
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintln(w, "Failed to decode JSON, make sure it follows documented format.")
    log.Printf("Failed to decode JSON: ", string(b), err)
  } else {
    date := data.DateOfBirth
    if validDate(date) {
      writeBirthday(username, date, w)
    } else {
      w.WriteHeader(http.StatusBadRequest)
      todayDate := time.Now().UTC().Format("2006-01-02 UTC")
      // Let's print today's date for clarity, as UTC time may be different from client. (assuming that user can be born yesterday and that server time is in sync with NTP or whatnot...)
      fmt.Fprintln(w, "Invalid date. Must follow format YYYY-MM-DD and be no later than today: ", todayDate)
      log.Printf("Invalid date. Must follow format YYYY-MM-DD and be no later than today: ", todayDate)
    }
  }
}

func main() {
  // We want this to be configurable to make it easier to develop/debug
  port := getPort()

  http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
    // Get username from path
    username := strings.TrimPrefix(r.URL.Path, "/hello/")

    // Validate username (only latin letters are supported)
    if IsLetter(username) {
      // Handle only GET and PUT HTTP methods
      switch r.Method {
        case http.MethodGet:
          handleGet(username, w)
        case http.MethodPut:
          handlePut(username, r, w)
        default:
          // Log undefined methods and give HTTP 501 status code
          w.WriteHeader(http.StatusNotImplemented)
          fmt.Fprintln(w, "Unsupported HTTP method. Only GET and PUT are supported.")
          log.Printf("Call to unsupported HTTP method: ", r.Method)
      }
    } else {
      // Log bad usernames (we don't log non /hello paths to reduce unnecessary log footprint - assuming "/hello" as an "obscure" authentication)
      w.WriteHeader(http.StatusBadRequest)
      fmt.Sprintf("Please use a valid username consisting of only latin letters")
    }
  })

  // Let's log when the app boots up
  log.Printf("app listening on port :%s", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}