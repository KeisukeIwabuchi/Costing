package totalcosting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	var master []Element
	master = append(master, first)
	master = append(master, input)
	master = append(master, output)
	master = append(master, last)

	var material, processing Cost
	material.InputOnAvg = false
	material.InputTiming = 0.0
	material.CMethod = AVG
	material.FirstCost = 206400
	material.InputCost = 717600

	processing.InputOnAvg = true
	processing.FirstCost = 161640
	processing.InputCost = 972360

	var costs []Cost
	costs = append(costs, material)
	costs = append(costs, processing)

	var box Box
	box.Master = master
	box.Costs = costs
	box.Run()

	actual := box.ProductAvgCost
	expected := 1300.0
	assert.Equal(t, expected, actual)
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

	var material Cost
	material.InputOnAvg = false
	material.InputTiming = 0.0
	material.CMethod = AVG
	material.FirstCost = 206400
	material.InputCost = 717600

	material.CalulateInputUnit(master)

	actual := material.Elements[0].Unit
	expected := 300
	assert.Equal(t, expected, actual)
}

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
