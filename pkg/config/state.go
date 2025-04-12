package config

import (
	"bytes"
	"sync"
	"time"
)

// AppState manages the application state with thread safety
type AppState struct {
	mu                  sync.RWMutex
	isRecording         bool
	transcribedText     string
	audioBuffer         bytes.Buffer
	partialTranscription string
	partialUpdateTime    time.Time
	// No longer used for voice commands
}

// NewAppState creates a new application state
func NewAppState(cfg *AppConfig) *AppState {
	return &AppState{
		isRecording: false,
	}
}

// IsRecording returns the current recording state
func (s *AppState) IsRecording() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRecording
}

// SetRecording sets the recording state
func (s *AppState) SetRecording(recording bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isRecording = recording
}

// GetTranscribedText returns the transcribed text
func (s *AppState) GetTranscribedText() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transcribedText
}

// SetTranscribedText sets the transcribed text
func (s *AppState) SetTranscribedText(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transcribedText = text
}

// GetAudioBuffer returns a copy of the current audio buffer
func (s *AppState) GetAudioBuffer() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.audioBuffer.Bytes()
}

// WriteToAudioBuffer writes data to the audio buffer
func (s *AppState) WriteToAudioBuffer(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audioBuffer.Write(data)
}

// ResetAudioBuffer clears the audio buffer
func (s *AppState) ResetAudioBuffer() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audioBuffer.Reset()
}

// GetPartialTranscription returns the current partial transcription
func (s *AppState) GetPartialTranscription() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.partialTranscription
}

// SetPartialTranscription sets the partial transcription
func (s *AppState) SetPartialTranscription(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.partialTranscription = text
}

// UpdatePartialTranscriptionTime updates the time of the last partial transcription
func (s *AppState) UpdatePartialTranscriptionTime() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.partialUpdateTime = time.Now()
}

// ShouldUpdatePartialTranscription checks if enough time has passed to update the transcription
func (s *AppState) ShouldUpdatePartialTranscription() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Use a very short interval (200ms) for much more responsive real-time updates
	return time.Since(s.partialUpdateTime) > 200*time.Millisecond
}

// Voice command functionality has been removed