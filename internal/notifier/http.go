package notifier

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"refurbed/assignment/internal/report"
	"sync"
	"sync/atomic"
	"time"
)

type httpNotifier struct {
	client      http.Client
	url         url.URL
	maxRequests int32
	mu          sync.Mutex
	payloads    []string
	ctx         context.Context
}

func NewHttpNotifier(timeout time.Duration, maxRequests int32, url url.URL, ctx context.Context) Notifier {
	return &httpNotifier{
		client:      http.Client{Timeout: timeout},
		url:         url,
		maxRequests: maxRequests,
		ctx:         ctx,
	}
}

func (h *httpNotifier) Enqueue(payload string) {
	h.mu.Lock()
	h.payloads = append(h.payloads, payload)
	h.mu.Unlock()
}

func (h *httpNotifier) Process() {
	var processing int32 = 0

	for {
		if processing == h.maxRequests {
			continue
		}

		h.mu.Lock()

		if len(h.payloads) == 0 {
			h.mu.Unlock()
			continue
		}

		payload := h.payloads[0]
		h.payloads = h.payloads[1:]

		atomic.AddInt32(&processing, 1)

		h.mu.Unlock()

		go h.notify(payload, &processing)
	}
}

func (h *httpNotifier) notify(payload string, processing *int32) {
	defer atomic.AddInt32(processing, -1)

	req, err := http.NewRequestWithContext(h.ctx, "POST", h.url.String(), bytes.NewBufferString(payload))

	if err != nil {
		report.Error(err)
		return
	}

	if _, err = h.client.Do(req); err != nil {
		report.Error(err)
	}
}
