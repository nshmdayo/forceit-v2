from mediapipe_server.main import SkeletonFrame, Joint, Vec3, encode_frame


def test_encode_frame_newline_delimited_json() -> None:
    frame = SkeletonFrame(timestamp_ms=1, joints=[Joint("j", Vec3(1, 2, 3), 0.5)])
    payload = encode_frame(frame)
    assert payload.endswith(b"\n")
    assert b'"timestamp_ms": 1' in payload
