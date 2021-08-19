package maker

import (
	"fmt"
	"math/rand"
	"time"
)

// TransNo 生成用户订单号，长度15位，组成结构如下：
// |00|111|22|333|44444|
//   |  |  |   |    |-- 当日的秒数
//   |  |  |   |------- 用户ID 最后三位，[001, 999]
//   |  |  |----------- 随机数，[10, 99]
//   |  |-------------- 当日是一年中的第几日，[000, 366]
//   |-- 年份最后2位
func TransNo(uid uint) string {
	now := time.Now()

	// 随机数
	rand.Seed(now.UnixNano())
	rn := rand.Intn(99-10+1) + 10

	// 当日的秒数
	sec := (now.Hour()*60+now.Minute())*60 + now.Second()

	// 重组用户ID，方便取最后三位
	userId := fmt.Sprintf("000%d", uid)

	return fmt.Sprintf(
		"%d%d%d%s%05d",
		now.Year()-2000,
		now.YearDay(),
		rn,
		userId[len(userId)-3:],
		sec,
	)
}
