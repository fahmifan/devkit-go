package mediastore_cmd

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"path"
	"path/filepath"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/devkit/pkg/core"
	"github.com/fahmifan/devkit/pkg/core/auth"
	"github.com/fahmifan/devkit/pkg/core/mediastore"
	"github.com/fahmifan/devkit/pkg/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaStoreCmd struct {
	*core.Ctx
}

type InternalSaveMultipartRequest struct {
	FileInfo  *multipart.FileHeader
	MediaType mediastore.MediaFileType
}

type InternalSaveMultipartResponse struct {
	ID uuid.UUID `json:"id"`
}

func (m *MediaStoreCmd) InternalSaveMultipart(ctx context.Context, req InternalSaveMultipartRequest) (InternalSaveMultipartResponse, error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return InternalSaveMultipartResponse{}, core.ErrInternalServer
	}

	if !authUser.Role.Can(auth.CreateMedia) {
		return InternalSaveMultipartResponse{}, connect.NewError(connect.CodePermissionDenied, nil)
	}

	fileInfo := req.FileInfo

	src, err := fileInfo.Open()
	if err != nil {
		return InternalSaveMultipartResponse{}, err
	}
	defer src.Close()

	ext := filepath.Ext(fileInfo.Filename)
	fileName := generateUniqueString() + ext
	dst := path.Join(m.RootDir, fileName)

	err = m.ObjectStorer.Store(ctx, dst, src)
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: m.ObjectStorer.Store")
		return InternalSaveMultipartResponse{}, err
	}

	publicURL := fmt.Sprintf("%s/%s", m.MediaServeBaseURL, fileName)

	mediaFile, err := mediastore.CreateMediaFile(mediastore.CreateMediaRequest{
		NewID:     uuid.New(),
		Now:       time.Now(),
		FileName:  fileName,
		FileType:  req.MediaType,
		Ext:       mediastore.Extension(ext),
		PublicURL: publicURL,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: mediastore.CreateMediaFile")
		return InternalSaveMultipartResponse{}, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err = mediastore.MediaFileWriter{}.Create(ctx, m.GormDB, &mediaFile)
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: mediastore.MediaFileWriter{}.Create")
		return InternalSaveMultipartResponse{}, core.ErrInternalServer
	}

	return InternalSaveMultipartResponse{
		ID: mediaFile.ID,
	}, nil
}

type InternalSaveRequest struct {
	Ext       mediastore.Extension
	Body      io.Reader
	MediaType mediastore.MediaFileType
}

func (m *MediaStoreCmd) InternalSave(ctx context.Context, tx *gorm.DB, req InternalSaveRequest) (InternalSaveMultipartResponse, error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return InternalSaveMultipartResponse{}, core.ErrInternalServer
	}

	if !authUser.Role.Can(auth.CreateMedia) {
		return InternalSaveMultipartResponse{}, connect.NewError(connect.CodePermissionDenied, nil)
	}

	ext := req.Ext
	fileName := generateUniqueString() + string(ext)
	dst := path.Join(m.RootDir, fileName)

	err := m.ObjectStorer.Store(ctx, dst, req.Body)
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSave: m.ObjectStorer.Store")
		return InternalSaveMultipartResponse{}, err
	}

	publicURL := fmt.Sprintf("%s/%s", m.MediaServeBaseURL, fileName)

	now := time.Now()

	mediaFile, err := mediastore.CreateMediaFile(mediastore.CreateMediaRequest{
		NewID:     uuid.New(),
		Now:       now,
		FileName:  fileName,
		FileType:  req.MediaType,
		Ext:       mediastore.Extension(ext),
		PublicURL: publicURL,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: mediastore.CreateMediaFile")
		return InternalSaveMultipartResponse{}, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if tx == nil {
		tx = m.GormDB
	}

	err = mediastore.MediaFileWriter{}.Create(ctx, tx, &mediaFile)
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: mediastore.MediaFileWriter{}.Create")
		return InternalSaveMultipartResponse{}, core.ErrInternalServer
	}

	return InternalSaveMultipartResponse{
		ID: mediaFile.ID,
	}, nil
}

func generateUniqueString() string {
	h := md5.New()
	randomNumber := fmt.Sprint(rand.Intn(10))
	timestamp := fmt.Sprint(time.Now().Unix())

	h.Write([]byte(randomNumber + timestamp))

	return hex.EncodeToString(h.Sum(nil))
}
