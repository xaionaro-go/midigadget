package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type pathValuer interface {
	PathValue(r *http.Request, name string) string
}

type Controller struct {
	midi       drivers.Out
	pathValuer pathValuer
}

func New(
	midi drivers.Out,
	pathValuer pathValuer,
) *Controller {
	return &Controller{
		midi:       midi,
		pathValuer: pathValuer,
	}
}

func httpError(
	resp http.ResponseWriter,
	statusCode int,
	format string,
	args ...any,
) {
	errorMessage := fmt.Sprintf(format, args...)
	resp.WriteHeader(statusCode)
	resp.Write([]byte(errorMessage))
}

func (ctrl *Controller) HTTPHandlerNoteOn(resp http.ResponseWriter, req *http.Request) {
	channelIDStr := ctrl.pathValuer.PathValue(req, "channelID")
	noteIDStr := ctrl.pathValuer.PathValue(req, "noteID")
	velocityStr := ctrl.pathValuer.PathValue(req, "velocity")

	channelID, err := strconv.ParseUint(channelIDStr, 10, 8)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to parse channel ID '%s': %v", channelIDStr, err)
		return
	}

	noteID, err := strconv.ParseUint(noteIDStr, 10, 8)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to parse note ID '%s': %v", noteIDStr, err)
		return
	}

	velocity, err := strconv.ParseUint(velocityStr, 10, 8)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to parse velocity '%s': %v", velocityStr, err)
		return
	}

	midiMsg := midi.NoteOn(uint8(channelID), uint8(noteID), uint8(velocity))
	err = ctrl.midi.Send(midiMsg)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to send MIDI message '%X': %v", midiMsg, err)
		return
	}
}

func (ctrl *Controller) HTTPHandlerNoteOff(resp http.ResponseWriter, req *http.Request) {
	channelIDStr := ctrl.pathValuer.PathValue(req, "channelID")
	noteIDStr := ctrl.pathValuer.PathValue(req, "noteID")

	channelID, err := strconv.ParseUint(channelIDStr, 10, 8)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to parse channel ID '%s': %v", channelIDStr, err)
		return
	}

	noteID, err := strconv.ParseUint(noteIDStr, 10, 8)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to parse note ID '%s': %v", noteIDStr, err)
		return
	}

	midiMsg := midi.NoteOff(uint8(channelID), uint8(noteID))
	err = ctrl.midi.Send(midiMsg)
	if err != nil {
		httpError(resp, http.StatusBadRequest, "unable to send MIDI message '%X': %v", midiMsg, err)
		return
	}
}
