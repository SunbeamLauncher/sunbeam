package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pomdtr/sunbeam/schemas"
	"github.com/pomdtr/sunbeam/types"
)

type PageGenerator func() (*types.Page, error)

func NewStaticGenerator(reader io.Reader) PageGenerator {
	var pageRef *types.Page
	return func() (*types.Page, error) {
		if pageRef != nil {
			return pageRef, nil
		}

		b, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		if err := schemas.Validate(b); err != nil {
			return nil, err
		}

		var page types.Page
		if err := json.Unmarshal(b, &page); err != nil {
			return nil, err
		}

		p, err := expandPage(page, nil)
		if err != nil {
			return nil, err
		}

		pageRef = p
		return p, nil
	}
}

func NewCommandGenerator(command *types.Command) PageGenerator {
	return func() (*types.Page, error) {
		output, err := command.Output(context.TODO())
		if err != nil {
			return nil, err
		}

		if err := schemas.Validate(output); err != nil {
			return nil, err
		}

		var page types.Page
		if err := json.Unmarshal(output, &page); err != nil {
			return nil, err
		}

		p, err := expandPage(page, nil)
		if err != nil {
			return nil, err
		}

		return p, nil
	}
}

func NewRequestGenerator(request *types.Request) PageGenerator {
	return func() (*types.Page, error) {
		req, err := http.NewRequest(request.Method, request.Url, strings.NewReader(request.Body))
		if err != nil {
			return nil, err
		}

		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		bs, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err := schemas.Validate(bs); err != nil {
			return nil, err
		}

		var page types.Page
		if err := json.Unmarshal(bs, &page); err != nil {
			return nil, err
		}

		p, err := expandPage(page, &url.URL{
			Scheme: res.Request.URL.Scheme,
			Host:   res.Request.URL.Host,
			Path:   path.Dir(res.Request.URL.Path),
		})

		if err != nil {
			return nil, err
		}

		return p, nil
	}
}
