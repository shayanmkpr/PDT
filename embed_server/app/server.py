from concurrent import futures
import grpc
from app.embedder import Embedder
from app.protos import embed_pb2, embed_pb2_grpc

class EmbedderServicer(embed_pb2_grpc.EmbedderServicer):
    def __init__(self):
        self.embedder = Embedder()

    def GetEmbeddings(self, request, context):
        vectors = self.embedder.embed_texts(request.texts)
        return embed_pb2.EmbedResponse(
            embeddings=[embed_pb2.Embedding(values=vec) for vec in vectors]
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    embed_pb2_grpc.add_EmbedderServicer_to_server(EmbedderServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC Embedder running on port 50051...")
    server.wait_for_termination()
