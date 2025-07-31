import spacy
from sentence_transformers import SentenceTransformer

# everything here is happening under the hood.
# even the different arrays are being embedded differently.
class Embedder:
    def __init__(self, model="all-MiniLM-L6-v2"):
        self.model = SentenceTransformer(model)

    def embed(self, texts):
        return self.model.encode(texts, convert_to_numpy=True).tolist()

class Chunker:
    def __init__(self):
        self.nlp = spacy.load("en_core_web_sm")

    def chunk(self, text, strategy="sentence"):
        doc = self.nlp(text)
        if strategy == "sentence":
            return [sent.text.strip() for sent in doc.sents]

        elif strategy == "clause":
            chunks, current = [], ""
            for token in doc:
                current += token.text_with_ws
                if token.is_sent_end or (
                    token.pos_ == "CCONJ" and token.text.lower() in {"but", "and", "or"}
                ) or (
                    token.dep_ == "mark" and token.text.lower() in {"because", "although", "since", "while"}
                ):
                    if current.strip():
                        chunks.append(current.strip())
                        current = ""
            if current.strip():
                chunks.append(current.strip())
            return chunks

        elif strategy == "semantic":
            return [sent.text.strip() for sent in doc.sents]

        return []

class TextProcessor:
    def __init__(self, model="all-MiniLM-L6-v2"):
        self.chunker = Chunker()
        self.embedder = Embedder(model)

    def process(self, text, strategy="semantic"):
        chunks = self.chunker.chunk(text, strategy)
        embeddings = self.embedder.embed(chunks)
        return {
            "chunks": chunks,
            "embeddings": embeddings,
            "strategy": strategy,
            "chunk_count": len(chunks),
        }

# Example usage
if __name__ == "__main__":
    text = ("This is the first sentence. And here is another one but this one is a bit too large "
            "so it might need more than a single vector. I am wondering if this really is working "
            "correctly? Let me test with another sentence to see how the chunking works.")

    processor = TextProcessor()
    result = processor.process(text, strategy="semantic")

    for i, results in enumerate(result):
        print(f'{i} " {results}')

    for chunk in result["chunks"]:
        print('-', chunk)

    print("Embeddings shape:", len(result["embeddings"]), "x", len(result["embeddings"][0]))
