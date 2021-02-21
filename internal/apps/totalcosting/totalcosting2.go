package totalcosting2

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
	NDBurden int         // 正常仕損の負担量
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

// Cost is PriceとUnitから費用を計算
func (e Element) Cost() float64 {
	return e.Price * float64(e.Unit)
}

// AddCost is costで指定された分の費用を加えたPriceを計算する
func (e *Element) AddCost(cost float64) {
	totalCost := e.Price*float64(e.Unit) + cost
	e.Price = totalCost / float64(e.Unit)
}

// IsBear is 仕損を負担するかの判定
// true: 負担する, false: 負担しない
func (e Element) IsBear(progress float64) bool {
	return e.Progress >= progress
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
	Elements    map[ElementType]Element
	CMethod     CalculationMethod
	DMethod     DefectiveProductMethod
	FirstCost   float64
	InputCost   float64
}

// CalulateInputUnit is 投入点と進捗度を比較して、進捗度が投入点以上なら投入
func (c *Cost) CalulateInputUnit(master map[ElementType]Element) {
	c.Elements = make(map[ElementType]Element, len(master))

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
func (c *Cost) CalulateConversionUnit(master map[ElementType]Element) {
	sumLeft := 0
	sumRight := 0
	c.Elements = make(map[ElementType]Element, len(master))

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

// GetFIFOOutputBurder is 正常仕損度外視法, 先入先出法の場合の完成品負担量を返す
func (c Cost) GetFIFOOutputBurder() int {
	unit := 0

	for _, e := range c.Elements {
		if e.Type == Input {
			unit += e.Unit
		}
		if e.Type == Last {
			unit -= e.Unit
		}
		if e.Type == NormalDefect {
			unit -= e.Unit
		}
	}

	return unit
}

// GetNormalDefectUnit is 正常仕損の数量を返す
func (c Cost) GetNormalDefectUnit() int {
	for _, e := range c.Elements {
		if e.Type == NormalDefect {
			return e.Unit
		}
	}

	return 0
}

// GetNormalDefectCost is 正常仕損の費用を返す
func (c Cost) GetNormalDefectCost() float64 {
	for _, e := range c.Elements {
		if e.Type == NormalDefect {
			return e.Cost()
		}
	}

	return 0
}

// GetTotalNDBurden is 負担量合計の計算
func (c Cost) GetTotalNDBurden() int {
	total := 0

	for _, e := range c.Elements {
		total += e.NDBurden
	}

	return total
}

// Box is 解く問題
type Box struct {
	Master           map[ElementType]Element
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
		total += c.Elements[Output].Price * float64(c.Elements[Output].Unit)
	}

	return total
}

// CalculationProductAvgCost is 完成品単位原価の計算
func (b Box) CalculationProductAvgCost() float64 {
	if b.Master[Output].Unit == 0.0 {
		return 0.0
	}
	return b.ProductTotalCost / float64(b.Master[Output].Unit)
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
		}

		// 平均的に投入
		if b.Costs[i].InputOnAvg {
			// 完成品換算量を計算
			b.Costs[i].CalulateConversionUnit(b.Master)
		}
	}

	// 正常仕損の扱い
	for i := 0; i < cCount; i++ {
		normalDefectProgress := b.Costs[i].Elements[NormalDefect].Progress

		// 度外視法
		if b.Costs[i].DMethod == Neglecting {
			// 月末仕掛品進捗度が正常仕損発生点を超えていれば両者負担
			for j := 0; j < len(b.Costs[i].Elements); j++ {
				// 完成品負担割合の計算
				// 先入先出法の場合
				if b.Costs[i].CMethod == FIFO {
					b.Costs[i].Elements[Output].NDBurden = b.Costs[i].GetFIFOOutputBurder()
				}

				// 平均法の場合
				if b.Costs[i].CMethod == AVG {
					b.Costs[i].Elements[Output].NDBurden = b.Costs[i].Elements[Output].Unit
				}

				// 月末仕掛品負担割合の計算
				if b.Costs[i].Elements[Last].Progress > normalDefectProgress {
					// 正常仕損発生点を超えているので負担
					b.Costs[i].Elements[Last].NDBurden = b.Costs[i].Elements[Last].Unit

					// 両者負担の場合は投入量から正常仕損量を控除
					b.Costs[i].Elements[Input].Unit -= b.Costs[i].Elements[Last].Unit
				} else {
					b.Costs[i].Elements[Last].NDBurden = 0
				}
			}
		}

		// 非度外視法
		if b.Costs[i].DMethod == Neglecting {
			// 月末仕掛品進捗度が正常仕損発生点を超えていれば両者負担
			for j := 0; j < len(b.Costs[i].Elements); j++ {
				elementType := b.Costs[i].Elements[j].Type

				if elementType != Output && elementType != Last {
					continue
				}

				if b.Costs[i].Elements[j].Progress > normalDefectProgress {
					// 正常仕損発生点を超えているので負担
					if b.Costs[i].CMethod == FIFO {
						// 先入先出法では完成品の負担量の計算に注意
						if elementType == Output {
							// 投入量から月末仕掛品量と正常仕損量を差し引いたのが
							// 完成品の負担量
							b.Costs[i].Elements[j].NDBurden = b.Costs[i].GetFIFOOutputBurder()
						}
						if elementType == Last {
							b.Costs[i].Elements[j].NDBurden = b.Costs[i].Elements[j].Unit
						}
					}
					if b.Costs[i].CMethod == AVG {
						// 平均法では負担量=仕掛品量
						b.Costs[i].Elements[j].NDBurden = b.Costs[i].Elements[j].Unit
					}
				} else {
					if elementType == Output {
						// 完成品が全部負担
						b.Costs[i].Elements[j].NDBurden = b.Costs[i].Elements[j].Unit
					} else {
						// 正常仕損発生点を超えていないので負担しない
						b.Costs[i].Elements[j].NDBurden = 0
					}
				}
			}
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
			totalCost := c.FirstCost + c.InputCost
			outputCost := totalCost - lastPrice*float64(totalUnit)
			c.Elements[Output].Price = outputCost / float64(c.Elements[Output].Unit)
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

	// 正常仕損の扱い
	for i := 0; i < cCount; i++ {
		// 度外視法
		if b.Costs[i].DMethod == Neglecting {
			// 月末仕掛品進捗度が正常仕損発生点を超えていれば両者負担

		}

		// 非度外視法の場合は
		if b.Costs[i].DMethod == NonNeglecting {

		}
	}

	// 月末仕掛品原価の計算
	b.EOTMTotalCost = b.CalculationEOFMCost()

	// 完成品原価の計算
	b.ProductTotalCost = b.CalculationProductCost()

	// 完成品単位原価の計算
	b.ProductAvgCost = b.CalculationProductAvgCost()
}

// Index is elementsの中からsearchで指定したElementTypeに一致する
// Elementを探してそのindexを返す
// 見つからなければ-1を返す
// searchに一致するものが複数あっても最初の1つしか返さないので注意
func Index(search ElementType, elements []Element) int {
	for i := 0; i < len(elements); i++ {
		if elements[i].Type == search {
			return i
		}
	}

	return -1
}

// GetNormalDefectProgress is 正常仕損の発生点を返す
func GetNormalDefectProgress(elements []Element) float64 {
	return elements[NormalDefect].Progress
}

// GetCountWithElementType is elementsの中にあるsearchに一致する要素の数を返す
func GetCountWithElementType(elements []Element, search []ElementType) int {
	count := 0
	elementsCount := len(elements)
	searchCount := len(search)

	for i := 0; i < elementsCount; i++ {
		for j := 0; j < searchCount; j++ {
			if elements[i].Type == search[j] {
				count++
			}
		}
	}

	return count
}
