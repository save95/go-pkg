package maker

import (
	newmaker "github.com/save95/go-utils/maker"
)

// TransNo 生成用户订单号，长度16-32位，组成结构如下：
// ｜00|111|22222|33333|44444|
//   ｜  |    |     |     |
//   ｜  |    |     |     `-- 当日的秒数（5位）
//   ｜  |    |      `------- 用户 ID 最后五位，[0001, 9999]（4位）
//   ｜  |    `-------------- 随机数补位（2-18位）
//   ｜  `------------------- 当日是一年中的第几日，[000, 366]（3位）
//    `---------------------- 年份最后2位（2位）
// Deprecated
func TransNo(uid uint) string {
	return newmaker.TransNoWith(uid, 16)
}

// TransNoWith 生成用户订单号，长度16-32位，组成结构如下：
// ｜00|111|22222|33333|44444|
//   ｜  |    |     |     |
//   ｜  |    |     |     `-- 当日的秒数（5位）
//   ｜  |    |      `------- 用户 ID 最后五位，[0001, 9999]（4位）
//   ｜  |    `-------------- 随机数补位（2-18位）
//   ｜  `------------------- 当日是一年中的第几日，[000, 366]（3位）
//    `---------------------- 年份最后2位（2位）
// Deprecated
func TransNoWith(uid, length uint) string {
	return newmaker.TransNoWith(uid, length)
}
