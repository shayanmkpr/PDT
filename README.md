### Promt Development Tools

## Tool 1:

# Git Diff of prompts(GDP):
    - Using Vector Embedding, we will calculate how much a part of a prompt has changed. If the change is more than a certain threshold (or on a spectrum) we will mark that line as changed, added, or removed.
        - This way, we will not mark unnecessary changes and only mark those that actually made a difference.
        - Given the model that the client is using and its embedding(GPT, Bard, ...) we should choose a different embedding for the git as well.
        - This is meaning based diff not literal diff. We are comparing different meanings, not words.

# Configuration:
    1. Embedding --> Should be chosen according to the Model that is being used.
    2. Fixed threshold or a spectrum.
    3. Sensitivity of the changes.

# To do:
    - Fyne for GUI and the client.
    - Compare the vects using Go and show in Fyne.
Test commit with correct email
