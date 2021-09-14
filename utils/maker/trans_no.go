package maker

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// TransNo 生成用户订单号，长度16-32位，组成结构如下：
// ｜00|111|22222|33333|44444|
//   ｜  |    |     |     |
//   ｜  |    |     |     `-- 当日的秒数（5位）
//   ｜  |    |      `------- 用户 ID 最后五位，[0001, 9999]（4位）
//   ｜  |    `-------------- 随机数补位（2-18位）
//   ｜  `------------------- 当日是一年中的第几日，[000, 366]（3位）
//    `---------------------- 年份最后2位（2位）
func TransNo(uid uint) string {
	return TransNoWith(uid, 16)
}

const (
	_symbols = "012345678901234567890123456789012345678901234567890123456789"
)

// TransNoWith 生成用户订单号，长度16-32位，组成结构如下：
// ｜00|111|22222|33333|44444|
//   ｜  |    |     |     |
//   ｜  |    |     |     `-- 当日的秒数（5位）
//   ｜  |    |      `------- 用户 ID 最后五位，[0001, 9999]（4位）
//   ｜  |    `-------------- 随机数补位（2-18位）
//   ｜  `------------------- 当日是一年中的第几日，[000, 366]（3位）
//    `---------------------- 年份最后2位（2位）
func TransNoWith(uid, length uint) string {
	baseLen := 14 // 除随机字符串之外的长度
	rnLen := int(math.Min(32, math.Max(16, float64(length)))) - baseLen

	bytes := make([]byte, 32)
	rand.Read(bytes)
	symbolsByteLength := byte(len(_symbols))
	for i, b := range bytes {
		bytes[i] = _symbols[b%symbolsByteLength]
	}
	rn := string(bytes)[0:rnLen]

	now := time.Now()

	// 当日的秒数
	sec := (now.Hour()*60+now.Minute())*60 + now.Second()

	// 用户ID转字符串
	userId := fmt.Sprintf("%06d", uid)

	return fmt.Sprintf(
		"%d%d%s%s%05d",
		now.Year()-2000,
		now.YearDay(),
		rn,
		userId[len(userId)-4:],
		sec,
	)
}
