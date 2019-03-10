package main

import (
  "testing"
  "strconv"
)

func TestGetPort(t *testing.T) {
  if port, _ := strconv.Atoi(getPort()); port != 9000 {
    t.Errorf("Port was not 9000. Maybe it's configured as env. variable for tests?")
  }
}

func TestValidDate(t *testing.T) {
  if validDate("dd") {
    t.Errorf("dd is not a valid date")
  }
  if validDate("20110101") {
    t.Errorf("20110101 is not a valid date")
  }
  if validDate("2011-13-01") {
    t.Errorf("2011-13-01 is not a valid date")
  }
  if validDate("2090-01-01") {
    t.Errorf("2090-01-01 is not a valid date")
  }
}

func TestBirthdayAlreadyHappenedThisYear(t *testing.T) {
  if birthdayAlreadyHappenedThisYear("2010-12-32") == true {
    t.Errorf("2010-12-32 birthday did not happen this year.")
  }
}

func TestDaysTill(t *testing.T) {
  if daysTill("2010-01-01") < 0 {
    t.Errorf("There should be a positive number of days till 2010 01 01")
  }
}
