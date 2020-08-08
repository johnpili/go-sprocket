package sprocket

import (
	"fmt"
	"net/http"
)

// VerifyHTTPGetHasCode a helper method that can be used to validate required HTTP response code
func VerifyHTTPGetHasCode(client *http.Client, url string, statusCode int) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != statusCode {
		return fmt.Errorf("Expecting %d but got %d", statusCode, response.StatusCode)
	}
	return nil
}
