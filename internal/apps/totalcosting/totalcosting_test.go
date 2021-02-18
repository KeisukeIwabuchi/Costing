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

func TestCost(t *testing.T) {
	testCases := []struct {
		E      Element
		Result float64
	}{
		{Element{Price: 100.0, Unit: 100}, 10000.0},
		{Element{Price: 100.0, Unit: 0}, 0.0},
		{Element{Price: 0.0, Unit: 100}, 0.0},
	}

	for _, testCase := range testCases {
		result := testCase.E.Cost()
		if result != testCase.Result {
			t.Errorf("Invalid result. testCase:%#v, actual:%f", testCase, result)
		}
	}
}

func TestAddCost(t *testing.T) {
	testCases := []struct {
		E        Element
		Argument float64
		Result   float64
	}{
		{Element{Price: 100.0, Unit: 100}, 10000.0, 200.0},
		{Element{Price: 0.0, Unit: 100}, 10000.0, 100.0},
	}

	for _, testCase := range testCases {
		testCase.E.AddCost(testCase.Argument)
		result := testCase.E.Price
		if result != testCase.Result {
			t.Errorf("Invalid result. testCase:%#v, actual:%f", testCase, result)
		}
	}
}

func TestIsBear(t *testing.T) {
	testCases := []struct {
		E        Element
		Argument float64
		Result   bool
	}{
		{Element{Progress: 0.5}, 1.0, false},
		{Element{Progress: 0.5}, 0.5, true},
		{Element{Progress: 0.5}, 0.2, true},
	}

	for _, testCase := range testCases {
		result := testCase.E.IsBear(testCase.Argument)
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

	expected := []int{
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

	expected := []int{
		180,
		1332,
		1440,
		72,
	}
	for i, e := range proccesing.Elements {
		actual := e.Unit
		assert.Equal(t, expected[i], actual)
	}

	expectedType := []ElementType{
		First,
		Input,
		Output,
		Last,
	}
	for i, e := range proccesing.Elements {
		actual := e.Type
		assert.Equal(t, expectedType[i], actual)
	}

	expectedProgress := []float64{
		0.6,
		0.0,
		1.0,
		0.3,
	}
	for i, e := range proccesing.Elements {
		actual := e.Progress
		assert.Equal(t, expectedProgress[i], actual)
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

func TestGetFIFOOutputBurder(t *testing.T) {
	first := Element{
		Type: First,
		Unit: 300,
	}
	input := Element{
		Type: Input,
		Unit: 1380,
	}
	output := Element{
		Type: Output,
		Unit: 1200,
	}
	normalDefect := Element{
		Type: NormalDefect,
		Unit: 240,
	}
	last := Element{
		Type: Last,
		Unit: 240,
	}

	var material Cost
	material.Elements = append(material.Elements, first)
	material.Elements = append(material.Elements, input)
	material.Elements = append(material.Elements, output)
	material.Elements = append(material.Elements, normalDefect)
	material.Elements = append(material.Elements, last)

	actual := material.GetFIFOOutputBurder()
	expected := 900
	assert.Equal(t, expected, actual)
}

func TestGetNormalDefectUnit(t *testing.T) {
	first := Element{
		Type: First,
		Unit: 300,
	}
	input := Element{
		Type: Input,
		Unit: 1380,
	}
	output := Element{
		Type: Output,
		Unit: 1320,
	}
	normalDefect := Element{
		Type: NormalDefect,
		Unit: 120,
	}
	last := Element{
		Type: Last,
		Unit: 240,
	}

	var material Cost
	material.Elements = append(material.Elements, first)
	material.Elements = append(material.Elements, input)
	material.Elements = append(material.Elements, output)
	material.Elements = append(material.Elements, normalDefect)
	material.Elements = append(material.Elements, last)

	actual := material.GetNormalDefectUnit()
	expected := 120
	assert.Equal(t, expected, actual)
}

func TestGetTotalNDBurden(t *testing.T) {
	first := Element{
		Type:     First,
		Unit:     300,
		NDBurden: 0,
	}
	input := Element{
		Type:     Input,
		Unit:     1380,
		NDBurden: 0,
	}
	output := Element{
		Type:     Output,
		Unit:     1320,
		NDBurden: 1020,
	}
	normalDefect := Element{
		Type:     NormalDefect,
		Unit:     120,
		NDBurden: 0,
	}
	last := Element{
		Type:     Last,
		Unit:     240,
		NDBurden: 240,
	}

	var material Cost
	material.Elements = append(material.Elements, first)
	material.Elements = append(material.Elements, input)
	material.Elements = append(material.Elements, output)
	material.Elements = append(material.Elements, normalDefect)
	material.Elements = append(material.Elements, last)

	actual := material.GetTotalNDBurden()
	expected := 1260
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

	actual := box.CalculationEOFMCost()
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

	actual := box.CalculationProductCost()
	expected := 1872000.0
	assert.Equal(t, expected, actual)
}

func TestCalculationProdutAvgCost(t *testing.T) {
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
	box.ProductTotalCost = box.CalculationProductCost()

	actual := box.CalculationProductAvgCost()
	expected := 1300.0
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
