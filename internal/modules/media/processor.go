package media

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"

	"github.com/nfnt/resize"
)

type ImageOptions struct {
	MaxWidth  uint
	MaxHeight uint
	Quality   int // 1 - 100
}

func ProcessImage(file io.Reader, opts ImageOptions) (io.Reader, int64, error) {
	// 1. decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, err
	}

	// 2. resize image using Thumbnail (maintains aspect ratio)
	// if pass 0 for width or height, it keeps the original scales
	newImg := resize.Thumbnail(opts.MaxWidth, opts.MaxHeight, img, resize.Lanczos3)

	// 3. encode to jpeg with specified compression quality
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImg, &jpeg.Options{Quality: opts.Quality})
	if err != nil {
		return nil, 0, err
	}

	return buf, int64(buf.Len()), nil
}
