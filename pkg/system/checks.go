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

package system

import (
	"os/exec"
)

// IsDockerAvailable checks if Docker is installed and the daemon is running
func IsDockerAvailable() (bool, string) {
	// Check if docker command exists
	_, err := exec.LookPath("docker")
	if err != nil {
		return false, "Docker command not found in PATH"
	}

	// Check if docker daemon is running
	cmd := exec.Command("docker", "ps")
	if err := cmd.Run(); err != nil {
		return false, "Docker daemon not running"
	}

	return true, "Docker is available"
}

// IsClaudeAvailable checks if Claude CLI is installed
func IsClaudeAvailable() (bool, string) {
	_, err := exec.LookPath("claude")
	if err != nil {
		return false, "Claude CLI not found in PATH"
	}

	// Verify it's actually executable
	cmd := exec.Command("claude", "--version")
	if err := cmd.Run(); err != nil {
		return false, "Found but not executable"
	}

	return true, "Claude CLI is available"
}
