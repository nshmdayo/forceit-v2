package sensor

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"forceit-v2/go/internal/domain"
)

type TCPJSONClient struct {
	Addr string
}

func (c TCPJSONClient) Stream(ctx context.Context, out chan<- domain.SkeletonFrame) error {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", c.Addr)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ctx.Err()
		}
		return fmt.Errorf("dial sensor server: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	tcpConn, _ := conn.(*net.TCPConn)
	const readDeadlineInterval = 500 * time.Millisecond
	var pending []byte

	for {
		if tcpConn != nil {
			if err := tcpConn.SetReadDeadline(time.Now().Add(readDeadlineInterval)); err != nil {
				return fmt.Errorf("set read deadline: %w", err)
			}
		}

		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			pending = append(pending, line...)
		}
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					continue
				}
			}
			if errors.Is(err, io.EOF) {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			}
			if errors.Is(err, context.Canceled) {
				return ctx.Err()
			}
			return err
		}

		line = pending
		pending = nil

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var f domain.SkeletonFrame
		if err := json.Unmarshal(line, &f); err != nil {
			continue
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- f:
		}
	}
}
