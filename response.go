package costgraph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maxResponseSize = 8 << 20

type APIError struct {
	StatusCode int
	Operation  string
	Code       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: CostGraph API returned HTTP %d %s", e.Operation, e.StatusCode, http.StatusText(e.StatusCode))
}

type requestError struct {
	operation string
	cause     error
}

func (e *requestError) Error() string {
	return e.operation + ": CostGraph API request failed"
}

func (e *requestError) Unwrap() error {
	return e.cause
}

func decode[T any](operation string, response *http.Response, err error) (T, error) {
	var value T
	if err != nil {
		return value, &requestError{operation: operation, cause: err}
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		apiError := &APIError{StatusCode: response.StatusCode, Operation: operation}
		var body struct {
			Code string `json:"code"`
		}
		if json.NewDecoder(io.LimitReader(response.Body, 4096)).Decode(&body) == nil {
			switch body.Code {
			case "not_found", "virtual_tag_exists", "virtual_tag_revision_conflict", "virtual_tag_exported", "virtual_tag_invalid", "virtual_tag_error":
				apiError.Code = body.Code
			}
		}
		return value, apiError
	}
	if response.StatusCode == http.StatusNoContent {
		return value, nil
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return value, &requestError{operation: operation, cause: err}
	}
	if len(body) > maxResponseSize {
		return value, fmt.Errorf("%s: CostGraph API response exceeds the size limit", operation)
	}
	var envelope struct {
		Status string `json:"status"`
	}
	if json.Unmarshal(body, &envelope) != nil || envelope.Status != "success" || json.Unmarshal(body, &value) != nil {
		return value, fmt.Errorf("%s: invalid CostGraph API response", operation)
	}
	return value, nil
}
