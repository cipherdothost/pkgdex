// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package sanitize_test

import (
	"errors"
	"html/template"
	"testing"

	"go.cipher.host/pkgdex/internal/sanitize"
)

func TestUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    template.HTML
		wantErr error
	}{
		{
			name:    "empty input",
			input:   "",
			want:    "",
			wantErr: sanitize.ErrEmptyUsageString,
		},
		{
			name:  "allowed elements with valid attributes",
			input: `<pre class="highlight" tabindex="0"><code class="language-go">fmt.Println("hello")</code></pre>`,
			want:  template.HTML(`<pre class="highlight" tabindex="0"><code class="language-go">fmt.Println(&#34;hello&#34;)</code></pre>`),
		},
		{
			name:  "allowed elements with valid styles",
			input: `<a href="https://example.com" style="color: blue; text-decoration: none; outline: none;">Link</a>`,
			want:  template.HTML(`<a href="https://example.com" style="color: blue; text-decoration: none; outline: none">Link</a>`),
		},
		{
			name:  "span with allowed attributes",
			input: `<span class="highlight" id="code-1">Hello</span>`,
			want:  template.HTML(`<span class="highlight" id="code-1">Hello</span>`),
		},
		{
			name:  "disallowed elements are stripped",
			input: `<script>alert('xss')</script><pre>Safe content</pre>`,
			want:  template.HTML(`<pre>Safe content</pre>`),
		},
		{
			name:  "disallowed attributes are stripped",
			input: `<pre onclick="alert('xss')" class="safe">Content</pre>`,
			want:  template.HTML(`<pre class="safe">Content</pre>`),
		},
		{
			name:  "disallowed styles are stripped",
			input: `<a style="position: absolute; color: blue;">Link</a>`,
			want:  template.HTML(`<a style="color: blue">Link</a>`),
		},
		{
			name:  "nested elements",
			input: `<pre class="highlight"><code class="go"><span class="keyword">func</span> main()</code></pre>`,
			want:  template.HTML(`<pre class="highlight"><code class="go"><span class="keyword">func</span> main()</code></pre>`),
		},
		{
			name:  "mixed content",
			input: `Regular text <pre>code</pre> more text`,
			want:  template.HTML(`Regular text <pre>code</pre> more text`),
		},
		{
			name:  "complex XSS attempt",
			input: `<pre><img src="x" onerror="alert(1)"><script>alert(2)</script><a href="javascript:alert(3)">click</a></pre>`,
			want:  template.HTML(`<pre><a href="javascript:alert(3)">click</a></pre>`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sanitize.Usage(tt.input)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Usage() error = nil, wantErr %v", tt.wantErr)
					return
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Usage() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				return
			}

			if err != nil {
				t.Errorf("Usage() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("Usage() = %q, want %q", got, tt.want)
			}
		})
	}
}
