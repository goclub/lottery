package lottery

import (
	xrand "github.com/goclub/rand"
)

type DrawCard struct {
	Threshold uint64 `note:"10次抽奖为一轮"`
	RoundAward uint64 `note:"一轮最多中3次奖"`
	Inventory uint64 `note:"计数:礼品库存"`

	count uint64 `note:"计数:每当抽奖行为产生时递增"`
	roundInventoryAdditions uint64 `note:"计数:当前这一轮新增的礼品数"`
}
func (v *DrawCard) Lottery() (bingo bool, err error)  {
	bingo, err = v.drawCardLottery() ; if err != nil {
	    return
	}
	if v.count >= v.Threshold {
		v.count = 0
		v.roundInventoryAdditions = 0
	}
	if bingo {
		v.Inventory--
		v.roundInventoryAdditions++
	}
	return
}
func (v *DrawCard) drawCardLottery() (bingo bool, err error) {
	v.count++
	if v.Inventory == 0{
		return false, nil
	}
	// 可中奖数
	canWinningNumber := v.RoundAward - v.roundInventoryAdditions
	if canWinningNumber == 0 {
		return false, nil
	}
	// 当前这一轮抽奖次数
	var chanceInRound uint64 = v.Threshold - v.count
	var random uint64
	random, err = xrand.RangeUint64(1, chanceInRound) ; if err != nil {
	    return
	}
	if random <= canWinningNumber {
		return true, nil
	}
	return false, nil
}