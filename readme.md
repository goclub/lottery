# lottery

> 由浅入深的实现**可配置**的概率和计数抽奖


## 概率

用文字描述需求: 实现`10%`中1元 `20%` 概率中2元 `70%` 概率中7元

用数据描述需求:

```js
{
  // 概率范围
  range: 100,
  proportions: [
      {
          proportion: 10,
          award: 1, // 1元,单位分
      },
      {
          proportion: 20,
          award: 2, // 2元,单位分
      },
      {
          proportion: 70,
          award: 7, // 7元,单位分
      },
  ]
}
```

这里的数据结构经过了设计,增加的 `range: 100` 是为了便于调整概率经度(百分之一,万分之一).


基于上述数据结构你可以自己思考或编码五分钟尝试去实现

....五分钟...

### 实现思路

将 `proportions` 转换为如下数据结构

```js
proportionRanges = [
  // proportion: 10
  {
      begin: 1, 
      end: 10,
  },
  // proportion: 20
  {
    begin: 11,
    end: 30,
  },
  // proportion: 70
  {
      begin: 31,
      end: 100,
  },
]
```

然后生成一个 1 到 100 之间的随tenThousand机数,遍历 `proportionRanges`

当 `randomUint >= item.begin && randomUint <= item.end` 时则返回奖项

## 实现

[Go版本函数](./probability.go?blob)

[Go版本调用](./probability_test.go?blob)

### 错误的实现方式

提供了正确的解法后,我列举一种错误的实现方式:

从逻辑的角度看分三步

1. 生成一个 1 到 100 之间的随机数,然后判断 `randomUint <= 10` 时返回1元,否则继续下一步
2. 再次生成一个 1 到 90 `(100-10) 之间的随机数,然后判断 `randomUint <= 20` 时返回2元,否则继续下一步
3. 直接返回7元,因为后面没有其他规则

这种方式虽然也能实现,但**概率的计算是不准确的**.通过数学中组合排列的方式去计算概率可以知道实际的概率并不是按 10% 20% 70%计算的.

## 计数

- 概率: **有10%的概率中奖**
- 计数: **每当遇到10的倍数时中奖** 

概率中奖有可能会出现连续抽10次一次都不中,但计数型确保了第10次一定中奖.计数型能更让奖项的发放更平均.

用文字描述需求: 每当遇到10的倍数时中奖

### 伪代码

用伪代码表达逻辑

```js
var data = {
    threshold: 10, // 每10次抽奖
    inventory: 0 // 礼品库存
    count: 0, //计数,每当抽奖行为产生时递增
}
function lottery() {
    data.count++
    if (data.count >= data.threshold) {
        data.inventory++
        data.count = 0 
    }
    if (data.inventory>0) {
        data.inventory--
        return true
    }
    return false
}
```

> 你可能会觉得只需要 `count` 而不需要 `inventory`.暂时先不解释原因,后面的章节会介绍到奖池的概念.

### 原子性

计数型包含了可变的数据 `count` `inventory`, 它们肯定是存储在数据库中,并且操作数据的时候要注意 [原子性](https://be.nimo.run/theory/atomicity)

用 redis-eval 实现满足原子性的逻辑:

```lua
local threshold = tonumber(ARGV[1])
-- 递增并获取递增后的值
local newCount = redis.call("INCR", "count")
-- 判断是否达到阈值 (注意这里要写 >= 而不是 = , 这样代码更健壮)
if (newCount >= threshold)
then
    -- 达到阈值则递增库存
    redis.call("INCR", "inventory")
    -- 并将计数归零
    redis.redis("SET", "count", 0)
end
-- 获取库存
local inventory = redis.call("GET", "inventory")
-- 确保库存是 number
if (inventory)
then
    inventory = tonumber(inventory)
else
    -- 如果库存不存在则视为0库存
    inventory = 0
end
-- 判断库存
if (inventory > 0 )
then
    -- 递减库存
    redis.call("DECR", "inventory")
    return "won"
end
-- 库存不足未中奖
return "miss"
```

## 抽牌

- 概率: **有10%的概率中奖**
- 计数: **每当遇到10的倍数时中奖**
- 抽牌: **每10次抽奖必定中奖,且最多只中1次**

我来举一个现实生活中的例子来解释这里的抽牌指的是什么:

桌面上有**红蓝黄**三张牌.你想要抽到**红牌**,于是闭上眼睛去抽牌:

- 在桌面上放 **红蓝黄** 三张牌
- **第1次**抽到了 **蓝牌**,此时桌上还有**红牌** **黄牌** 
- **第2次**抽到了 **黄牌**,此时桌上还有**红牌**
- **第3次**抽到了 **红牌**
- 你重新**将三张牌放回到了桌上**,还是想要抽到**红牌**
- **第4次**抽到了 **黄牌**,此时桌上还有**红牌** **蓝牌**
- **第5次**抽到了 **红牌**
- **第6次**抽到了 **蓝牌**

> 我们将视桌子上放三张牌到牌抽完视为一轮 `round`
> 第1/2/3次抽奖为一轮,第4/5/6次抽奖为一轮

计数型是控制第10次中奖,而抽牌式则是让中奖更加平均.可能是1~10中的任意一次中奖.

> 那么每 10 次抽奖为一轮 `round` 
 
用文字描述需求: **每10次抽奖必定中奖,且最多中1次**

伪代码:


```js
var data = {
    threshold: 10, // 10次抽奖为一轮
    roundAward: 1, // 一轮最多中1次
    roundInventoryAdditions: 0, // 计数:当前这一轮新增的礼品数
    count: 0, //计数:每当抽奖行为产生时递增
    inventory: 0 // 计数:礼品库存
}
function lottery(data) {
    coreLottery(data)
    if (data.count >= data.threshold) {
        data.count = 0
        data.roundInventoryAdditions = 0
    }
}
function coreLottery(data) {
    data.count++
    // 每轮的中奖次数不能超过 roundAward
    if (data.roundInventoryAdditions >= data.roundAward) {
        return false
    }
    // 保底操作:
    // 如果规则是每10次抽奖最多中 1 次,那么第10次必中
    //  mustWinningTime = 10 - 1 + 1
    //  mustWinningTime =      9 + 1
    //  mustWinningTime =         10

    // 如果规则是每10次抽奖最多中 2 次,那么第9和第10次必中
    //  mustWinningTime = 10 - 2 + 1
    //  mustWinningTime =      8 + 1
    //  mustWinningTime =         9
    
    var mustWinningTime = data.threshold - roundAward
    if (data.count >= mustWinning) {
        data.inventory++
        data.roundInventoryAdditions++
        return true
    }
    return false
}
```

redis 版本 @grifree 来实现