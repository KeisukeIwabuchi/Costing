package totalcosting

// Element はBOX図の構成要素を想定
type Element struct {
	Price float64
	Unit  int
	Type  ElementType
}

// Elements はElementの集合体を想定、BOX図の左側と右側を表現するために使う
type Elements []Element

// ElementType はElementの種別を表す
type ElementType int

// Elementの種別
const (
	other ElementType = iota
	First
	Input
	Output
	Last
)

// CalculationMethod is 月末仕掛品の計算方法
type CalculationMethod int

// 月末仕掛品の計算方法(先入先出法 or 平均法)
const (
	unknown CalculationMethod = iota
	FIFO
	AVG
)

// BOX is 仕掛品のBOX図
type BOX struct {
	Left  Elements
	Right Elements

	CMethod CalculationMethod
}

// CalculateUnitPrice is Calculate Unit Price
func (b BOX) CalculateUnitPrice() float64 {
	switch b.CMethod {
	case FIFO:
		return b.UnitPriceWithFIFO()
	case AVG:
		return b.UnitPriceWithAVG()
	}

	// default
	return b.UnitPriceWithAVG()
}

// CalculateAverageUnitPrice is return average price
func CalculateAverageUnitPrice(e Elements, filter ...ElementType) float64 {
	var sumPrice float64
	var sumUnit int

	contains := func(arr []ElementType, e_type ElementType) bool {
		for _, v := range arr {
			if e_type == v {
				return true
			}
		}
		return false
	}

	for _, v := range e {
		if !contains(filter, v.Type) {
			continue
		}

		sumPrice += v.Price
		sumUnit += v.Unit
	}

	return sumPrice / float64(sumUnit)
}

// UnitPriceWithFIFO is Calculate Unit Price with FIFO
func (b BOX) UnitPriceWithFIFO() float64 {
	return CalculateAverageUnitPrice(b.Left, Input)
}

// UnitPriceWithAVG is Calculate Unit Price with Average Method
func (b BOX) UnitPriceWithAVG() float64 {
	return CalculateAverageUnitPrice(b.Left, First, Input)
}
