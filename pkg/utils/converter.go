package converter

func ConvertStorePrice(value *float64) *int64 {

	if value == nil {
		return nil
	}

	result := int64(*value * 100)

	return &result
}

func ConvertToDisplayPrice(value *int64) *float64 {
	if value == nil {
		return nil
	}

	result := float64(*value) / 100

	return &result
}
