package controllers

import (
	"log"

	"github.com/auroralaboratories/pulse"
	reflections "gopkg.in/oleiade/reflections.v1"
)

/*
type Sink struct {
	Name        string  `json:"name"` // the name reported by PulseAudio
	Index       uint8   `json:"index"`
	Description string  `json:"description"` // device.description property
	Volume      uint8   `json:"volume"`
	Muted       bool    `json:"muted"`
	Volumes     []uint8 `json:"volumes"`
}
*/

type Sink struct {
	//Client             *Client `json:"Client"`  // I don't know what *Client is
	CardIndex          int                    `json:"CardIndex"`
	Description        string                 `json:"Description"`
	DriverName         string                 `json:"DriverName"`
	Index              int                    `json:"Index"`
	ModuleIndex        int                    `json:"ModuleIndex"`
	MonitorSourceIndex int                    `json:"MonitorSourceIndex"`
	MonitorSourceName  string                 `json:"MonitorSourceName"`
	Muted              bool                   `json:"Muted"`
	Name               string                 `json:"Name"`
	NumFormats         int                    `json:"NumFormats"`
	NumPorts           int                    `json:"NumPorts"`
	NumVolumeSteps     int                    `json:"NumVolumeSteps"`
	State              int                    `json:"State"` // in pulse, SinkState is its own type
	Channels           int                    `json:"Channels"`
	CurrentVolumeStep  int                    `json:"CurrentVolumeStep"`
	VolumeFactor       float64                `json:"VolumeFactor"`
	properties         map[string]interface{} `json:"properties"`
}

type Sinks []Sink

type SinkInput pulse.Source
type SinkInputs []pulse.Source

/*
[{"BaseVolumeStep":0,
"CardIndex":0,
"Channels":2,
"Client":{"ID":"ce8e1abe-b365-4c17-841a-1dda4843a763","Name":"test-client-get-sinks","Server":"","OperationTimeout":5000000000},
"CurrentVolumeStep":65536,
"Description":"Monitor of Built-in Audio Analog Stereo",
"DriverName":"module-alsa-card.c",
"Index":0,
"ModuleIndex":6,
"MonitorOfSinkIndex":0,
"MonitorOfSinkName":"alsa_output.pci-0000_00_1b.0.analog-stereo",
"Muted":true,
"Name":"alsa_output.pci-0000_00_1b.0.analog-stereo.monitor",
"NumFormats":1,
"NumPorts":0,
"NumVolumeSteps":65537,
"State":2,
"VolumeFactor":1}]
*/

/*
type SinkInput struct {
	Name        string  `json:"name"` // the name reported by PulseAudio
	Index       uint8   `json:"index"`
	Description string  `json:"description"` // media.name property
	Volume      uint8   `json:"volume"`
	Muted       bool    `json:"muted"`
	Volumes     []uint8 `json:"volumes"`
}
*/

func marshalSinks(psinks []pulse.Sink) (osinks Sinks, err error) {
	for _, s := range psinks {
		x, err := marshalSink(s)
		if err != nil {
			log.Panic(err)
		}
		osinks = append(osinks, x)
	}
	return osinks, nil
}

func marshalSink(ps pulse.Sink) (os Sink, err error) {
	fields, _ := reflections.Fields(ps)
	for _, field := range fields {
		value, err := reflections.GetField(&ps, field)
		if err != nil {
			log.Panic(err)
		}
		_ = reflections.SetField(&os, field, value)
	}
	return os, nil

}

func GetSinks() ([]Sink, error) {
	//log.Print("in GetSinks")
	client, err := pulse.NewClient(`test-client-get-sinks`)
	if err != nil {
		log.Fatal(err)
	}
	psinks, err := client.GetSinks()
	if err != nil {
		log.Fatal(err)
	}
	sinks, err := marshalSinks(psinks)
	return sinks, nil
}

func GetSinkInputs() ([]pulse.Source, error) {
	log.Print("in GetSinkInputs")
	client, err := pulse.NewClient(`test-client-get-sinks`)
	if err != nil {
		log.Fatal(err)
	}
	inputs, err := client.GetSources()
	if err != nil {
		return nil, err
	}
	log.Print(inputs)
	return inputs, nil

}

func (s Sink) ToggleMute() error {
	client, err := pulse.NewClient(`test-client-get-sinks`)
	if err != nil {
		log.Fatal(err)
	}
	pSinks, err := client.GetSinks()
	if err != nil {
		log.Fatal(err)
	}

	for _, ps := range pSinks {
		if ps.Description == s.Description {
			err := ps.ToggleMute()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func (s Sink) SetMute(muted bool) error {
	client, err := pulse.NewClient(`test-client-get-sinks`)
	if err != nil {
		log.Fatal(err)
	}
	pSinks, err := client.GetSinks()
	if err != nil {
		log.Fatal(err)
	}

	for _, ps := range pSinks {
		if ps.Description == s.Description {
			if muted {
				err = ps.Mute()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				err = ps.Unmute()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return nil
}

// SetVolumeFactor sets the volume of all channels of this sink to a factor of the maximum volume (0.0 <= v <= 1.0).
func (s Sink) SetNewVolume(volume float64) error {
	client, err := pulse.NewClient(`test-client-get-sinks`)
	if err != nil {
		log.Fatal(err)
	}
	pSinks, err := client.GetSinks()
	if err != nil {
		log.Fatal(err)
	}

	for _, ps := range pSinks {
		if ps.Description == s.Description {
			//curVolSteps := float64(ps.CurrentVolumeStep)
			numVolSteps := float64(ps.NumVolumeSteps)
			volFactor := volume / numVolSteps
			log.Print("Setting volume factor to ", volFactor)
			err = ps.SetVolume(volFactor)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}
