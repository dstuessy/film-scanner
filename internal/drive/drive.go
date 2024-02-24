package drive

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/dstuessy/film-scanner/internal/auth"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GetContext() context.Context {
	return context.Background()
}

func GetDriveFileService(token *oauth2.Token, ctx context.Context) (*drive.Service, error) {
	if token.Expiry.Before(time.Now()) {
		return nil, new(auth.TokenExpiredError)
	}
	return drive.NewService(ctx, option.WithTokenSource(auth.OauthConf.TokenSource(ctx, token)))
}

func CreateFolder(srv *drive.Service, name string, parentId string) (*drive.File, error) {
	f := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
	}

	if parentId != "" {
		f.Parents = []string{parentId}
	}

	r, err := srv.Files.Create(f).Do()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func FindFolder(srv *drive.Service, name string) (*drive.File, error) {
	q := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and trashed=false", name)

	files, err := srv.Files.List().
		PageSize(1).
		Q(q).
		Fields("files(id, name)").
		Spaces("drive").
		Do()
	if err != nil {
		return nil, err
	}

	if len(files.Files) == 0 {
		return nil, nil
	}

	return files.Files[0], nil
}

func ListFiles(srv *drive.Service, parentId string) (*drive.FileList, error) {
	q := "mimeType='image/jpeg' or mimeType='application/vnd.google-apps.folder' and trashed=false"

	if parentId != "" {
		q = fmt.Sprintf("%s and '%s' in parents", q, parentId)
	}

	fmt.Println(q)

	files, err := srv.Files.List().
		PageSize(10).
		Q(q).
		Fields("nextPageToken, files(id, name)").
		Spaces("drive").
		Do()
	if err != nil {
		return nil, err
	}

	return files, nil
}

func SaveImage(srv *drive.Service, img []byte, name string, parentId string) (*drive.File, error) {
	f := &drive.File{
		Name:     name,
		MimeType: "image/jpeg",
	}

	if parentId != "" {
		f.Parents = []string{parentId}
	}

	r, err := srv.Files.Create(f).Media(bytes.NewReader(img)).Do()
	if err != nil {
		return nil, err
	}

	return r, nil
}
