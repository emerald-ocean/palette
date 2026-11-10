package recolour

import (
	"image"
	"testing"
)

func BenchmarkRampImage(b *testing.B) {
	r, err := New(testStops(b, testDark, testGreen, testFoam))
	if err != nil {
		b.Fatalf("New() error = %v", err)
	}
	src := image.NewYCbCr(image.Rect(0, 0, 512, 512), image.YCbCrSubsampleRatio420)
	for i := range src.Y {
		src.Y[i] = uint8(i * 2654435761 >> 13)
	}
	for i := range src.Cb {
		src.Cb[i] = uint8(i * 40503 >> 5)
		src.Cr[i] = uint8(i * 69069 >> 7)
	}
	b.ReportAllocs()
	for b.Loop() {
		r.Image(src)
	}
}
