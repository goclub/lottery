package lottery_test

import (
	"github.com/goclub/lottery"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestProbability_Lottery(t *testing.T) {
	// 运行一百次观察结果
	awardCount := map[int]int{}
	for i:=0;i<100;i++ {
		award, err := lottery.Probability{
			Range: 100,
			Proportions: []lottery.Proportion{
				{
					Proportion: 10,
					Award: 1,
				},
				{
					Proportion: 20,
					Award: 2,
				},
				{
					Proportion: 70,
					Award: 7,
				},
			},
		}.Lottery() ; assert.NoError(t, err)
		awardUint64 := award.(int)
		awardCount[awardUint64] = awardCount[awardUint64]+1
	}
	log.Print("awardCount: ", awardCount)
	// 允许出现5的误差
	assert.Greater(t, awardCount[1], 5)
	assert.Less(t, awardCount[1], 15)
	// 允许出现10的误差
	assert.Greater(t, awardCount[2], 10)
	assert.Less(t, awardCount[2], 30)
	// 允许出现20的误差
	assert.Greater(t, awardCount[7], 50)
	assert.Less(t, awardCount[7], 90)
}

