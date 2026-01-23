package media

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // register PNG decoder for image.Decode
	"io"

	"github.com/nfnt/resize"
)

type imageOptions struct {
	MaxWidth  uint
	MaxHeight uint
	Quality   int // 1 - 100
}

func processImage(file io.Reader, opts imageOptions) (io.Reader, int64, error) {
	// 1. decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, err
	}

	// 2. resize only if the image is actually larger thn the max limits
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()

	if dx <= 0 || dy <= 0 {
		return nil, 0, fmt.Errorf("invalid image dimensions: %dx%d", dx, dy)
	}

	width := uint(dx)
	height := uint(dy)

	var finalImg image.Image
	if width > opts.MaxWidth || height > opts.MaxHeight {
		// if pass 0 for width or height, it keeps the original scales
		finalImg = resize.Thumbnail(opts.MaxWidth, opts.MaxHeight, img, resize.Lanczos3)
	} else {
		finalImg = img
	}

	// 3. encode to jpeg with specified compression quality
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, finalImg, &jpeg.Options{Quality: opts.Quality})
	if err != nil {
		return nil, 0, err
	}

	return buf, int64(buf.Len()), nil
}
