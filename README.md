### Promt Development Tools

## Tool 1:

# Git Diff of prompts(GDP):
- Using Vector Embedding, we will calculate how much a part of a prompt has changed. If the change is more than a certain threshold (or on a spectrum) we will mark that line as changed, added, or removed.
    - This way, we will not mark unnecessary changes and only mark those that actually made a difference.
    - Given the model that the client is using and its embedding(GPT, Bard, ...) we should choose a different embedding for the git as well.
    - ...
# Configuratoin:
   
    2. Fixed threshold or a spectrum.
    3. Sensitivity of the changes.
    
# What to do?
    1. Use clauses chunks (sub-sentences)
    2. Use something like at first for the demo text-embedding-ada-002 but smaller and easier to work with.
    3. Use the static threhsold at first. Work on the spectrum later.

