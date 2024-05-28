package loggingtest

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monorepo/common/logging"
)

// IgnoreFieldValue is used to ignore a specific field in AssertLogs.
const IgnoreFieldValue = "<ignored value>"

// Recorder is a helper to test a logger.
type Recorder struct {
	*safeBuffer
	expectedTSFieldName string
	fixedExpectedFields logging.Fields
}

// NewRecorder creates a new logger recorder.
func NewRecorder(expectedTSFieldName string) *Recorder {
	return &Recorder{
		safeBuffer: &safeBuffer{
			Buffer: bytes.NewBuffer(nil),
		},
		expectedTSFieldName: expectedTSFieldName,
	}
}

// Writer returns the output writer to use in your logger.
func (r *Recorder) Writer() io.Writer {
	return r.safeBuffer
}

// ExpectFields always check the given fields for all log checked by AssertLogs().
func (r *Recorder) ExpectFields(fields logging.Fields) {
	r.fixedExpectedFields = fields
}

// AssertLog checks the given log line has been only written by the logger.
func (r *Recorder) AssertLog(t *testing.T, expectedUserFields logging.Fields) {
	r.AssertLogs(t, []logging.Fields{
		expectedUserFields,
	})
}

// AssertLogs checks the given logs has been written by the logger.
func (r *Recorder) AssertLogs(t *testing.T, listExpectedUserFields []logging.Fields) {
	var ms []logging.Fields

	jsonDecoder := json.NewDecoder(r.safeBuffer)
	for {
		var fields logging.Fields

		err := jsonDecoder.Decode(&fields)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		ms = append(ms, fields)
	}

	if !assert.Equal(t, len(listExpectedUserFields), len(ms), "number of log lines incorrect") {
		return
	}

	for idx := 0; idx < len(ms); idx++ {
		m, expectedUserFields := ms[idx], listExpectedUserFields[idx]

		ts, ok := m[r.expectedTSFieldName]
		if ok {
			d, err := time.Parse(time.RFC3339, ts.(string))
			assert.NoError(t, err)
			assert.WithinDuration(t, time.Now(), d, 5*time.Second)

			delete(m, r.expectedTSFieldName)
		} else {
			t.Errorf("missing %q field", r.expectedTSFieldName)
		}

		expectedFields := r.computeExpectedFields(expectedUserFields)

		for fk, fv := range expectedFields {
			if fv == IgnoreFieldValue {
				if _, ok := m[fk]; ok {
					delete(m, fk)
				} else {
					t.Errorf("missing %q field", fk)
				}

				delete(expectedFields, fk)
			}
		}

		assert.Equal(t, expectedFields, m)
	}
}

func (r *Recorder) computeExpectedFields(expectedUserFields logging.Fields) logging.Fields {
	expectedFields := make(logging.Fields)

	for fn, fv := range r.fixedExpectedFields {
		expectedFields[fn] = fv
	}

	for fn, fv := range expectedUserFields {
		expectedFields[fn] = fv
	}

	return expectedFields
}

type safeBuffer struct {
	*bytes.Buffer
	m sync.Mutex
}

func (s *safeBuffer) Read(p []byte) (n int, err error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.Buffer.Read(p)
}

func (s *safeBuffer) Write(p []byte) (n int, err error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.Buffer.Write(p)
}

func (s *safeBuffer) Len() int {
	s.m.Lock()
	defer s.m.Unlock()

	return s.Buffer.Len()
}

func (s *safeBuffer) Bytes() []byte {
	s.m.Lock()
	defer s.m.Unlock()

	return s.Buffer.Bytes()
}
