package totalcosting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAverageUnitPrice(t *testing.T) {
	first := Element{
		Cost: []float64{10000},
		Unit: 100,
		Type: First,
	}
	input := Element{
		Cost: []float64{50000},
		Unit: 500,
		Type: Input,
	}

	var elements Elements
	elements = append(elements, first)
	elements = append(elements, input)

	actual := CalculateAverageUnitPrice(elements, Input)
	expected := 100.0
	assert.Equal(t, expected, actual)

	actual = CalculateAverageUnitPrice(elements, First, Input)
	expected = 100.0
	assert.Equal(t, expected, actual)
}

func TestUnitPriceWithFIFO(t *testing.T) {
	first := Element{
		Cost: []float64{10000},
		Unit: 100,
		Type: First,
	}
	input := Element{
		Cost: []float64{50000},
		Unit: 500,
		Type: Input,
	}

	var elements Elements
	elements = append(elements, first)
	elements = append(elements, input)

	var box BOX
	box.Left = elements

	actual := box.UnitPriceWithFIFO()
	expected := CalculateAverageUnitPrice(elements, Input)
	assert.Equal(t, expected, actual)
}

func TestUnitPriceWithAVG(t *testing.T) {
	first := Element{
		Cost: []float64{10000},
		Unit: 100,
		Type: First,
	}
	input := Element{
		Cost: []float64{50000},
		Unit: 500,
		Type: Input,
	}

	var elements Elements
	elements = append(elements, first)
	elements = append(elements, input)

	var box BOX
	box.Left = elements

	actual := box.UnitPriceWithAVG()
	expected := CalculateAverageUnitPrice(elements, First, Input)
	assert.Equal(t, expected, actual)
}

func TestRun(t *testing.T) {
	first := Element{
		Type:     First,
		Cost:     []float64{206400, 161640},
		Unit:     300,
		Progress: 0.6,
	}
	input := Element{
		Type: Input,
		Cost: []float64{717600, 972360},
		Unit: 1380,
	}
	output := Element{
		Type: Output,
		Unit: 1440,
	}
	last := Element{
		Type:     Last,
		Unit:     240,
		Progress: 0.3,
	}

	var left, right Elements
	left = append(left, first)
	left = append(left, input)
	right = append(right, output)
	right = append(right, last)

	var box BOX
	box.Left = left
	box.Right = right
	box.CMethod = AVG

	box.Run()

	actual := box.ProductAvgPrice
	expected := 1300.0
	assert.Equal(t, expected, actual)
}
