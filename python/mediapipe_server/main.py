import json
import socket
import time
from dataclasses import asdict, dataclass


@dataclass
class Vec3:
    x: float
    y: float
    z: float


@dataclass
class Joint:
    name: str
    position: Vec3
    confidence: float


@dataclass
class SkeletonFrame:
    timestamp_ms: int
    joints: list[Joint]


def encode_frame(frame: SkeletonFrame) -> bytes:
    payload = asdict(frame)
    return (json.dumps(payload) + "\n").encode("utf-8")


def serve(host: str = "127.0.0.1", port: int = 50051) -> None:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server:
        server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        server.bind((host, port))
        server.listen(1)
        conn, _ = server.accept()
        with conn:
            while True:
                now = int(time.time() * 1000)
                frame = SkeletonFrame(
                    timestamp_ms=now,
                    joints=[
                        Joint("right_wrist", Vec3(0.1, 0.2, 0.0), 0.9),
                        Joint("left_wrist", Vec3(-0.1, 0.2, 0.0), 0.9),
                    ],
                )
                conn.sendall(encode_frame(frame))
                time.sleep(1 / 30)


if __name__ == "__main__":
    serve()
