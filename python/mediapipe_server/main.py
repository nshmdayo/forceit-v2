import json
import socket
import time
from collections.abc import Callable
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


def build_frame(now_ms: int) -> SkeletonFrame:
    return SkeletonFrame(
        timestamp_ms=now_ms,
        joints=[
            Joint("right_wrist", Vec3(0.1, 0.2, 0.0), 0.9),
            Joint("left_wrist", Vec3(-0.1, 0.2, 0.0), 0.9),
        ],
    )


def validate_fps(fps: float) -> float:
    if fps <= 0:
        raise ValueError(f"fps must be positive, got {fps}")
    return fps


def stream_frames(
    conn: socket.socket,
    fps: float,
    frame_builder: Callable[[int], SkeletonFrame] = build_frame,
) -> None:
    validated_fps = validate_fps(fps)
    frame_interval_seconds = 1 / validated_fps
    while True:
        now_ms = int(time.time() * 1000)
        frame = frame_builder(now_ms)
        conn.sendall(encode_frame(frame))
        time.sleep(frame_interval_seconds)


def serve(
    host: str = "127.0.0.1",
    port: int = 50051,
    fps: float = 30,
    frame_builder: Callable[[int], SkeletonFrame] = build_frame,
) -> None:
    validated_fps = validate_fps(fps)
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server:
        server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        server.bind((host, port))
        server.listen(1)
        conn, _ = server.accept()
        with conn:
            stream_frames(conn=conn, fps=validated_fps, frame_builder=frame_builder)


if __name__ == "__main__":
    HOST = "127.0.0.1"
    PORT = 50051
    FPS = 30.0
    serve(host=HOST, port=PORT, fps=FPS)
