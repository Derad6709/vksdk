package callback_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestRetryCounterFromContext(t *testing.T) {
	const retryCounter = 123
	ctx := context.WithValue(
		context.Background(),
		internal.CallbackRetryCounterKey,
		retryCounter,
	)
	assert.Equal(t, retryCounter, callback.RetryCounterFromContext(ctx))
}

func TestRetryAfter(t *testing.T) {
	code := http.StatusGone
	date := time.Now().Add(time.Minute * 5)

	cb := callback.NewCallback()
	cb.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		callback.RetryAfter(
			ctx,
			code,
			date,
		)
	})

	jsonStr := []byte(`{"type": "message_new","object": {}}`)

	req, err := http.NewRequest("POST", "/callback", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cb.HandleFunc)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusText(code)+"\n", rr.Body.String())
	assert.Equal(t, code, rr.Code)
	assert.Equal(t, date.Format(http.TimeFormat), rr.Header().Get("Retry-After"))
}