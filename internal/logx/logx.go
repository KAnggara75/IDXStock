/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/gpl-3.0.html>.
 *
 * @author KAnggara75 on Fri 27/06/25 07.44
 * @project IDXStock logx
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/logx
 */

package logx

import (
	"context"
	"github.com/KAnggara75/IDXStock/internal/app"
	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/rs/zerolog"
)

var emptyTraceIDPlaceholder = config.Get40Space()

func getTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceIDValue := ctx.Value("trace_id")
	if traceIDValue != nil {
		if traceID, ok := traceIDValue.(string); ok {
			return traceID
		}
	}
	return ""
}

func withTraceID(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	traceID := getTraceIDFromContext(ctx)
	if traceID != "" {
		return event.Str("trace_id", traceID)
	}
	return event.Str("trace_id", emptyTraceIDPlaceholder)
}

func Debug(msg string) {
	app.NewLogger().Debug().Str("trace_id", emptyTraceIDPlaceholder).Msg(msg)
}

func Debugf(format string, args ...interface{}) {
	app.NewLogger().Debug().Str("trace_id", emptyTraceIDPlaceholder).Msgf(format, args...)
}

func DebugCtx(ctx context.Context, msg string) {
	withTraceID(ctx, app.NewLogger().Debug()).Msg(msg)
}

func DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	withTraceID(ctx, app.NewLogger().Debug()).Msgf(format, args...)
}

func Info(msg string) {
	app.NewLogger().Info().Str("trace_id", emptyTraceIDPlaceholder).Msg(msg)
}

func Infof(format string, args ...interface{}) {
	app.NewLogger().Info().Str("trace_id", emptyTraceIDPlaceholder).Msgf(format, args...)
}

func InfoCtx(ctx context.Context, msg string) {
	withTraceID(ctx, app.NewLogger().Info()).Msg(msg)
}

func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	withTraceID(ctx, app.NewLogger().Info()).Msgf(format, args...)
}

func Warn(msg string) {
	app.NewLogger().Warn().Str("trace_id", emptyTraceIDPlaceholder).Msg(msg)
}

func Warnf(format string, args ...interface{}) {
	app.NewLogger().Warn().Str("trace_id", emptyTraceIDPlaceholder).Msgf(format, args...)
}

func WarnCtx(ctx context.Context, msg string) {
	withTraceID(ctx, app.NewLogger().Warn()).Msg(msg)
}

func WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	withTraceID(ctx, app.NewLogger().Warn()).Msgf(format, args...)
}

func Error(msg string) {
	app.NewLogger().Error().Str("trace_id", emptyTraceIDPlaceholder).Msg(msg)
}

func Errorf(format string, args ...interface{}) {
	app.NewLogger().Error().Str("trace_id", emptyTraceIDPlaceholder).Msgf(format, args...)
}

func ErrorCtx(ctx context.Context, msg string) {
	withTraceID(ctx, app.NewLogger().Error()).Msg(msg)
}

func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	withTraceID(ctx, app.NewLogger().Error()).Msgf(format, args...)
}

func Fatal(msg string) {
	app.NewLogger().Fatal().Str("trace_id", emptyTraceIDPlaceholder).Msg(msg)
}

func Fatalf(format string, args ...interface{}) {
	app.NewLogger().Fatal().Str("trace_id", emptyTraceIDPlaceholder).Msgf(format, args...)
}

func FatalCtx(ctx context.Context, msg string) {
	withTraceID(ctx, app.NewLogger().Fatal()).Msg(msg)
}

func FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	withTraceID(ctx, app.NewLogger().Error()).Msgf(format, args...)
}
