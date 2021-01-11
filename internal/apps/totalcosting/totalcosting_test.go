package totalcosting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAverageUnitPrice(t *testing.T) {
	first := Element{
		Cost: 10000,
		Unit: 100,
		Type: First,
	}
	input := Element{
		Cost: 50000,
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
		Cost: 10000,
		Unit: 100,
		Type: First,
	}
	input := Element{
		Cost: 50000,
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
		Cost: 10000,
		Unit: 100,
		Type: First,
	}
	input := Element{
		Cost: 50000,
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
