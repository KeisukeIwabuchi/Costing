package totalcosting

// ElementType はBox図の要素の種別を表す
type ElementType int

// Box図のElementの種別
// (月初仕掛品, 投入, 完成品, 月末仕掛品, 正常仕損, 異常仕損, 正常減損, 異常減損)
const (
	First ElementType = iota
	Input
	Output
	Last
	NormalDefect
	AbnormalDefect
	NormalImpairment
	AbnormalImpairment
)

// Element はBOX図の構成要素を想定
type Element struct {
	Type     ElementType // 種別
	Price    float64     // 単価
	Unit     int         // 数量
	Progress float64     // 加工進捗度
}

// IsLeftElement is ElementTypeがBox図左側の要素かを確認する
// true: Left, false: Right
func (e Element) IsLeftElement() bool {
	leftElement := []ElementType{First, Input}

	for _, v := range leftElement {
		if e.Type == v {
			return true
		}
	}

	return false
}

// CalculationMethod is 月末仕掛品の計算方法
type CalculationMethod int

// 月末仕掛品の計算方法(先入先出法 or 平均法)
const (
	FIFO CalculationMethod = iota
	AVG
)

// DefectiveProductMethod 正常仕損の計算方法
type DefectiveProductMethod int

// 正常仕損の計算方法(度外視法 or 非度外視法)
const (
	Neglecting DefectiveProductMethod = iota
	NonNeglecting
)

// Cost is 仕掛品のBOX図
type Cost struct {
	InputOnAvg  bool
	InputTiming float64
	Elements    []Element
	CMethod     CalculationMethod
	DMethod     DefectiveProductMethod
	FirstCost   float64
	InputCost   float64
}

// CalulateInputUnit is 投入点と進捗度を比較して、進捗度が投入点以上なら投入
func (c *Cost) CalulateInputUnit(master []Element) {
	c.Elements = make([]Element, len(master))

	for i, m := range master {
		var element Element

		element.Type = m.Type

		if element.Type == Output {
			element.Progress = 1.0
		} else {
			element.Progress = m.Progress
		}

		if m.Progress < c.InputTiming {
			element.Unit = 0
		} else {
			element.Unit = m.Unit
		}

		c.Elements[i] = element
	}
}

// CalulateConversionUnit is 完成品換算量を計算
func (c *Cost) CalulateConversionUnit(master []Element) {
	sumLeft := 0
	sumRight := 0
	c.Elements = make([]Element, len(master))

	for i, m := range master {
		var element Element

		element.Type = m.Type

		if element.Type == Input {
			c.Elements[i] = element
			continue
		}

		if element.Type == Output {
			element.Progress = 1.0
			element.Unit = m.Unit
		} else {
			element.Progress = m.Progress
			element.Unit = int(float64(m.Unit) * m.Progress)
		}

		c.Elements[i] = element

		// 投入量計算のための処理
		// Unitの計算後にやること
		if element.IsLeftElement() {
			sumLeft += element.Unit
		} else {
			sumRight += element.Unit
		}
	}

	for i := 0; i < len(c.Elements); i++ {
		if c.Elements[i].Type == Input {
			c.Elements[i].Unit = sumRight - sumLeft
		}
	}
}

// GetPriceFIFO is 先入先出法での月末仕掛品平均単価を返す
func (c Cost) GetPriceFIFO() float64 {
	for _, e := range c.Elements {
		if e.Type == Input {
			return c.InputCost / float64(e.Unit)
		}
	}

	return 0.0
}

// GetPriceAVG is 平均法での月末仕掛品平均単価を返す
func (c Cost) GetPriceAVG() float64 {
	totalUnit := 0

	for _, e := range c.Elements {
		if e.IsLeftElement() {
			totalUnit += e.Unit
		}
	}

	return (c.FirstCost + c.InputCost) / float64(totalUnit)
}

// Box is 解く問題
type Box struct {
	Master           []Element
	Costs            []Cost
	ProductTotalCost float64
	ProductAvgCost   float64
	EOTMTotalCost    float64
}

// CalculationEOFMCost is 月末仕掛品原価の計算
func (b Box) CalculationEOFMCost() float64 {
	total := 0.0

	for _, c := range b.Costs {
		for _, e := range c.Elements {
			if e.Type == Last {
				total += e.Price * float64(e.Unit)
			}
		}
	}

	return total
}

// CalculationProductCost is 完成品原価の計算
func (b Box) CalculationProductCost() float64 {
	total := 0.0

	for _, c := range b.Costs {
		for _, e := range c.Elements {
			if e.Type == Output {
				total += e.Price * float64(e.Unit)
			}
		}
	}

	return total
}

// CalculationProductAvgCost is 完成品単位原価の計算
func (b Box) CalculationProductAvgCost() float64 {
	for _, e := range b.Master {
		if e.Type == Output {
			return b.ProductTotalCost / float64(e.Unit)
		}
	}

	return 0.0
}

// Add is test
func Add(e *[]Element, timing float64, master []Element) Element {
	for _, m := range master {
		var element Element

		element.Type = m.Type

		if element.Type == Output {
			element.Progress = 1.0
		} else {
			element.Progress = m.Progress
		}

		if m.Progress < timing {
			element.Unit = 0
		} else {
			element.Unit = m.Unit
		}

		*e = append(*e, element)
		return element
	}
	return Element{}
}

// Run is culcurate answer
func (b *Box) Run() {
	// 数量の計算
	cCount := len(b.Costs)
	for i := 0; i < cCount; i++ {
		// // 定点で投入
		if !b.Costs[i].InputOnAvg {
			// 投入量の計算
			b.Costs[i].CalulateInputUnit(b.Master)

			// n := len(b.Master)
			// b.Costs[i].Elements = make([]Element, 10)
			// for j := 0; j < n; j++ {
			// 	var element Element

			// 	element.Type = b.Master[j].Type

			// 	if element.Type == Output {
			// 		element.Progress = 1.0
			// 	} else {
			// 		element.Progress = b.Master[j].Progress
			// 	}

			// 	if b.Master[j].Progress < b.Costs[i].InputTiming {
			// 		element.Unit = 0
			// 	} else {
			// 		element.Unit = b.Master[j].Unit
			// 	}

			// 	b.Costs[i].Elements[j] = element
			// }
		}

		// 平均的に投入
		if b.Costs[i].InputOnAvg {
			// 完成品換算量を計算
			b.Costs[i].CalulateConversionUnit(b.Master)
		}
	}

	// 月末仕掛品原価の計算
	for _, c := range b.Costs {
		// 先入先出法
		if c.CMethod == FIFO {
			lastPrice := c.GetPriceFIFO()
			totalUnit := 0

			// 完成品以外の平均単価を代入
			for i := 0; i < len(c.Elements); i++ {
				if c.Elements[i].IsLeftElement() || c.Elements[i].Type == Output {
					continue
				}

				c.Elements[i].Price = lastPrice
				totalUnit += c.Elements[i].Unit
			}

			// 差額で完成品の平均単価を計算
			for i := 0; i < len(c.Elements); i++ {
				if c.Elements[i].Type != Output {
					continue
				}

				totalCost := c.FirstCost + c.InputCost
				outputCost := totalCost - lastPrice*float64(totalUnit)
				c.Elements[i].Price = outputCost / float64(c.Elements[i].Unit)
			}
		}

		// 平均法
		if c.CMethod == AVG {
			lastPrice := c.GetPriceAVG()

			for i := 0; i < len(c.Elements); i++ {
				if c.Elements[i].IsLeftElement() {
					continue
				}

				c.Elements[i].Price = lastPrice
			}
		}
	}

	// 月末仕掛品原価の計算
	b.EOTMTotalCost = b.CalculationEOFMCost()

	// 完成品原価の計算
	b.ProductTotalCost = b.CalculationProductCost()

	// 完成品単位原価の計算
	b.ProductAvgCost = b.CalculationProductAvgCost()
}
