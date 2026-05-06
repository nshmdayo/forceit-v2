package sensor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"

	"forceit-v2/go/internal/domain"
)

type TCPJSONClient struct {
	Addr string
}

func (c TCPJSONClient) Stream(ctx context.Context, out chan<- domain.SkeletonFrame) error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("dial sensor server: %w", err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		var f domain.SkeletonFrame
		if err := json.Unmarshal(scanner.Bytes(), &f); err != nil {
			continue
		}
		out <- f
	}
	return scanner.Err()
}
