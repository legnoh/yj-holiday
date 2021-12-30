package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Holiday struct {
	Date time.Time
	Name string
}

type HolidayJSON struct {
	Name      string `json:"name"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

var (
	calendar     = ics.NewCalendarFor("yahoojapan.holiday.tool.legnoh.lkj.io")
	icsFilePath  = "./htdocs/yahoojapan/holidays.ics"
	jsonFilePath = "./htdocs/yahoojapan/holidays.json"
)

func main() {

	yjHolidays := []Holiday{}

	calendar.SetXWRCalName("YJ holidays")
	calendar.SetColor("#FF2968")
	calendar.SetXWRTimezone("Asia/Tokyo")
	calendar.SetTzid("Asia/Tokyo")

	events, err := getEvents()
	if err != nil {
		panic(err)
	}

	// Event processing...
	for _, v := range events {
		yjHolidays = addEvent(v, yjHolidays)

		if v.Date.Month() == time.January && v.Date.Day() == 1 {
			newYearHolidays := []time.Time{
				time.Date(v.Date.Year(), time.January, 2, 0, 0, 0, 0, v.Date.Location()),
				time.Date(v.Date.Year(), time.January, 3, 0, 0, 0, 0, v.Date.Location()),
				time.Date(v.Date.Year(), time.January, 4, 0, 0, 0, 0, v.Date.Location()),
				time.Date(v.Date.Year(), time.December, 29, 0, 0, 0, 0, v.Date.Location()),
				time.Date(v.Date.Year(), time.December, 30, 0, 0, 0, 0, v.Date.Location()),
				time.Date(v.Date.Year(), time.December, 31, 0, 0, 0, 0, v.Date.Location()),
			}
			for _, h := range newYearHolidays {
				nyh := Holiday{Date: h, Name: "年末年始休暇"}
				yjHolidays = addEvent(nyh, yjHolidays)
			}
		}

		if v.Date.Weekday() == time.Saturday {
			yjHolidays = addEvent(getSpecialHoliday(v.Date, yjHolidays), yjHolidays)
		}
	}

	// Create ICS
	calendar_bytes := []byte(calendar.Serialize())
	icsFile, err := os.Create(icsFilePath)
	if err != nil {
		panic(err)
	}
	defer icsFile.Close()
	if err := ioutil.WriteFile(icsFilePath, calendar_bytes, 0666); err != nil {
		panic(err)
	}

	// Create JSON
	json_bytes, err := json.Marshal(yjHolidays)
	if err != nil {
		panic(err)
	}
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	if err := ioutil.WriteFile(jsonFilePath, json_bytes, 0666); err != nil {
		panic(err)
	}
	fmt.Println("success to write json/ics file")
}

func addEvent(v Holiday, yjHolidays []Holiday) []Holiday {

	uuidObj, _ := uuid.NewUUID()
	event := calendar.AddEvent(uuidObj.String())
	event.SetSummary(v.Name)

	event.SetStartAt(v.Date)
	event.SetEndAt(v.Date.AddDate(0, 0, 1))
	event.SetTimeTransparency(ics.TransparencyTransparent)
	event.SetDtStampTime(time.Now())

	return append(yjHolidays, v)
}

func getSpecialHoliday(date time.Time, yjHolidays []Holiday) Holiday {

	for {
		date = date.AddDate(0, 0, -1)
		if !isHoliday(date, yjHolidays) {
			return Holiday{Date: date, Name: "振替特別休日"}
		}
	}
}

func isHoliday(date time.Time, yjHolidays []Holiday) bool {
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return true
	}
	for _, v := range yjHolidays {
		if date.Year() == v.Date.Year() && date.Month() == v.Date.Month() && date.Day() == v.Date.Day() {
			return true
		}
	}
	return false
}

func getEvents() ([]Holiday, error) {

	var (
		csvUrl        = "https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv"
		holidays      []Holiday
		csvDateFormat = "2006/1/2 -07:00"
	)

	resp, err := http.Get(csvUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str, _, err := transform.String(japanese.ShiftJIS.NewDecoder(), string(body))
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(strings.NewReader(str))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, v := range records {
		name := v[1]
		date, err := time.Parse(csvDateFormat, v[0]+" +09:00")
		if err != nil {
			continue
		}
		holiday := Holiday{Name: name, Date: date}
		holidays = append(holidays, holiday)
	}
	return holidays, nil
}

func (h Holiday) MarshalJSON() ([]byte, error) {
	v, err := json.Marshal(HolidayJSON{
		Name:      h.Name,
		Date:      h.Date.Format("2006-01-02"),
		StartTime: strconv.Itoa(int(h.Date.Unix())),
		EndTime:   strconv.Itoa(int(h.Date.Add(time.Second * 86399).Unix())),
	})
	return v, err
}
