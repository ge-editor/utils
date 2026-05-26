# locale package

実際に使うのは？

- UTF-8
- UTF-8 with BOM
- Shift_JIS（日本）
- Windows-1252（欧米）


| 地域    | 代表例                |
| ----- | ------------------ |
| 日本    | Shift_JIS / EUC-JP |
| 中国    | GB2312 / GBK       |
| 韓国    | EUC-KR             |
| 西欧    | ISO-8859-1         |
| ロシア   | KOI8-R             |
| DOS時代 | CP437 / CP932 など   |


| locale | 追加           |
| ------ | ------------ |
| ja_JP  | Shift_JIS    |
| zh_CN  | GB18030      |
| zh_TW  | Big5         |
| ko_KR  | EUC-KR       |
| ru_RU  | KOI8-R       |
| en_US  | Windows-1252 |

---

```go
package encoding

import (
	"os"
	"strings"
)

func DetectLocaleEncodings() []string {
	encs := []string{"utf-8"}

	locale := getLocale()

	switch {
	case strings.HasPrefix(locale, "ja"):
		encs = append(encs, "shift_jis")
	case strings.HasPrefix(locale, "zh_cn"):
		encs = append(encs, "gb18030")
	case strings.HasPrefix(locale, "zh_tw"):
		encs = append(encs, "big5")
	case strings.HasPrefix(locale, "ko"):
		encs = append(encs, "euc-kr")
	case strings.HasPrefix(locale, "ru"):
		encs = append(encs, "koi8-r")
	}

	return encs
}

func getLocale() string {
	l := os.Getenv("LC_ALL")
	if l == "" {
		l = os.Getenv("LANG")
	}
	return strings.ToLower(l)
}
```

---
