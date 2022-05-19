package app

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/kakao/varlog/internal/varlogadm"
	"github.com/kakao/varlog/internal/varlogadm/mrmanager"
	"github.com/kakao/varlog/internal/varlogadm/snmanager"
	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/log"
)

func InitCLI() *cli.App {
	app := &cli.App{
		Name:    "varlogadm",
		Usage:   "run varlog admin server",
		Version: "v0.0.1",
	}
	cli.VersionFlag = &cli.BoolFlag{Name: "version"}
	app.Commands = append(app.Commands, initStartCommand())
	return app
}

func initStartCommand() *cli.Command {
	return &cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start [flags]",
		Action:  start,
		Flags: []cli.Flag{
			flagClusterID.StringFlag(false, types.ClusterID(1).String()),
			flagListen.StringFlag(false, varlogadm.DefaultListenAddress),
			flagReplicationFactor.UintFlag(false, varlogadm.DefaultReplicationFactor),
			flagMetadataRepository.StringSliceFlag(true, nil),
			flagInitMRConnRetryCount.IntFlag(false, mrmanager.DefaultInitialMRConnectRetryCount),
			flagInitMRConnRetryBackoff.DurationFlag(false, mrmanager.DefaultInitialMRConnectRetryBackoff),
			flagSNWatcherRPCTimeout.DurationFlag(false, varlogadm.DefaultWatcherRPCTimeout),
			flagMRConnTimeout.DurationFlag(false, mrmanager.DefaultMRConnTimeout),
			flagMRCallTimeout.DurationFlag(false, mrmanager.DefaultMRCallTimeout),
			flagLogDir.StringFlag(false, ""),
			flagLogToStderr.BoolFlag(),
		},
	}
}

func start(c *cli.Context) error {
	clusterID, err := types.ParseClusterID(c.String(flagClusterID.Name))
	if err != nil {
		return err
	}
	logger, err := newLogger(c)
	if err != nil {
		return err
	}
	logger = logger.Named("adm").With(zap.Uint32("cid", uint32(clusterID)))
	defer func() {
		_ = logger.Sync()
	}()

	mrMgr, err := mrmanager.New(context.TODO(),
		mrmanager.WithAddresses(c.StringSlice(flagMetadataRepository.Name)...),
		mrmanager.WithInitialMRConnRetryCount(c.Int(flagInitMRConnRetryCount.Name)),
		mrmanager.WithInitialMRConnRetryBackoff(c.Duration(flagInitMRConnRetryBackoff.Name)),
		mrmanager.WithMRManagerConnTimeout(c.Duration(flagMRConnTimeout.Name)),
		mrmanager.WithMRManagerCallTimeout(c.Duration(flagMRCallTimeout.Name)),
		mrmanager.WithClusterID(clusterID),
		mrmanager.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	snMgr, err := snmanager.New(context.TODO(),
		snmanager.WithClusterID(clusterID),
		snmanager.WithClusterMetadataView(mrMgr.ClusterMetadataView()),
		snmanager.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	opts := []varlogadm.Option{
		varlogadm.WithLogger(logger),
		varlogadm.WithListenAddress(c.String(flagListen.Name)),
		varlogadm.WithReplicationFactor(c.Uint(flagReplicationFactor.Name)),
		varlogadm.WithMetadataRepositoryManager(mrMgr),
		varlogadm.WithStorageNodeManager(snMgr),
		varlogadm.WithWatcherOptions(
			varlogadm.WithWatcherRPCTimeout(c.Duration(flagSNWatcherRPCTimeout.Name)),
		),
	}
	return Main(opts, logger)
}

func newLogger(c *cli.Context) (*zap.Logger, error) {
	logtostderr := c.Bool(flagLogToStderr.Name)
	logdir := c.String(flagLogDir.Name)
	if !logtostderr && len(logdir) == 0 {
		return zap.NewNop(), nil
	}

	logOpts := []log.Option{
		log.WithHumanFriendly(),
		log.WithZapLoggerOptions(zap.AddStacktrace(zap.DPanicLevel)),
	}
	if !logtostderr {
		logOpts = append(logOpts, log.WithoutLogToStderr())
	}
	if len(logdir) != 0 {
		absDir, err := filepath.Abs(logdir)
		if err != nil {
			return nil, err
		}
		logOpts = append(logOpts, log.WithPath(filepath.Join(absDir, "varlogadm.log")))
	}
	return log.New(logOpts...)
}
