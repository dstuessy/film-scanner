package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type TokenExpiredError struct{}

func (m *TokenExpiredError) Error() string {
	return "Token Expired"
}

func getContext() context.Context {
	return context.Background()
}

func getDriveService(token *oauth2.Token, ctx context.Context) (*drive.Service, error) {
	if token.Expiry.Before(time.Now()) {
		return nil, new(TokenExpiredError)
	}
	return drive.NewService(ctx, option.WithTokenSource(oauthConf.TokenSource(ctx, token)))
}

func listFiles(srv *drive.Service, parentId string) (*drive.FileList, error) {
	q := "mimeType='application/vnd.google-apps.folder'"

	if parentId != "" {
		q = fmt.Sprintf("%s and '%s' in parents", q, parentId)
	}

	files, err2 := srv.Files.List().
		PageSize(10).
		Q(q).
		Fields("nextPageToken, files(id, name)").
		Spaces("drive").
		Do()
	if err2 != nil {
		return nil, err2
	}

	return files, nil
}
