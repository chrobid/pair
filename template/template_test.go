package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func ExampleT_CoAuthor_present() {
	tt := New()

	template := `Add secret message

  Co-Authored-By: Alice <alice@example.com>`
	var buffer bytes.Buffer
	buffer.WriteString(template)
	tt.ReadFrom(&buffer)

	coAuthor, present := tt.CoAuthor()
	fmt.Println(coAuthor)
	fmt.Println(present)
	// Output:
	// Alice <alice@example.com>
	// true
}

func ExampleT_CoAuthor_absent() {
	tt := New()

	template := `Add secret message`
	var buffer bytes.Buffer
	buffer.WriteString(template)
	tt.ReadFrom(&buffer)

	coAuthor, present := tt.CoAuthor()
	fmt.Println(coAuthor)
	fmt.Println(present)
	// Output:
	// false
}

func TestT(t *testing.T) {
	t.Run("CoAuthor", func(t *testing.T) {
		testcases := []struct {
			templatePath string
			coAuthor     string
			present      bool
		}{
			{
				"simple.txt",
				"Alice <alice@example.com>",
				true,
			},
			{
				"none.txt",
				"",
				false,
			},
			{
				"double.txt", // only one co-author is supported
				"Alice <alice@example.com>",
				true,
			},
		}

		for _, tc := range testcases {
			tt := New()

			f, err := os.Open(filepath.Join("testdata", tc.templatePath))
			if err != nil {
				t.Fatalf("Missing test data: %s", tc.templatePath)
			}
			tt.ReadFrom(f)

			coAuthor, present := tt.CoAuthor()
			if present != tc.present {
				t.Errorf("Co-author detection failed for %s", tc.templatePath)
			}
			if coAuthor != tc.coAuthor {
				t.Errorf("Expected co-author to be '%s', was '%s'", tc.coAuthor, coAuthor)
			}
		}

	})

	t.Run("ReadFrom", func(t *testing.T) {
		tt := New()

		expected := `Capitalized, short (50 chars or less) summary

More detailed explanatory text, if necessary.  Wrap it to about 72
characters or so.  In some contexts, the first line is treated as the
subject of an email and the rest of the text as the body.  The blank
line separating the summary from the body is critical (unless you omit
the body entirely); tools like rebase can get confused if you run the
two together.

Write your commit message in the imperative: "Fix bug" and not "Fixed bug"
or "Fixes bug."  This convention matches up with commit messages generated
by commands like git merge and git revert.

Further paragraphs come after blank lines.

- Bullet points are okay, too

- Typically a hyphen or asterisk is used for the bullet, followed by a
  single space, with blank lines in between, but conventions vary here

- Use a hanging indent`

		var buffer bytes.Buffer
		buffer.WriteString(expected)

		n, err := tt.ReadFrom(&buffer)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n != int64(len(expected)) {
			t.Errorf("Expected %d bytes to be read, got %d", int64(len(expected)), n)
		}

		if tt.content.String() != expected {
			t.Errorf("Expected template content to be '%v', got '%v'", expected, tt.content.String())
		}
	})

	t.Run("WriteTo", func(t *testing.T) {
		var out bytes.Buffer
		tt := New()

		expected := `lala`
		var buffer bytes.Buffer
		buffer.WriteString(expected)

		tt.ReadFrom(&buffer)
		n, err := tt.WriteTo(&out)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n != int64(len(expected)) {
			t.Errorf("Expected %d bytes to be read, got %d", int64(len(expected)), n)
		}
		if out.String() != expected {
			t.Errorf("Expected output to be '%v', got '%v'", expected, out.String())
		}
	})
}
