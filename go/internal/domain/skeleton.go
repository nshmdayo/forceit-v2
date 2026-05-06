package domain

type Joint struct {
	Name       string  `json:"name"`
	Position   Vec3    `json:"position"`
	Confidence float64 `json:"confidence"`
}

type SkeletonFrame struct {
	TimestampMs int64   `json:"timestamp_ms"`
	Joints      []Joint `json:"joints"`
}
