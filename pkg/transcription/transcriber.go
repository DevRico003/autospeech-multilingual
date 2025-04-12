package transcription

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tarasowski/autospeech/pkg/audio"
	"github.com/tarasowski/autospeech/pkg/config"
)

// Transcriber handles speech-to-text transcription
type Transcriber struct {
	cfg   *config.AppConfig
	state *config.AppState
}

// NewTranscriber creates a new transcription service
func NewTranscriber(cfg *config.AppConfig, state *config.AppState) *Transcriber {
	return &Transcriber{
		cfg:   cfg,
		state: state,
	}
}

// TranscribeAudio transcribes the current audio buffer
func (t *Transcriber) TranscribeAudio() (string, error) {
	audioData := t.state.GetAudioBuffer()
	if len(audioData) == 0 {
		log.Println("No audio data captured")
		return "", fmt.Errorf("no audio data captured")
	}

	log.Printf("Captured %d bytes of audio data", len(audioData))

	// Save audio to a temporary WAV file
	tmpDir, err := audio.CreateTempDir("speech-reco")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	wavFile := filepath.Join(tmpDir, "recording.wav")
	log.Printf("Saving audio to temporary WAV file: %s", wavFile)
	if err := audio.SaveAsWav(audioData, wavFile); err != nil {
		return "", err
	}

	// Try different transcription methods
	log.Println("Starting transcription...")
	
	// Try vosk first
	transcript, err := t.transcribeWithVosk(wavFile)
	if err == nil && transcript != "" {
		return transcript, nil
	}
	
	// Then try system commands
	transcript, err = t.transcribeWithSystemCommand(wavFile)
	if err == nil && transcript != "" {
		return transcript, nil
	}
	
	// Return a default message if all methods fail
	return "Speech recognition failed. Please install Vosk or other speech recognition software.", nil
}

// QuickTranscribe performs a fast transcription on partial audio
func (t *Transcriber) QuickTranscribe(audioData []byte) (string, error) {
	// Create a temporary WAV file
	tmpDir, err := audio.CreateTempDir("speech-reco-rt")
	if err != nil {
		log.Printf("Failed to create temp directory for real-time transcription: %v", err)
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	wavFile := filepath.Join(tmpDir, "partial.wav")
	if err := audio.SaveAsWav(audioData, wavFile); err != nil {
		log.Printf("Failed to save partial audio: %v", err)
		return "", err
	}

	// Add debug info about the audio file
	log.Printf("Created temporary WAV file for real-time transcription: %s (size: %d bytes)", 
		wavFile, len(audioData))

	// Try to transcribe with Vosk (fastest option)
	transcript, err := t.transcribeWithVosk(wavFile)
	if err == nil && transcript != "" {
		log.Printf("Vosk real-time transcription: '%s'", transcript)
		return transcript, nil
	} else if err != nil {
		log.Printf("Vosk transcription failed: %v", err)
	} else {
		log.Printf("Vosk returned empty transcription")
	}

	// Fallback to other methods if available
	transcript, err = t.transcribeWithSystemCommand(wavFile)
	if err == nil && transcript != "" {
		log.Printf("System command real-time transcription: '%s'", transcript)
		return transcript, nil
	} else if err != nil {
		log.Printf("System command transcription failed: %v", err)
	} else {
		log.Printf("System command returned empty transcription")
	}

	return "", fmt.Errorf("real-time transcription failed")
}

// transcribeWithVosk uses the Vosk speech recognition toolkit
func (t *Transcriber) transcribeWithVosk(wavFile string) (string, error) {
	// Try to find vosk-transcribe in various locations
	possiblePaths := []string{
		"./vosk-transcribe",                         // Current directory
		"../vosk-transcribe",                        // Parent directory
		"/usr/local/bin/vosk-transcribe",            // System-wide install
		"/usr/bin/vosk-transcribe",                  // Alternative system location
		filepath.Join(os.Getenv("HOME"), "vosk-transcribe"), // User's home directory
	}
	
	voskCmd := ""
	found := false
	
	// Check each possible path
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			voskCmd = path
			found = true
			log.Printf("Found vosk-transcribe at %s", voskCmd)
			break
		}
	}
	
	// If not found in specific locations, try PATH lookup
	if !found {
		var err error
		cmd := "vosk-transcribe"
		voskCmd, err = exec.LookPath(cmd)
		if err != nil {
			log.Printf("vosk-transcribe not found in PATH: %v", err)
			return "", fmt.Errorf("vosk-transcribe not found, run setup-vosk-small.sh to install")
		}
		log.Printf("Found vosk-transcribe in PATH at %s", voskCmd)
	}
	
	// Run vosk-transcribe with the WAV file
	log.Printf("Running Vosk transcription with: %s", voskCmd)
	cmd := exec.Command(voskCmd, wavFile)
	
	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("vosk transcription failed: %v, output: %s", err, string(output))
	}
	
	// Return the transcribed text
	transcription := strings.TrimSpace(string(output))
	if transcription == "" {
		return "", fmt.Errorf("no transcription output from vosk")
	}
	
	return transcription, nil
}

// transcribeWithSystemCommand tries other system speech recognition tools
func (t *Transcriber) transcribeWithSystemCommand(wavFile string) (string, error) {
	// Try different system commands
	cmds := []struct {
		name    string
		command string
		args    []string
	}{
		// Try system speech-to-text commands that might be available
		{"speech-recognition", "speech-recognition", []string{wavFile}},
		{"speech-to-text", "speech-to-text", []string{wavFile}},
	}
	
	for _, cmdInfo := range cmds {
		path, err := exec.LookPath(cmdInfo.command)
		if err == nil {
			log.Printf("Found system speech recognition tool: %s at %s", cmdInfo.name, path)
			cmd := exec.Command(path, cmdInfo.args...)
			
			log.Printf("Running %s...", cmdInfo.name)
			output, err := cmd.CombinedOutput()
			if err == nil && len(output) > 0 {
				log.Printf("Transcription result from %s: %s", cmdInfo.name, string(output))
				return string(output), nil
			} else if err != nil {
				log.Printf("%s failed: %v", cmdInfo.name, err)
			}
		}
	}
	
	return "", fmt.Errorf("no system speech recognition tools found")
}

// ExtractTextFromJSON extracts text field from JSON output
func ExtractTextFromJSON(jsonStr string) string {
	log.Println("Extracting text from JSON output...")
	// Look for "text" field in JSON
	textFieldIdx := strings.Index(jsonStr, "\"text\":")
	if textFieldIdx == -1 {
		log.Println("No text field found in JSON")
		return ""
	}
	
	// Move past the field name and colon
	textStart := textFieldIdx + 7
	
	// Find the opening quote
	for textStart < len(jsonStr) && jsonStr[textStart] != '"' {
		textStart++
	}
	if textStart >= len(jsonStr) {
		log.Println("No opening quote found in JSON text field")
		return ""
	}
	textStart++ // Move past the opening quote
	
	// Find the closing quote (accounting for escaped quotes)
	textEnd := textStart
	for textEnd < len(jsonStr) {
		if jsonStr[textEnd] == '"' && (textEnd == 0 || jsonStr[textEnd-1] != '\\') {
			break
		}
		textEnd++
	}
	
	if textEnd >= len(jsonStr) {
		log.Println("No closing quote found in JSON text field")
		return ""
	}
	
	extracted := jsonStr[textStart:textEnd]
	log.Printf("Extracted text from JSON: %s", extracted)
	return extracted
}