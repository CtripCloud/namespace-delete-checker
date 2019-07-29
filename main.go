package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ctripcloud/namespace-delete-check/cfg"
	"github.com/ctripcloud/namespace-delete-check/handlers"
	"github.com/ctripcloud/namespace-delete-check/k8s"
	"github.com/ctripcloud/namespace-delete-check/logger"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

type SvrParameters struct {
	port                int    // webhook server port
	certFile            string // path to the x509 certificate for https
	keyFile             string // path to the x509 private key matching `certFile`
	cfgFile             string // path to configuration file
	graceShutdownPeriod int    // seconds to wait before exit
	kubeTokenFile       string // path to token used to setup k8s client
}

var params SvrParameters

func init() {
	flag.IntVar(&params.port, "port", 443, "Webhook server port.")
	flag.IntVar(&params.graceShutdownPeriod, "grace-shutdown-period", 10, "Time to wait(seconds) before shutting down the server forcely.")
	flag.StringVar(&params.certFile, "tlsCertFile", "/etc/namespace-delete-check/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&params.keyFile, "tlsKeyFile", "/etc/namespace-delete-check/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&params.cfgFile, "cfgFile", "/etc/namespace-delete-check/config/config.json", "File containing the webhook configuration.")

	flag.StringVar(&logger.LogDir, "logdir", "/var/log/namespace-delete-check", "Log directory of server.")
	flag.IntVar(&logger.LogRemain, "logRemain", 14, "Days to remain logs in logDir.")

	flag.StringVar(&k8s.KubeConfigFile, "kubeConfigFile", "", "Kubernetest config file used to connect kube-api-server.")
}

func main() {
	flag.Parse()
	stopCh := SetupSignalHandler()

	logger.MustInit()
	defer logger.Close()

	cfg.MustInit(params.cfgFile)
	k8s.MustInit()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", params.port),
		Handler: handlers.NewHandler(),
	}

	go func() {
		if err := srv.ListenAndServeTLS(params.certFile, params.keyFile); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatalf("Failed to listen and serve webhook server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout (graceShutdownPeriod)
	<-stopCh
	logrus.Infof("Shutdown Server (within %v seconds)...", params.graceShutdownPeriod)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.graceShutdownPeriod)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server shutdown failed")
	}
	logrus.Info("Server shutdown successfully")
}

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
