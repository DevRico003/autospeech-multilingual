package audio

import (
	"log"
	"os"
	"path/filepath"

	"github.com/tarasowski/autospeech/pkg/config"
)

// SaveAsWav converts raw audio data to a WAV file
func SaveAsWav(audioData []byte, outputFile string) error {
	log.Println("Creating WAV file...")
	
	// Create a simple WAV header (44 bytes)
	header := make([]byte, 44)

	// RIFF header
	copy(header[0:4], []byte("RIFF"))
	// File size (to be filled later)
	// WAVE header
	copy(header[8:12], []byte("WAVE"))
	// fmt chunk
	copy(header[12:16], []byte("fmt "))
	// Chunk size: 16 for PCM
	header[16] = 16
	header[17] = 0
	header[18] = 0
	header[19] = 0
	// Audio format: 1 for PCM
	header[20] = 1
	header[21] = 0
	// Number of channels: 1 for mono
	header[22] = byte(config.Channels)
	header[23] = 0
	// Sample rate (16000)
	sampleRateBytes := uint32(config.SampleRate)
	header[24] = byte(sampleRateBytes)
	header[25] = byte(sampleRateBytes >> 8)
	header[26] = byte(sampleRateBytes >> 16)
	header[27] = byte(sampleRateBytes >> 24)
	// Byte rate: SampleRate * NumChannels * BitsPerSample/8
	byteRate := uint32(config.SampleRate * config.Channels * 16 / 8)
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)
	// Block align: NumChannels * BitsPerSample/8
	blockAlign := config.Channels * 16 / 8
	header[32] = byte(blockAlign)
	header[33] = byte(blockAlign >> 8)
	// Bits per sample
	header[34] = 16
	header[35] = 0
	// data chunk
	copy(header[36:40], []byte("data"))
	// Data size
	dataSize := uint32(len(audioData))
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	// Set file size in header (file size - 8 bytes for RIFF header)
	fileSize := uint32(len(header) + len(audioData) - 8)
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		return err
	}

	// Write header and audio data to file
	file, err := os.Create(outputFile)
	if err != nil {
		log.Printf("Error creating WAV file: %v", err)
		return err
	}
	defer file.Close()

	if _, err := file.Write(header); err != nil {
		log.Printf("Error writing WAV header: %v", err)
		return err
	}
	if _, err := file.Write(audioData); err != nil {
		log.Printf("Error writing WAV data: %v", err)
		return err
	}

	log.Println("WAV file created successfully")
	return nil
}