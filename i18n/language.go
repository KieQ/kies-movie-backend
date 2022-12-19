package i18n

import (
	"context"
	"kies-movie-backend/constant"
	"strings"
)

func Translate(ctx context.Context, index SentenceIndex) string {
	var lang = ""
	if ctx.Value(constant.Language) != nil {
		if val, ok := ctx.Value(constant.Language).(string); ok {
			lang = val
		}
	}
	if val, exist := translate[index]; exist {
		switch strings.ToLower(lang) {
		case "", "en":
			return val.English
		case "zh":
			return val.Chinese
		default:
			return "unknown language"
		}
	} else {
		return "unknown sentence index"
	}
}
