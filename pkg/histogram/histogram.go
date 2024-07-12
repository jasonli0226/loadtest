package histogram

import (
	"fmt"
	"sort"
	"strings"
)

type Histogram struct {
	bins     []int
	binWidth float64
	maxBins  int
}

func NewHistogram(data []float64, maxBins int) *Histogram {
	if len(data) == 0 {
		return &Histogram{maxBins: maxBins}
	}

	sort.Float64s(data)
	min, max := data[0], data[len(data)-1]
	binWidth := (max - min) / float64(maxBins)

	bins := make([]int, maxBins)
	for _, v := range data {
		binIndex := int((v - min) / binWidth)
		if binIndex >= maxBins {
			binIndex = maxBins - 1
		}
		bins[binIndex]++
	}

	return &Histogram{
		bins:     bins,
		binWidth: binWidth,
		maxBins:  maxBins,
	}
}

func (h *Histogram) String() string {
	if len(h.bins) == 0 {
		return "No data available for histogram"
	}

	maxCount := 0
	for _, count := range h.bins {
		if count > maxCount {
			maxCount = count
		}
	}

	var sb strings.Builder
	sb.WriteString("Response Time Histogram:\n")
	for i, count := range h.bins {
		barLength := int(float64(count) / float64(maxCount) * 50)
		bar := strings.Repeat("#", barLength)
		sb.WriteString(fmt.Sprintf("%6.2f ms - %6.2f ms | %s (%d)\n",
			float64(i)*h.binWidth, float64(i+1)*h.binWidth, bar, count))
	}
	return sb.String()
}
