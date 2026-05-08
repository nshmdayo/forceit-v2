import pytest

from mediapipe_server.main import (
    Joint,
    SkeletonFrame,
    Vec3,
    encode_frame,
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


def test_validate_fps_accepts_positive_value() -> None:
    assert validate_fps(30.0) == 30.0
