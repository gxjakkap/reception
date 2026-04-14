// Copyright 2026 Jakkaphat Chalermphanaphan

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"regexp"
	"strings"
)

var roleMentionRegex = regexp.MustCompile(`^<@&(\d+)>$`)

func ExtractRoleID(roleStr string) string {
	roleStr = strings.TrimSpace(roleStr)

	matches := roleMentionRegex.FindStringSubmatch(roleStr)
	if len(matches) == 2 {
		return matches[1]
	}

	isDigits := true
	for _, c := range roleStr {
		if c < '0' || c > '9' {
			isDigits = false
			break
		}
	}
	if isDigits && len(roleStr) > 0 {
		return roleStr
	}

	return ""
}

var customEmojiRegex = regexp.MustCompile(`^<a?:(\w+):(\d+)>$`)

func NormalizeEmoji(emojiStr string) string {
	emojiStr = strings.TrimSpace(emojiStr)

	matches := customEmojiRegex.FindStringSubmatch(emojiStr)
	if len(matches) == 3 {
		// Return name:id format which matches discordgo APIName()
		return matches[1] + ":" + matches[2]
	}

	// For Unicode emojis, it's already the literal character
	return emojiStr
}

func StringPtr(s string) *string {
	return &s
}
