package middlewares

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/p000ic/webapi"
)

func TestZeroValueRecoveryConcurrentRequests(t *testing.T) {
	host := webapi.NewHost(webapi.Config{DisableAutoReport: true}, &Recovery{})
	if err := host.AddEndpoint(http.MethodGet, "/panic", func(*webapi.Context) {
		panic("test panic")
	}); err != nil {
		t.Fatalf("AddEndpoint() error = %v", err)
	}

	const requests = 16
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			recorder := httptest.NewRecorder()
			host.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/panic", nil))
			if recorder.Code != http.StatusInternalServerError {
				t.Errorf("ServeHTTP() status = %d; want %d", recorder.Code, http.StatusInternalServerError)
			}
		}()
	}
	wg.Wait()
}
