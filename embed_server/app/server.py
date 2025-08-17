import grpc
from concurrent import futures
import time

from app.protos import embed_pb2, embed_pb2_grpc
from app.embedder import TextProcessor

class EmbedderServicer(embed_pb2_grpc.EmbedderServicer):
    def __init__(self):
        self.processor = TextProcessor()

    def GetEmbeddings(self, request, context):
        try:
            results = self.processor.process(request.text, strategy=request.strategy)
            
            # Get data from results dictionary based on TextProcessor's output format
            chunks = results["chunks"]
            embeddings = [
                embed_pb2.Embedding(values=embedding) 
                for embedding in results["embeddings"]
            ]
            # Generate sequential indices since TextProcessor doesn't provide them
            indices = list(range(len(chunks)))

            return embed_pb2.EmbedResponse(
                chunks=chunks,
                embeddings=embeddings,
                indices=indices
            )
            print("Server working")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error processing request: {str(e)}")
            return embed_pb2.EmbedResponse()

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    embed_pb2_grpc.add_EmbedderServicer_to_server(EmbedderServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    try:
        while True:
            time.sleep(86400)  # one day
    except KeyboardInterrupt:
        server.stop(0)
        print("\nServer stopped gracefully")

if __name__ == '__main__':
    serve()
