package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jba/muxpatterns"
	"github.com/xaionaro-go/midigadget/pkg/controller"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	defer midi.CloseDriver()

	drv, err := rtmididrv.New()
	if err != nil {
		log.Panic(err)
	}

	listenAddr := flag.String("listen-addr", ":3333", "")
	midiPort := flag.String("midi-port", "", "MIDI port name")
	flag.Parse()

	_ = drv
	out, err := midi.FindOutPort(*midiPort)
	if err != nil {
		log.Panicf("unable to find MIDI port '%s': %v. Available ports: %s", *midiPort, err, midi.GetOutPorts().String())
		return
	}

	err = out.Open()
	if err != nil {
		log.Panicf("unable to open MIDI port '%s': %v", *midiPort, err)
		return
	}

	mux := muxpatterns.NewServeMux()
	ctrl := controller.New(out, mux)

	mux.HandleFunc("/noteOn/{channelID}/{noteID}/{velocity}", ctrl.HTTPHandlerNoteOn)
	mux.HandleFunc("/noteOff/{channelID}/{noteID}", ctrl.HTTPHandlerNoteOff)

	log.Println("Listening on", *listenAddr)
	err = http.ListenAndServe(*listenAddr, mux)
	log.Panic(err)
}
