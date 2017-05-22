package ast

func convertSliceMap(i interface{}) []m {
	var m []m
	for _, mm := range i.([]interface{}) {
		m = append(m, convertMap(mm))
	}

	return m
}

func convertMap(i interface{}) m {
	return m(i.(map[string]interface{}))
}

func convertString(i interface{}) string {
	return i.(string)
}

func convertInt(i interface{}) int {
	return int(i.(float64))
}

func convertBool(i interface{}) bool {
	return i.(bool)
}
