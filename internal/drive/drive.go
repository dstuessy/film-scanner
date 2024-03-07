package drive

import (
	"bytes"
	"context"
	"fmt"

	"github.com/dstuessy/film-scanner/internal/auth"

	"golang.org/x/oauth2"
	gdrive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const DriveDirName = "Open Scanner"

func GetContext() context.Context {
	return context.Background()
}

func GetDriveFileService(token *oauth2.Token, ctx context.Context) (*gdrive.Service, error) {
	return gdrive.NewService(ctx, option.WithTokenSource(auth.OauthConf.TokenSource(ctx, token)))
}

func CreateFolder(srv *gdrive.Service, name string, parentId string) (*gdrive.File, error) {
	f := &gdrive.File{
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

func FindFolder(srv *gdrive.Service, name string) (*gdrive.File, error) {
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

func GetFile(srv *gdrive.Service, id string) (*gdrive.File, error) {
	file, err := srv.Files.Get(id).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GetWorkspaceDir(srv *gdrive.Service) (*gdrive.File, error) {
	q := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and trashed=false", DriveDirName)

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

func ListFiles(srv *gdrive.Service, parentId string) (*gdrive.FileList, error) {
	q := "(mimeType='image/jpeg' or mimeType='application/vnd.google-apps.folder') and trashed=false"

	if parentId != "" {
		q = fmt.Sprintf("%s and '%s' in parents", q, parentId)
	}

	fmt.Println(q)

	files, err := srv.Files.List().
		PageSize(10).
		Q(q).
		Fields("nextPageToken, files(id, name, thumbnailLink, iconLink, mimeType)").
		Spaces("drive").
		Do()
	if err != nil {
		return nil, err
	}

	return files, nil
}

func SaveImage(srv *gdrive.Service, img []byte, name string, parentId string) (*gdrive.File, error) {
	f := &gdrive.File{
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

func DeleteFile(srv *gdrive.Service, id string) error {
	_, err := srv.Files.Update(id, &gdrive.File{Trashed: true}).Do()
	return err
}
