package fasttld

import (
	"log"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

var idnaM *idna.Profile = idna.New(idna.MapForLookup(), idna.Transitional(true), idna.BidiRule())

type runeSlice []rune

const alphabets = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"

// Obtained from IETF RFC 3490
const labelSeparators string = "\u002e\u3002\uff0e\uff61"

var labelSeparatorsRuneSlice runeSlice = runeSlice(labelSeparators)

const controlChars string = "\u0000\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008\t\n\v\f\r\u000e\u000f" +
	"\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\u001c\u001d\u001e\u001f"

const whitespace string = controlChars + " \u0085\u0086\u00a0\u1680\u200b\u200c\u200d\uFEFF"

var whitespaceRuneSlice runeSlice = runeSlice(whitespace)

const invalidHostNameChars = controlChars + " !\u0085\u0086\u00a0\u1680\u200b\u200c\u200d\u2025\uFEFF\uff1a"

var invalidHostNameCharsRuneSlice runeSlice = runeSlice(invalidHostNameChars)

const validHostNameChars = "-." + numbers + alphabets + "\u3002\uff0e\uff61"

var validHostNameCharsRuneSlice runeSlice = runeSlice(validHostNameChars)

const endOfHostWithPortDelimiters string = `/\?#`

var endOfHostWithPortDelimitersSet asciiSet = makeASCIISet(endOfHostWithPortDelimiters)

const endOfHostDelimiters string = endOfHostWithPortDelimiters + ":"

var endOfHostDelimitersSet asciiSet = makeASCIISet(endOfHostDelimiters)

// Characters that cannot appear in UserInfo
const invalidUserInfoChars string = endOfHostWithPortDelimiters + "[]"

var invalidUserInfoCharsSet asciiSet = makeASCIISet(invalidUserInfoChars)

// For extracting URL scheme.
var schemeFirstCharSet asciiSet = makeASCIISet(alphabets)
var schemeRemainingCharSet asciiSet = makeASCIISet(alphabets + "+-.0123456789")

// getSchemeEndIndex checks if string s begins with a URL Scheme and
// returns its last index. Returns -1 if no Scheme exists.
func getSchemeEndIndex(s string) int {
	var colon bool
	var slashCount int

	for i := 0; i < len(s); i++ {
		// first character
		if i == 0 {
			// expecting schemeFirstCharSet or slash
			if schemeFirstCharSet.contains(s[i]) {
				continue
			}
			if s[i] == '/' || s[i] == '\\' {
				slashCount++
				continue
			}
			return -1
		}
		// second character onwards
		// if no slashes yet, look for schemeRemainingCharSet or colon
		// otherwise look for slashes
		if slashCount == 0 {
			if !colon {
				if schemeRemainingCharSet.contains(s[i]) {
					continue
				}
				if s[i] == ':' {
					colon = true
					continue
				}
			}
			if s[i] == '/' || s[i] == '\\' {
				slashCount++
				continue
			}
			return -1
		}
		// expecting only slashes
		if s[i] == '/' || s[i] == '\\' {
			slashCount++
			continue
		}
		if slashCount < 2 {
			return -1
		}
		return i
	}
	if slashCount >= 2 {
		return len(s)
	}
	return -1
}

// asciiSet is a 32-byte value, where each bit represents the presence of a
// given ASCII character in the set. The 128-bits of the lower 16 bytes,
// starting with the least-significant bit of the lowest word to the
// most-significant bit of the highest word, map to the full range of all
// 128 ASCII characters. The 128-bits of the upper 16 bytes will be zeroed,
// ensuring that any non-ASCII character will be reported as not in the set.
// This allocates a total of 32 bytes even though the upper half
// is unused to avoid bounds checks in asciiSet.contains.
type asciiSet [8]uint32

// makeASCIISet creates a set of ASCII characters.
//
// Similar to strings.makeASCIISet but skips input validation.
func makeASCIISet(chars string) (as asciiSet) {
	// all characters in chars are expected to be valid ASCII characters
	for _, c := range chars {
		as[c/32] |= 1 << (c % 32)
	}
	return as
}

// contains reports whether c is inside the set.
//
// same as strings.contains.
func (as *asciiSet) contains(c byte) bool {
	return (as[c/32] & (1 << (c % 32))) != 0
}

// indexAnyASCII returns the index of the first instance of any Unicode code point
// from asciiSet in s, or -1 if no Unicode code point from asciiSet is present in s.
//
// Similar to strings.IndexAny but takes in an asciiSet instead of a string
// and skips input validation.
func indexAnyASCII(s string, as asciiSet) int {
	for i, b := range []byte(s) {
		if as.contains(b) {
			return i
		}
	}
	return -1
}

// hasInvalidCharsOrConsecutiveLabelSeparators checks s for
// invalid runes or consecutive label separators
func hasInvalidCharsOrConsecutiveLabelSeparators(s string) bool {
	var isLabelSeparator bool
	for _, c := range s {
		if runeBinarySearch(c, labelSeparatorsRuneSlice) {
			if isLabelSeparator {
				return true
			}
			isLabelSeparator = true
		} else {
			isLabelSeparator = false
		}
		if runeBinarySearch(c, invalidHostNameCharsRuneSlice) {
			return true
		}
	}
	return false
}

// indexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
//
// chars is assumed to be sorted by integer value in ascending order.
//
// Similar to strings.IndexAny but uses runeSlice
// and skips input validation.
func indexAny(s string, chars runeSlice) int {
	for i, c := range s {
		if runeBinarySearch(c, chars) {
			return i
		}
	}
	return -1
}

// runeBinarySearch returns true if target exists in sortedRunes
// otherwise it returns false.
//
// sortedRunes must be already sorted by integer value in ascending order.
func runeBinarySearch(target rune, sortedRunes runeSlice) bool {
	var low int
	high := len(sortedRunes) - 1

	for low <= high {
		median := (low + high) / 2

		if sortedRunes[median] < target {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	return low != len(sortedRunes) && sortedRunes[low] == target
}

// lastIndexAny returns the index of the last instance of any Unicode code
// point from chars in s, or -1 if no Unicode code point from chars is
// present in s.
//
// Similar to strings.LastIndexAny but skips input validation and uses runeSlice.
func lastIndexAny(s string, chars runeSlice) int {
	for i := len(s); i > 0; {
		var lowerBound int
		if i > 4 {
			// minimises size of slice to search
			// rune is an alias of int32 with maximum size of 4 bytes
			lowerBound = i - 4
		}
		r, size := utf8.DecodeLastRuneInString(s[lowerBound:i])
		i -= size
		if runeBinarySearch(r, chars) {
			return i
		}
	}
	return -1
}

// reverse reverses a slice of strings in-place.
func reverse(input []string) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}

// sepSize returns byte length of an sep rune, given the rune's first byte.
func sepSize(r byte) int {
	// r is the first byte of any of the runes in labelSeparators
	if r == 46 {
		// First byte of '.' is 46
		// size of '.' is 1
		return 1
	}
	// First byte of any label separator other than '.' is not 46
	// size of separator is 3
	return 3
}

// formatAsPunycode formats s as punycode.
func formatAsPunycode(s string) string {
	asPunyCode, err := idnaM.ToASCII(s)
	if err != nil {
		log.Println(strings.SplitAfterN(err.Error(), "idna: invalid label", 2)[0])
		return ""
	}
	return asPunyCode
}

// indexLastByteBefore returns the index of the last instance of byte b
// before any byte in notAfterCharsSet, otherwise -1
func indexLastByteBefore(s string, b byte, notAfterCharsSet asciiSet) int {
	if firstNotAfterCharIdx := indexAnyASCII(s, notAfterCharsSet); firstNotAfterCharIdx != -1 {
		return strings.LastIndexByte(s[0:firstNotAfterCharIdx], b)
	}
	return strings.LastIndexByte(s, b)
}
