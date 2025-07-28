Not sure to do the python server first or the Go server?

General Map:
                ┌────────────┐
                │  Text In   │
                └─────┬──────┘
                      │
         ┌────────────▼────────────┐
         │    Go Server (Main)     │
         │  - Chunking (by meaning)│
         │  - Batching all chunks  │
         │  - Sends to embedding svc│
         └────────────┬────────────┘
                      │
         ┌────────────▼────────────┐
         │ Python Embedding Server │
         │ - FastAPI or gRPC       │
         │ - Preloaded SBERT model │
         │ - Returns embeddings[]  │
         └────────────┬────────────┘
                      │
         ┌────────────▼────────────┐
         │  Go Server (Diff Logic) │
         │  - Cosine similarity    │
         │  - Thresholding         │
         │  - Diff labeling        │
         └────────────┬────────────┘
                      │
              ┌───────▼────────┐
              │   Output / UI  │
              │  - Git-like log│
              │  - CLI/JSON/UI │
              └────────────────┘
python server:
    []1. test the embedding and the difference.
    []2. Handle the chunking and the diffing.
    []3. Regardless of whether we are going to handle all the text in python or not, we will need to have the python server to understand and process chunks. So that in case the Go server failed, the python server handles it.

Go Server:
    to python:
        []1. How to do the chunking?
        []2. How to batch them?
        []3. Send through grpc.
    from python:
        []4. Similarity
        []5. Thresholding
        []6. Diffing
    
