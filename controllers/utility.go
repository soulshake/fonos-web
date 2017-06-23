package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	//"github.com/convox/console/httperr"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	log.Print("in Dashboard")
	data := Data{
		Body: "hello world",
	}

	RenderHTML(w, r, "dashboard", "index", data)
}

func ShowSinks(w http.ResponseWriter, r *http.Request) {
	log.Print("in ShowSinks")
	data := Data{}

	sinks, err := GetSinks()
	if err != nil {
		data.Error = err
		log.Print(err)
	}

	data.Sinks = sinks

	// Serve HTML or JSON
	switch accept := HTMLorJSON(r); accept {
	case "json":
		json, err := json.Marshal(sinks)
		if err != nil {
			log.Print(err)
		}
		w.Write([]byte(string(json)))
	default:
		log.Print("Need to serve HTML")
		RenderHTML(w, r, "dashboard", "sinks", data)
	}
}

func ShowPaWebControl(w http.ResponseWriter, r *http.Request) {
	data := Data{}
	sinks, err := GetSinks()
	if err != nil {
		data.Error = err
		log.Print(err)
	}
	data.Sinks = sinks
	RenderHTML(w, r, "dashboard", "pawebcontrol", data)
}

func ShowVolumes(w http.ResponseWriter, r *http.Request) {
	data := Data{}

	log.Print("in ShowVolumes")
	sinks, err := GetSinks()
	if err != nil {
		data.Error = err
		log.Print(err)
	}
	data.Sinks = sinks

	// Serve HTML or JSON
	switch accept := HTMLorJSON(r); accept {
	case "json":
		json, err := json.Marshal(sinks)
		if err != nil {
			log.Print(err)
			data.Error = err
		}
		w.Write([]byte(string(json)))
	default:
		RenderHTML(w, r, "dashboard", "volumes", data)
	}
}

func ShowSinkInputs(w http.ResponseWriter, r *http.Request) {
	data := Data{}

	inputs, err := GetSinkInputs()
	if err != nil {
		data.Error = err
	}
	data.SinkInputs = inputs

	// Serve HTML or JSON
	switch accept := HTMLorJSON(r); accept {
	case "json":
		json, err := json.Marshal(inputs)
		if err != nil {
			log.Print(err)
		}
		w.Write([]byte(string(json)))
	default:
		RenderHTML(w, r, "dashboard", "sink-inputs", data)
	}
}

func HTMLorJSON(r *http.Request) string {
	for k, v := range r.Header {
		if k == "Accept" {
			if strings.Join(v, "") == "application/json" {
				return "json"
			}
		}
	}
	return "html"
}

func ShowPaCmd(w http.ResponseWriter, r *http.Request) {
	data := Data{}
	i, o, err := get_cmd(w, r)
	if err != nil {
		data.Error = err
		log.Print(err)
	} else {
		data.Input = i
		data.Output = o
	}

	RenderHTML(w, r, "dashboard", "pacmd", data)
}

func getSinkById(str_id string) (*Sink, error) {
	sinks, err := GetSinks()
	if err != nil {
		log.Print(err)
	}
	for _, s := range sinks {
		id, err := strconv.Atoi(str_id)
		if err != nil {
			log.Print(err)
		}
		if s.Index == id {
			return &s, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("couldn't find sink with ID %s", str_id))
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	data := Data{}
	sinks, err := GetSinks()
	if err != nil {
		log.Panic(err)
		data.Error = err
	}
	data.Sinks = sinks

	str_id := r.FormValue("id")
	if str_id != "" {
		str_id = strings.TrimPrefix(str_id, "s")

		muted := r.FormValue("muted")
		str_volume := r.FormValue("volume")

		s, err := getSinkById(str_id)
		if err != nil {
			log.Print(err)
		}

		if str_volume != "" {
			nv, err := strconv.ParseFloat(str_volume, 64)
			if err != nil {
				log.Print(err)
			}
			s.SetNewVolume(nv)
		}
		if muted != "" {
			log.Print(fmt.Printf("Got request to change ID %s from volume %d to volume %s\n", str_id, s.CurrentVolumeStep, str_volume))
			s.SetNewVolume(float64(s.CurrentVolumeStep))

			if muted == "1" {
				log.Print("we should be muted")
				s.SetMute(true)
			} else {
				log.Print("we should NOT be muted")
				s.SetMute(false)
			}
		}
	}
	RenderHTML(w, r, "dashboard", "pawebcontrol", data)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/assets/img/favicon.ico", http.StatusSeeOther)
}

/*
func errorHandler(err error, w http.ResponseWriter, r *http.Request) {
	switch terr := err.(type) {
	//case nil:
	//http.Redirect(w, r, "/grid/err500", 303)
	case *httperr.Error:
		if terr.Server() {
			reqId := r.Header.Get("X-Request-Id")
			fmt.Printf("request_id=%s error=%q\n", reqId, terr.Error())
			for _, frame := range terr.Stack {
				fmt.Printf("ns=grid request_id=%s filename=%s method=%s line=%d\n",
					reqId,
					frame.Filename,
					frame.Method,
					frame.Line)
			}
			Unknown(w, r)
		} else {
			message := []byte(err.Error())
			SetFlash(w, "error", message)
		}
		return
	case error:
		reqId := r.Header.Get("X-Request-Id")
		fmt.Printf("request_id=%s error=%q\n", reqId, terr.Error())
		Unknown(w, r)
	}
}
*/

var list = map[string]string{
	"+":          "amixer set Master playback 5%+",
	"-":          "amixer set Master playback 5%-",
	"help":       "pacmd --help",
	"list-cards": "pacmd --list-cards",
	"list-sinks": "pacmd --list-sinks",
	"t start":    "transmission-remote -tall --start",
}

func get_cmd(w http.ResponseWriter, r *http.Request) (input string, output string, err error) {

	cmd := strings.TrimSpace(r.FormValue("cmd"))
	result, err := run(cmd)
	if err != nil {
		return "", "", err
	}
	return cmd, result, nil
}

func run(cmd string) (string, error) {
	if len(cmd) != 0 {
		if value, ok := list[cmd]; ok && value != "" {
			cmd = list[cmd]
		} else {
			cmd = "pacmd " + cmd
		}
		result, err := execute(cmd)
		if err != nil {
			return "", err
		}
		return result, nil
	}
	return "", errors.New("No command provided!")
}

func execute(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]
	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(out), nil
}
