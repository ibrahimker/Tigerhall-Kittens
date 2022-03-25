// Package logging is a package for logrus wrapper
package logging

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
	"google.golang.org/grpc/metadata"
)

var (
	isTerminal = term.IsTerminal(int(os.Stdout.Fd()))
	isTest     = strings.HasSuffix(os.Args[0], ".test")
	flagMu     = sync.Mutex{}
)

// NewLogger creates a new logrus logger with a fluentd formatter that will
// work natively with stackdriver logging
func NewLogger() *logrus.Logger {
	// Prepare a new logger
	logger := logrus.New()

	logger.Level = logrus.InfoLevel
	if isTest {
		// testing.Verbose can panic if Init hasn't been called yet,
		// e.g. if NewLogger is used as part of a global declaration.
		testing.Init()

		// Also parse the test flags, to be able to query -test.v. This
		// won't affect release binaries, because of the isTest check.
		// Do this behind a mutex to avoid data races, and don't parse
		// the flags twice.
		flagMu.Lock()
		if !flag.Parsed() {
			flag.Parse()
		}
		flagMu.Unlock()

		if !testing.Verbose() {
			// Keep the tests quiet, unless -test.v is used.
			logger.Level = logrus.FatalLevel
		}
	}

	logger.Out = os.Stdout
	if isTerminal {
		logger.Formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		}
	} else {
		logger.Formatter = &FluentdFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	logger.SetReportCaller(true)
	return logger
}

// WithError takes an error and logger and returns a standardised error logger
func WithError(err error, logger logrus.FieldLogger) *logrus.Entry {
	return logger.WithError(err).WithField("stacktrace", fmt.Sprintf("%+v", err))
}

// NewTestLogger is a wrapper to create logger for testing
func NewTestLogger() *logrus.Entry {
	return NewLogger().WithField("env", "testing")
}

// NewHandlerLogger is a wrapper to initiate logger at handler level
func NewHandlerLogger(ctx context.Context, logger *logrus.Entry, handlerName string, handlerRequest interface{}) (*logrus.Entry, context.Context) {
	fields := logrus.Fields{
		"handler-name": handlerName,
	}
	if handlerRequest != nil {
		fields["handler-req"] = handlerRequest
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields["handler-metadata"] = md
	}
	logger = WithContext(ctx, logger).WithFields(fields)
	ctx = WithLogger(ctx, logger)
	logger.Infof("Start Handler %s", handlerName)
	return logger, ctx
}

// NewServiceLogger is a wrapper to initiate logger at service level
func NewServiceLogger(ctx context.Context, serviceName string, fields logrus.Fields) *logrus.Entry {
	if fields != nil {
		fields["svc-name"] = serviceName
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields["svc-metadata"] = md
	}
	logger := FromContext(ctx).WithFields(fields)
	logger.Infof("Start Service %s", serviceName)
	return logger
}

// NewRepoLogger is a wrapper to initiate logger at repository level
func NewRepoLogger(ctx context.Context, repoName string, fields logrus.Fields) *logrus.Entry {
	if fields != nil {
		fields["repo-name"] = repoName
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields["repo-metadata"] = md
	}
	logger := FromContext(ctx).WithFields(fields)
	logger.Debugf("Start Repository %s", repoName)
	return logger
}
