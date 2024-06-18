package pagination

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFromURLQuery(t *testing.T) {
	tests := map[string]struct {
		offset  string
		limit   string
		want    Pagination
		wantErr bool
	}{
		"when offset cannot be converted to int, returns error": {
			offset:  "pouet",
			limit:   "100",
			wantErr: true,
		},
		"when limit cannot be converted to int, returns error": {
			offset:  "100",
			limit:   "pouet",
			wantErr: true,
		},
		"when offset is negative, returns error": {
			offset:  "-2",
			limit:   "100",
			wantErr: true,
		},
		"when limit is negative, returns error": {
			offset:  "100",
			limit:   "-2",
			wantErr: true,
		},
		"nominal": {
			offset: "0",
			limit:  "100",
			want: Pagination{
				Offset: 0,
				Limit:  100,
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var page Pagination
			var err error
			api := gin.Default()
			api.GET("/", func(context *gin.Context) {
				page, err = GetFromURLQuery(context)
			})
			url := fmt.Sprintf("/?offset=%s&&limit=%s", tt.offset, tt.limit)

			r := httptest.NewRequest(http.MethodGet, url, bytes.NewReader(nil))
			rw := httptest.NewRecorder()
			api.ServeHTTP(rw, r)

			assert.Equal(t, tt.want, page)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestPageableToLabelSlice(t *testing.T) {
	t.Run("When data is not a slice, returns empty slice and error", func(t *testing.T) {
		sampleData := Pageable{}
		sampleData.Data = int64(1)
		data, err := PageableToSlice[Label](sampleData)
		assert.Error(t, err)
		assert.Equal(t, []Label{}, data)
	})

	t.Run("When data is a slice of various types, returns empty slice and error", func(t *testing.T) {
		var customData interface{}
		interData := []byte(`["toto", 1]`)
		err := json.Unmarshal(interData, &customData)
		require.NoError(t, err)
		sampleData := Pageable{
			Data: customData,
		}
		data, err := PageableToSlice[Label](sampleData)
		assert.Error(t, err)
		assert.Equal(t, []Label{}, data)
	})

	t.Run("When data is a slice of Labels, returns no error and slice of label", func(t *testing.T) {
		labelPageable := MockPageableLabel("my label", "my second label")
		data, err := PageableToSlice[Label](labelPageable)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(data))
		assert.True(t, data[0].Label.Valid)
		assert.Equal(t, "my label", data[0].Label.String)
		assert.True(t, data[1].Label.Valid)
		assert.Equal(t, "my second label", data[1].Label.String)

	})
}
func TestDefaultPagination(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {

		out := Default()
		expect := Pageable{
			Offset: 0,
			Limit:  500,
		}

		assert.Equal(t, expect.Limit, out.Limit)
		assert.Equal(t, expect.Offset, out.Offset)
	})
}
func TestBuildPageable(t *testing.T) {
	t.Run("nominal label", func(t *testing.T) {
		data := []Label{
			{
				Label: NullEmptyString{
					sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
		}

		out := BuildPageable[Label](Pagination{}, 1, data)
		expect := Pageable{
			Offset: 0,
			Limit:  0,
			Total:  1,
			Data:   data,
		}

		assert.Equal(t, expect.Data, out.Data)
		assert.Equal(t, expect.Limit, out.Limit)
		assert.Equal(t, expect.Total, out.Total)
		assert.Equal(t, expect.Offset, out.Offset)
	})

	t.Run("nominal label with specified offset and limit", func(t *testing.T) {
		data := []Label{
			{
				Label: NullEmptyString{
					sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
		}

		out := BuildPageable[Label](Pagination{
			Limit:  12,
			Offset: 47,
		}, 1, data)

		expect := Pageable{
			Offset: 47,
			Limit:  12,
			Total:  1,
			Data:   data,
		}

		assert.Equal(t, expect.Data, out.Data)
		assert.Equal(t, expect.Limit, out.Limit)
		assert.Equal(t, expect.Total, out.Total)
		assert.Equal(t, expect.Offset, out.Offset)
	})
}
