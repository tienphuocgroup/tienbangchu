package converter

import (
	"fmt"
	"strings"
	"sync"
)

// TurboVietnameseConverter provides the fastest possible number-to-text conversion
// for Vietnamese with minimal allocations and optimized algorithms
type TurboVietnameseConverter struct {
	// Static maps for faster lookups and reduced allocations
	units      [10]string
	tens       [10]string
	scales     [7]string
	specialMap map[int]string
	
	// Pre-allocated buffer pool to avoid repeated allocations in high-performance scenarios
	bufferPool *sync.Pool
}

// NewTurboConverter creates a new instance of the ultra-optimized Vietnamese converter
func NewTurboConverter() NumberConverter {
	conv := &TurboVietnameseConverter{
		// Using arrays instead of slices to avoid heap allocations
		units: [10]string{
			"", "một", "hai", "ba", "bốn", "năm", "sáu", "bảy", "tám", "chín",
		},
		tens: [10]string{
			"", "mười", "hai mươi", "ba mươi", "bốn mươi", "năm mươi",
			"sáu mươi", "bảy mươi", "tám mươi", "chín mươi",
		},
		scales: [7]string{
			"", "nghìn", "triệu", "tỷ", "nghìn tỷ", "triệu tỷ", "tỷ tỷ",
		},
		// Pre-compute special cases for faster access
		specialMap: map[int]string{
			1: "mốt",   // Special case for "một" in tens position
			4: "tư",    // Special case for "bốn" in tens position
			5: "lăm",   // Special case for "năm" in tens position
		},
		// Create string builder pool for reuse
		bufferPool: &sync.Pool{
			New: func() interface{} {
				// Pre-allocate builders with sufficient capacity
				sb := &strings.Builder{}
				sb.Grow(100)
				return sb
			},
		},
	}
	
	return conv
}

// Convert converts a number to Vietnamese text
func (c *TurboVietnameseConverter) Convert(number int64) (string, error) {
	return c.ConvertWithCurrency(number, "đồng")
}

// ConvertWithCurrency converts a number to Vietnamese text with specified currency
func (c *TurboVietnameseConverter) ConvertWithCurrency(number int64, currency string) (string, error) {
	// Handle validation with pre-checks
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

	// Get a pre-allocated string builder from the pool
	sb := c.bufferPool.Get().(*strings.Builder)
	sb.Reset() // Clear any previous content
	defer func() {
		// Return to pool when done
		c.bufferPool.Put(sb)
	}()
	
	// Direct, stack-based processing of digits
	// This approach avoids both array creation and sorting
	// Using 5 as that's the max needed for 15 digits
	var groups [5]int
	var groupCount int
	
	// Extract groups of 3 digits with direct arithmetic
	temp := number
	for temp > 0 {
		groups[groupCount] = int(temp % 1000)
		temp /= 1000
		groupCount++
	}
	
	// Process each group from highest to lowest without recursion
	firstGroup := true
	for i := groupCount - 1; i >= 0; i-- {
		group := groups[i]
		
		// Skip zero groups unless it's the only group
		if group == 0 {
			if groupCount == 1 {
				sb.WriteString("không")
			}
			continue
		}
		
		// Add space between groups except for the first one
		if !firstGroup {
			sb.WriteRune(' ')
		}
		
		// Convert three-digit group with direct string concat
		c.appendGroup(sb, group, i, firstGroup)
		
		// Add appropriate scale suffix
		if i > 0 {
			sb.WriteRune(' ')
			sb.WriteString(c.scales[i])
		}
		
		firstGroup = false
	}
	
	// Add currency if specified
	if currency != "" {
		sb.WriteRune(' ')
		sb.WriteString(currency)
	}
	
	// Return the result - applying any final normalization
	result := sb.String()
	
	// The only normalization needed in practice is mươi một -> mươi mốt
	// This is more efficient than a full string replacement
	if strings.Contains(result, "mươi một") {
		result = strings.ReplaceAll(result, "mươi một", "mươi mốt")
	}
	
	return result, nil
}

// appendGroup directly appends a 3-digit group conversion to the string builder
func (c *TurboVietnameseConverter) appendGroup(sb *strings.Builder, group int, scale int, isFirst bool) {
	// Split digits for direct access (more efficient than multiple divisions)
	hundreds := group / 100
	remainder := group % 100
	tens := remainder / 10
	units := remainder % 10
	
	// Process hundreds place
	if hundreds > 0 {
		sb.WriteString(c.units[hundreds])
		sb.WriteString(" trăm")
		
		// Only add connective words if needed
		if remainder > 0 {
			sb.WriteRune(' ')
			if tens == 0 {
				// Special case for numbers like x01 to x09
				sb.WriteString("lẻ")
				sb.WriteRune(' ')
				sb.WriteString(c.units[units])
				return
			}
		} else {
			// No remaining digits
			return
		}
	} else if !isFirst && remainder > 0 {
		// Handle cases like x,001 where x > 0
		sb.WriteString("không trăm")
		sb.WriteRune(' ')
		if tens == 0 {
			sb.WriteString("lẻ")
			sb.WriteRune(' ')
			sb.WriteString(c.units[units])
			return
		}
	}
	
	// Process tens place with special cases
	if tens > 1 {
		// 20-99
		sb.WriteString(c.units[tens])
		sb.WriteString(" mươi")
		if units > 0 {
			sb.WriteRune(' ')
			// Special cases handled via map for better performance
			if special, exists := c.specialMap[units]; exists {
				sb.WriteString(special)
			} else {
				sb.WriteString(c.units[units])
			}
		}
	} else if tens == 1 {
		// 10-19
		sb.WriteString("mười")
		if units > 0 {
			sb.WriteRune(' ')
			if units == 5 {
				sb.WriteString("lăm")
			} else {
				sb.WriteString(c.units[units])
			}
		}
	} else {
		// 1-9
		sb.WriteString(c.units[units])
	}
}
