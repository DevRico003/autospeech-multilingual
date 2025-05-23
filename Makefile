.PHONY: build run clean deps setup-small setup-medium

BINARY_NAME=autospeech

build:
	go build -o $(BINARY_NAME) ./cmd/autospeech

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -f speech-reco.log

# Install dependencies
deps:
	go mod tidy

# Setup small Vosk model (recommended for most use cases)
setup-small:
	./setup-vosk-small.sh

# Setup medium Vosk model (better accuracy but requires more resources)
setup-medium:
	./setup-vosk-medium.sh
