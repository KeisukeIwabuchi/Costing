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

// LeftElementType はBox図左側のElementの種別を表す
type LeftElementType int

// Box図左側のElementの種別
// (月初仕掛品, 当月投入)
const (
	First LeftElementType = iota
	Input
)

// RightElementType はBox図右側のElementの種別を表す
type RightElementType int

// Box図左側のElementの種別
// (完成品, 月末仕掛品, 正常仕損, 異常仕損, 正常減損, 異常減損)
const (
	Output RightElementType = iota
	Last
	NormalDefect
	AbnormalDefect
	NormalImpairment
	AbnormalImpairment
)

// Element はBOX図の構成要素を想定
type Element struct {
	Price    float64 // 単価
	Unit     int     // 数量
	Progress float64 // 加工進捗度
}

// LeftElement はBox図の左側のElement
type LeftElement struct {
	Type  LeftElementType
	Value Element
}

// RightElement はBox図の右側のElement
type RightElement struct {
	Type  RightElementType
	Value Element
}

// Box is 仕掛品のBOX図
type Box struct {
	InputTiming float64
	Left        []LeftElement
	Right       []RightElement
	CMethod     CalculationMethod
	DMethod     DefectiveProductMethod
	FirstCost   float64
	InputCost   float64
}

// TotalCosting is 解く問題
type TotalCosting struct {
	Boxes          []Box
	ProductAvgCost float64
	EOTMCost       float64
}

// Run is culcurate answer
func (t TotalCosting) Run() {
	// var product, processing Box
	// var totalUnitLeft, totalUnitRight int

	// for _, v := range t.boxes[0].Left {
	// 	// BOX図の左側はFirstかInputのみ
	// 	if v.Type != First && v.Type != Input {
	// 		panic("Invalid Left Type")
	// 	}

	// 	// 材料の計算
	// 	product.Left = append(product.Left, v)

	// 	// 加工費の計算
	// 	element := v
	// 	if v.Type == First {
	// 		// 加工費は加工進捗度をかけて完成品換算量を求める
	// 		// 当月投入量は差分で求めるのでここでは計算しない
	// 		element.Unit = int(float64(element.Unit) * element.Progress)

	// 		// 当月投入量のUnitは加算したくないのでLeftではここで足す
	// 		totalUnitLeft += element.Unit
	// 	}
	// 	processing.Left = append(processing.Left, element)
	// }

	// for _, v := range t.boxes[0].Right {
	// 	// BOX図の右側はFirstとInput以外の全て
	// 	if v.Type == First || v.Type == Input {
	// 		panic("Invalid Right Type")
	// 	}

	// 	// 材料の計算
	// 	product.Right = append(product.Right, v)

	// 	// 加工費の計算
	// 	element := v
	// 	if v.Type != Output {
	// 		// 加工費は加工進捗度をかけて完成品換算量を求める
	// 		// 完成品はそのまま
	// 		element.Unit = int(float64(element.Unit) * element.Progress)
	// 	}
	// 	processing.Right = append(processing.Right, element)

	// 	// 完成品のUnitも加算したいのでRightではここで足す
	// 	totalUnitRight += element.Unit
	// }

	// // ここまでの計算結果でtotalUnitLeftの方が大きければ何かが間違っている
	// if totalUnitLeft > totalUnitRight {
	// 	panic("Invalid totalUnitLeft")
	// }

	// for _, v := range processing.Left {
	// 	if v.Type == Input {
	// 		v.Unit = totalUnitRight - totalUnitLeft
	// 	}
	// }
}

// // CalculateUnitPrice is Calculate Unit Price
// func (b Box) CalculateUnitPrice() float64 {
// 	switch b.CMethod {
// 	case FIFO:
// 		return b.UnitPriceWithFIFO()
// 	case AVG:
// 		return b.UnitPriceWithAVG()
// 	}

// 	// default
// 	return b.UnitPriceWithAVG()
// }

// // CalculateAverageUnitPrice is return average price
// func CalculateAverageUnitPrice(e Elements, filter ...ElementType) float64 {
// 	var sumCost float64
// 	var sumUnit int

// 	contains := func(arr []ElementType, e_type ElementType) bool {
// 		for _, v := range arr {
// 			if e_type == v {
// 				return true
// 			}
// 		}
// 		return false
// 	}

// 	for _, v := range e {
// 		if !contains(filter, v.Type) {
// 			continue
// 		}

// 		for _, cost := range v.Cost {
// 			sumCost += cost
// 		}
// 		sumUnit += v.Unit
// 	}

// 	return sumCost / float64(sumUnit)
// }

// // UnitPriceWithFIFO is Calculate Unit Price with FIFO
// func (b Box) UnitPriceWithFIFO() float64 {
// 	return CalculateAverageUnitPrice(b.Left, Input)
// }

// // UnitPriceWithAVG is Calculate Unit Price with Average Method
// func (b Box) UnitPriceWithAVG() float64 {
// 	return CalculateAverageUnitPrice(b.Left, First, Input)
// }

// // CreateVirtualLeft is Calculat BOX.VirtualLeft
// func (b Box) CreateVirtualLeft() {
// 	b.VirtualLeft = b.Left

// 	input, _ := GetElement(b.Left, Input)
// 	last, _ := GetElement(b.Right, Last)
// 	normalDefect, result := GetElement(b.Left, NormalDefect)

// 	// 正常仕損がなければ終了
// 	if !result {
// 		return
// 	}

// 	// 度外視法
// 	if b.DMethod == Neglecting {
// 		// 月末仕掛品の加工進捗度が仕損発生点を超えている場合
// 		if normalDefect.Progress >= last.Progress {
// 			input.Unit -= normalDefect.Unit
// 		}
// 	}

// 	// 非度外視法(特に必要な処理無し)
// }

// // CreateVirtualRight is Calculat BOX.VirtualRight
// func (b Box) CreateVirtualRight() {
// 	b.VirtualRight = b.Right

// 	output, _ := GetElement(b.Right, Output)
// 	last, _ := GetElement(b.Right, Last)
// 	normalDefect, result := GetElement(b.Left, NormalDefect)

// 	// 正常仕損がなければ終了
// 	if !result {
// 		return
// 	}

// 	// 度外視法
// 	if b.DMethod == Neglecting {
// 		// 月末仕掛品の加工進捗度が仕損発生点を超えていない場合(完成品が全部負担)
// 		if normalDefect.Progress < last.Progress {
// 			output.Unit += normalDefect.Unit
// 		}
// 	}

// 	// 非度外視法(特に必要な処理無し)
// }

// // GetElement is return Element
// func GetElement(elements Elements, search ElementType) (Element, bool) {
// 	for _, element := range elements {
// 		if element.Type == search {
// 			return element, true
// 		}
// 	}

// 	return Element{}, false
// }

// // Contains is ckeck ElementType
// func Contains(arr []ElementType, eType ElementType) bool {
// 	for _, v := range arr {
// 		if eType == v {
// 			return true
// 		}
// 	}
// 	return false
// }
