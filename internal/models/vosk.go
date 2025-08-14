package models

import (
	"encoding/json"
	"fmt"

	vosk "github.com/alphacep/vosk-api/go"
)

type Opts struct {
	ModelPath string
}

type VoskModel struct {
	opts  Opts
	model *vosk.VoskModel
}

func NewVoskModel(opts Opts) (*VoskModel, error) {
	vosk.SetLogLevel(0)
	model, err := vosk.NewModel(opts.ModelPath)
	if err != nil {
		return nil, fmt.Errorf("vosk model init: %w", err)
	}

	return &VoskModel{
		opts:  opts,
		model: model,
	}, nil
}

const (
	sampleRate = 16000.0
)

type recognitionResult struct {
	Text string `json:"text"`
}

func (t *VoskModel) Transcribe(input []byte) (string, error) {
	rec, err := vosk.NewRecognizer(t.model, sampleRate)
	if err != nil {
		return "", fmt.Errorf("new recognizer: %w", err)
	}
	defer rec.Free()

	rec.AcceptWaveform(input)

	var r recognitionResult
	err = json.Unmarshal([]byte(rec.FinalResult()), &r)
	if err != nil {
		return "", fmt.Errorf("unmarshal recognition result: %w", err)
	}

	return r.Text, nil
}
