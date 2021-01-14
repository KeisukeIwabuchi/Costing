package totalcosting

// Element はBOX図の構成要素を想定
type Element struct {
	Cost     float64
	Price    float64
	Unit     int
	Progress float64
	Type     ElementType
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
	NormalDefect
	AbnormalDefect
	NormalImpairment
	AbnormalImpairment
)

// CalculationMethod is 月末仕掛品の計算方法
type CalculationMethod int

// 月末仕掛品の計算方法(先入先出法 or 平均法)
const (
	unknown CalculationMethod = iota
	FIFO
	AVG
)

// DefectiveProductMethod 正常仕損の計算方法
type DefectiveProductMethod int

// 正常仕損の計算方法(度外視法 or 非度外視法)
const (
	Neglecting DefectiveProductMethod = iota
	NonNeglecting
)

// BOX is 仕掛品のBOX図
type BOX struct {
	Left         Elements
	Right        Elements
	VirtualLeft  Elements
	VirtualRight Elements
	CMethod      CalculationMethod
	DMethod      DefectiveProductMethod
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
	var sumCost float64
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

		sumCost += v.Cost
		sumUnit += v.Unit
	}

	return sumCost / float64(sumUnit)
}

// UnitPriceWithFIFO is Calculate Unit Price with FIFO
func (b BOX) UnitPriceWithFIFO() float64 {
	return CalculateAverageUnitPrice(b.Left, Input)
}

// UnitPriceWithAVG is Calculate Unit Price with Average Method
func (b BOX) UnitPriceWithAVG() float64 {
	return CalculateAverageUnitPrice(b.Left, First, Input)
}

// CreateVirtualLeft is Calculat BOX.VirtualLeft
func (b BOX) CreateVirtualLeft() {
	b.VirtualLeft = b.Left

	for _, value := range b.Left {
		if value.Type != NormalDefect {
			continue
		}

		if b.DMethod == Neglecting {

		}
		if b.DMethod == NonNeglecting {

		}
	}
}

// GetElement is return Element
func GetElement(elements Elements, search ElementType) (Element, bool) {
	for _, element := range elements {
		if element.Type == search {
			return element, true
		}
	}

	return Element{}, false
}

// Contains is ckeck ElementType
func Contains(arr []ElementType, eType ElementType) bool {
	for _, v := range arr {
		if eType == v {
			return true
		}
	}
	return false
}
