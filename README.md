# Speech Recognition Tool

A Go application that transcribes speech from your microphone and displays it in the terminal. The transcribed text is copied to clipboard for easy pasting.

## Features

- Real-time speech transcription displayed in terminal
- Transcribed text automatically copied to clipboard
- Minimalist system tray interface with Start/Stop Recording options
- No distracting notifications or text in system tray - all output displayed in terminal only

## Requirements

- Go 1.21+
- PortAudio development libraries
- Python 3.8+ (for Vosk)
- C compiler (for CGO dependencies)

## Installation

### 1. Install dependencies

#### Ubuntu/Debian
```bash
sudo apt-get install -y libportaudio2 libx11-dev xorg-dev libxtst-dev libpng-dev libayatana-appindicator3-dev python3-pip python3-dev python3-venv unzip
```

#### macOS
```bash
brew install portaudio python unzip
```

### 2. Install Vosk

Run the included setup script for the small model (recommended for lower-resource machines):
```bash
./setup-vosk-small.sh
```

For better accuracy with the medium model (~500MB, requires more resources):
```bash
./setup-vosk-medium.sh
```

### 3. Build and run

```bash
make build
./autospeech
```

Or you can build and run in one step:

```bash
make run
```

## Usage

1. Start the application
2. Click the system tray icon to access controls
3. Click "Start Recording" to begin speech recognition
4. Speak clearly into your microphone - you'll see real-time transcription in the terminal
5. Click "Stop Recording" when finished
6. The transcribed text will be copied to clipboard automatically
7. Use Ctrl+V or your system's paste shortcut to paste the text where needed

## Models

The application uses the Vosk speech recognition toolkit with different model options:

- **Small Model** (~40MB): Fast with decent accuracy (default)
- **Medium Model** (~500MB): Better accuracy but requires more resources

To switch between models, use the provided make commands:
- For small model: `make setup-small`
- For medium model: `make setup-medium`

## Moving to Binary Distribution

If you want to distribute the compiled binary:

1. First ensure the Vosk model and transcription tool are set up on the target system:
   ```bash
   make setup-small   # or setup-medium for better accuracy
   ```

2. Copy the compiled binary to the desired location
   ```bash
   sudo cp autospeech /usr/local/bin/
   ```

3. The application will find the Vosk transcription tool automatically, as it's installed in standard system locations by the setup script.

## License

MIT