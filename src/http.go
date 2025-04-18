package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func extractErrorRate(reader io.Reader) int {
	re := regexp.MustCompile(`(\d+)]]$`)
	body, err := io.ReadAll(reader)
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading HTTP body: %s", err))
		return 0
	}
	str := string(body)
	matches := re.FindStringSubmatch(str)
	value, err := strconv.Atoi(matches[1])
	if err == nil {
		return value
	}
	return 0
}

func probeHTTP(target string, w http.ResponseWriter, module Module) (success bool) {
	config := module.HTTP

	client := &http.Client{
		Timeout: module.Timeout,
	}

	requestURL, err := url.JoinPath(
		os.Getenv("SENTRY_URL"),
		"/api/0/projects/",
		os.Getenv("SENTRY_ORG"),
		target,
		"/stats/",
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating request URL for target %s: %s", target, err))
	}

	slog.Info(fmt.Sprintf("Request URL: %s", requestURL))
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating request for target %s: %s", target, err))
		return
	}

	caser := cases.Title(language.English)
	for key, value := range config.Headers {
		if caser.String(key) == "Host" {
			request.Host = value
			continue
		} else if caser.String(key) == "Authorization" {
			// Skip auth headers. We are addiing it manually
			continue
		}

		request.Header.Set(key, value)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SENTRY_AUTH_TOKEN")))

	resp, err := client.Do(request)
	// Err won't be nil if redirects were turned off. See https://github.com/golang/go/issues/3795
	if err != nil && resp == nil {
		slog.Warn(fmt.Sprintf("Error for HTTP request to %s: %s", target, err))
	} else {
		defer resp.Body.Close()
		if len(config.ValidStatusCodes) != 0 {
			if slices.Contains(config.ValidStatusCodes, resp.StatusCode) {
				success = true
			}
			if slices.Contains(config.ValidStatusCodes, resp.StatusCode) {
				success = true
			}
		} else if 200 <= resp.StatusCode && resp.StatusCode < 300 {
			success = true
		}
		if success {
			fmt.Fprintf(w, "probe_sentry_error_received %d\n", extractErrorRate(resp.Body))
		}
	}
	if resp == nil {
		resp = &http.Response{}
	}

	fmt.Fprintf(w, "probe_sentry_status_code %d\n", resp.StatusCode)
	fmt.Fprintf(w, "probe_sentry_content_length %d\n", resp.ContentLength)

	return
}
