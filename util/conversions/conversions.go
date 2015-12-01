package conversions

import (
	"fmt"
	"strconv"
	"strings"
)

func IntForValue(value interface{}) (out int, err error) {
	switch value := value.(type) { // shadow variable
	case string:
		i, err := strconv.ParseFloat(value, 32)
		if err == nil {
			out = int(i)
		} else {
			fmt.Println(err)
			out = -1
		}

	case float32:
		out = int(value)
	case float64:
		out = int(value)
	case int:
		out = int(value)
	case bool:
		out = 0
		if value {
			out = 1
		}
	default:
		return 0, fmt.Errorf("float32ForValue(): %v is not a "+
			"number or Boolean", value)
	}
	return out, nil
}

func Float32ForValue(value interface{}) (out float32, err error) {
	switch value := value.(type) { // shadow variable
	case string:
		i, err := strconv.ParseFloat(value, 32)
		if err == nil {
			out = float32(i)
		} else {
			fmt.Println(err)
			out = -1
		}
	case float32:
		out = value
	case float64:
		out = float32(value)
	case int:
		out = float32(value)
	case bool:
		out = 0
		if value {
			out = 1
		}
	default:
		return 0, fmt.Errorf("float32ForValue(): %v is not a "+"number or Boolean", value)
	}
	return out, nil
}

func Float64ForValue(value interface{}) (out float64, err error) {
	switch value := value.(type) { // shadow variable
	case string:
		i, err := strconv.ParseFloat(value, 32)
		if err == nil {
			out = float64(i)
		} else {
			fmt.Println(err)
			out = -1
		}
	case float32:
		out = float64(value)
	case float64:
		out = value
	case int:
		out = float64(value)
	case bool:
		out = 0
		if value {
			out = 1
		}
	default:
		return 0, fmt.Errorf("float32ForValue(): %v is not a "+"number or Boolean", value)
	}
	return out, nil
}

func String2Array16(value string, atype string, out []int16) (err error) {
	separator := string(",")
	str := strings.Replace(value, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	imageParts := strings.Split(str, separator)
	i := 0
	for _, voxel := range imageParts {
		tmp, _ := strconv.Atoi(voxel)
		out[i] = int16(tmp)
		i++
	}
	return
}

func String2Array32(value string, atype string, out []int32) (err error) {
	separator := string(",")
	str := strings.Replace(value, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	imageParts := strings.Split(str, separator)
	i := 0
	for _, voxel := range imageParts {
		tmp, _ := strconv.Atoi(voxel)
		out[i] = int32(tmp)
		i++
	}
	return
}
