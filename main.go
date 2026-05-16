package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	foodCapacity float32 = 1.0
	energyPerDay float32 = -0.5
	hungerPerDay float32 = 0.5
	agePerDay    int     = 1
	maxEnergy    float32 = 25.0
	maxHunger    float32 = -10.0

	saveFile = "save.json"
)

var (
	renameFlag = flag.String("rename", "", "name your tamagochi")
	feedFlag   = flag.Bool("feed", false, "feed your tamagochi")
	sleepFlag  = flag.Bool("sleep", false, "put your tamagochi to sleep")
	statusFlag = flag.Bool("status", false, "show your tamagochi's status")
)

type tamagochi struct {
	Name       string    `json:"name"`
	Hunger     float32   `json:"hunger"`
	Energy     float32   `json:"energy"`
	Age        int       `json:"age"`
	SeenAt     time.Time `json:"seenAt"`
	IsSleeping bool      `json:"IsSleeping"`
}

func (t *tamagochi) feed() {
	if t.IsSleeping {
		fmt.Printf("you can't feed %s while its sleeping\n", t.Name)
		return
	}

	if t.Hunger+foodCapacity >= maxHunger && t.Hunger+foodCapacity <= -maxHunger {
		t.Hunger += foodCapacity
		fmt.Printf("%s has been fed\n", t.Name)
		save(*t)
	} else {
		fmt.Printf("%s is full!\n", t.Name)
	}
}

func (t *tamagochi) rename(name string) {
	lastName := t.Name
	t.Name = name
	fmt.Printf("you have renamed %s to %s\n", lastName, t.Name)
	save(*t)
}

func (t *tamagochi) sleep() {
	if t.IsSleeping {
		t.IsSleeping = false
		fmt.Printf("%s is not longer sleeping\n", t.Name)

		daysPassed := int(time.Since(t.SeenAt).Hours() / 24)

		if daysPassed > 0 {
			var energy float32 = float32(daysPassed) * energyPerDay
			var hunger float32 = float32(daysPassed) * hungerPerDay
			if t.Energy+energy <= maxEnergy && t.Energy+energy >= -maxEnergy {
				t.Energy += energy
			}
			if t.Hunger+hunger >= maxHunger && t.Hunger+hunger <= -maxHunger {
				t.Hunger += hunger
			}
			t.Age += daysPassed
		}

		save(*t)
		return
	}

	t.IsSleeping = true
	fmt.Printf("%s is sleeping\n", t.Name)

	save(*t)
}

func (t tamagochi) status() {
	fmt.Printf("name: %s\nhunger: %v\nenergy: %v\nage: %v\nis sleeping: %v\n", t.Name, t.Hunger, t.Energy, t.Age, t.IsSleeping)
}

func load() (tamagochi, error) {
	data, err := os.ReadFile(saveFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("you don't have a tamagochi yet")
			pet := tamagochi{
				Name:       "",
				Hunger:     5.0,
				Energy:     5.0,
				Age:        0,
				SeenAt:     time.Now(),
				IsSleeping: false,
			}
			pet.sleep()
			return pet, nil
		}
		return tamagochi{}, fmt.Errorf("failed to read file %w", err)
	}

	var pet tamagochi
	if err := json.Unmarshal(data, &pet); err != nil {
		return tamagochi{}, fmt.Errorf("failed to decode: %w", err)
	}
	return pet, nil
}

func save(pet tamagochi) error {
	data, err := json.MarshalIndent(pet, "", "	")
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}
	return os.WriteFile(saveFile, data, 0644)
}

func main() {
	flag.Parse()

	pet, err := load()
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
		pet.status()
	default:
		flag.Usage()
	}
}
