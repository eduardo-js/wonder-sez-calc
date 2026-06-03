package httpx_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wonderlic-calc/backend/internal/httpx"
)

type testReq struct {
	Name   string  `json:"name"   binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Mode   string  `json:"mode"   binding:"omitempty,oneof=a b c"`
}

func newJSONContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestBindJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantOK     bool
		wantStatus int
		wantCode   string
		wantFields map[string]string
	}{
		{
			name:   "valid full",
			body:   `{"name":"foo","amount":1.5,"mode":"a"}`,
			wantOK: true,
		},
		{
			name:   "valid without optional mode",
			body:   `{"name":"foo","amount":10}`,
			wantOK: true,
		},
		{
			name:       "missing name",
			body:       `{"amount":5}`,
			wantOK:     false,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeValidation,
			wantFields: map[string]string{"name": "required"},
		},
		{
			// amount=0 is the float64 zero value; `required` fires before `gt`.
			name:       "amount zero",
			body:       `{"name":"foo","amount":0}`,
			wantOK:     false,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeValidation,
			wantFields: map[string]string{"amount": "required"},
		},
		{
			name:       "amount negative",
			body:       `{"name":"foo","amount":-1}`,
			wantOK:     false,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeValidation,
			wantFields: map[string]string{"amount": "gt"},
		},
		{
			name:       "mode not in set",
			body:       `{"name":"foo","amount":1,"mode":"z"}`,
			wantOK:     false,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeValidation,
			wantFields: map[string]string{"mode": "oneof"},
		},
		{
			name:       "malformed JSON",
			body:       `{not valid json`,
			wantOK:     false,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeBadRequest,
		},
		{
			name:       "wrong type amount as string",
			body:       `{"name":"foo","amount":"notanumber"}`,
			wantOK:     false,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeBadRequest,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, w := newJSONContext(tc.body)

			var dst testReq
			ok := httpx.BindJSON(c, &dst)

			if ok != tc.wantOK {
				t.Errorf("BindJSON returned %v, want %v", ok, tc.wantOK)
			}

			if tc.wantOK {
				if w.Code != 0 && w.Code != http.StatusOK {
					t.Errorf("expected no response written on success, got status %d body %s", w.Code, w.Body.String())
				}
				if c.IsAborted() {
					t.Error("expected context not aborted on success")
				}
				return
			}

			if !c.IsAborted() {
				t.Error("expected context aborted on failure")
			}
			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tc.wantStatus)
			}

			var body struct {
				Error struct {
					Code   string            `json:"code"`
					Fields map[string]string `json:"fields"`
				} `json:"error"`
			}
			if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			if body.Error.Code != tc.wantCode {
				t.Errorf("code: got %q, want %q", body.Error.Code, tc.wantCode)
			}

			for k, v := range tc.wantFields {
				if got, ok2 := body.Error.Fields[k]; !ok2 {
					t.Errorf("fields missing key %q", k)
				} else if got != v {
					t.Errorf("fields[%q]: got %q, want %q", k, got, v)
				}
			}
		})
	}
}

// TestBindJSON_SuccessPopulatesStruct verifies the success path populates the dst struct.
func TestBindJSON_SuccessPopulatesStruct(t *testing.T) {
	t.Parallel()
	body, _ := json.Marshal(map[string]any{"name": "test", "amount": 5.0})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var dst testReq
	ok := httpx.BindJSON(c, &dst)

	if !ok {
		t.Fatal("expected BindJSON to return true")
	}
	if dst.Name != "test" {
		t.Errorf("Name: got %q, want %q", dst.Name, "test")
	}
	if dst.Amount != 5.0 {
		t.Errorf("Amount: got %v, want 5.0", dst.Amount)
	}
}
