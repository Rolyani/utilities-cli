package transform

import "testing"

func TestApplyAddComma(t *testing.T) {
	in := []byte("apple\nbanana\npear")
	got := string(Apply(in, AddComma, "", "", ""))
	want := "apple,\nbanana,\npear,"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplyPrefix(t *testing.T) {
	in := []byte("apple\nbanana")
	got := string(Apply(in, Prefix, "> ", "", ""))
	want := "> apple\n> banana"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplySuffix(t *testing.T) {
	in := []byte("apple\nbanana")
	got := string(Apply(in, Suffix, "", ";", ""))
	want := "apple;\nbanana;"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplySplitSpace(t *testing.T) {
	in := []byte("hello world again")
	got := string(Apply(in, Split, "", "", "space"))
	want := "hello \nworld \nagain"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplySplitComma(t *testing.T) {
	in := []byte("a,b,c")
	got := string(Apply(in, Split, "", "", "comma"))
	want := "a,\nb,\nc"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplySplitBoth(t *testing.T) {
	in := []byte("a, b, c")
	got := string(Apply(in, Split, "", "", "both"))
	want := "a,\n \nb,\n \nc"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFirstNLinesAllLines(t *testing.T) {
	in := []byte("one\ntwo\nthree")
	got := FirstNLines(in, 10)
	want := "one\ntwo\nthree"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFirstNLinesLimited(t *testing.T) {
	in := []byte("one\ntwo\nthree\nfour")
	got := FirstNLines(in, 2)
	want := "one\ntwo"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
