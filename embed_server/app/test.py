from embedder import ChunkSplitter

text = "hello this is shayan. i am testing my app and I wanted to make sure tha this is working correctly with the clause not just the snetences but I am not sure about it."
splitter = ChunkSplitter()

print(splitter.chunk_text_advanced(text, strategy = "clause"))
