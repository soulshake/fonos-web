package controllers

import (
	"log"
)

type Volumes []Volume

type Volume struct {
	Name string `json:"name"`
}

func GetVolumes() (volumes Volumes, err error) {
	sinks, err := GetSinks()
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range sinks {
		volumes = append(volumes, Volume{Name: s.Name})
	}
	return volumes, nil
}
