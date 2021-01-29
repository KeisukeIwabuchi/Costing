package totalcosting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLeftElement(t *testing.T) {
	testCases := []struct {
		E      Element
		Result bool
	}{
		{Element{Type: First}, true},
		{Element{Type: Input}, true},
		{Element{Type: Output}, false},
		{Element{Type: Last}, false},
		{Element{Type: NormalDefect}, false},
		{Element{Type: AbnormalDefect}, false},
		{Element{Type: NormalImpairment}, false},
		{Element{Type: AbnormalImpairment}, false},
	}

	for _, testCase := range testCases {
		result := testCase.E.IsLeftElement()
		if result != testCase.Result {
			t.Errorf("Invalid result. testCase:%#v, actual:%t", testCase, result)
		}
	}
}

func TestCalulateInputUnit(t *testing.T) {
	first := Element{
		Type:     First,
		Unit:     300,
		Progress: 0.6,
	}
	input := Element{
		Type: Input,
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

	var master []Element
	master = append(master, first)
	master = append(master, input)
	master = append(master, output)
	master = append(master, last)

	var material Cost
	material.InputOnAvg = false
	material.InputTiming = 0.0

	material.CalulateInputUnit(master)

	expected := []float64{
		300,
		1380,
		1440,
		240,
	}
	for i, e := range material.Elements {
		actual := e.Unit
		assert.Equal(t, expected[i], actual)
	}
}

func TestCalulateConversionUnit(t *testing.T) {
	first := Element{
		Type:     First,
		Unit:     300,
		Progress: 0.6,
	}
	input := Element{
		Type: Input,
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

	var master []Element
	master = append(master, first)
	master = append(master, input)
	master = append(master, output)
	master = append(master, last)

	var proccesing Cost
	proccesing.InputOnAvg = true

	proccesing.CalulateConversionUnit(master)

	expected := []float64{
		180,
		1332,
		1440,
		72,
	}
	for i, e := range proccesing.Elements {
		actual := e.Unit
		assert.Equal(t, expected[i], actual)
	}
}

func TestGetPriceFIFO(t *testing.T) {
	input := Element{
		Type: Input,
		Unit: 1380,
	}

	var material Cost
	material.Elements = append(material.Elements, input)
	material.FirstCost = 206400
	material.InputCost = 717600

	actual := material.GetPriceFIFO()
	expected := 520.0
	assert.Equal(t, expected, actual)
}

func TestGetPriceAVG(t *testing.T) {
	first := Element{
		Type:     First,
		Unit:     300,
		Progress: 0.6,
	}
	input := Element{
		Type: Input,
		Unit: 1380,
	}

	var material Cost
	material.Elements = append(material.Elements, first)
	material.Elements = append(material.Elements, input)
	material.FirstCost = 206400
	material.InputCost = 717600

	actual := material.GetPriceAVG()
	expected := 550.0
	assert.Equal(t, expected, actual)
}

func TestCalculationEOFMCost(t *testing.T) {
	materialLast := Element{
		Type:  Last,
		Price: 550.0,
		Unit:  240,
	}
	processingLast := Element{
		Type:  Last,
		Price: 750.0,
		Unit:  72,
	}

	var material, processing Cost
	material.Elements = append(material.Elements, materialLast)
	processing.Elements = append(processing.Elements, processingLast)

	var costs []Cost
	costs = append(costs, material)
	costs = append(costs, processing)

	var box Box
	box.Costs = append(box.Costs, material)
	box.Costs = append(box.Costs, processing)
	box.CalculationEOFMCost()

	actual := box.EOTMTotalCost
	expected := 186000.0
	assert.Equal(t, expected, actual)
}

func TestCalculationProdutCost(t *testing.T) {
	materialProduct := Element{
		Type:  Output,
		Price: 550.0,
		Unit:  1440,
	}
	processingProduct := Element{
		Type:  Output,
		Price: 750.0,
		Unit:  1440,
	}

	var material, processing Cost
	material.Elements = append(material.Elements, materialProduct)
	processing.Elements = append(processing.Elements, processingProduct)

	var costs []Cost
	costs = append(costs, material)
	costs = append(costs, processing)

	var box Box
	box.Master = append(box.Master, materialProduct)
	box.Costs = append(box.Costs, material)
	box.Costs = append(box.Costs, processing)
	box.CalculationProductCost()

	actual := box.ProductTotalCost
	expected := 1872000.0
	assert.Equal(t, expected, actual)

	actual = box.ProductAvgCost
	expected = 1300.0
	assert.Equal(t, expected, actual)
}

func TestRun(t *testing.T) {
	first := Element{
		Type:     First,
		Unit:     300,
		Progress: 0.6,
	}
	input := Element{
		Type: Input,
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

	var box Box

	box.Master = append(box.Master, first)
	box.Master = append(box.Master, input)
	box.Master = append(box.Master, output)
	box.Master = append(box.Master, last)

	var material, processing Cost
	material.InputOnAvg = false
	material.InputTiming = 0.0
	material.CMethod = AVG
	material.FirstCost = 206400
	material.InputCost = 717600

	processing.InputOnAvg = true
	processing.CMethod = AVG
	processing.FirstCost = 161640
	processing.InputCost = 972360

	box.Costs = append(box.Costs, material)
	box.Costs = append(box.Costs, processing)

	box.Run()

	actual := box.ProductAvgCost
	expected := 1300.0
	assert.Equal(t, expected, actual)
}
