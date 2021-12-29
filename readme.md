# lottery

> 由浅入深的实现**可配置**的概率和计数抽奖


## 简单的概率

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

**实现思路**

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

[Go版本实现](./probability.go?blob)

[Go版本调用](./probability_test.go?blob)

提供了正确的解法后,我列举一种错误的实现方式:

从逻辑的角度看分三步

1. 生成一个 1 到 100 之间的随机数,然后判断 `randomUint <= 10` 时返回1元,否则继续下一步
2. 再次生成一个 1 到 90 `(100-10) 之间的随机数,然后判断 `randomUint <= 20` 时返回2元,否则继续下一步
3. 直接返回7元,因为后面没有其他规则

这种方式虽然也能实现,但**概率的计算是不准确的**.通过数学中组合排列的方式去计算概率可以知道实际的概率并不是按 10% 20% 70%计算的.