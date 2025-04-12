package audio

import (
	"log"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"

	"github.com/tarasowski/autospeech/pkg/config"
)

// Recorder manages audio recording from a microphone
type Recorder struct {
	state      *config.AppState
	cfg        *config.AppConfig
	callback   func([]byte)
	stopChan   chan struct{}
}

// NewRecorder creates a new audio recorder
func NewRecorder(state *config.AppState, cfg *config.AppConfig) *Recorder {
	return &Recorder{
		state:    state,
		cfg:      cfg,
		stopChan: make(chan struct{}),
	}
}

// StartRecording begins audio recording with the given data callback
func (r *Recorder) StartRecording(dataCallback func([]byte)) error {
	r.callback = dataCallback
	r.state.ResetAudioBuffer()
	r.state.SetRecording(true)
	
	log.Println("Starting audio recording...")
	
	// Create an audio input stream
	inputStream, err := portaudio.OpenDefaultStream(
		config.Channels, 
		0, 
		float64(config.SampleRate), 
		config.FramesPerBuffer,
		r.audioInputCallback,
	)
	
	if err != nil {
		log.Printf("Failed to open audio stream: %v", err)
		r.state.SetRecording(false)
		return err
	}
	defer inputStream.Close()

	if err := inputStream.Start(); err != nil {
		log.Printf("Failed to start audio stream: %v", err)
		r.state.SetRecording(false)
		return err
	}

	log.Println("Audio recording started successfully")
	
	// Keep recording until stopped
	for r.state.IsRecording() {
		select {
		case <-r.stopChan:
			log.Println("Recording stopped by request")
			r.state.SetRecording(false)
		case <-time.After(100 * time.Millisecond):
			// Continue recording
		}
	}

	log.Println("Stopping audio recording...")
	inputStream.Stop()
	return nil
}

// StopRecording stops the current recording session
func (r *Recorder) StopRecording() {
	if r.state.IsRecording() {
		r.stopChan <- struct{}{}
	}
}

// audioInputCallback processes incoming audio data
func (r *Recorder) audioInputCallback(in []int16) {
	// Convert audio samples to bytes and store in buffer
	buf := make([]byte, len(in)*2)
	for i, sample := range in {
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	
	// Store in global buffer
	r.state.WriteToAudioBuffer(buf)
	
	// Call the data callback if provided
	if r.callback != nil {
		r.callback(buf)
	}
}

// CreateTempDir creates a temporary directory for audio processing
func CreateTempDir(prefix string) (string, error) {
	tmpDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		log.Printf("Failed to create temp directory: %v", err)
		return "", err
	}
	return tmpDir, nil
}