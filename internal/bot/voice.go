package bot

import (
	"bytes"
	"fmt"
	"os/exec"

	"copium-bot/internal/models"
)

type Transcriber interface {
	Transcribe(input []byte) (string, error)
}

type Voice struct {
	transcriber Transcriber
}

func NewVoice(transcriber Transcriber) *Voice {
	return &Voice{
		transcriber: transcriber,
	}
}

func (v *Voice) Process(r models.BotRequest) (models.BotResponse, error) {
	if r.Message.Voice == nil {
		return models.BotResponse{}, fmt.Errorf("voice message is empty")
	}

	wav, err := convertToWav(r.Message.Voice)
	if err != nil {
		return models.BotResponse{}, fmt.Errorf("convert to wav: %w", err)
	}

	result, err := v.transcriber.Transcribe(wav)
	if err != nil {
		return models.BotResponse{}, fmt.Errorf("transcribe: %w", err)
	}

	if len(result) == 0 {
		return models.BotResponse{}, fmt.Errorf("empty result")
	}

	return models.BotResponse{
		ChatID:  r.Message.ChatID,
		ReplyTo: r.Message.ID,
		Text:    result,
	}, nil
}

func convertToWav(ogg []byte) ([]byte, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0", // Read OGG from stdin
		"-f", "wav", // Output format = WAV
		"-ar", "16000", // Sample rate (16kHz)
		"-ac", "1", // Mono channel
		"-c:a", "pcm_s16le", // 16-bit PCM encoding
		"pipe:1", // Write WAV to stdout
	)

	cmd.Stdin = bytes.NewReader(ogg)

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
