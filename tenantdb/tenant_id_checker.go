package tenantdb

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
)

var numberChars = map[rune]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

var lowercaseAlphabetChars = map[rune]bool{
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'g': true,
	'h': true,
	'i': true,
	'j': true,
	'k': true,
	'l': true,
	'm': true,
	'n': true,
	'o': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'u': true,
	'v': true,
	'w': true,
	'x': true,
	'z': true,
}

// SanitizeTenantId validate tenantId and sanitize it
func SanitizeTenantId(ctx context.Context, tenantId string) (newId string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SanitizeTenantId")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	tenantId = strings.TrimSpace(tenantId)

	if len(tenantId) <= 5 {
		err = fmt.Errorf("tenant id must have at least 5 characters")
		return
	}

	str := bytes.Buffer{}
	defer str.Reset()

	for i, char := range tenantId {
		// first 3 chars must not be number
		_, isNumber := numberChars[char]
		_, isAlphabet := lowercaseAlphabetChars[char]

		if i >= 0 && i <= 2 {
			if isNumber {
				err = fmt.Errorf("first 3 characters on tenant id must not be number")
				return
			}
		}

		if isNumber || isAlphabet {
			str.WriteRune(char)
		}
	}

	if str.String() != tenantId {
		return "", fmt.Errorf("tenant id contains unpermitted character")
	}

	return str.String(), nil
}
