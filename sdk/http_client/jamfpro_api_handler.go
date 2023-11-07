// jamfpro_api_handler.go
/* ------------------------------Summary----------------------------------------
This is a api handler module for the http_client to accommodate specifics of
jamf's api(s). It handles the encoding (marshalling) and decoding (unmarshalling)
of data. It also sets the correct content headers for the various http methods.

This module integrates with the http_client logger for wrapped error handling
for human readable return codes. It also supports the http_clients debugMode for
verbose logging.

The logic of this module is defined as follows:
Classic API:

For requests (GET, POST, PUT, DELETE):
- Encoding (Marshalling): Use XML format.
For responses (GET, POST, PUT):
- Decoding (Unmarshalling): Use XML format.
For responses (DELETE):
- Handle response codes as response body lacks anything useful.
Headers
- Set content header as application/xml

JamfPro API:

For requests (GET, POST, PUT, DELETE):
- Encoding (Marshalling): Use JSON format.
For responses (GET, POST, PUT):
- Decoding (Unmarshalling): Use JSON format.
For responses (DELETE):
- Handle response codes as response body lacks anything useful.
Headers
- Set content header as application/json
*/
package http_client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// Endpoint constants represent the URL suffixes used for Jamf API token interactions.
const (
	BaseDomain              = ".jamfcloud.com"                // BaseDomain represents the base domain for the jamf instance.
	OAuthTokenEndpoint      = "/api/oauth/token"              // OAuthTokenEndpoint: The endpoint to obtain an OAuth token.
	BearerTokenEndpoint     = "/api/v1/auth/token"            // BearerTokenEndpoint: The endpoint to obtain a bearer token.
	TokenRefreshEndpoint    = "/api/v1/auth/keep-alive"       // TokenRefreshEndpoint: The endpoint to refresh an existing token.
	TokenInvalidateEndpoint = "/api/v1/auth/invalidate-token" // TokenInvalidateEndpoint: The endpoint to invalidate an active token.
)

// ClassicApiHandler handles the specifics of the Classic API.
type ClassicApiHandler struct {
	logger    Logger
	debugMode bool
}

// JamfProApiHandler handles the specifics of the JamfPro API.
type JamfProApiHandler struct {
	logger    Logger
	debugMode bool
}

// UnknownApiHandler provides default behavior for unrecognized API types.
type UnknownApiHandler struct {
	logger    Logger
	debugMode bool
}

// SetLogger assigns a logger instance to the ClassicApiHandler.
func (h *ClassicApiHandler) SetLogger(logger Logger) {
	h.logger = logger
}

// SetLogger assigns a logger instance to the JamfProApiHandler.
func (h *JamfProApiHandler) SetLogger(logger Logger) {
	h.logger = logger
}

// SetLogger assigns a logger instance to the UnknownApiHandler.
func (h *UnknownApiHandler) SetLogger(logger Logger) {
	h.logger = logger
}

func (h *ClassicApiHandler) SetDebugMode(debug bool) {
	h.debugMode = debug
}

func (h *JamfProApiHandler) SetDebugMode(debug bool) {
	h.debugMode = debug
}

func (h *UnknownApiHandler) SetDebugMode(debug bool) {
	h.debugMode = debug
}

// ConstructAPIResourceEndpoint returns the full URL for a Jamf API resource endpoint path.
func (c *Client) ConstructAPIResourceEndpoint(endpointPath string) string {
	url := fmt.Sprintf("https://%s%s%s", c.InstanceName, BaseDomain, endpointPath)
	c.logger.Info("Request will be made to API Resource URL:", "URL", url)
	return url
}

// ConstructAPIAuthEndpoint returns the full URL for a Jamf API auth endpoint path.
func (c *Client) ConstructAPIAuthEndpoint(endpointPath string) string {
	url := fmt.Sprintf("https://%s%s%s", c.InstanceName, BaseDomain, endpointPath)
	c.logger.Info("Request will be made to API authentication URL:", "URL", url)
	return url
}

// APIHandler is an interface for encoding, decoding, and determining content types for different API implementations.
// It encapsulates behavior for encoding and decoding requests and responses.
type APIHandler interface {
	MarshalRequest(body interface{}, method string) ([]byte, error)
	MarshalMultipartRequest(fields map[string]string, files map[string]string) ([]byte, string, error) // New method for multipart
	UnmarshalResponse(resp *http.Response, out interface{}) error
	GetContentTypeHeader(method string) string
	SetLogger(logger Logger)
	SetDebugMode(debug bool)
	GetAcceptHeader(url string) string
}

// GetContentTypeHeader for ClassicApiHandler always returns XML as the content type.
func (h *ClassicApiHandler) GetContentTypeHeader(method string) string {
	return "application/xml"
}

// GetContentTypeHeader for JamfProApiHandler always returns JSON as the content type.
func (h *JamfProApiHandler) GetContentTypeHeader(method string) string {
	return "application/json"
}

func (h *UnknownApiHandler) GetContentTypeHeader(method string) string {
	// For an unknown API handler, defaults to JSON handling behavior.
	return "application/json"
}

// GetAPIHandler determines the appropriate APIHandler based on the endpoint.
// It identifies the type of API (Classic, JamfPro, or Unknown) and returns the corresponding handler.
func GetAPIHandler(endpoint string, logLevel LogLevel) APIHandler {
	var handler APIHandler
	if strings.Contains(endpoint, "/JSSResource") {
		handler = &ClassicApiHandler{}
	} else if strings.Contains(endpoint, "/api") {
		handler = &JamfProApiHandler{}
	} else {
		handler = &UnknownApiHandler{}
	}
	logger := NewDefaultLogger()
	logger.SetLevel(logLevel)
	handler.SetLogger(logger)
	return handler
}

// MarshalRequest encodes the request body in XML format for the Classic API.
func (h *ClassicApiHandler) MarshalRequest(body interface{}, method string) ([]byte, error) {
	data, err := xml.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Log the full request body for POST and PUT methods
	if method == "POST" || method == "PUT" {
		h.logger.Trace("Full Request Body:", "Body", string(data))
	}

	return data, nil
}

// UnmarshalResponse decodes the response body from XML format for the Classic API.
func (h *ClassicApiHandler) UnmarshalResponse(resp *http.Response, out interface{}) error {
	// Handle DELETE method
	if resp.Request.Method == "DELETE" {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		} else {
			return fmt.Errorf("DELETE request failed with status code: %d", resp.StatusCode)
		}
	}
	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("Failed reading response body", "error", err)
		return err
	}

	// Log the raw response body and headers
	h.logger.Trace("Raw HTTP Response for Classic API:", string(bodyBytes))
	h.logger.Debug("Unmarshaling response for Classic API", "status", resp.Status)

	// Log headers when in debug mode
	h.logger.Debug("HTTP Response Headers:", resp.Header)

	// Check the Content-Type header
	contentType := resp.Header.Get("Content-Type")

	// If content type is HTML
	if strings.Contains(contentType, "text/html") {
		errMsg := extractErrorMessageFromHTML(string(bodyBytes))
		h.logger.Warn("Received HTML content", "error_message", errMsg, "status_code", resp.StatusCode)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    errMsg,
		}
	}

	// If content type is XML
	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml;charset=UTF-8") {
		// Check the status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			h.logger.Error("Received non-success status code", "status_code", resp.StatusCode)
			return fmt.Errorf("received non-success status code: %d", resp.StatusCode)
		}

		// Try to unmarshal the XML response
		err = xml.Unmarshal(bodyBytes, out)
		if err != nil {
			// If unmarshalling fails, check if the content might be HTML
			if strings.Contains(string(bodyBytes), "<html>") {
				errMsg := extractErrorMessageFromHTML(string(bodyBytes))
				h.logger.Warn("Received HTML content instead of expected XML", "error_message", errMsg, "status_code", resp.StatusCode)
				return fmt.Errorf(errMsg)
			}

			h.logger.Error("Failed to unmarshal XML response", "error", err)
			return fmt.Errorf("failed to unmarshal XML response: %v", err)
		}
	} else {
		// If the content type is neither XML nor HTML
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	return nil
}

// MarshalRequest encodes the request body in JSON format for the JamfPro API.
func (h *JamfProApiHandler) MarshalRequest(body interface{}, method string) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		h.logger.Error("Failed marshaling request for JamfPro API", "error", err)
		return nil, err
	}

	// If in debug mode and the method is either POST (Create) or PUT (Update), log the full request body
	if h.debugMode {
		h.logger.Debug("Marshaling request for JamfPro API", "method", method)
		if method == "POST" || method == "PUT" || method == "PATCH" {
			h.logger.Debug("Full Request Body for JamfPro API:", string(data))
		}
	}

	return data, nil
}

// UnmarshalResponse decodes the response body from JSON format for the JamfPro API.
func (h *JamfProApiHandler) UnmarshalResponse(resp *http.Response, out interface{}) error {
	// Handle DELETE method
	if resp.Request.Method == "DELETE" {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		} else {
			return fmt.Errorf("DELETE request failed with status code: %d", resp.StatusCode)
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("Failed reading response body for JamfPro API", "error", err)
		return err
	}

	// Log the raw response body and headers
	h.logger.Trace("Raw HTTP Response for Jamf Pro API:", string(bodyBytes))
	h.logger.Debug("Unmarshaling response for Jamf Pro API", "status", resp.Status)

	// Log headers when in debug mode
	h.logger.Debug("HTTP Response Headers:", resp.Header)

	// Check the Content-Type header
	contentType := resp.Header.Get("Content-Type")

	// If content type is HTML
	if strings.Contains(contentType, "text/html") {
		errMsg := extractErrorMessageFromHTML(string(bodyBytes))
		h.logger.Warn("Received HTML content", "error_message", errMsg, "status_code", resp.StatusCode)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    errMsg,
		}
	}

	// If content type is JSON
	if strings.Contains(contentType, "application/json") {
		// Check the status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			h.logger.Error("Received non-success status code", "status_code", resp.StatusCode)
			return fmt.Errorf("received non-success status code: %d", resp.StatusCode)
		}

		// Try to unmarshal the JSON response
		err = json.Unmarshal(bodyBytes, out)
		if err != nil {
			// If unmarshalling fails, check if the content might be HTML
			if strings.Contains(string(bodyBytes), "<html>") {
				errMsg := extractErrorMessageFromHTML(string(bodyBytes))
				h.logger.Warn("Received HTML content instead of expected JSON", "error_message", errMsg, "status_code", resp.StatusCode)
				return fmt.Errorf(errMsg)
			}

			h.logger.Error("Failed to unmarshal JSON response", "error", err)
			return fmt.Errorf("failed to unmarshal JSON response: %v", err)
		}
	} else {
		// If the content type is neither JSON nor HTML
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	return nil
}

// MarshalRequest returns an error since the API type is unsupported.
func (h *UnknownApiHandler) MarshalRequest(body interface{}, method string) ([]byte, error) {
	h.logger.Warn("Attempted to marshal request for an unsupported API type")

	return nil, fmt.Errorf("unsupported API type")
}

// UnmarshalResponse returns an error since the API type is unsupported.
func (h *UnknownApiHandler) UnmarshalResponse(resp *http.Response, out interface{}) error {
	h.logger.Warn("Attempted to unmarshal response for an unsupported API type", "status", resp.Status)
	return fmt.Errorf("unsupported API type")
}

// MarshalMultipartRequest marshals a request with multipart form data for the JamfPro API.
func (h *JamfProApiHandler) MarshalMultipartRequest(fields map[string]string, files map[string]string) ([]byte, string, error) {
	h.logger.Debug("Marshaling multipart request for JamfPro API")
	return MarshalMultipartFormData(fields, files)
}

// MarshalMultipartFormData takes a map with form fields and file paths and returns the encoded body and content type.
func MarshalMultipartFormData(fields map[string]string, files map[string]string) ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the simple fields to the form data
	for field, value := range fields {
		if err := writer.WriteField(field, value); err != nil {
			return nil, "", err
		}
	}

	// Add the files to the form data
	for formField, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(formField, filepath)
		if err != nil {
			return nil, "", err
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, "", err
		}
	}

	// Close the writer before returning
	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body.Bytes(), contentType, nil
}

// MarshalMultipartRequest is not supported by the Classic API and will return an error.
func (h *ClassicApiHandler) MarshalMultipartRequest(fields map[string]string, files map[string]string) ([]byte, string, error) {
	return nil, "", fmt.Errorf("multipart request not supported by the Classic API")
}

// MarshalMultipartRequest is not supported for an unknown API type and will return an error.
func (h *UnknownApiHandler) MarshalMultipartRequest(fields map[string]string, files map[string]string) ([]byte, string, error) {
	return nil, "", fmt.Errorf("multipart request not supported for an unknown API type")
}

// GetAcceptHeader returns the correct Accept header based on the URL endpoint.
// If the URL matches certain patterns that are exceptions, a specific Accept header is returned.
func (h *JamfProApiHandler) GetAcceptHeader(url string) string {
	// List of endpoints that should have the "image/*" Accept header
	imageEndpoints := []string{
		"/api/v1/branding-images/download/",
		"/api/v1/icon/download/",
	}

	// Check if the requested URL is one of the image endpoints
	for _, endpoint := range imageEndpoints {
		if strings.Contains(url, endpoint) {
			return "image/*" // Return "image/*" for image download endpoints
		}
	}

	// Default to "application/json" for all other endpoints
	return "application/json"
}

// GetAcceptHeader returns the Accept header value for requests made through the Classic API.
// Since the Classic API does not have exceptions like the Jamf Pro API for certain endpoints,
// it defaults to the standard content type, which is "application/xml".
func (h *ClassicApiHandler) GetAcceptHeader(url string) string {
	return h.GetContentTypeHeader("") // Returns the default "application/xml"
}

// GetAcceptHeader returns the Accept header value for requests made through an unknown API type.
// Since the handler does not know the specifics of the API, it defaults to "application/json",
// which is a reasonable default for modern web APIs.
func (h *UnknownApiHandler) GetAcceptHeader(url string) string {
	return h.GetContentTypeHeader("") // Returns the default "application/json"
}
