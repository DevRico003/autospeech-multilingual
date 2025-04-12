#!/bin/bash

# Install pip and Python dev packages
echo "Installing Python dependencies..."
sudo apt-get update
sudo apt-get install -y python3-pip python3-dev python3-venv

# Create virtual environment
echo "Creating virtual environment..."
python3 -m venv ~/vosk-env

# Activate virtual environment and install vosk
echo "Installing Vosk..."
source ~/vosk-env/bin/activate
pip install vosk sounddevice

# Download a small model
echo "Downloading Vosk model..."
mkdir -p ~/vosk-models
cd ~/vosk-models
if [ ! -d "vosk-model-small-en-us-0.15" ]; then
  curl -LO https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip
  unzip vosk-model-small-en-us-0.15.zip
  rm vosk-model-small-en-us-0.15.zip
fi

# Create a direct vosk-transcribe script in the project directory
echo "Creating Vosk transcription script..."
cat > ./vosk-transcribe << 'EOL'
#!/bin/bash
# Activate the virtual environment and run the transcription script
source ~/vosk-env/bin/activate
python3 - "$@" << 'PYCODE'
#!/usr/bin/env python3
import sys
import json
import os
from vosk import Model, KaldiRecognizer, SetLogLevel
import wave

SetLogLevel(-1)  # Disable debug messages

if len(sys.argv) < 2:
    print("Usage: vosk-transcribe input.wav")
    sys.exit(1)

wav_file = sys.argv[1]

# Find model in the user's home directory
model_path = os.path.expanduser("~/vosk-models/vosk-model-small-en-us-0.15")
if not os.path.exists(model_path):
    print(f"Error: Model not found at {model_path}")
    sys.exit(1)

model = Model(model_path)

# Open the WAV file
wf = wave.open(wav_file, "rb")
if wf.getnchannels() != 1 or wf.getsampwidth() != 2 or wf.getcomptype() != "NONE":
    print("Audio file must be WAV format mono PCM.")
    sys.exit(1)

# Create recognizer
rec = KaldiRecognizer(model, wf.getframerate())
rec.SetWords(True)

# Process audio
results = []
while True:
    data = wf.readframes(4000)
    if len(data) == 0:
        break
    if rec.AcceptWaveform(data):
        part_result = json.loads(rec.Result())
        results.append(part_result)

part_result = json.loads(rec.FinalResult())
results.append(part_result)

# Extract text from results
full_text = " ".join([res.get("text", "") for res in results if "text" in res])
print(full_text)
PYCODE
EOL

chmod +x ./vosk-transcribe

# Copy to /usr/local/bin for system-wide access
echo "Installing vosk-transcribe to system path..."
sudo cp ./vosk-transcribe /usr/local/bin/

echo "Setup complete! You can now run the speech recognition app."