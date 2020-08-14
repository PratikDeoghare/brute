package brute_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pratikdeoghare/brute"
)

func TestRouter(t *testing.T) {

	routerTests := []struct {
		pattern        string
		reqPath        string
		expectedParams brute.Params
		expectedCode   int
	}{
		{
			pattern:        "/",
			reqPath:        "/",
			expectedParams: nil,
			expectedCode:   200,
		},
		{
			pattern:        "/a",
			reqPath:        "/a",
			expectedParams: nil,
			expectedCode:   200,
		},
		{
			pattern:        "/a/b",
			reqPath:        "/a/b",
			expectedParams: nil,
			expectedCode:   200,
		},
		{
			pattern:        "/:a/b",
			reqPath:        "/a/b",
			expectedParams: brute.Params{{"a", "a"}},
			expectedCode:   200,
		},
		{
			pattern:        "/:a/:b",
			reqPath:        "/a/b",
			expectedParams: brute.Params{{"a", "a"}, {"b", "b"}},
			expectedCode:   200,
		},
		{
			pattern:        "/a/:boo",
			reqPath:        "/a/woogie",
			expectedParams: brute.Params{{"boo", "woogie"}},
			expectedCode:   200,
		},
		{
			pattern:        "/a/:boo",
			reqPath:        "/b/woogie",
			expectedParams: nil,
			expectedCode:   404,
		},
	}

	for _, routerTest := range routerTests {
		pattern := routerTest.pattern
		router := brute.New(brute.Endpoint{
			Method:  http.MethodGet,
			Pattern: pattern,
			Handler: handlerAssertParams(t, pattern, routerTest.expectedParams),
		})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}

}

func handlerAssertParams(t *testing.T, s string, expectedParams brute.Params) brute.Handle {
	return func(w http.ResponseWriter, r *http.Request, params brute.Params) {
		if len(params) != len(expectedParams) {
			t.Errorf("params mismatch expected: %s, got: %s", expectedParams, params)
		}

		for i, param := range params {
			if param != expectedParams[i] {
				t.Errorf("params mismatch expected: %s, got: %s", expectedParams, params)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
