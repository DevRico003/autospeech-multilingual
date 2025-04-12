#!/bin/bash

# This script downloads and sets up the medium-sized Vosk model for English
# and updates the transcription script to use it

echo "Setting up Vosk medium model..."

# Create directory for models if it doesn't exist
mkdir -p ~/vosk-models

# Go to the models directory
cd ~/vosk-models

# Download the medium model (about 500MB)
echo "Downloading medium Vosk model (this may take a while)..."
curl -LO https://alphacephei.com/vosk/models/vosk-model-en-us-0.22.zip

# Extract the model
echo "Extracting model..."
unzip vosk-model-en-us-0.22.zip

# Clean up zip file
rm vosk-model-en-us-0.22.zip

# Update the vosk-transcribe.py script to use the medium model
echo "Updating transcription script to use medium model..."
sed -i 's/vosk-model-small-en-us-0.15/vosk-model-en-us-0.22/' ~/vosk-transcribe.py

echo "Setup complete! The medium Vosk model is now installed."
echo "Your speech recognition app will use the more accurate model now."