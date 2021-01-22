package totalcosting

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

// Box is 解く問題
type Box struct {
	Master         []Element
	Costs          []Cost
	ProductAvgCost float64
	EOTMCost       float64
}

// Run is culcurate answer
func (b Box) Run() {
	for _, cost := range b.Costs {
		// 定点で投入
		if !cost.InputOnAvg {
			// 投入量の計算
			cost.CalulateInputUnit(b.Master)
		}

		// 平均的に投入
		if cost.InputOnAvg {
			// 完成品換算量を計算
			cost.CalulateConversionUnit(b.Master)
		}
	}
}

// CalulateInputUnit is 投入点と進捗度を比較して、進捗度が投入点以上なら投入
func (c Cost) CalulateInputUnit(master []Element) {
	for _, m := range master {
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

		c.Elements = append(c.Elements, element)
	}
}

// CalulateConversionUnit is 完成品換算量を計算
func (c Cost) CalulateConversionUnit(master []Element) {
	var sumLeft, sumRight int

	for _, m := range master {
		var element Element

		element.Type = m.Type

		if element.Type == Input {
			c.Elements = append(c.Elements, element)
			continue
		}

		if element.Type == Output {
			element.Progress = 1.0
		} else {
			element.Progress = m.Progress
		}

		element.Unit = int(float64(m.Unit) * m.Progress)

		c.Elements = append(c.Elements, element)

		// 投入量計算のための処理
		// Unitの計算後にやること
		if element.IsLeftElement() {
			sumLeft += element.Unit
		} else {
			sumRight += element.Unit
		}
	}

	for _, e := range c.Elements {
		if e.Type == Input {
			e.Unit = sumRight - sumLeft
		}
	}
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
