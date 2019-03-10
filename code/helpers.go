package main

import (
  "fmt"
  "log"
  "os"
  "math"
  "strconv"
  "strings"
  "time"
)

func getPort() string {
  port := os.Getenv("PORT")
  if len(port) == 0 {
    // Port 9000 is an easy choice for local development as it wouldn't require root permissions
    return "9000"
  }
  return port
}

func getCurrentTime() (int, int, int) {
  currentTime := time.Now().UTC()
  currentYear, _ := strconv.Atoi(currentTime.Format("2006"))
  currentMonth, _ := strconv.Atoi(currentTime.Format("01"))
  currentDay, _ := strconv.Atoi(currentTime.Format("02"))
  return currentYear, currentMonth, currentDay
}

func Date(year, month, day int) time.Time {
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func validDate(date string) bool {
  t, err := time.Parse("2006-01-02", date)
  if err != nil {
    log.Println(err)
    // Return false if date format is invalid
    return false
  }
  currentYear, currentMonth, currentDay := getCurrentTime()
  birthdayYear, _ := strconv.Atoi(t.Format("2006"))
  birthdayMonth, _ := strconv.Atoi(t.Format("01"))
  birthdayDay, _ := strconv.Atoi(t.Format("02"))
  // Return false if date is after yesterday
  if Date(birthdayYear, birthdayMonth, birthdayDay).Unix() > Date(currentYear, currentMonth, currentDay).Unix() {
    return false
  }
  return true
}

func daysBetween(date string, targetYear int) int {
  dates := strings.Split(date, "-")
  _, currentMonth, currentDay := getCurrentTime()
  birthdayDay, _ := strconv.Atoi(dates[2])
  birthdayMonth, _ := strconv.Atoi(dates[1])

  t := Date(targetYear, birthdayMonth, birthdayDay)
  hours := time.Now().UTC().Sub(t).Hours()
  // Check if it's today
  if (currentMonth == birthdayMonth) && (currentDay == birthdayDay) {
    return 0
  } else {
    if int(math.Abs(hours)) < 24 {
      return 1
    } else {
      return int(math.Abs(hours)/24)
    }
  }
}

func birthdayAlreadyHappenedThisYear(date string) bool {
  _, currentMonth, currentDay := getCurrentTime()
  dates := strings.Split(date, "-")
  birthdayDay, _ := strconv.Atoi(dates[2])
  birthdayMonth, _ := strconv.Atoi(dates[1])

  if birthdayMonth > int(currentMonth) {
    return false
  } else if birthdayMonth < int(currentMonth) {
    return true
  } else if birthdayDay > int(currentDay) {
    return false
  } else {
    return true
  }
}

func daysTill(date string) int {
  currentYear, _, _ := getCurrentTime()
  bdate := date
  targetYear := currentYear
  if birthdayAlreadyHappenedThisYear(bdate) == true {
    targetYear = currentYear + 1
  }
  dbetween := daysBetween(bdate, targetYear)
  return dbetween
}

func templateMessage(date string) string {
  days := daysTill(date)
  if days > 0 {
    return fmt.Sprintf("Your birthday is in %d day(s)", days)
  } else {
    return "Happy birthday!"
  }
}