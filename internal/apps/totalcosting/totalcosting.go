package totalcosting

// Material は材料費を想定しているけど加工費のことも考えると不適切なネーミング
type Material struct {
	Price float64
	Unit  int
}

// WIP means Work in Process
type WIP struct {
	First  Material
	Input  Material
	Output Material
	End    Material

	CET CalcEndType
}

// CalcEndType is 月末仕掛品の計算方法
type CalcEndType int

// 月末仕掛品の計算方法(先入先出法 or 平均法)
const (
	unknown CalcEndType = iota
	FIFO
	AVG
)

// CalculateUnitPrice is Calculate Unit Price
func (w WIP) CalculateUnitPrice() float64 {
	switch w.CET {
	case FIFO:
		return w.UnitPriceWithFIFO()
	case AVG:
		return w.UnitPriceWithAVG()
	}

	// default
	return w.UnitPriceWithAVG()
}

// CalculateAverageUnitPrice is return average price
func CalculateAverageUnitPrice(materials ...Material) float64 {
	var sumPrice float64
	var sumUnit int

	for _, value := range materials {
		sumPrice += value.Price
		sumUnit += value.Unit
	}

	return sumPrice / float64(sumUnit)
}

// UnitPriceWithFIFO is Calculate Unit Price with FIFO
func (w WIP) UnitPriceWithFIFO() float64 {
	return CalculateAverageUnitPrice(w.Input)
}

// UnitPriceWithAVG is Calculate Unit Price with Average Method
func (w WIP) UnitPriceWithAVG() float64 {
	return CalculateAverageUnitPrice(w.First, w.Input)
}
