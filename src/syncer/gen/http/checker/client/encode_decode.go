// Code generated by goa v3.13.2, DO NOT EDIT.
//
// checker HTTP client encoders and decoders
//
// Command:
// $ goa gen syncer/design

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	checker "syncer/gen/checker"
	checkerviews "syncer/gen/checker/views"

	goahttp "goa.design/goa/v3/http"
)

// BuildGetRequest instantiates a HTTP request object with method and path set
// to call the "checker" service "get" endpoint
func (c *Client) BuildGetRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetCheckerPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("checker", "get", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetResponse returns a decoder for responses returned by the checker
// get endpoint. restoreBody controls whether the response body should be
// restored after having been read.
func DecodeGetResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("checker", "get", err)
			}
			p := NewGetSyncOK(&body)
			view := "default"
			vres := &checkerviews.Sync{Projected: p, View: view}
			if err = checkerviews.ValidateSync(vres); err != nil {
				return nil, goahttp.ErrValidationError("checker", "get", err)
			}
			res := checker.NewSync(vres)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("checker", "get", resp.StatusCode, string(body))
		}
	}
}
