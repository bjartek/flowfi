package bot

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httputil"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func HttpGet[T any](ctx context.Context, url string, zapLogger *zap.Logger) (T, error) {
	var m T
	request, err := retryablehttp.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return m, err
	}
	request.Header.Add("Content-Type", "application/json")

	dump, err := httputil.DumpRequestOut(request.Request, false)
	if err != nil {
		zapLogger.With(zap.Error(err)).Warn("could not dump request")
	}

	zapLogger.Debug(string(dump))

	client := retryablehttp.NewClient()
	client.Logger = retryablehttp.LeveledLogger(&LeveledZap{zapLogger})

	res, err := client.Do(request)
	if err != nil {
		zapLogger.With(zap.Error(err)).Warn("could not perform http request")
		return m, errors.Wrap(err, "could not perfom http request")
	}
	zapLogger = zapLogger.With(zap.String("status", res.Status))

	responseDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		zapLogger.With(zap.Error(err)).Warn("could not dump request")
	}
	zapLogger.Debug(string(responseDump))
	zapLogger.Debug("http status")
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		zapLogger.With(zap.Error(err)).Warn("could not read body")
		return m, errors.Wrap(err, "could not read body")
	}
	zapLogger = zapLogger.With(zap.String("result", string(body)))

	result, err := parseJSON[T](body)
	if err != nil {
		zapLogger.With(zap.Error(err)).Warn("could not marshal body as json")
		return m, errors.Wrap(err, "could not marshal json")
	}
	return result, nil
}

func parseJSON[T any](s []byte) (T, error) {
	var r T

	if err := json.Unmarshal(s, &r); err != nil {
		return r, err
	}

	return r, nil
}

/*
func toJSON(T any) ([]byte, error) {
	return json.Marshal(T)
}
*/

type LeveledZap struct {
	logger *zap.Logger
}

func (l *LeveledZap) Error(msg string, keysAndValues ...interface{}) {
	l.logger.With(fields(keysAndValues)...).Error(msg)
}

func (l *LeveledZap) Info(msg string, keysAndValues ...interface{}) {
	l.logger.With(fields(keysAndValues)...).Info(msg)
}

func (l *LeveledZap) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.With(fields(keysAndValues)...).Debug(msg)
}

func (l *LeveledZap) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.With(fields(keysAndValues)...).Warn(msg)
}

func fields(keysAndValues []interface{}) []zap.Field {
	fields := []zap.Field{}

	for i := 0; i < len(keysAndValues)-1; i += 2 {
		fields = append(fields, zap.Any(keysAndValues[i].(string), keysAndValues[i+1]))
	}

	return fields
}
