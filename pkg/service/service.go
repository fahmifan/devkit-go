package service

import (
	"context"
	"database/sql"

	"connectrpc.com/connect"
	"github.com/coocood/freecache"
	"github.com/fahmifan/devkit/pkg/config"
	"github.com/fahmifan/devkit/pkg/core"
	"github.com/fahmifan/devkit/pkg/core/auth/auth_cmd"
	"github.com/fahmifan/devkit/pkg/core/mediastore/mediastore_cmd"
	"github.com/fahmifan/devkit/pkg/jobqueue"
	"github.com/fahmifan/devkit/pkg/jobqueue/outbox"
	"github.com/fahmifan/devkit/pkg/mailer"
	autogradv1 "github.com/fahmifan/devkit/pkg/pb/devkit/v1"
	"github.com/fahmifan/devkit/pkg/pb/devkit/v1/devkitv1connect"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type Service struct {
	coreCtx *core.Ctx

	*auth_cmd.AuthCmd
	*mediastore_cmd.MediaStoreCmd

	outboxService *outbox.OutboxService
}

var _ devkitv1connect.DevkitServiceHandler = &Service{}

func NewService(
	gormDB *gorm.DB,
	sqlDB *sql.DB,
	jwtKey string,
	debug bool,
	mediaCfg core.MediaConfig,
	senderEmail string,
	mailer mailer.Mailer,
	cache *freecache.Cache,
) *Service {
	outboxService := outbox.NewOutboxService(gormDB, sqlDB, config.Debug())

	coreCtx := &core.Ctx{
		GormDB:         gormDB,
		JWTKey:         jwtKey,
		MediaConfig:    mediaCfg,
		SenderEmail:    senderEmail,
		AppLink:        config.WebBaseURL(),
		LogoURL:        config.BaseURL() + "/logo.png",
		Mailer:         mailer,
		OutboxEnqueuer: outboxService,
		SqlDB:          sqlDB,
		Debug:          debug,
		Cache:          cache,
		Flight:         &singleflight.Group{},
	}

	return &Service{
		coreCtx:       coreCtx,
		outboxService: outboxService,
		AuthCmd:       &auth_cmd.AuthCmd{Ctx: coreCtx},
		MediaStoreCmd: &mediastore_cmd.MediaStoreCmd{Ctx: coreCtx},
	}
}

func (service *Service) Ping(ctx context.Context, req *connect.Request[autogradv1.Empty]) (*connect.Response[autogradv1.PingResponse], error) {
	return &connect.Response[autogradv1.PingResponse]{
		Msg: &autogradv1.PingResponse{
			Message: "pong",
		},
	}, nil
}

func (service *Service) RunOutboxService(debug bool) error {
	return service.outboxService.Run(debug)
}

func (service *Service) StopOutboxService() {
	service.outboxService.Stop()
}

func (service *Service) RegisterJobHandlers() {
	handlers := []jobqueue.JobHandler{}

	outbox.RegisterHandlers(service.coreCtx.GormDB, service.coreCtx.SqlDB, service.coreCtx.Debug, handlers)
}
