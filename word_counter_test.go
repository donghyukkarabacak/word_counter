package word_counter

import (
	"strings"
	"testing"
)

func TestCountBasic(t *testing.T) {
	got, err := Count(strings.NewReader("the quick brown fox the lazy dog"), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]int{
		"the":   2,
		"quick": 1,
		"brown": 1,
		"fox":   1,
		"lazy":  1,
		"dog":   1,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d words, want %d: %v", len(got), len(want), got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("count[%q] = %d, want %d", k, got[k], v)
		}
	}
}

func TestCountFoldCase(t *testing.T) {
	got, err := Count(strings.NewReader("The the THE tHe"), Options{FoldCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected one key, got %v", got)
	}
	if got["the"] != 4 {
		t.Errorf("count[the] = %d, want 4", got["the"])
	}
}

func TestCountNoFoldCase(t *testing.T) {
	got, err := Count(strings.NewReader("The the THE"), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["The"] != 1 || got["the"] != 1 || got["THE"] != 1 {
		t.Errorf("expected distinct counts, got %v", got)
	}
}

func TestCountStopWords(t *testing.T) {
	stop := map[string]struct{}{"the": {}, "a": {}}
	got, err := Count(strings.NewReader("the cat sat on a mat"), Options{StopWords: stop})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["the"]; ok {
		t.Errorf("stop word 'the' should be excluded: %v", got)
	}
	if _, ok := got["a"]; ok {
		t.Errorf("stop word 'a' should be excluded: %v", got)
	}
	if got["cat"] != 1 || got["sat"] != 1 || got["on"] != 1 || got["mat"] != 1 {
		t.Errorf("unexpected counts: %v", got)
	}
}

func TestCountStopWordsAfterFold(t *testing.T) {
	// Stop words are matched after folding, so a lowercase stop word
	// suppresses an uppercase occurrence.
	stop := map[string]struct{}{"the": {}}
	got, err := Count(strings.NewReader("The THE the"), Options{FoldCase: true, StopWords: stop})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestCountPunctuation(t *testing.T) {
	got, err := Count(strings.NewReader("hello, world! hello... world?"), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["hello"] != 2 || got["world"] != 2 {
		t.Errorf("unexpected counts: %v", got)
	}
}

func TestCountApostrophes(t *testing.T) {
	got, err := Count(strings.NewReader("don't don't 'quoted' o'clock"), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["don't"] != 2 {
		t.Errorf("count[don't] = %d, want 2", got["don't"])
	}
	if got["quoted"] != 1 {
		t.Errorf("count[quoted] = %d, want 1", got["quoted"])
	}
	if got["o'clock"] != 1 {
		t.Errorf("count[o'clock] = %d, want 1", got["o'clock"])
	}
}

func TestCountEmpty(t *testing.T) {
	got, err := Count(strings.NewReader(""), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestCountUnicode(t *testing.T) {
	got, err := Count(strings.NewReader("café Café naïveté"), Options{FoldCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["café"] != 2 {
		t.Errorf("count[café] = %d, want 2", got["café"])
	}
	if got["naïveté"] != 1 {
		t.Errorf("count[naïveté] = %d, want 1", got["naïveté"])
	}
}

func TestCountDigits(t *testing.T) {
	got, err := Count(strings.NewReader("version 2 version 2 version 3"), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["version"] != 3 || got["2"] != 2 || got["3"] != 1 {
		t.Errorf("unexpected counts: %v", got)
	}
}

func TestCountMultiline(t *testing.T) {
	got, err := Count(strings.NewReader("one two\nthree two\n"), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["two"] != 2 || got["one"] != 1 || got["three"] != 1 {
		t.Errorf("unexpected counts: %v", got)
	}
}

func TestCountNilStopWords(t *testing.T) {
	// A nil StopWords map must not panic.
	got, err := Count(strings.NewReader("a b c"), Options{StopWords: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 words, got %v", got)
	}
}
