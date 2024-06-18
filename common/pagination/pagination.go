package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Default returns a default pagination with offset 0 and limit DefaultLimit500
func Default() Pagination {
	return Pagination{
		Offset: 0,
		Limit:  DefaultLimit500,
	}
}

// GetFromURLQuery gets page (limit and offset) from url query
func GetFromURLQuery(c *gin.Context) (Pagination, error) {
	key := "offset"
	offset, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(DefaultOffset)))
	if err != nil {
		return Pagination{}, BadRequestValueError{Key: key, Err: err}
	}

	if offset < 0 {
		return Pagination{}, BadRequestValueError{Key: key, Err: fmt.Errorf("offset (%d) cannot be negative", offset)}
	}

	key = "limit"
	limit, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(DefaultLimit500)))
	if err != nil {
		return Pagination{}, BadRequestValueError{Key: key, Err: err}
	}

	if limit < 0 {
		return Pagination{}, BadRequestValueError{Key: key, Err: fmt.Errorf("limit (%d) cannot be negative", limit)}
	}

	return Pagination{Offset: offset, Limit: limit}, nil
}

// BuildPageable builds pageable from entity
func BuildPageable[T any](page Pagination, total int64, data []T) Pageable {
	return Pageable{
		Limit:  page.Limit,
		Offset: page.Offset,
		Total:  total,
		Data:   data,
	}
}

// PageableToSlice casts Data interface field (from json Unmarshalling of Pageable) to slice of type T
func PageableToSlice[T any](pageable Pageable) ([]T, error) {
	sliceInterface, ok := pageable.Data.([]interface{})
	if !ok {
		return []T{}, errors.New("unable to cast data field of pageable to a slice of interface")
	}

	var data []T
	for _, value := range sliceInterface {
		var sample T
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return []T{}, fmt.Errorf("unable to encode JSON elements for %T datatype", sample)
		}
		err = json.Unmarshal(jsonBytes, &sample)
		if err != nil {
			return []T{}, fmt.Errorf("unable to encode JSON elements for %T datatype", sample)
		}
		data = append(data, sample)
	}
	return data, nil
}
