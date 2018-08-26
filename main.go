package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
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

func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("$GOPATH/src/github.com/csvwolf/offwork")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		viper.SetDefault("time", "19:00")
		f, _ := os.Create(os.Getenv("GOPATH") + "/src/github.com/csvwolf/offwork/config.json")
		defer f.Close()
		fmt.Println("初始化中……")
		viper.WriteConfig()
		viper.ReadInConfig()
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "offwork"
	app.Usage = "啥时候下班呀"

	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				prompt := promptui.Prompt{
					Label: "下班时间 >>>> ",
				}

				result, _ := prompt.Run()

				viper.Set("time", result)
				viper.WriteConfig()
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		date := time.Now().Local()
		day := date.Weekday().String()
		timing := strings.Split(viper.GetString("time"), ":")

		if len(timing) != 2 {
			fmt.Println("请按照格式：hour:minute 输出，比如 19:00")
		}

		hour, _ := strconv.Atoi(timing[0])
		minute, _ := strconv.Atoi(timing[1])

		targetTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
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
