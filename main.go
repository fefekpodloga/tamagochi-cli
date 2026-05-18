package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	renameFlag      = flag.String("rename", "", "name your tamagochi")
	feedFlag        = flag.Bool("feed", false, "feed your tamagochi")
	sleepFlag       = flag.Bool("sleep", false, "put your tamagochi to sleep")
	statusFlag      = flag.Bool("status", false, "show your tamagochi's status")
	debugStatusFlag = flag.Bool("debug", false, "status in debug mode")
)

func main() {
	flag.Parse()

	var pet Tamagochi
	err := pet.load()
	if err != nil {
		fmt.Println("error loading data: ", err)
		os.Exit(1)
	}

	switch {
	case *renameFlag != "":
		pet.rename(*renameFlag)
	case *feedFlag:
		pet.feed()
	case *sleepFlag:
		pet.sleep()
	case *statusFlag:
		pet.status(false)
	case *debugStatusFlag:
		pet.status(true)
	default:
		flag.Usage()
	}
}
