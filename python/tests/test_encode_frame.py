import pytest

from mediapipe_server.main import (
    Joint,
    SkeletonFrame,
    Vec3,
    encode_frame,
    stream_frames,
    validate_fps,
)


def test_encode_frame_newline_delimited_json() -> None:
    frame = SkeletonFrame(timestamp_ms=1, joints=[Joint("j", Vec3(1, 2, 3), 0.5)])
    payload = encode_frame(frame)
    assert payload.endswith(b"\n")
    assert b'"timestamp_ms": 1' in payload


@pytest.mark.parametrize("fps", [0, -1, -30.0])
def test_validate_fps_rejects_non_positive_values(fps: float) -> None:
    with pytest.raises(ValueError, match="fps must be positive"):
        validate_fps(fps)


@pytest.mark.parametrize("fps", ["30", None])
def test_validate_fps_rejects_non_numeric_values(fps: object) -> None:
    with pytest.raises(ValueError, match="fps must be a number"):
        validate_fps(fps)  # type: ignore[arg-type]


def test_validate_fps_accepts_positive_value() -> None:
    assert validate_fps(30.0) == 30.0


def test_stream_frames_breaks_on_broken_pipe(monkeypatch: pytest.MonkeyPatch) -> None:
    class BrokenConn:
        def sendall(self, data: bytes) -> None:
            raise BrokenPipeError

    def fail_sleep(_seconds: float) -> None:
        raise AssertionError("sleep should not be called after disconnect")

    monkeypatch.setattr("mediapipe_server.main.time.sleep", fail_sleep)

    fixed_frame = SkeletonFrame(timestamp_ms=123, joints=[])

    stream_frames(conn=BrokenConn(), fps=30.0, frame_builder=lambda _now_ms: fixed_frame)
