// Package util provides utility functions for image processing and color manipulation.
package util

import (
	"fmt"
	"image"
	"image/color"
	"sort"
)

// Color extraction constants
const (
	MaxUniqueColors = 256 // Maximum number of unique colors to extract
)

// GetImageColors extracts all unique colors from an image.
// Returns a map where keys are colors and values are empty structs (for set behavior).
// Returns an error if the image is nil or invalid.
func GetImageColors(img image.Image) (map[color.Color]struct{}, error) {
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return make(map[color.Color]struct{}), nil
	}

	colors := make(map[color.Color]struct{})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			colors[c] = struct{}{}

			// Prevent excessive memory usage
			if len(colors) > MaxUniqueColors {
				return colors, fmt.Errorf("image has more than %d unique colors", MaxUniqueColors)
			}
		}
	}

	return colors, nil
}

// GetImageColorsLimited extracts up to maxColors unique colors from an image.
// Returns a slice of colors, limited to the specified maximum.
func GetImageColorsLimited(img image.Image, maxColors int) ([]color.Color, error) {
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	if maxColors <= 0 {
		return nil, fmt.Errorf("maxColors must be positive, got %d", maxColors)
	}

	colorMap, err := GetImageColors(img)
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	colors := make([]color.Color, 0, len(colorMap))
	for c := range colorMap {
		if len(colors) >= maxColors {
			break
		}
		colors = append(colors, c)
	}

	return colors, nil
}

// GetDominantColors extracts the most common colors from an image.
// Returns up to maxColors colors, sorted by frequency.
func GetDominantColors(img image.Image, maxColors int) ([]color.Color, error) {
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	if maxColors <= 0 {
		return nil, fmt.Errorf("maxColors must be positive, got %d", maxColors)
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return []color.Color{}, nil
	}

	// Count color occurrences
	colorCounts := make(map[color.Color]int)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			colorCounts[c]++
		}
	}

	// Convert to slice for sorting
	type colorCount struct {
		color color.Color
		count int
	}

	counts := make([]colorCount, 0, len(colorCounts))
	for c, count := range colorCounts {
		counts = append(counts, colorCount{c, count})
	}

	// Sort by frequency (descending)
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	// Extract top colors
	limit := maxColors
	if limit > len(counts) {
		limit = len(counts)
	}

	colors := make([]color.Color, limit)
	for i := 0; i < limit; i++ {
		colors[i] = counts[i].color
	}

	return colors, nil
}

// ColorToRGBA converts a color.Color to RGBA values.
// Returns 8-bit RGB and alpha values.
func ColorToRGBA(c color.Color) (r, g, b, a uint8) {
	if c == nil {
		return 0, 0, 0, 0
	}

	r32, g32, b32, a32 := c.RGBA()
	// RGBA() returns 16-bit values, convert to 8-bit
	r = uint8(r32 >> 8)
	g = uint8(g32 >> 8)
	b = uint8(b32 >> 8)
	a = uint8(a32 >> 8)
	return
}

// ColorToHex converts a color to a hex string.
// Returns format "#RRGGBB" or "#RRGGBBAA" if alpha < 255.
func ColorToHex(c color.Color) string {
	if c == nil {
		return "#000000"
	}

	r, g, b, a := ColorToRGBA(c)

	if a == 255 {
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a)
}

// HexToColor converts a hex string to a color.
// Supports formats: "#RGB", "#RRGGBB", "#RRGGBBAA".
func HexToColor(hex string) (color.Color, error) {
	if len(hex) == 0 {
		return nil, fmt.Errorf("hex string cannot be empty")
	}

	// Remove # prefix if present
	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r, g, b, a uint8 = 0, 0, 0, 255

	switch len(hex) {
	case 3: // #RGB
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r *= 17 // Convert 0-15 to 0-255
		g *= 17
		b *= 17
	case 6: // #RRGGBB
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	case 8: // #RRGGBBAA
		fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
	default:
		return nil, fmt.Errorf("invalid hex color format: %s", hex)
	}

	return color.NRGBA{R: r, G: g, B: b, A: a}, nil
}

// ColorsEqual compares two colors for equality.
// Uses RGBA values for comparison.
func ColorsEqual(c1, c2 color.Color) bool {
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// GetImageSize returns the width and height of an image.
func GetImageSize(img image.Image) (width, height int, err error) {
	if img == nil {
		return 0, 0, fmt.Errorf("image cannot be nil")
	}

	bounds := img.Bounds()
	width = bounds.Dx()
	height = bounds.Dy()

	return width, height, nil
}

// IsValidImageSize checks if the given dimensions are valid for an image.
func IsValidImageSize(width, height int) bool {
	const (
		minSize = 1
		maxSize = 4096
	)
	return width >= minSize && width <= maxSize &&
		height >= minSize && height <= maxSize
}

// ClampInt clamps a value between min and max.
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat32 clamps a float32 value between min and max.
func ClampFloat32(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
