import spacy
from sentence_transformers import SentenceTransformer
from typing import List, Dict, Tuple

class Embedder:
    def __init__(self, model_name="all-MiniLM-L6-v2"):
        self.model = SentenceTransformer(model_name)
    
    def embed_texts(self, texts):
        return self.model.encode(texts, convert_to_numpy=True).tolist()

class ChunkSplitter:
    def __init__(self):
        self.nlp = spacy.load("en_core_web_sm")
    
    def chunk_text_advanced(self, text: str, strategy: str = "sentence") -> Dict[str, List[str]]:
        """
        Advanced chunking with different strategies.
        
        Args:
            text: Input text to chunk
            strategy: "sentence", "clause", or "semantic"
        
        Returns:
            Dictionary with 'texts' key containing list of text chunks
        """
        doc = self.nlp(text)
        chunks = []
        
        if strategy == "sentence":
            # Simple sentence-based chunking
            chunks = [sent.text.strip() for sent in doc.sents]
            
        elif strategy == "clause":
            # Split on clauses and coordinating conjunctions
            current_chunk = ""
            
            for token in doc:
                current_chunk += token.text_with_ws
                
                # Split on sentence boundaries or major conjunctions
                if (token.is_sent_end or 
                    (token.pos_ == "CCONJ" and token.text.lower() in ["but", "and", "or"]) or
                    (token.dep_ == "mark" and token.text.lower() in ["because", "although", "since", "while"])):
                    
                    if current_chunk.strip():
                        chunks.append(current_chunk.strip())
                        current_chunk = ""
            
            # Add remaining text
            if current_chunk.strip():
                chunks.append(current_chunk.strip())
                
        elif strategy == "semantic":
            # Group sentences by semantic similarity (simplified approach)
            sentences = [sent.text.strip() for sent in doc.sents]
            chunks = sentences  # For now, just return sentences
            
        return {"texts": [chunk for chunk in chunks if chunk]}

class TextProcessor:
    """Combined class that chunks text and generates embeddings"""
    
    def __init__(self, embedding_model="all-MiniLM-L6-v2"):
        self.chunker = ChunkSplitter()
        self.embedder = Embedder(embedding_model)
    
    def process_text(self, text: str, strategy: str = "sentence") -> Dict[str, any]:
        """
        Process text by chunking and embedding.
        
        Args:
            text: Input text to process
            strategy: Chunking strategy ("sentence", "clause", or "semantic")
        
        Returns:
            Dictionary containing chunks, embeddings, and metadata
        """
        # Get chunks
        chunk_result = self.chunker.chunk_text_advanced(text, strategy)
        chunks = chunk_result["texts"]
        
        # Generate embeddings for chunks
        embeddings = self.embedder.embed_texts(chunks)
        
        return {
            "chunks": chunks,
            "embeddings": embeddings,
            "chunk_count": len(chunks),
            "strategy": strategy
        }
    
    def get_chunks_only(self, text: str, strategy: str = "sentence") -> List[str]:
        """
        Get just the text chunks (for compatibility with existing embed_texts calls).
        
        Args:
            text: Input text to chunk
            strategy: Chunking strategy
        
        Returns:
            List of text chunks
        """
        chunk_result = self.chunker.chunk_text_advanced(text, strategy)
        return chunk_result["texts"]

# Example usage and integration patterns
if __name__ == "__main__":
    # Initialize processor
    processor = TextProcessor()
    
    sample_text = """This is the first sentence. And here is another one but this one is a bit too large so it might need more than a single vector. I am wondering if this really is working correctly? Let me test with another sentence to see how the chunking works."""
    
    # Method 1: Full processing (chunks + embeddings)
    result = processor.process_text(sample_text, strategy="clause")
    print("Full processing result:")
    print(f"Chunks: {result['chunks']}")
    print(f"Embeddings shape: {len(result['embeddings'])} x {len(result['embeddings'][0])}")
    print(f"Strategy used: {result['strategy']}")
    
    print("\n" + "="*50 + "\n")
    
    # Method 2: Just get chunks for existing embedder workflow
    chunks = processor.get_chunks_only(sample_text, strategy="sentence")
    embeddings = processor.embedder.embed_texts(chunks)
    print("Separate chunking and embedding:")
    print(f"Chunks: {chunks}")
    print(f"Embeddings shape: {len(embeddings)} x {len(embeddings[0])}")
    
    print("\n" + "="*50 + "\n")
    
    # Method 3: Direct integration with your existing classes
    chunker = ChunkSplitter()
    embedder = Embedder()
    
    chunk_result = chunker.chunk_text_advanced(sample_text, strategy="clause")
    texts_to_embed = chunk_result["texts"]  # This is what goes into embed_texts
    final_embeddings = embedder.embed_texts(texts_to_embed)
    
    print("Direct integration:")
    print(f"Chunks to embed: {texts_to_embed}")
    print(f"Final embeddings shape: {len(final_embeddings)} x {len(final_embeddings[0])}")
