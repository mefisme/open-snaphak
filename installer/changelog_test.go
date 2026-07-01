package main

import (
	"strings"
	"testing"
)

// TestFormatChangelogFullNotesNoTruncation locks the changelog contract: the newest release's notes are
// printed IN FULL (every line, no matter how many) with no "... full notes" link, while older releases are
// listed compactly with a link. Regression guard for the closed-beta report where the console showed only
// links (a 404 into the private repo) instead of the changelog.
func TestFormatChangelogFullNotesNoTruncation(t *testing.T) {
	// A newest release whose notes are far longer than the old 12-line cap.
	var lines []string
	for i := 1; i <= 30; i++ {
		lines = append(lines, "- change number "+itoa(i))
	}
	newestBody := strings.Join(lines, "\n")

	list := []ghRelease{
		{TagName: "v0.2.0", Prerelease: true, PublishedAt: "2026-07-01T00:00:00Z", Body: newestBody,
			HTMLURL: "https://github.com/snaphak/open-snaphak/releases/tag/v0.2.0"},
		{TagName: "v0.1.0", Prerelease: true, PublishedAt: "2026-06-28T00:00:00Z", Body: "- older note",
			HTMLURL: "https://github.com/snaphak/open-snaphak/releases/tag/v0.1.0"},
	}

	out := formatChangelog(list, "v0.2.0")

	// Every line of the newest release must appear -- nothing truncated.
	for _, ln := range lines {
		if !strings.Contains(out, ln) {
			t.Errorf("newest release note line missing from output: %q", ln)
		}
	}
	// No "full notes" link, and the newest release's own page link must NOT be emitted (it's printed in full).
	if strings.Contains(out, "full notes:") {
		t.Error("output still truncates with a 'full notes:' link -- the newest notes must print in full")
	}
	if strings.Contains(out, "releases/tag/v0.2.0") {
		t.Error("the newest release should print its notes verbatim, not a link to its page")
	}
	// The installed marker + the newest headline.
	if !strings.Contains(out, "Latest release: v0.2.0 (beta)") {
		t.Errorf("missing latest-release headline; got:\n%s", out)
	}
	if !strings.Contains(out, "<- you have this") {
		t.Error("missing the installed-version marker")
	}
	// The OLDER release is a compact link, not full notes.
	if !strings.Contains(out, "Earlier releases:") {
		t.Error("missing the 'Earlier releases:' section")
	}
	if !strings.Contains(out, "releases/tag/v0.1.0") {
		t.Error("older release should be listed with its GitHub link")
	}
}

// TestFormatChangelogSingleRelease: one release -> full notes, no "Earlier releases" section.
func TestFormatChangelogSingleRelease(t *testing.T) {
	list := []ghRelease{
		{TagName: "v0.1.0-beta.4", Prerelease: true, PublishedAt: "2026-07-01T12:34:56Z",
			Body: "Changes since v0.1.0-beta.3:\n\n- fix the changelog command"},
	}
	out := formatChangelog(list, "")
	if !strings.Contains(out, "- fix the changelog command") {
		t.Errorf("full notes missing; got:\n%s", out)
	}
	if strings.Contains(out, "Earlier releases:") {
		t.Error("a single release must not render an 'Earlier releases' section")
	}
	if strings.Contains(out, "2026-07-01T12:34:56Z") {
		t.Error("date should be trimmed to YYYY-MM-DD")
	}
	if !strings.Contains(out, "2026-07-01") {
		t.Error("expected the trimmed date")
	}
}

// TestReleaseHeadline covers the pure one-line formatter.
func TestReleaseHeadline(t *testing.T) {
	r := ghRelease{TagName: "v1.0.0", Prerelease: false, PublishedAt: "2026-08-15T09:00:00Z"}
	if got := releaseHeadline(r, "v1.0.0"); got != "v1.0.0   2026-08-15   <- you have this" {
		t.Errorf("stable installed headline = %q", got)
	}
	rb := ghRelease{TagName: "v1.1.0-beta.1", Prerelease: true, PublishedAt: "2026-09-01T00:00:00Z"}
	if got := releaseHeadline(rb, "v1.0.0"); got != "v1.1.0-beta.1 (beta)   2026-09-01" {
		t.Errorf("beta non-installed headline = %q", got)
	}
}

// itoa avoids importing strconv for a one-off in the table above.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var d []byte
	for n > 0 {
		d = append([]byte{byte('0' + n%10)}, d...)
		n /= 10
	}
	return string(d)
}
