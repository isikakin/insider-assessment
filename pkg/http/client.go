package http

import (
	"assesment/pkg/log"
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type ClientDecorator interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
type httpClientDecorator struct {
	httpClient http.Client
	log        log.Logger
}

func (h *httpClientDecorator) Do(ctx context.Context, req *http.Request) (*http.Response, error) {

	var (
		err        error
		response   *http.Response
		field      log.Field
		contents   []byte
		statusCode = 500
	)

	req.Header.Add("Content-Type", "application/json")
	if ctx.Value("correlationId") != nil {
		req.Header.Add("X-Correlation-ID", ctx.Value("correlationId").(string))
		req.Header.Add("X-CorrelationId", ctx.Value("correlationId").(string))
		req.Header.Add("CorrelationId", ctx.Value("correlationId").(string))
		req.Header.Add(fiber.HeaderXRequestID, ctx.Value("correlationId").(string))
	}

	if req.Body != nil {
		contents, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(contents))
	}

	reqRelativeUrl := req.URL.Path + "?" + req.URL.RawQuery
	field = log.Field{Message: "Starting request",
		RequestBody:    string(contents),
		Duration:       0,
		HttpStatusCode: 102,
		Url:            reqRelativeUrl,
		HostName:       req.URL.Host,
		HttpMethod:     req.Method}
	h.log.Request(ctx, field)
	timer := time.Now()

	response, err = h.httpClient.Do(req)

	elapsed := time.Since(timer)
	contents = []byte{}
	if response != nil && response.Body != nil {
		contents, _ = ioutil.ReadAll(response.Body)
		response.Body = ioutil.NopCloser(bytes.NewReader(contents))
		statusCode = response.StatusCode
	}
	field = log.Field{Message: "Finished request",
		ResponseBody:   string(contents),
		Duration:       elapsed.Seconds(),
		HttpStatusCode: statusCode,
		Url:            reqRelativeUrl,
		HostName:       req.URL.Host,
		HttpMethod:     req.Method}

	if statusCode >= 200 && statusCode < 300 {
		h.log.Response(ctx, field)
	} else {
		h.log.ResponseWithLevel(ctx, field, logrus.WarnLevel)
	}

	return response, err
}

func NewHttpClientDecorator(log log.Logger, timeOutInSecond int) ClientDecorator {
	timeout := time.Second * time.Duration(timeOutInSecond)
	return &httpClientDecorator{
		httpClient: http.Client{Timeout: timeout},
		log:        log,
	}
}
