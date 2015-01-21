package vcf

import (
	"strconv"
	"strings"
)

func infoToMap(info string) map[string]interface{} {
	infoMap := make(map[string]interface{})
	fields := strings.Split(info, ";")
	for _, field := range fields {
		if strings.Contains(field, "=") {
			split := strings.Split(field, "=")
			fieldName, fieldValue := split[0], split[1]
			if fieldName == "AC" || fieldName == "AF" {
				infoMap[fieldName] = strings.Split(fieldValue, ",")
			} else {
				infoMap[fieldName] = fieldValue
			}
		} else {
			infoMap[field] = true
		}
	}
	return infoMap
}

func buildInfoSubFields(variant *Variant) {
	info := variant.Info
	variant.Depth = parseIntFromInfoMap("DP", info)
	variant.AlleleCount = parseIntSliceFromInfoMap("AC", info)
	variant.AlleleFrequency = parseFloatSliceFromInfoMap("AF", info)
	variant.AncestralAllele = parseStringFromInfoMap("AA", info)
	variant.TotalAlleles = parseIntFromInfoMap("AN", info)
	variant.End = parseIntFromInfoMap("END", info)
	variant.MAPQ0Reads = parseIntFromInfoMap("MQ0", info)
	variant.NumberOfSamples = parseIntFromInfoMap("NS", info)
	variant.MappingQuality = parseFloatFromInfoMap("MQ", info)
	variant.Cigar = parseStringFromInfoMap("CIGAR", info)
	variant.InDBSNP = parseBoolFromInfoMap("DB", info)
	variant.InHapmap2 = parseBoolFromInfoMap("H2", info)
	variant.InHapmap3 = parseBoolFromInfoMap("H3", info)
	variant.IsSomatic = parseBoolFromInfoMap("SOMATIC", info)
	variant.IsValidated = parseBoolFromInfoMap("VALIDATED", info)
	variant.In1000G = parseBoolFromInfoMap("1000G", info)
	variant.BaseQuality = parseFloatFromInfoMap("BQ", info)
	variant.StrandBias = parseFloatFromInfoMap("SB", info)
}

func parseIntFromInfoMap(key string, info map[string]interface{}) *int {
	if value, found := info[key]; found {
		if str, ok := value.(string); ok {
			intvalue, err := strconv.Atoi(str)
			if err == nil {
				return &intvalue
			}
		}
	}
	return nil
}

func parseIntSliceFromInfoMap(key string, info map[string]interface{}) []int {
	if valueSlice, found := info[key]; found {
		strValueSlice := valueSlice.([]string)
		returnSlice := make([]int, len(strValueSlice))
		for i, value := range strValueSlice {
			convertedValue, err := strconv.Atoi(value)
			if err == nil {
				returnSlice[i] = convertedValue
			}
		}
		return returnSlice
	}

	return nil
}

func parseStringFromInfoMap(key string, info map[string]interface{}) *string {
	if value, found := info[key]; found {
		if str, ok := value.(string); ok {
			return &str
		}
	}
	return nil
}

func parseFloatFromInfoMap(key string, info map[string]interface{}) *float64 {
	if value, found := info[key]; found {
		if str, ok := value.(string); ok {
			convertedValue, err := strconv.ParseFloat(str, 64)
			if err == nil {
				return &convertedValue
			}
		}
	}
	return nil
}

func parseFloatSliceFromInfoMap(key string, info map[string]interface{}) []float64 {
	if valueSlice, found := info[key]; found {
		strValueSlice := valueSlice.([]string)
		returnSlice := make([]float64, len(strValueSlice))
		for i, value := range strValueSlice {
			convertedValue, err := strconv.ParseFloat(value, 64)
			if err == nil {
				returnSlice[i] = convertedValue
			}
		}
		return returnSlice
	}

	return nil
}

func parseBoolFromInfoMap(key string, info map[string]interface{}) *bool {
	if value, found := info[key]; found {
		if b, ok := value.(bool); ok {
			return &b
		}
	}
	return nil
}

func splitMultipleAltInfos(info map[string]interface{}, numberOfAlternatives int) []map[string]interface{} {
	maps := make([]map[string]interface{}, 0, 2)
	separator := ","

	for key, v := range info {
		if value, ok := v.(string); ok {
			if strings.Contains(value, separator) {
				alternatives := strings.Split(value, separator)
				for position, alt := range alternatives {
					maps = insertMapSlice(maps, position, key, alt)
				}
			} else {
				for i := 0; i < numberOfAlternatives; i++ {
					maps = insertMapSlice(maps, i, key, value)
				}
			}
		} else {
			maps = insertMapSlice(maps, 0, key, v)
		}
	}

	return maps
}

func insertMapSlice(maps []map[string]interface{}, position int, key string, alt interface{}) []map[string]interface{} {
	if len(maps) <= position {
		for i := len(maps); i <= position; i++ {
			maps = append(maps, make(map[string]interface{}))
		}
	}
	maps[position][key] = alt
	return maps
}
