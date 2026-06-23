package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed assets/BeVietnamPro-Bold.ttf
var fontBoldTTF []byte

//go:embed assets/BeVietnamPro-Regular.ttf
var fontRegularTTF []byte

const (
	ogWidth  = 1200
	ogHeight = 630
)

var (
	ogFontBold    *opentype.Font
	ogFontRegular *opentype.Font
	ogFontInitErr error
	ogFontOnce    sync.Once

	ogCache   = map[string][]byte{}
	ogCacheMu sync.RWMutex
)

func initOGFonts() error {
	ogFontOnce.Do(func() {
		b, err := opentype.Parse(fontBoldTTF)
		if err != nil {
			ogFontInitErr = fmt.Errorf("parse bold font: %w", err)
			return
		}
		r, err := opentype.Parse(fontRegularTTF)
		if err != nil {
			ogFontInitErr = fmt.Errorf("parse regular font: %w", err)
			return
		}
		ogFontBold = b
		ogFontRegular = r
	})
	return ogFontInitErr
}

type ogParams struct {
	Title    string
	Subtitle string // e.g. "Lotus AI · 19/05/2026"
}

func renderOGImage(p ogParams) ([]byte, error) {
	if err := initOGFonts(); err != nil {
		return nil, err
	}

	bg := color.RGBA{0x0f, 0x14, 0x19, 0xff}     // --bg dark
	accent := color.RGBA{0x7e, 0xe8, 0xa3, 0xff} // --accent green
	textCol := color.RGBA{0xe6, 0xed, 0xf3, 0xff}
	dim := color.RGBA{0x8b, 0x94, 0x9e, 0xff}

	img := image.NewRGBA(image.Rect(0, 0, ogWidth, ogHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Accent bar on the left
	draw.Draw(img, image.Rect(0, 0, 16, ogHeight), &image.Uniform{accent}, image.Point{}, draw.Src)

	// Brand: "🌸 Lotus AI" — drop the emoji (font may lack glyph) and use a colored dot instead
	dotX, dotY, dotR := 80, 90, 14
	drawCircle(img, dotX, dotY, dotR, accent)
	if err := drawText(img, "Lotus AI", 110, 100, 36, ogFontBold, textCol); err != nil {
		return nil, err
	}

	// Title — large, multi-line wrap
	title := strings.TrimSpace(p.Title)
	if title == "" {
		title = "Lotus AI — VN Stock Market Toolkit"
	}
	titleFace, err := opentype.NewFace(ogFontBold, &opentype.FaceOptions{Size: 56, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil, err
	}
	defer titleFace.Close()
	lines := wrapText(title, titleFace, ogWidth-160)
	if len(lines) > 4 {
		lines = lines[:4]
		lines[3] = truncateTail(lines[3], titleFace, ogWidth-160, "…")
	}
	y := 230
	for _, ln := range lines {
		if err := drawTextWithFace(img, ln, 80, y, titleFace, textCol); err != nil {
			return nil, err
		}
		y += 72
	}

	// Subtitle (dim) near bottom
	if sub := strings.TrimSpace(p.Subtitle); sub != "" {
		if err := drawText(img, sub, 80, ogHeight-60, 26, ogFontRegular, dim); err != nil {
			return nil, err
		}
	}

	// Bottom-right CTA
	cta := "lotusai.servehttp.com"
	face, err := opentype.NewFace(ogFontRegular, &opentype.FaceOptions{Size: 24, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil, err
	}
	defer face.Close()
	w := textWidth(cta, face)
	if err := drawTextWithFace(img, cta, ogWidth-80-w, ogHeight-60, face, accent); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawCircle(img *image.RGBA, cx, cy, r int, c color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, c)
			}
		}
	}
}

func drawText(img *image.RGBA, s string, x, y int, size float64, f *opentype.Font, c color.Color) error {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return err
	}
	defer face.Close()
	return drawTextWithFace(img, s, x, y, face, c)
}

func drawTextWithFace(img *image.RGBA, s string, x, y int, face font.Face, c color.Color) error {
	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{c},
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
	return nil
}

func textWidth(s string, face font.Face) int {
	d := &font.Drawer{Face: face}
	return d.MeasureString(s).Ceil()
}

func wrapText(s string, face font.Face, maxWidth int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	current := words[0]
	for _, w := range words[1:] {
		cand := current + " " + w
		if textWidth(cand, face) <= maxWidth {
			current = cand
		} else {
			lines = append(lines, current)
			current = w
		}
	}
	lines = append(lines, current)
	return lines
}

func truncateTail(s string, face font.Face, maxWidth int, tail string) string {
	if textWidth(s, face) <= maxWidth {
		return s
	}
	runes := []rune(s)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		cand := strings.TrimRight(string(runes), " ") + tail
		if textWidth(cand, face) <= maxWidth {
			return cand
		}
	}
	return tail
}

// ogImageFor returns cached PNG bytes for the given cache key, building it with builder on miss.
func ogImageFor(key string, builder func() (ogParams, error)) ([]byte, error) {
	ogCacheMu.RLock()
	if b, ok := ogCache[key]; ok {
		ogCacheMu.RUnlock()
		return b, nil
	}
	ogCacheMu.RUnlock()

	params, err := builder()
	if err != nil {
		return nil, err
	}
	png, err := renderOGImage(params)
	if err != nil {
		return nil, err
	}
	ogCacheMu.Lock()
	ogCache[key] = png
	ogCacheMu.Unlock()
	return png, nil
}

// formatVnDate renders a time as "02/01/2006" without importing the heavy formatter elsewhere.
func formatVnDate(t time.Time) string { return t.Format("02/01/2006") }
