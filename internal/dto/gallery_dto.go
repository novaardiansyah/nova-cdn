package dto

const (
	MaxUploadSize   = 10 * 1024 * 1024 // 10MB
	UploadDirBase   = "public"
	DefaultImageDir = "gallery"
	ModelPrefix     = "App\\Models\\"
)

var AllowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type UploadInput struct {
	Dir         string
	Description string
	IsPrivate   bool
	SubjectID   *uint
	SubjectType *string
	UserID      uint
}
