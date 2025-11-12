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

package tui

// This file re-exports colors from the style package for backwards compatibility
import "github.com/christopher/mcl/pkg/tui/style"

// Primary Colors
var (
	PurpleHaze   = style.PurpleHaze
	CrimsonPulse = style.CrimsonPulse
	SunsetGlow   = style.SunsetGlow
)

// Accent Colors (Maestro-Specific)
var (
	OceanTide  = style.OceanTide
	OceanSurge = style.OceanSurge
	OceanDepth = style.OceanDepth
	OceanAbyss = style.OceanAbyss
	HotPink    = style.HotPink
	NeonGreen  = style.NeonGreen
)

// Grayscale
var (
	GhostWhite = style.GhostWhite
	SilverMist = style.SilverMist
	DimGray    = style.DimGray
	DeepSpace  = style.DeepSpace
)

// Focus
var (
	FocusedBorder   = style.FocusedBorder
	UnfocusedBorder = style.UnfocusedBorder
)

// Re-export functions
var (
	GetOceanTideShade = style.GetOceanTideShade
	GetDaemonShade    = style.GetDaemonShade
)
