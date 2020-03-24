package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/romchi/zabbix.mdadm/mdadm"
)

type devices struct {
	MDName string `json:"{#MD.NAME}"`
}

func main() {
	statsCommand := flag.NewFlagSet("stats", flag.ExitOnError)

	statsDeviceName := statsCommand.String("name", "", `Device "name" to get stats (Required)`)

	if len(os.Args) < 2 {
		fmt.Println("[discovery, stats] - required one command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "discovery":
		discoveryMD()
	case "stats":
		statsCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if statsCommand.Parsed() {
		if len(*statsDeviceName) < 1 {
			statsCommand.PrintDefaults()
			os.Exit(0)
		}
		statsRaidController(*statsDeviceName)
	}
}

func statsRaidController(rName string) {
	raids, err := mdadm.MdadmDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range raids {
		if rName == c.Name {
			i, err := c.RaidStats()
			if err != nil {
				log.Fatal(err)
			}
			m, err := c.MDStats()
			if err != nil {
				log.Fatal(err)
			}
			i.MD = m
			ji, _ := json.Marshal(i)
			fmt.Print(string(ji))
		}
	}
}

func discoveryMD() {
	result := []devices{}

	mdDevices, err := mdadm.MdadmDevicesList()
	if err != nil {
		log.Fatal(err)
	}

	for _, md := range mdDevices {
		result = append(result, devices{MDName: md})
	}

	r, _ := json.Marshal(result)

	fmt.Println(string(r))
}
