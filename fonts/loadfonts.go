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

package fonts

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadFont(dc *gg.Context, fontBytes []byte, size float64) error {
	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	dc.SetFontFace(face)
	return nil
}

func LoadInterRegular(dc *gg.Context, size float64) error {
	return loadFont(dc, InterRegular, size)
}

func LoadInterBold(dc *gg.Context, size float64) error {
	return loadFont(dc, InterBold, size)
}

func LoadInterLight(dc *gg.Context, size float64) error {
	return loadFont(dc, InterLight, size)
}

func LoadInterMedium(dc *gg.Context, size float64) error {
	return loadFont(dc, InterMedium, size)
}
