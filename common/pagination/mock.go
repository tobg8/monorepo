package pagination

// MockPageableLabel creates a valid pageable label
func MockPageableLabel(labels ...string) Pageable {
	var pageable Pageable
	pageable.Offset = 0
	pageable.Limit = 999999
	var data []interface{}
	for _, value := range labels {
		newData := map[string]interface{}{
			"label": value,
		}
		data = append(data, newData)
	}
	pageable.Total = int64(len(data))
	pageable.Data = data
	return pageable
}
