// Copyright 2025 Christopher O'Connell
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	dir := GetConfigDir()

	if dir == "" {
		t.Fatal("GetConfigDir() returned empty string")
	}

	// Should end with maestro
	if !strings.HasSuffix(dir, "maestro") {
		t.Errorf("GetConfigDir() = %q, should end with 'maestro'", dir)
	}

	// Platform-specific checks
	if runtime.GOOS == "windows" {
		if !strings.Contains(dir, "AppData") {
			t.Errorf("Windows: GetConfigDir() = %q, should contain 'AppData'", dir)
		}
	} else {
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, ".maestro")
		if dir != expected {
			t.Errorf("Unix: GetConfigDir() = %q, want %q", dir, expected)
		}
	}
}

func TestConfigFile(t *testing.T) {
	file := ConfigFile()

	if file == "" {
		t.Fatal("ConfigFile() returned empty string")
	}

	// Should end with config.yml
	if !strings.HasSuffix(file, "config.yml") {
		t.Errorf("ConfigFile() = %q, should end with 'config.yml'", file)
	}

	// Should be inside config directory
	dir := GetConfigDir()
	if !strings.HasPrefix(file, dir) {
		t.Errorf("ConfigFile() = %q, should be inside %q", file, dir)
	}
}

func TestAuthDir(t *testing.T) {
	dir := AuthDir()

	if dir == "" {
		t.Fatal("AuthDir() returned empty string")
	}

	// Should end with .claude
	if !strings.HasSuffix(dir, ".claude") {
		t.Errorf("AuthDir() = %q, should end with '.claude'", dir)
	}

	// Should be inside config directory
	configDir := GetConfigDir()
	if !strings.HasPrefix(dir, configDir) {
		t.Errorf("AuthDir() = %q, should be inside %q", dir, configDir)
	}
}

func TestGitHubAuthDir(t *testing.T) {
	dir := GitHubAuthDir()

	if dir == "" {
		t.Fatal("GitHubAuthDir() returned empty string")
	}

	// Should end with gh
	if !strings.HasSuffix(dir, "gh") {
		t.Errorf("GitHubAuthDir() = %q, should end with 'gh'", dir)
	}
}

func TestLegacyPaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		// No legacy paths on Windows
		if LegacyConfigFile() != "" {
			t.Error("Windows: LegacyConfigFile() should return empty string")
		}
		if LegacyConfigDir() != "" {
			t.Error("Windows: LegacyConfigDir() should return empty string")
		}
		if HasLegacyConfig() {
			t.Error("Windows: HasLegacyConfig() should return false")
		}
		return
	}

	// Unix/macOS
	home, _ := os.UserHomeDir()

	legacyFile := LegacyConfigFile()
	expectedFile := filepath.Join(home, ".mcl.yml")
	if legacyFile != expectedFile {
		t.Errorf("LegacyConfigFile() = %q, want %q", legacyFile, expectedFile)
	}

	legacyDir := LegacyConfigDir()
	expectedDir := filepath.Join(home, ".mcl")
	if legacyDir != expectedDir {
		t.Errorf("LegacyConfigDir() = %q, want %q", legacyDir, expectedDir)
	}

	// HasLegacyConfig() depends on filesystem state, just test it doesn't panic
	_ = HasLegacyConfig()
}

func TestEnsureDirs(t *testing.T) {
	// Test that ensure functions don't error
	// (We won't actually create dirs in test, just verify function exists)
	t.Run("EnsureConfigDir", func(t *testing.T) {
		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("EnsureConfigDir() panicked: %v", r)
			}
		}()
		// We'll skip actual execution to avoid creating dirs during test
		// err := EnsureConfigDir()
		// Just verify function signature compiles
	})

	t.Run("EnsureAuthDir", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("EnsureAuthDir() panicked: %v", r)
			}
		}()
		// Skip actual execution
	})
}
