package transform

import "testing"

func TestApply(t *testing.T) {
	tests := []struct {
		name        string
		in          []byte
		op          Operation
		prefix      string
		suffix      string
		splitChoice string
		want        string
	}{
		{
			name: "add comma",
			in:   []byte("apple\nbanana\npear"),
			op:   AddComma,
			want: "apple,\nbanana,\npear,",
		},
		{
			name:   "prefix",
			in:     []byte("apple\nbanana"),
			op:     Prefix,
			prefix: "> ",
			want:   "> apple\n> banana",
		},
		{
			name:   "suffix",
			in:     []byte("apple\nbanana"),
			op:     Suffix,
			suffix: ";",
			want:   "apple;\nbanana;",
		},
		{
			name:        "split space",
			in:          []byte("hello world again"),
			op:          Split,
			splitChoice: "space",
			want:        "hello \nworld \nagain",
		},
		{
			name:        "split comma",
			in:          []byte("a,b,c"),
			op:          Split,
			splitChoice: "comma",
			want:        "a,\nb,\nc",
		},
		{
			name:        "split both",
			in:          []byte("a, b, c"),
			op:          Split,
			splitChoice: "both",
			want:        "a,\n \nb,\n \nc",
		},
		{
			name: "unknown operation returns input unchanged",
			in:   []byte("keep me"),
			op:   Operation("unknown"),
			want: "keep me",
		},
		{
			name: "preserve trailing newline when adding comma",
			in:   []byte("apple\nbanana\n"),
			op:   AddComma,
			want: "apple,\nbanana,\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(Apply(tt.in, tt.op, tt.prefix, tt.suffix, tt.splitChoice))
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFirstNLines(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		n    int
		want string
	}{
		{
			name: "all lines returned when n is large",
			in:   []byte("one\ntwo\nthree"),
			n:    10,
			want: "one\ntwo\nthree",
		},
		{
			name: "limited number of lines",
			in:   []byte("one\ntwo\nthree\nfour"),
			n:    2,
			want: "one\ntwo",
		},
		{
			name: "single line",
			in:   []byte("one"),
			n:    1,
			want: "one",
		},
		{
			name: "zero lines requested",
			in:   []byte("one\ntwo"),
			n:    0,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstNLines(tt.in, tt.n)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
