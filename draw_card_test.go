package lottery_test

import (
	"fmt"
	"github.com/goclub/lottery"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDrawCard_Lottery(t *testing.T) {
	data := lottery.DrawCard{
		Threshold: 3, // 每3次抽奖
		RoundAward: 1, // 有一次中奖
		Inventory: 10, // 共发放10个奖品
	}
	bingoCount := 0
	for i := 0; i <	33 ; i++ {
		bingo, err := data.Lottery() ; assert.NoError(t, err)
		if bingo {
			bingoCount++
		}
		fmt.Println(i+1, bingo)
	}
	log.Print("bingoCount", bingoCount)
}
