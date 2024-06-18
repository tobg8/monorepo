package pagination

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"time"
)

// Pagination to use this struct for all endpoint in the project that require paging
type Pagination struct {
	Offset int
	Limit  int
}

// Pageable describes a generic model
type Pageable struct {
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int64       `json:"total"`
	Data   interface{} `json:"data"`
}

// ResponseError encapsulates error in message to send to HTTP client
type ResponseError struct {
	Message string `json:"message"`
}

// NotFoundError is returned by repository when an entity is not found by its ID
type NotFoundError struct {
	Entity interface{}
}

// Error implement err interface
func (e NotFoundError) Error() string {
	return fmt.Sprintf("%T not found", e.Entity)
}

// RepositoryError defines errors happening at repository level
type RepositoryError struct {
	Usecase   string
	VersionID uint64
	Err       error
}

func (e RepositoryError) Error() string {
	return fmt.Sprintf("err in repository for usecase %s: '%v'", e.Usecase, e.Err)
}

// DeletePeriodError defines errors happening at repository level while deleting period
type DeletePeriodError struct {
	Usecase   string
	VersionID uint64
	PeriodID  uint64
	Err       error
}

func (e DeletePeriodError) Error() string {
	return fmt.Sprintf("err in deleting period for usecase %s: '%v' for periodID %d", e.Usecase, e.Err, e.PeriodID)
}

// RowsAffectedError defines errors when the number of affected rows by Update is not the one expected
type RowsAffectedError struct {
	Usecase      string
	VersionID    uint64
	AffectedRows int
	ExpectedRows int
}

func (e RowsAffectedError) Error() string {
	return fmt.Sprintf("err in repository for usecase %s: %d affected rows, %d was expected",
		e.Usecase,
		e.AffectedRows,
		e.ExpectedRows)
}

// BadRequestKeyError defines errors for bad requests when Key is not found in url
type BadRequestKeyError struct {
	Key string
}

func (e BadRequestKeyError) Error() string {
	return fmt.Sprintf("bad request: %q is not found in url", e.Key)
}

// BadRequestValueError defines errors for bad requests when Value is not valid
type BadRequestValueError struct {
	Key   string
	Value interface{}
	Err   error
}

func (e BadRequestValueError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("bad request: %v is not a valid value for key %q", e.Value, e.Key)
	}
	return fmt.Sprintf("bad request, value from key %q could not be parsed: %q", e.Key, e.Err.Error())
}

// MissingQueryParameterError defines errors when URL parameter is missing
type MissingQueryParameterError struct {
	Key string
}

func (e MissingQueryParameterError) Error() string {
	if e.Key == "" {
		return "missing key in url query"
	}
	return fmt.Sprintf("missing key %q in query string", e.Key)
}

// Label describes label entity
type Label struct {
	Label NullEmptyString `json:"label"`
}

// LabelRow describes label row entity for scan
type LabelRow struct {
	ID    int
	Label NullEmptyString
}

// NullBool encapsulates sql null boolean with custom marshalling/unmarshalling
type NullBool struct {
	sql.NullBool
}

// IsEmpty returns true if models.NullBool is either not valid or false
func (nb NullBool) IsEmpty() bool {
	if !nb.Valid || !nb.Bool {
		return true
	}
	return false
}

// MarshalJSON marshals models.NullBool datatype
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON unmarshal models.NullBool datatype
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nb.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nb.Bool)
	nb.Valid = err == nil
	return err
}

// MapToBool returns NullBool boolean value pointer if valid, nil otherwise
func (nb NullBool) MapToBool() *bool {
	if !nb.Valid {
		return nil
	}
	value := nb.Bool
	return &value
}

// MapForRequest returns boolean value (integer) if valid, nil otherwise
func (nb NullBool) MapForRequest() interface{} {
	if !nb.Valid {
		return nil
	}

	if nb.Bool {
		return 1
	}
	return 0
}

// NullInt encapsulates sql null int with custom marshalling/unmarshalling
type NullInt struct {
	sql.NullInt64
}

// MarshalJSON marshals models.NullInt datatype
func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid || ni.Int64 == 0 {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON unmarshal models.NullInt datatype
func (ni *NullInt) UnmarshalJSON(b []byte) error {
	var f float64
	err := json.Unmarshal(b, &f)
	ni.Int64 = int64(f)
	ni.Valid = err == nil
	if string(b) == "null" {
		ni.Valid = false
	}
	return err
}

// Scan scans int to models.NullInt datatype
func (ni *NullInt) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*ni = NullInt{}
	}

	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	ni.Int64 = i.Int64
	if i.Int64 == 0 {
		ni.Valid = false
	} else {
		ni.Valid = true
	}
	return nil
}

// MapToInt64 returns NullInt integer value pointer if valid, nil otherwise
func (ni NullInt) MapToInt64() *int64 {
	if !ni.Valid {
		return nil
	}
	value := ni.Int64
	return &value
}

// MapToInt32 returns NullInt integer value pointer if valid and in range, nil otherwise
func (ni NullInt) MapToInt32() (*int32, error) {
	if !ni.Valid {
		return nil, nil
	}

	if ni.Int64 > math.MaxInt32 {
		return nil, fmt.Errorf("could not convert NullInt to *int32, value %d too high", ni.Int64)
	}

	value := int32(ni.Int64)
	return &value, nil
}

// NullFloat encapsulates sql null float with custom marshalling/unmarshalling
type NullFloat struct {
	sql.NullFloat64
}

// MarshalJSON marshals models.NullFloat datatype
func (nf NullFloat) MarshalJSON() ([]byte, error) {
	if !nf.Valid || nf.Float64 == 0.0 {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON unmarshal models.NullFloat datatype
func (nf *NullFloat) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = err == nil
	if string(b) == "null" {
		nf.Valid = false
	}
	return err
}

// Scan scans numbers to NullFloat datatype
func (nf *NullFloat) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nf = NullFloat{}
		return nil
	}

	var i sql.NullFloat64
	if err := i.Scan(value); err != nil {
		return err
	}

	if math.IsNaN(i.Float64) {
		nf.Float64 = 0.
		nf.Valid = false
		return nil
	}

	i.Float64 = float64(int(i.Float64*math.Pow10(4))) / math.Pow10(4)
	nf.Float64 = i.Float64
	if i.Float64 == 0 {
		nf.Valid = false
	} else {
		nf.Valid = true
	}
	return nil
}

// IsEmpty returns true if models.NullFloat is either not valid or empty
func (nf NullFloat) IsEmpty() bool {
	if !nf.Valid || nf.Float64 == 0 {
		return true
	}
	return false
}

// IsEmpty returns true if models.NullInt is either not valid or empty
func (ni NullInt) IsEmpty() bool {
	if !ni.Valid || ni.Int64 == 0 {
		return true
	}
	return false
}

// MapToFloat64 returns NullFloat value pointer if valid, nil otherwise
func (nf NullFloat) MapToFloat64() *float64 {
	if !nf.Valid {
		return nil
	}

	value := nf.Float64
	return &value
}

// MapToFloat32 returns NullFloat value pointer if valid and in range, nil otherwise
func (nf NullFloat) MapToFloat32() (*float32, error) {
	if !nf.Valid {
		return nil, nil
	}

	if nf.Float64 > math.MaxFloat32 {
		return nil, fmt.Errorf("could not convert nullFloat to *float32, value %f too high", nf.Float64)
	}
	value := float32(nf.Float64)
	return &value, nil
}

// NullString encapsulates sql null string with custom marshalling/unmarshalling
type NullString struct {
	sql.NullString
}

// MarshalJSON marshals models.NullString datatype
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid || ns.String == "" {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON unmarshal models.NullString datatype
func (ns *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == "" {
		ns.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

// Scan scans any variable types to models.NullString datatype
func (ns *NullString) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*ns = NullString{}
	}

	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}

	ns.NullString = i
	if i.Valid && i.String == "" {
		ns.Valid = false
	}
	return nil
}

// IsEmpty returns true if models.NullString is either not valid or empty
func (ns NullString) IsEmpty() bool {
	if !ns.Valid || ns.String == "" {
		return true
	}
	return false
}

// MapToString returns NullString value pointer if valid, nil otherwise
func (ns NullString) MapToString() *string {
	if !ns.Valid {
		return nil
	}

	value := ns.String
	return &value
}

// NullEmptyString encapsulates sql null string with custom marshalling/unmarshalling to allow empty string
type NullEmptyString struct {
	sql.NullString
}

// MarshalJSON marshals models.NullEmptyString datatype
func (ns NullEmptyString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal("")
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON unmarshals models.NullEmptyString datatype
func (ns *NullEmptyString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

// NullTime encapsulates sql null time with custom marshalling/unmarshalling
type NullTime struct {
	sql.NullTime
}

// MarshalJSON marshals models.NullTime datatype
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid || nt.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time.Format("2006-01-02"))
}

// UnmarshalJSON unmarshal models.NullTime datatype
// When using Unmarshal method for null time, layout must be either "YYYY-MM-DD" or RFC3339
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}

	var err error
	nt.Time, err = time.Parse(`"2006-01-02"`, string(b))
	if err == nil {
		nt.Valid = true
		return nil
	}

	/* Try RFC3339 format layout */
	nt.Time, err = time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err == nil {
		nt.Valid = true
	}
	return err
}

// Scan method scans time.Time value from database and forces timezone to UTC
// Do not use NullTime if you want to scan time from database with database's timezone
func (nt *NullTime) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NullTime{}
	}

	var i sql.NullTime
	if err := i.Scan(value); err != nil {
		return err
	}

	zone, _ := i.Time.Zone()
	if zone != "UTC" {
		nt.NullTime.Time = time.Date(i.Time.Year(), i.Time.Month(), i.Time.Day(), 0, 0, 0, 0, time.UTC)
		nt.NullTime.Valid = true
	} else {
		nt.NullTime = i
	}

	if nt.Time.IsZero() {
		nt.Valid = false
	}

	return nil
}

// MapToString returns NullTime string value pointer (RFC3339) if valid , nil otherwise
func (nt NullTime) MapToString() *string {
	if !nt.Valid {
		return nil
	}

	value := nt.Time.Format(time.RFC3339)
	return &value
}

// AfterOrEqual returns true if variable is equal or after arg.
func (nt NullTime) AfterOrEqual(t NullTime) bool {
	if !nt.Time.IsZero() && t.Time.IsZero() {
		return false
	}
	return nt.Time == t.Time || nt.Time.After(t.Time)
}

// BeforeOrEqual returns true if variable is equal or before arg.
func (nt NullTime) BeforeOrEqual(t NullTime) bool {
	if nt.Time.IsZero() && !t.Time.IsZero() {
		return false
	}
	if !nt.Time.IsZero() && t.Time.IsZero() {
		return true
	}

	return nt.Time == t.Time || nt.Time.Before(t.Time)
}

// JSONNullInt64 encapsulates sql null int with marshalling/unmarshalling
type JSONNullInt64 struct {
	sql.NullInt64
}

// MarshalJSON marshals models.JSONNullInt64 datatype
func (v JSONNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON unmarshals JSONNullInt64 datatype
func (v *JSONNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

// MapForRequest converts JSONNullInt64 to value if valid, nil otherwise
func (v JSONNullInt64) MapForRequest() interface{} {
	if v.Valid {
		return v.Int64
	}
	return nil
}

// JSONNullFloat64 encapsulates sql null float with marshalling/unmarshalling
type JSONNullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON marshals JSONNullFloat64 datatype
func (v JSONNullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON unmarshals JSONNullFloat64 datatype
func (v *JSONNullFloat64) UnmarshalJSON(data []byte) error {
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

// MapForRequest converts JSONNullFloat64 to value if valid, nil otherwise
func (v JSONNullFloat64) MapForRequest() interface{} {
	if v.Valid {
		return v.Float64
	}
	return nil
}
