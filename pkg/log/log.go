package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

const correlationIdKey string = "loggerWithCorrelation"

type Logger interface {
	Info(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	Exception(ctx context.Context, message string, error error)
	WithCorrelationId(ctx context.Context, id string) context.Context
	Request(ctx context.Context, withFields Field)
	Response(ctx context.Context, withFields Field)
	ResponseWithLevel(ctx context.Context, withFields Field, level logrus.Level)
}

type logger struct {
	logRus   *logrus.Entry
	logLevel logrus.Level
}

func (l *logger) Info(ctx context.Context, message string) {

	l.withContext(ctx).WithFields(logrus.Fields{
		"DateTime": time.Now()}).Info(message)
}

func (l *logger) Warn(ctx context.Context, message string) {
	l.withContext(ctx).WithFields(logrus.Fields{
		"DateTime": time.Now()}).Warn(message)
}

func (l *logger) Exception(ctx context.Context, message string, error error) {
	l.withContext(ctx).WithFields(logrus.Fields{
		"DateTime":  time.Now(),
		"Exception": error}).Error(message)
}

func (l *logger) WithCorrelationId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIdKey, l.withContext(ctx).WithFields(logrus.Fields{"CorrelationId": id}))
}

func (l *logger) withContext(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(correlationIdKey)
	if logger == nil {
		return l.logRus
	}
	var logEntry = (logger.(*logrus.Entry))
	logEntry.Logger.SetLevel(l.logLevel)

	return logEntry
}

func (l *logger) Request(ctx context.Context, withFields Field) {

	var fields = logrus.Fields{
		"DateTime":       time.Now(),
		"Url":            withFields.Url,
		"HttpStatusCode": 102,
		"Duration":       0,
		"RequestBody":    withFields.RequestBody,
		"ResponseBody":   "",
		"HttpMethod":     withFields.HttpMethod,
		"Headers":        withFields.Headers,
		"HostName":       withFields.HostName,
	}

	for key, value := range withFields.Extra {
		fields[key] = value
	}

	l.withContext(ctx).WithFields(fields).Info(withFields.Message)

}

func (l *logger) Response(ctx context.Context, withFields Field) {

	var fields = logrus.Fields{
		"DateTime":       time.Now(),
		"Url":            withFields.Url,
		"HttpStatusCode": withFields.HttpStatusCode,
		"Duration":       withFields.Duration,
		"RequestBody":    withFields.RequestBody,
		"ResponseBody":   withFields.ResponseBody,
		"HttpMethod":     withFields.HttpMethod,
		"Headers":        withFields.Headers,
		"HostName":       withFields.HostName,
	}

	for key, value := range withFields.Extra {
		fields[key] = value
	}

	l.withContext(ctx).WithFields(fields).Info(withFields.Message)

}

func (l *logger) ResponseWithLevel(ctx context.Context, withFields Field, level logrus.Level) {

	var fields = logrus.Fields{
		"DateTime":       time.Now(),
		"Url":            withFields.Url,
		"HttpStatusCode": withFields.HttpStatusCode,
		"Duration":       withFields.Duration,
		"RequestBody":    withFields.RequestBody,
		"ResponseBody":   withFields.ResponseBody,
		"HttpMethod":     withFields.HttpMethod,
		"HostName":       withFields.HostName,
	}

	for key, value := range withFields.Extra {
		fields[key] = value
	}

	l.withContext(ctx).WithFields(fields).Logln(level, withFields.Message)

}

func NewLogger() Logger {
	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	return &logger{logRus: logrus.NewEntry(log), logLevel: logrus.InfoLevel}
}

type Field struct {
	Url            string
	HostName       string
	HttpStatusCode int
	Duration       float64
	RequestBody    string
	ResponseBody   string
	HttpMethod     string
	Message        string
	Headers        interface{}
	Extra          map[string]interface{}
}
