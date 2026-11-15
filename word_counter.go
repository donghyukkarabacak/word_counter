// Package word_counter provides a small utility for counting word
// frequencies in text, with optional case folding and stop-word exclusion.
//
// The design goal is to keep the API concrete and dependency-free so the
// library can be dropped into any Go program without pulling in modules.
package word_counter

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Options controls how words are normalised before counting.
type Options struct {
	// FoldCase lowercases every word before counting, so "The" and "the"
	// are counted together.
	FoldCase bool

	// StopWords is a set of words to ignore. Membership is checked after
	// case folding when FoldCase is true, so callers should supply
	// lowercase stop words in that case.
	StopWords map[string]struct{}
}

// Count reads all words from r and returns a map of word to frequency.
//
// A word is a maximal run of Unicode letters and digits. Apostrophes are
// kept inside words so that "don't" stays one token, but a leading or
// trailing apostrophe is trimmed. Everything else (punctuation, whitespace,
// symbols) separates words.
//
// The returned map is never nil, even when r yields no words; callers can
// mutate it without a nil check.
func Count(r io.Reader, opts Options) (map[string]int, error) {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(r)
	// Allow long lines. The default 64 KiB token limit is easy to hit with
	// generated text, and a word counter should not care about line length.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		for _, w := range splitWords(scanner.Text()) {
			if opts.FoldCase {
				w = strings.ToLower(w)
			}
			if opts.StopWords != nil {
				if _, stop := opts.StopWords[w]; stop {
					continue
				}
			}
			counts[w]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return counts, nil
}

// splitWords splits a single line into word tokens. It is deliberately
// permissive about apostrophes: they are allowed inside a word but stripped
// from the edges, which matches how people read contractions and possessives.
func splitWords(line string) []string {
	var words []string
	var b strings.Builder

	flush := func() {
		if b.Len() == 0 {
			return
		}
		w := b.String()
		b.Reset()
		// Trim edge apostrophes that were captured while scanning.
		w = strings.Trim(w, "'")
		if w != "" {
			words = append(words, w)
		}
	}

	for _, r := range line {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case r == '\'' && b.Len() > 0:
			// Keep an apostrophe only when it is between word characters.
			b.WriteRune(r)
		default:
			flush()
		}
	}
	flush()
	return words
}
