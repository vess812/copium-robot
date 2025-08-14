package bot

import (
	"bytes"
	"fmt"
	"os/exec"

	"copium-bot/internal/domain"
)

type Model interface {
	Transcribe(input []byte) (string, error)
}

type Transcriber struct {
	model Model
}

func NewTranscriber(model Model) *Transcriber {
	return &Transcriber{
		model: model,
	}
}

func (v *Transcriber) Process(r domain.Request) (domain.Response, error) {
	if r.Message.Voice == nil && r.Message.VideoNote == nil {
		return domain.Response{}, fmt.Errorf("voice and video note are empty")
	}

	var data []byte
	switch {
	case r.Message.Voice != nil:
		data = r.Message.Voice
	case r.Message.VideoNote != nil:
		data = r.Message.VideoNote
	}

	wav, err := ffmpegConvertToWav(data)
	if err != nil {
		return domain.Response{}, fmt.Errorf("convert to wav: %w", err)
	}

	result, err := v.model.Transcribe(wav)
	if err != nil {
		return domain.Response{}, fmt.Errorf("transcribe: %w", err)
	}

	if len(result) == 0 {
		return domain.Response{}, fmt.Errorf("empty result")
	}

	return domain.Response{
		ChatID:  r.Message.ChatID,
		ReplyTo: r.Message.ID,
		Text:    result,
	}, nil
}

func ffmpegConvertToWav(data []byte) ([]byte, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0", // Read data from stdin
		"-vn",       // Disable video processing
		"-f", "wav", // Output format = WAV
		"-ar", "16000", // Sample rate (16kHz)
		"-ac", "1", // Mono channel
		"-c:a", "pcm_s16le", // 16-bit PCM encoding
		"pipe:1", // Write WAV to stdout
	)

	cmd.Stdin = bytes.NewReader(data)

	var wavBuf bytes.Buffer
	cmd.Stdout = &wavBuf

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return wavBuf.Bytes(), nil
}
