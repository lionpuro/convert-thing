package files

var imageFormats = []string{
	"jpeg",
	"jpg",
	"png",
	"webp",
}

func Formats() map[string][]string {
	return map[string][]string{
		"image": imageFormats,
	}
}
