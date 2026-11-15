# word_counter

Counts word frequencies in text, with optional case folding and stop-word exclusion.

```go
package main

import (
	"fmt"
	"strings"

	wc "word_counter"
)

func main() {
	stop := map[string]struct{}{"the": {}, "and": {}}
	counts, err := wc.Count(strings.NewReader("The cat and the hat."), wc.Options{
		FoldCase:  true,
		StopWords: stop,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(counts)
	// map[cat:1 hat:1]
}
```

## Why this exists

Most word-count snippets are a `strings.Fields` call and a map, which silently
loses contractions and treats punctuation as part of the word. This library
keeps apostrophes inside words, splits on everything else, and lets callers
opt into case folding and stop-word filtering without a separate pass over the
text. The trade-off is that it does not attempt stemming or lemmatisation: it
counts surface forms, which is predictable and cheap, at the cost of not
merging "run" and "running".

## API

- `Count(r io.Reader, opts Options) (map[string]int, error)`
- `Options{FoldCase bool, StopWords map[string]struct{}}`

`Count` always returns a non-nil map, even for empty input, so callers can
write to it without a nil check.

## The awkward edge

Stop words are matched *after* case folding. If you set `FoldCase: true`,
supply lowercase stop words; an uppercase entry such as `"The"` will never
match. This is the one place where the two options interact, and getting it
backwards is the usual first bug.
