package i18n

import (
	"context"
	"strings"
)

const ContextLanguage = "language"

func Translate(ctx context.Context, index SentenceIndex) string {
	var lang = ""
	if ctx.Value(ContextLanguage) != nil {
		if val, ok := ctx.Value(ContextLanguage).(string); ok {
			lang = val
		}
	}
	if val, exist := translate[index]; exist {
		switch strings.ToLower(lang) {
		case "", "en":
			return val.English
		case "zh-cn":
			return val.Chinese
		default:
			return "unknown language"
		}
	} else {
		return "unknown sentence index"
	}
}
