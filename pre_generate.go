package lottery

type PreGenerate struct {
	Total uint64 // 总数
	AwardTotal uint64 // 奖品数
}

func (v *PreGenerate) Gen(callback func(number uint64, award bool)) {
	stepLength := float64(v.AwardTotal)/float64(v.Total)
	var awardCount uint64 = 0
	var n uint64 = 1
	for ;n<=v.Total;n++ {
		cursor := stepLength*float64(n)
		award := func() bool {
			if (cursor - float64(awardCount)) >= 1 {
				return true
			}
			return false
		}()
		if award {
			awardCount++
		}
		callback(n, award)
	}
	// 防御编程:防止少生成 (从测试结果看这个判断是多余的,但是理论上加上此判断更安心)
	if awardCount != v.AwardTotal {
		panic("goclub/lottery: awardCount != v.AwardTotal")
	}
}