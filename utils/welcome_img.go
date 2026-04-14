// Copyright 2026 Jakkaphat Chalermphanaphan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/gxjakkap/reception/fonts"
)

func GenerateWelcomeImage(username string, guildName string, avatarURL string, bgImg image.Image, textColorHex string) ([]byte, error) {
	const (
		W = 1024
		H = 450
	)

	dc := gg.NewContext(W, H)

	// Draw Background
	if bgImg != nil {
		// Resize/Crop background to fit
		dc.DrawImageAnchored(bgImg, W/2, H/2, 0.5, 0.5)
	} else {
		// Default solid background
		dc.SetRGB(0.15, 0.15, 0.15)
		dc.Clear()
	}

	// Draw Avatar
	if avatarURL != "" {
		resp, err := http.Get(avatarURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			avatarImg, _, err := image.Decode(resp.Body)
			if err == nil {
				// Create circular avatar
				radius := 100.0
				x, y := float64(W/2), float64(H/2-50)
				
				// Draw shadow/outer circle
				dc.SetRGBA(0, 0, 0, 0.5)
				dc.DrawCircle(x, y, radius+5)
				dc.Fill()

				dc.DrawCircle(x, y, radius)
				dc.Clip()
				dc.DrawImageAnchored(avatarImg, int(x), int(y), 0.5, 0.5)
				dc.ResetClip()
			}
		}
	}

	// Set Text Color
	if textColorHex == "" {
		textColorHex = "#ffffff"
	}
	dc.SetHexColor(textColorHex)

	// Load Fonts
	fontBold, err := truetype.Parse(fonts.InterBold)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bold font: %w", err)
	}
	fontRegular, err := truetype.Parse(fonts.InterRegular)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regular font: %w", err)
	}

	// Welcome Text
	faceBold := truetype.NewFace(fontBold, &truetype.Options{Size: 60})
	dc.SetFontFace(faceBold)
	dc.DrawStringAnchored("WELCOME", W/2, H/2+100, 0.5, 0.5)

	// Username Text
	faceRegular := truetype.NewFace(fontRegular, &truetype.Options{Size: 40})
	dc.SetFontFace(faceRegular)
	dc.DrawStringAnchored(username, W/2, H/2+160, 0.5, 0.5)

	// Guild Name Text
	faceSmall := truetype.NewFace(fontRegular, &truetype.Options{Size: 24})
	dc.SetFontFace(faceSmall)
	dc.DrawStringAnchored(fmt.Sprintf("to %s", guildName), W/2, H/2+210, 0.5, 0.5)

	// Output
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// Helper to download image (optional, can be refactored into a separate util if needed frequently)
func fetchImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	return img, err
}
