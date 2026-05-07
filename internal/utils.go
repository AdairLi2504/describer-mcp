package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	// register all kinds of decoder
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	//
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/image/draw"

	"github.com/mark3labs/mcp-go/mcp"
)

func ConstructTextToCallToolResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

func EncodeImage(imageSource string, toBase64 bool, compress bool) (string, error) {
	var img image.Image
	// Transform the image to image.Image for the next step
	// Or just use the original data
	// TODO: transfer to base64 directly when just need to transfet to base64
	switch classifyImageSourceType(imageSource) {
	case TypeUnknown:
		return "", fmt.Errorf("unsupported image source type")
	case TypeURL:
		// Only when the image is base64, it can be compressed.
		if !toBase64 && !compress {
			return imageSource, nil
		} else if !toBase64 && compress {
			return "", fmt.Errorf("impossbile to compress image url")
		} else {
			resp, err := http.Get(imageSource)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return "", fmt.Errorf("request failed: %d %s", resp.StatusCode, resp.Status)
			}
			contentType := resp.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "image/") {
				// still attempt to sniff the content in case the header is wrong
				var sampled bytes.Buffer
				if _, err := sampled.ReadFrom(resp.Body); err != nil {
					return "", err
				}
				detected := http.DetectContentType(sampled.Bytes())
				if !strings.HasPrefix(detected, "image/") {
					return "", fmt.Errorf("response type is not image: header=%s detected=%s", contentType, detected)
				}
				img, _, err = image.Decode(bytes.NewReader(sampled.Bytes()))
				if err != nil {
					return "", fmt.Errorf("failed to decode image: %w; header=%s detected=%s", err, contentType, detected)
				}
			} else {
				// header says image/* -- read body into buffer so we can both sniff and decode reliably
				var bodyBuf bytes.Buffer
				if _, err := bodyBuf.ReadFrom(resp.Body); err != nil {
					return "", err
				}
				img, _, err = image.Decode(bytes.NewReader(bodyBuf.Bytes()))
				if err != nil {
					detected := http.DetectContentType(bodyBuf.Bytes())
					return "", fmt.Errorf("failed to decode image: %w; header=%s detected=%s", err, contentType, detected)
				}
			}
		}
	case TypeLocalPath:
		file, err := os.Open(imageSource)
		if err != nil {
			return "", err
		}
		defer file.Close()
		img, _, err = image.Decode(file)
		if err != nil {
			return "", fmt.Errorf("failed to decode local image: %w", err)
		}
	case TypeImageBase64:
		if !compress {
			return imageSource, nil
		} else {
			_, b64Data, _ := strings.Cut(imageSource, ",")
			imgBytes, err := base64.StdEncoding.DecodeString(b64Data)
			if err != nil {
				return "", err
			}
			img, _, err = image.Decode(bytes.NewReader(imgBytes))
			if err != nil {
				return "", err
			}
		}
	}
	if compress {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		maxDim := 1024
		scale := float64(maxDim) / math.Max(float64(w), float64(h))
		// When the max side of the original image is lower than 1024, it is no need to scale
		if scale < 1.0 {
			newW := int(math.Round(float64(w) * scale))
			newH := int(math.Round(float64(h) * scale))
			imgNew := image.NewRGBA(image.Rect(0, 0, newW, newH))
			draw.CatmullRom.Scale(imgNew, imgNew.Bounds(), img, img.Bounds(), draw.Over, nil)
			img = imgNew
		}
	}
	var imgResBuf bytes.Buffer
	if err := png.Encode(&imgResBuf, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgResBuf.Bytes()), nil
}

type URIType int

const (
	TypeUnknown     URIType = iota
	TypeURL                 // http/https
	TypeImageBase64         // image base64 data URI
	TypeLocalPath
)

func classifyImageSourceType(imageSource string) URIType {
	// Trim spaces and control characters to avoid parse issues
	imageSource = strings.TrimSpace(imageSource)

	isValidPath := func(path string) bool {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return false
		}
		return !fileInfo.IsDir()
	}

	if isValidPath(imageSource) {
		return TypeLocalPath
	}

	u, err := url.Parse(imageSource)
	if err != nil {
		return TypeUnknown
	}

	switch u.Scheme {
	case "data":
		if strings.HasPrefix(u.Opaque, "image/") && strings.Contains(u.Opaque, ";base64") {
			return TypeImageBase64
		}
		return TypeUnknown
	case "http", "https":
		return TypeURL
	// TODO: add ftp and file protocol support
	default:
		return TypeUnknown
	}
}
