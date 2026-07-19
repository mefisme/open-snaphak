package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// User modding content (overrides, prefabs, custom strings) lives in the app-data folder alongside the install
// record -- one consolidated place, %LOCALAPPDATA%\snapmap-plus\ (appDataDir). This is the right home for mod
// data: it is reliably writable (unlike a game folder under Program Files, where non-elevated writes fail or get
// redirected), it survives a game uninstall / "verify integrity of game files" / reinstall, and it is out of the
// home directory so its name never collides with a repository clone. The backend reads the same tree; its startup
// scaffolder (dllmain.c) mirrors userContentSubdirs -- keep the two in sync.

// userContentSubdirs is the content tree the installer scaffolds and the backend reads. Mirrors dllmain.c's subs[].
var userContentSubdirs = []string{"strings", "overrides", "prefabs"}

// oldUserContentDir is where content lived before this version: %USERPROFILE%\snaphak\ (a path reused from the
// original tool, historically in the home root). Empty if USERPROFILE is unset.
func oldUserContentDir() string {
	base := os.Getenv("USERPROFILE")
	if base == "" {
		return ""
	}
	return filepath.Join(base, "snaphak")
}

// oldAppDataDir is the pre-rename app-data folder, %LOCALAPPDATA%\open-snaphak (install record + token).
func oldAppDataDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return ""
	}
	return filepath.Join(base, "open-snaphak")
}

func oldRecordPath() string {
	if d := oldAppDataDir(); d != "" {
		return filepath.Join(d, "install.json")
	}
	return ""
}

func oldTokenPath() string {
	if d := oldAppDataDir(); d != "" {
		return filepath.Join(d, "token")
	}
	return ""
}

// ensureUserDataTree creates the content folders (overrides / prefabs / strings) under the app-data dir if a
// fresh profile lacks them, so the disk-backed features work on a clean install instead of silently no-opping.
func ensureUserDataTree() {
	dir := appDataDir()
	if dir == "" {
		return
	}
	for _, sub := range userContentSubdirs {
		os.MkdirAll(filepath.Join(dir, sub), 0o755)
	}
}

// migrateUserData scaffolds the content tree and, one time, folds a user's existing content forward from the
// old home-root %USERPROFILE%\snaphak\ folder into the app-data dir. Copy-not-move: the old folder is left
// untouched as a backup, and a file already present at the destination is never overwritten -- so this is
// idempotent (a re-run copies nothing) and never clobbers newer content. Best-effort; never fails the install.
func migrateUserData() {
	dir := appDataDir()
	if dir == "" {
		return
	}
	ensureUserDataTree()

	old := oldUserContentDir()
	if old == "" || sameFile(old, dir) {
		return
	}
	if n := copyTreeMissing(old, dir); n > 0 {
		fmt.Printf("  ~ copied your saved content (overrides / prefabs) into the new location (%d file(s)).\n", n)
		fmt.Printf("    It now lives in %s; your old folder was left untouched as a backup.\n", dir)
	}
}

// copyTreeMissing recursively copies every file under src into dst, preserving the relative layout, but only
// when the destination file does NOT already exist (it never overwrites). Returns the number of files copied.
// Best-effort: an unreadable/unwritable file (or a missing src) is skipped, not fatal.
func copyTreeMissing(src, dst string) int {
	copied := 0
	filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil { // unreadable entry, or src doesn't exist -> nothing to migrate
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel, rerr := filepath.Rel(src, path)
		if rerr != nil {
			return nil
		}
		target := filepath.Join(dst, rel)
		if _, err := os.Stat(target); err == nil {
			return nil // already present -> never overwrite
		}
		if os.MkdirAll(filepath.Dir(target), 0o755) != nil {
			return nil
		}
		if copyFile(path, target) == nil {
			copied++
		}
		return nil
	})
	return copied
}
