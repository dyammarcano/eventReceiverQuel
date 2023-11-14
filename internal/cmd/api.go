package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dyammarcano/template-go/internal/helpers"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"runtime/trace"
)

type (
	ExternalAPIResponse struct {
		Args    map[string]interface{} `json:"args"`
		Headers map[string]interface{} `json:"headers"`
		Origin  string                 `json:"origin"`
		URL     string                 `json:"url"`
	}

	ExternalAPIResponseError struct {
		Error string `json:"error"`
	}
)

func CallExternalAPI(cmd *cobra.Command, _ []string) error {
	defer trace.StartRegion(cmd.Context(), "CallExternalAPI").End()
	// clean console
	fmt.Fprintf(cmd.OutOrStdout(), "\033[H\033[2J")

	apiUrl, err := helpers.GetString("url")
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(cmd.Context(), "GET", apiUrl, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	defer trace.StartRegion(cmd.Context(), "client.Do").End()
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer trace.StartRegion(cmd.Context(), "response.Body.Close").End()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		var errResponse ExternalAPIResponseError
		if err = json.Unmarshal(data, &errResponse); err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), errResponse.Error)
		return nil
	}

	var apiResponse ExternalAPIResponse

	defer trace.StartRegion(cmd.Context(), "json.Unmarshal").End()
	if err = json.Unmarshal(data, &apiResponse); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "your public ip: %s", apiResponse.Origin)
	return nil
}
