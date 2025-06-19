package converter

import (
	"fmt"
	"strings"
)

type NumberConverter interface {
	Convert(number int64) (string, error)
	ConvertWithCurrency(number int64, currency string) (string, error)
}

type vietnameseConverter struct {
	units     []string
	tens      []string
	scales    []string
	zeroWords map[int]string
}

func NewVietnameseConverter() NumberConverter {
	return &vietnameseConverter{
		units: []string{
			"", "một", "hai", "ba", "bốn", "năm", "sáu", "bảy", "tám", "chín",
		},
		tens: []string{
			"", "mười", "hai mươi", "ba mươi", "bốn mươi", "năm mươi",
			"sáu mươi", "bảy mươi", "tám mươi", "chín mươi",
		},
		scales: []string{
			"", "nghìn", "triệu", "tỷ", "nghìn tỷ", "triệu tỷ", "tỷ tỷ",
		},
		zeroWords: map[int]string{
			1: "lẻ",
			2: "không trăm",
		},
	}
}

func (vc *vietnameseConverter) Convert(number int64) (string, error) {
	return vc.ConvertWithCurrency(number, "đồng")
}

func (vc *vietnameseConverter) ConvertWithCurrency(number int64, currency string) (string, error) {
	if number < 0 {
		return "", fmt.Errorf("negative numbers not supported")
	}

	if number > 999999999999999 {
		return "", fmt.Errorf("number too large (max: 999,999,999,999,999)")
	}

	if number == 0 {
		if currency != "" {
			return "không " + currency, nil
		}
		return "không", nil
	}

	groups := vc.splitIntoGroups(number)

	var parts []string
	groupCount := len(groups)

	// Handle single digit case
	if groupCount == 1 && groups[0] < 10 {
		result := vc.units[groups[0]]
		if currency != "" {
			result += " " + currency
		}
		return result, nil
	}

	for i, group := range groups {
		scaleIndex := groupCount - i - 1
		if group == 0 {
			// Only read zero group if it's the lowest group (units)
			if i != groupCount-1 {
				continue
			}
		}

		groupText := vc.convertThreeDigitGroup(group, scaleIndex, i == 0)

		if groupText != "" && group != 0 {

			if scaleIndex > 0 && scaleIndex < len(vc.scales) {
				groupText += " " + vc.scales[scaleIndex]
			}
			parts = append(parts, groupText)
		} else if group == 0 && i == groupCount-1 {
			// Only add zero for the lowest group if all others are zero
			if len(parts) == 0 {
				parts = append(parts, "không")
			}
		}
	}

	if len(parts) == 0 {
		if currency != "" {
			return "không " + currency, nil
		}
		return "không", nil
	}

	result := strings.Join(parts, " ")
	result = vc.normalizeVietnamese(result)

	if currency != "" {
		result += " " + currency
	}

	return result, nil
}

func (vc *vietnameseConverter) splitIntoGroups(number int64) []int {
	var groups []int
	
	for number > 0 {
		groups = append([]int{int(number % 1000)}, groups...)
		number /= 1000
	}
	
	return groups
}

func (vc *vietnameseConverter) convertThreeDigitGroup(group int, scaleIndex int, isFirst bool) string {
	if group == 0 {
		return ""
	}

	hundreds := group / 100
	remainder := group % 100
	tens := remainder / 10
	units := remainder % 10

	var parts []string

	// Hundreds
	if hundreds > 0 {
		parts = append(parts, vc.units[hundreds]+" trăm")
	} else if !isFirst && (tens > 0 || units > 0) {
		parts = append(parts, "không trăm")
	}

	// Tens/Units
	if tens > 1 {
		parts = append(parts, vc.units[tens]+" mươi")
		if units == 1 {
			parts = append(parts, "mốt")
		} else if units == 4 {
			parts = append(parts, "tư")
		} else if units == 5 {
			parts = append(parts, "lăm")
		} else if units > 0 {
			parts = append(parts, vc.units[units])
		}
	} else if tens == 1 {
		parts = append(parts, "mười")
		if units == 5 {
			parts = append(parts, "lăm")
		} else if units > 0 {
			parts = append(parts, vc.units[units])
		}
	} else if tens == 0 && units > 0 {
		if hundreds > 0 {
			parts = append(parts, "lẻ")
		}
		parts = append(parts, vc.units[units])
	}

	return strings.Join(parts, " ")
}

// Removed: now handled in convertThreeDigitGroup
func (vc *vietnameseConverter) convertTensAndUnits(tens, units, scaleIndex int, hasHundreds bool) string {
	return ""
}

func (vc *vietnameseConverter) getUnitWord(digit int, isStandalone bool, scaleIndex int) string {
	if digit == 0 {
		return ""
	}
	
	if digit == 4 {
		if isStandalone || scaleIndex > 0 {
			return "bốn"
		}
		return "tư"
	}
	
	return vc.units[digit]
}

func (vc *vietnameseConverter) normalizeVietnamese(text string) string {
	words := strings.Fields(text)
	
	var normalized []string
	for i, word := range words {
		if word == "một" && i > 0 && i < len(words)-1 {
			prevWord := words[i-1]
			if strings.HasSuffix(prevWord, "mười") && prevWord != "mười" {
				normalized = append(normalized, "mốt")
				continue
			}
		}
		normalized = append(normalized, word)
	}
	
	return strings.Join(normalized, " ")
}
