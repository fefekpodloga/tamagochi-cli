package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	renameFlag      = flag.String("rename", "", "name your tamagotchi")
	feedFlag        = flag.Bool("feed", false, "feed your tamagotchi")
	sleepFlag       = flag.Bool("sleep", false, "put your tamagotchi to sleep")
	statusFlag      = flag.Bool("status", false, "show your tamagotchi's status")
	debugStatusFlag = flag.Bool("debug", false, "status in debug mode")
)

func main() {
	flag.Parse()

	var pet Tamagotchi
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
