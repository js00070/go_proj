package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	_ "github.com/go-sql-driver/mysql"
)

type record struct {
	rank     int
	city     string
	AQI      int
	airlevel string //空气质量
	source   string // 首要污染物
	pm25     int
	pm10     int
	CO       float64
	NO2      int
	O3_1h    int
	O3_8h    int
	SO2      int
	year     int
	month    int
	day      int
	hour     int
	//id       int
}

var dataTime string

func getTime() string {
	doc, err := goquery.NewDocument("http://pm25.in/rank")
	if err != nil {
		panic("err")
	}
	timestr := strings.SplitN(doc.Find("div.time").Text(), "：", 2)[1]
	timestr = strings.Split(timestr, "\n")[0]
	return timestr
}

func getData() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/airdata")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//fmt.Println("helloworld!")
	doc, err := goquery.NewDocument("http://pm25.in/rank")
	if err != nil {
		panic("err")
	}
	timestr := strings.SplitN(doc.Find("div.time").Text(), "：", 2)[1]
	timestr = strings.Split(timestr, "\n")[0]
	timestrlist := strings.Split(timestr, " ")
	date := strings.Split(timestrlist[0], "-")
	clock := strings.Split(timestrlist[1], ":")
	year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	day, _ := strconv.Atoi(date[2])
	hour, _ := strconv.Atoi(clock[0])
	var recordlist = make([]record, 0)
	doc.Find("div.table").Each(func(i int, s1 *goquery.Selection) {
		s1.Find("tbody").Each(func(j int, s2 *goquery.Selection) {
			s2.Find("tr").Each(func(k int, s3 *goquery.Selection) {
				templist := make([]string, 16)
				s3.Find("td").Each(func(l int, s4 *goquery.Selection) {
					templist[l] = s4.Text()
				})
				rank, err := strconv.Atoi(templist[0])
				AQI, err := strconv.Atoi(templist[2])
				pm25, err := strconv.Atoi(templist[5])
				pm10, err := strconv.Atoi(templist[6])
				CO, err := strconv.ParseFloat(templist[7], 64)
				NO2, err := strconv.Atoi(templist[8])
				O3_1h, err := strconv.Atoi(templist[9])
				O3_8h, err := strconv.Atoi(templist[10])
				SO2, err := strconv.Atoi(templist[11])
				if err != nil {
					panic("abc")
				}

				//db.Exec("set city 'utf8mb4';")
				//fmt.Println(templist[4])
				stmt, err := db.Prepare("INSERT INTO airtable(rank,city,AQI,airlevel,source,pm25,pm10,CO,NO2,O3_1h,O3_8h,SO2,year,month,day,hour) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
				if err != nil {
					fmt.Println("insert data failed:", err.Error())
					return
				}
				//fmt.Println(templist[1])
				res, err := stmt.Exec(rank, templist[1], AQI, templist[3], templist[4], pm25, pm10, CO, NO2, O3_1h, O3_8h, SO2, year, month, day, hour)
				if err != nil {
					panic(err)
				}
				id, err := res.LastInsertId()
				if err != nil {
					panic(err)
				}
				fmt.Println(id)
				recordlist = append(recordlist,
					record{rank, templist[1], AQI, templist[3], templist[4], pm25, pm10, CO, NO2, O3_1h, O3_8h, SO2, year, month, day, hour})
				//fmt.Println(recordlist[0])
			})
		})

	})

	//fmt.Println(html)
}

func main() {
	dataTime = getTime()
	for {
		fmt.Println(dataTime)
		time.Sleep(1 * time.Minute)
		now := getTime()
		if dataTime != now {
			getData()
			dataTime = now
		}
	}
}
