# Compile .proto
python -m grpc_tools.protoc \
       -I=app/protos \
       --python_out=app/protos \
       --grpc_python_out=app/protos \
       app/protos/embed.proto

# for linux:
sed -i '' 's/^import embed_pb2/from . import embed_pb2/' app/protos/embed_pb2_grpc.py

# for mac:
# sed -i '' 's/^import embed_pb2/from . import embed_pb2/' app/protos/embed_pb2_grpc.py
