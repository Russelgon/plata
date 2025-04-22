package log

import (
	"context"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)

	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Debugf(template string, args ...interface{})

	InfoCtx(ctx context.Context, msg string, fields ...zap.Field)
	ErrorCtx(ctx context.Context, msg string, fields ...zap.Field)
	WarnCtx(ctx context.Context, msg string, fields ...zap.Field)
	DebugCtx(ctx context.Context, msg string, fields ...zap.Field)
}
type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewZapLogger() Logger {
	l, _ := zap.NewProduction()
	return &zapLogger{
		logger: l,
		sugar:  l.Sugar(),
	}
}

func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.sugar.Infof(template, args...)
}

func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.sugar.Errorf(template, args...)
}

func (z *zapLogger) Warnf(template string, args ...interface{}) {
	z.sugar.Warnf(template, args...)
}

func (z *zapLogger) Debugf(template string, args ...interface{}) {
	z.sugar.Debugf(template, args...)
}

func (z *zapLogger) InfoCtx(_ context.Context, msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) ErrorCtx(_ context.Context, msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) WarnCtx(_ context.Context, msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) DebugCtx(_ context.Context, msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}
