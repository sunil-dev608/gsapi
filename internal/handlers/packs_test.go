package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPacksHandler_getPacksForItems(t *testing.T) {
	type args struct {
		numberOfItems int
	}
	tests := []struct {
		name       string
		configData []int
		args       args
		want       map[int]int
	}{
		{
			name:       "success: 1,000,000,000 items == 5000*20000",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 1000000000,
			},
			want: map[int]int{
				5000: 200000,
			},
		},
		{
			name:       "success: 10000 items == 5000*2",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 10000,
			},
			want: map[int]int{
				5000: 2,
			},
		},
		{
			name:       "success: 7000 items == 5000*1 + 2000*1",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 7000,
			},
			want: map[int]int{
				5000: 1,
				2000: 1,
			},
		},
		{
			name:       "success: 11000 items == 5000*2 + 1000*1",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 11000,
			},
			want: map[int]int{
				5000: 2,
				1000: 1,
			},
		},
		{
			name:       "success: 11000 items == 10000*1 + 1000*1",
			configData: []int{10000, 5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 11000,
			},
			want: map[int]int{
				10000: 1,
				1000:  1,
			},
		},
		{
			name:       "success: 4500 items == 2000*2 + 500*1",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 4500,
			},
			want: map[int]int{
				2000: 2,
				500:  1,
			},
		},
		{
			name:       "success: 20 items == 250*1",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 20,
			},
			want: map[int]int{
				250: 1,
			},
		},
		{
			name:       "success: 451 items == 250*2",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: 451,
			},
			want: map[int]int{
				250: 2,
			},
		},
		{
			name:       "success: -1 items == ",
			configData: []int{5000, 2000, 1000, 500, 250},
			args: args{
				numberOfItems: -1,
			},
			want: map[int]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &PacksHandler{
				packSizes: tt.configData,
			}
			if got := ph.getPacksForItems(tt.args.numberOfItems); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PacksHandler.getPacksForItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacksHandler_GetPacksForItems(t *testing.T) {
	configData := []int{5000, 2000, 1000, 500, 250}
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			requestBody: `{
				"items": 12000
			}`,
			expectedCode: http.StatusOK,
			expectedBody: `{
				"5000": 2,
				"2000": 1
			}`,
		},
		{
			name: "error: bad request",
			requestBody: `{
				"items": xyz
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message": "bad request"}`,
		},
		{
			name: "error: bad request",
			requestBody: `{
				"items": -1
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message": "items must be greater than 0"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postHandler := func(w http.ResponseWriter, req *http.Request) {
				ph := &PacksHandler{
					packSizes: configData,
				}
				ph.GetPacksForItems(echo.New().NewContext(req, w))
			}
			req := httptest.NewRequest("POST", "/api/v1/get-packs", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(postHandler)

			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			assert.Equal(t, tt.expectedCode, rr.Code, tt.name)

			// Check the response body is what we expect.
			assert.JSONEq(t, tt.expectedBody, rr.Body.String(), tt.name)

		})
	}
}
