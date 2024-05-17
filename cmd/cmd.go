package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coocood/freecache"
	"github.com/fahmifan/devkit/pkg/config"
	"github.com/fahmifan/devkit/pkg/core"
	"github.com/fahmifan/devkit/pkg/dbconn"
	"github.com/fahmifan/devkit/pkg/fs"
	"github.com/fahmifan/devkit/pkg/httpsvc"
	"github.com/fahmifan/devkit/pkg/logs"
	"github.com/fahmifan/devkit/pkg/mailer/smtp"
	"github.com/fahmifan/devkit/pkg/service"
	"github.com/spf13/cobra"
)

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "devkit",
		Short: "Devkit is a auto grader for programming assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(serverCmd())

	return rootCmd.Execute()
}

func mustInitService() *service.Service {
	gormDB := dbconn.MustPostgres(config.Debug())
	mediaCfg := core.MediaConfig{
		RootDir:      config.FileUploadPath(),
		ObjectStorer: fs.NewLocalStorage(),
	}
	smtpConfig := config.SMTPConfig()

	mailer, err := smtp.NewSmtpClient(&smtpConfig)
	if err != nil {
		logs.Err(err, "init mailer failed:")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("get sql db failed:", err)
	}

	debug := config.Debug()

	cacheSize := 10 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)

	svc := service.NewService(
		gormDB,
		sqlDB,
		config.JWTKey(),
		debug,
		mediaCfg,
		config.SenderEmail(),
		mailer,
		cache,
	)

	return svc
}

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run devkit server",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := mustInitService()

			ctx := context.Background()

			server := httpsvc.NewServer(
				config.Port(),
				httpsvc.WithService(service),
				httpsvc.WithJWTKey(config.JWTKey()),
			)

			go func() {
				logs.Info("run server")
				server.Run()
			}()

			go func() {
				logs.Info("run outbox service")
				service.RegisterJobHandlers()
				service.RunOutboxService(config.Debug())
			}()

			// Wait for a signal to quit:
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
			<-signalChan

			logs.Info("stopping server")
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()
			server.Stop(ctx)

			logs.Info("stopping outbox service")
			service.StopOutboxService()
			logs.Info("outbox service stopped")

			return nil
		},
	}
}
