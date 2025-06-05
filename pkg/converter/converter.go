package converter

// The NumberConverter interface is defined in vietnamese.go

// NewConverter creates and returns the optimal Vietnamese number converter implementation
// This is the main entry point for applications using this library
func NewConverter() NumberConverter {
	return NewTurboConverter()
}
