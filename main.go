package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

func getFormatDuation(nanoseconds int64) (int64, int64, int64) {
	hours := nanoseconds / int64(time.Hour)
	nanoseconds = nanoseconds - hours*int64(time.Hour)
	minutes := nanoseconds / int64(time.Minute)
	nanoseconds = nanoseconds - minutes*int64(time.Minute)
	seconds := nanoseconds / int64(time.Second)
	return hours, minutes, seconds
}

func main() {
	app := cli.NewApp()
	app.Name = "offwork"
	app.Usage = "啥时候下班呀"
	app.Action = func(c *cli.Context) error {
		date := time.Now().Local()
		day := date.Weekday().String()
		targetTime := time.Date(date.Year(), date.Month(), date.Day(), 19, 0, 0, 0, date.Location())
		if day == "Saturday" || day == "Sunday" {
			fmt.Println("已经是周末了谢谢")
		} else {
			// TODO
			duration := targetTime.Sub(date)
			dn := duration.Nanoseconds()
			if dn <= 0 {
				fmt.Println("明天再来上班吧，下班了")
			} else {
				h, m, s := getFormatDuation(dn)
				fmt.Printf("还有 %d:%d:%d 下班\n", h, m, s)
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
