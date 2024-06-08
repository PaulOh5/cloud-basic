package network

import (
	"io"
	"net/http"
)

func getResponseBody(resp *http.Response) string {
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	resp.Body.Close()
	return string(output)
}
