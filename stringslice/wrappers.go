package stringslice

import (
	"strings"
	"unicode"
)

// Contains returns a filter function that calls strings.Contains
// with the element and the given argument.
func Contains(substr string) MatchFunc {
	return func(s string) bool { return strings.Contains(s, substr) }
}

// ContainsAny returns a filter function that calls strings.ContainsAny
// with the element and the given argument.
func ContainsAny(chars string) MatchFunc {
	return func(s string) bool { return strings.ContainsAny(s, chars) }
}

// ContainsRune returns a filter function that calls strings.ContainsRune
// with the element and the given argument.
func ContainsRune(r rune) MatchFunc {
	return func(s string) bool { return strings.ContainsRune(s, r) }
}

// EqualFold returns a filter function that calls strings.EqualFold
// with the element and the given argument.
func EqualFold(t string) MatchFunc {
	return func(s string) bool { return strings.EqualFold(s, t) }
}

// HasPrefix returns a filter function that calls strings.HasPrefix
// with the element and the given argument.
func HasPrefix(prefix string) MatchFunc {
	return func(s string) bool { return strings.HasPrefix(s, prefix) }
}

// HasSuffix returns a filter function that calls strings.HasSuffix
// with the element and the given argument.
func HasSuffix(suffix string) MatchFunc {
	return func(s string) bool { return strings.HasSuffix(s, suffix) }
}

// Repeat returns a map function that calls strings.Repeat
// With the element and the given argument.
func Repeat(count int) MapFunc {
	return func(s string) string { return strings.Repeat(s, count) }
}

// Replace returns a map function that calls strings.Replace
// With the element and the given argument.
func Replace(old, new string, n int) MapFunc {
	return func(s string) string { return strings.Replace(s, old, new, n) }
}

// ReplaceAll returns a map function that calls strings.ReplaceAll
// With the element and the given argument.
func ReplaceAll(old, new string) MapFunc {
	return func(s string) string { return strings.ReplaceAll(s, old, new) }
}

// ToLowerSpecial returns a map function that calls strings.ToLowerSpecial
// With the element and the given argument.
func ToLowerSpecial(c unicode.SpecialCase) MapFunc {
	return func(s string) string { return strings.ToLowerSpecial(c, s) }
}

// ToTitleSpecial returns a map function that calls strings.ToTitleSpecial
// With the element and the given argument.
func ToTitleSpecial(c unicode.SpecialCase) MapFunc {
	return func(s string) string { return strings.ToTitleSpecial(c, s) }
}

// ToUpperSpecial returns a map function that calls strings.ToUpperSpecial
// With the element and the given argument.
func ToUpperSpecial(c unicode.SpecialCase) MapFunc {
	return func(s string) string { return strings.ToUpperSpecial(c, s) }
}

// Trim returns a map function that calls strings.Trim
// With the element and the given argument.
func Trim(cutset string) MapFunc {
	return func(s string) string { return strings.Trim(s, cutset) }
}

// TrimFunc returns a map function that calls strings.TrimFunc
// With the element and the given argument.
func TrimFunc(f func(rune) bool) MapFunc {
	return func(s string) string { return strings.TrimFunc(s, f) }
}

// TrimLeft returns a map function that calls strings.TrimLeft
// With the element and the given argument.
func TrimLeft(cutset string) MapFunc {
	return func(s string) string { return strings.TrimLeft(s, cutset) }
}

// TrimLeftFunc returns a map function that calls strings.TrimLeftFunc
// With the element and the given argument.
func TrimLeftFunc(f func(rune) bool) MapFunc {
	return func(s string) string { return strings.TrimLeftFunc(s, f) }
}

// TrimPrefix returns a map function that calls strings.TrimPrefix
// With the element and the given argument.
func TrimPrefix(prefix string) MapFunc {
	return func(s string) string { return strings.TrimPrefix(s, prefix) }
}

// TrimRight returns a map function that calls strings.TrimRight
// With the element and the given argument.
func TrimRight(cutset string) MapFunc {
	return func(s string) string { return strings.TrimRight(s, cutset) }
}

// TrimRightFunc returns a map function that calls strings.TrimRightFunc
// With the element and the given argument.
func TrimRightFunc(f func(rune) bool) MapFunc {
	return func(s string) string { return strings.TrimRightFunc(s, f) }
}

// TrimSuffix returns a map function that calls strings.TrimSuffix
// With the element and the given argument.
func TrimSuffix(suffix string) MapFunc {
	return func(s string) string { return strings.TrimSuffix(s, suffix) }
}
