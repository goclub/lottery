package lottery

import xrand "github.com/goclub/rand"
import xerr "github.com/goclub/error"

func (rule Probability) Lottery() (award interface{}, err error) {
	// 验证
	if rule.Range == 0 {
		return nil, xerr.New("goclub/lottery: Probability{}.Lottery() Probability{} Range can not be zero")
	}
	var proportionTotal uint64 = 0
	for _, proportion := range rule.Proportions {
		proportionTotal += proportion.Proportion
	}
	if proportionTotal != rule.Range {
		return nil, xerr.New("goclub/lottery: Probability{}.Lottery() Probability{} Range and Proportions unmatched, Proportions[*].Proportion total must equal Range")
	}
	// 计算
	var cursorBegin uint64 = 1
	var proportionRanges []Range
	for _, item := range rule.Proportions {
		proportionRanges = append(proportionRanges, Range{
			Begin: cursorBegin,
			End: cursorBegin + item.Proportion-1,
		})
		cursorBegin = cursorBegin + item.Proportion
	}
	// 判断
	randomUint, err := xrand.RangeUint64(1, rule.Range) ; if err != nil {
	    return
	}
	for index, item := range proportionRanges {
		if randomUint >= item.Begin && randomUint <= item.End {
			return rule.Proportions[index].Award, nil
		}
	}
	return nil, xerr.New("goclub/lottery: Unknown error, Probability{}.Lottery can not match any proportion")
}

// 概率配置数据结构
type Probability struct {
	Range uint64
	Proportions []Proportion
}
type Range struct {
	Begin uint64
	End uint64
}
type Proportion struct {
	Proportion uint64
	Award interface{}
}