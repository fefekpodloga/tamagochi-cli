package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	foodCapacity  float32 = 0.5
	energyPerDay  float32 = 0.5
	hungerPerDay  float32 = 0.2
	agePerDay     int     = 1
	maxEnergy     float32 = 25
	maxHunger     float32 = 10
	defaultName   string  = "noname"
	defaultHunger float32 = 5.0
	defaultEnergy float32 = 10.0
	defaultAge    int     = 0

	saveFileName string = "save.json"
)

type Tamagochi struct {
	Name       string    `json:"name"`
	Hunger     float32   `json:"hunger"`
	Energy     float32   `json:"energy"`
	Age        int       `json:"age"`
	SeenAt     time.Time `json:"seenAt"`
	IsSleeping bool      `json:"IsSleeping"`
}

func (t *Tamagochi) timeUpdate() {
	daysPassed := int(time.Since(t.SeenAt).Hours() / 24)

	if daysPassed == 0 {
		return
	}

	energy := energyPerDay * float32(daysPassed)
	if t.IsSleeping {
		if t.Energy+energy < maxEnergy {
			t.Energy += energy
		}
	} else {
		if t.Energy-energy > -maxEnergy {
			t.Energy -= energy
		}
	}

	hunger := hungerPerDay * float32(daysPassed)
	if t.Hunger+hunger < maxHunger {
		t.Hunger += hunger
	}

	t.save(true)
}

func (t *Tamagochi) sleep() {
	if t.IsSleeping {
		t.IsSleeping = false
		fmt.Printf("%s is not longer sleeping\n", t.Name)
		t.save(true)
		return
	}

	t.IsSleeping = true
	fmt.Printf("%s has just gone to sleep\n", t.Name)
	t.save(true)
}

func (t *Tamagochi) rename(name string) {
	lastName := t.Name
	t.Name = name
	fmt.Printf("%s has been renamed to %s\n", lastName, t.Name)
	t.save(false)
}

func (t *Tamagochi) feed() {
	if t.IsSleeping {
		fmt.Printf("you can't feed %s while sleeping\n", t.Name)
		return
	}

	if t.Hunger-foodCapacity > -maxHunger {
		t.Hunger -= foodCapacity
		fmt.Printf("you have fed %s\n", t.Name)
	} else {
		fmt.Printf("%s is full!\n", t.Name)
	}

	t.save(false)
}

func (t Tamagochi) status(debug bool) {
	if !debug {
		fmt.Printf("name: %s\nhunger: %v\nenergy: %v\nage: %v\nis sleeping: %v\n", t.Name, t.Hunger, t.Energy, t.Age, t.IsSleeping)
	} else {
		fmt.Printf("name: %s\nhunger: %v\nenergy: %v\nage: %v\nis sleeping: %v\nseen at: %v\n", t.Name, t.Hunger, t.Energy, t.Age, t.IsSleeping, t.SeenAt)
	}
}

func (t *Tamagochi) save(saveTime bool) error {
	if saveTime {
		t.SeenAt = time.Now()
	}

	data, err := json.MarshalIndent(t, "", "	")
	if err != nil {
		return fmt.Errorf("failed to encode: %w\n", err)
	}
	return os.WriteFile(saveFileName, data, 0644)
}

func (t *Tamagochi) load() error {
	data, err := os.ReadFile(saveFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("you don't have a tamagochi yet")

			t.Name = defaultName
			t.Hunger = defaultHunger
			t.Energy = defaultEnergy
			t.Age = defaultAge
			t.SeenAt = time.Now()
			t.IsSleeping = false

			t.save(false) // time is being set above
			return nil
		}
		return fmt.Errorf("failed to red file %w", err)
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("failed to decode %w", err)
	}
	return nil
}
