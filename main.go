package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/sirupsen/logrus"
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

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	yjHolidays := []Holiday{}

	calendar.SetXWRCalName("YJ holidays")
	calendar.SetColor("#FF2968")
	calendar.SetXWRTimezone("Asia/Tokyo")
	calendar.SetTzid("Asia/Tokyo")

	events, err := getEvents()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	// Event processing...
	for _, v := range events {
		yjHolidays = addEvent(v, yjHolidays)

		// 年末年始休暇
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
				nyh := Holiday{Date: h, Name: "年末年始休日"}
				yjHolidays = addEvent(nyh, yjHolidays)
			}
		}

		// 土曜日の場合は前の営業日を探してそれも休日にする
		if v.Date.Weekday() == time.Saturday {
			yjHolidays = addEvent(getSpecialHoliday(v.Date, yjHolidays), yjHolidays)
		}
	}

	// Create ICS
	calendar_string := calendar.Serialize()

	// go-ical経由で付与できないCOLOR要素を付与する
	re := regexp.MustCompile(`COLOR:#FF2968`)
	calendar_string = re.ReplaceAllString(calendar_string, "X-APPLE-CALENDAR-COLOR:#FF2968")

	calendar_bytes := []byte(calendar_string)
	icsFile, err := os.Create(icsFilePath)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	defer icsFile.Close()
	if err := os.WriteFile(icsFilePath, calendar_bytes, 0666); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	// Create JSON
	json_bytes, err := json.Marshal(yjHolidays)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	if err := os.WriteFile(jsonFilePath, json_bytes, 0666); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	fmt.Println("success to write json/ics file")
}

func addEvent(v Holiday, yjHolidays []Holiday) []Holiday {

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	dtProperty := &ics.KeyValues{Key: "VALUE", Value: []string{"DATE"}}
	uid := base64.StdEncoding.EncodeToString([]byte(v.Name + v.Date.Format("20060102")))

	log.Infof("add: %s", v)
	event := calendar.AddEvent(uid)
	event.SetAllDayStartAt(v.Date, dtProperty)
	event.SetAllDayEndAt(v.Date, dtProperty)
	event.SetSummary(v.Name)
	event.SetTimeTransparency(ics.TransparencyOpaque)
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
		csvDateFormat = "2006/1/2-07:00"
	)

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	resp, err := http.Get(csvUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
		date, err := time.Parse(csvDateFormat, v[0]+"-00:00")
		if err != nil {
			log.Warnf("info: %s is not parsable date. skipped...", v[0])
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
