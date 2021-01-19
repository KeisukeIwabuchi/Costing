package totalcosting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestCalculateAverageUnitPrice(t *testing.T) {
// 	first := Element{
// 		Cost: []float64{10000},
// 		Unit: 100,
// 		Type: First,
// 	}
// 	input := Element{
// 		Cost: []float64{50000},
// 		Unit: 500,
// 		Type: Input,
// 	}

// 	var elements Elements
// 	elements = append(elements, first)
// 	elements = append(elements, input)

// 	actual := CalculateAverageUnitPrice(elements, Input)
// 	expected := 100.0
// 	assert.Equal(t, expected, actual)

// 	actual = CalculateAverageUnitPrice(elements, First, Input)
// 	expected = 100.0
// 	assert.Equal(t, expected, actual)
// }

// func TestUnitPriceWithFIFO(t *testing.T) {
// 	first := Element{
// 		Cost: []float64{10000},
// 		Unit: 100,
// 		Type: First,
// 	}
// 	input := Element{
// 		Cost: []float64{50000},
// 		Unit: 500,
// 		Type: Input,
// 	}

// 	var elements Elements
// 	elements = append(elements, first)
// 	elements = append(elements, input)

// 	var box BOX
// 	box.Left = elements

// 	actual := box.UnitPriceWithFIFO()
// 	expected := CalculateAverageUnitPrice(elements, Input)
// 	assert.Equal(t, expected, actual)
// }

// func TestUnitPriceWithAVG(t *testing.T) {
// 	first := Element{
// 		Cost: []float64{10000},
// 		Unit: 100,
// 		Type: First,
// 	}
// 	input := Element{
// 		Cost: []float64{50000},
// 		Unit: 500,
// 		Type: Input,
// 	}

// 	var elements Elements
// 	elements = append(elements, first)
// 	elements = append(elements, input)

// 	var box BOX
// 	box.Left = elements

// 	actual := box.UnitPriceWithAVG()
// 	expected := CalculateAverageUnitPrice(elements, First, Input)
// 	assert.Equal(t, expected, actual)
// }

func TestRun(t *testing.T) {
	first := LeftElement{
		Type: First,
		Value: Element{
			Unit:     300,
			Progress: 0.6,
		},
	}
	input := LeftElement{
		Type: Input,
		Value: Element{
			Unit: 1380,
		},
	}
	output := RightElement{
		Type: Output,
		Value: Element{
			Unit: 1440,
		},
	}
	last := RightElement{
		Type: Last,
		Value: Element{
			Unit:     240,
			Progress: 0.3,
		},
	}

	var materialBox, processingBox Box
	materialBox.InputTiming = 0.0
	materialBox.Left = append(materialBox.Left, first)
	materialBox.Left = append(materialBox.Left, input)
	materialBox.Right = append(materialBox.Right, output)
	materialBox.Right = append(materialBox.Right, last)
	materialBox.CMethod = AVG
	materialBox.FirstCost = 206400
	materialBox.InputCost = 717600

	processingBox = materialBox
	processingBox.FirstCost = 161640
	processingBox.InputCost = 972360

	var costing TotalCosting
	costing.Boxes = append(costing.Boxes, materialBox)
	costing.Run()

	actual := costing.ProductAvgCost
	expected := 1300.0
	assert.Equal(t, expected, actual)
}
