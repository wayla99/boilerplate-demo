// Package player_spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package player_spec

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error     *string `json:"error,omitempty"`
	ErrorCode *string `json:"error_code,omitempty"`
}

// PlayerRequest defines model for PlayerRequest.
type PlayerRequest struct {
	Address *string `json:"address,omitempty"`
	Name    string  `json:"name"`
}

// PlayerResponse defines model for PlayerResponse.
type PlayerResponse struct {
	Address  *string `json:"address,omitempty"`
	Name     *string `json:"name,omitempty"`
	PlayerId *int    `json:"player_id,omitempty"`
}

// PostPlayerJSONRequestBody defines body for PostPlayer for application/json ContentType.
type PostPlayerJSONRequestBody = PlayerRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create player
	// (POST /player)
	PostPlayer(c *fiber.Ctx) error
	// Get player
	// (GET /player/{player_id})
	GetPlayerPlayerId(c *fiber.Ctx, playerId int) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// PostPlayer operation middleware
func (siw *ServerInterfaceWrapper) PostPlayer(c *fiber.Ctx) error {

	return siw.Handler.PostPlayer(c)
}

// GetPlayerPlayerId operation middleware
func (siw *ServerInterfaceWrapper) GetPlayerPlayerId(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "player_id" -------------
	var playerId int

	err = runtime.BindStyledParameterWithOptions("simple", "player_id", c.Params("player_id"), &playerId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter player_id: %w", err).Error())
	}

	return siw.Handler.GetPlayerPlayerId(c, playerId)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(m)
	}

	router.Post(options.BaseURL+"/player", wrapper.PostPlayer)

	router.Get(options.BaseURL+"/player/:player_id", wrapper.GetPlayerPlayerId)

}

type PostPlayerRequestObject struct {
	Body *PostPlayerJSONRequestBody
}

type PostPlayerResponseObject interface {
	VisitPostPlayerResponse(ctx *fiber.Ctx) error
}

type PostPlayer201TextResponse string

func (response PostPlayer201TextResponse) VisitPostPlayerResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "text/plain")
	ctx.Status(201)

	_, err := ctx.WriteString(string(response))
	return err
}

type PostPlayer400JSONResponse ErrorResponse

func (response PostPlayer400JSONResponse) VisitPostPlayerResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type GetPlayerPlayerIdRequestObject struct {
	PlayerId int `json:"player_id"`
}

type GetPlayerPlayerIdResponseObject interface {
	VisitGetPlayerPlayerIdResponse(ctx *fiber.Ctx) error
}

type GetPlayerPlayerId200JSONResponse PlayerResponse

func (response GetPlayerPlayerId200JSONResponse) VisitGetPlayerPlayerIdResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetPlayerPlayerId400JSONResponse ErrorResponse

func (response GetPlayerPlayerId400JSONResponse) VisitGetPlayerPlayerIdResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Create player
	// (POST /player)
	PostPlayer(ctx context.Context, request PostPlayerRequestObject) (PostPlayerResponseObject, error)
	// Get player
	// (GET /player/{player_id})
	GetPlayerPlayerId(ctx context.Context, request GetPlayerPlayerIdRequestObject) (GetPlayerPlayerIdResponseObject, error)
}

type StrictHandlerFunc func(ctx *fiber.Ctx, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// PostPlayer operation middleware
func (sh *strictHandler) PostPlayer(ctx *fiber.Ctx) error {
	var request PostPlayerRequestObject

	var body PostPlayerJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.PostPlayer(ctx.UserContext(), request.(PostPlayerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostPlayer")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(PostPlayerResponseObject); ok {
		if err := validResponse.VisitPostPlayerResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetPlayerPlayerId operation middleware
func (sh *strictHandler) GetPlayerPlayerId(ctx *fiber.Ctx, playerId int) error {
	var request GetPlayerPlayerIdRequestObject

	request.PlayerId = playerId

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetPlayerPlayerId(ctx.UserContext(), request.(GetPlayerPlayerIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetPlayerPlayerId")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetPlayerPlayerIdResponseObject); ok {
		if err := validResponse.VisitGetPlayerPlayerIdResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xTzW7bMAx+FYPb0aidbSff1q0ochgWdMchKDSJSVTEkkrRQ43A7z5QcrL8OD0MPewS",
	"KCSlj9+Pd6B9G7xDxxGaHUS9wVal4x2RpweMwbuIUgjkAxJbTG2Uthy4DwgNRCbr1jCUufOovcGJ9lDu",
	"K/7XE2qWC4ut6pEe8LnDyJdIyhjCGCexnGqvoBA+d5bQQPMzTy1fQb5G8h+gSwjp0UdrjrrWMa6RpvhL",
	"ybqVT9OWt9L7evfte/F5MYcSfiNF6x00MLupb2oB8AGdChYa+JhKJQTFm7RllcETE5/FFD6KrXdzAw0s",
	"fOTMGrJGGPnWm14mtXeMLl1SIWytTteqpyjw+2zI6T3hChp4V/0NTzUmpzo1czi1gqnDVMiCp5U/1LMz",
	"cMYXFiL2DNZg1GQDZzky00LULPDlpsj/Z1BehGEop69qQsVoithpjTGuuu22F30/1fWbyXH6FU3scqtM",
	"QQe1Sohd2yrqoYEvab2Rl9BS6yhpHgtLmR79rnaH0A2y0RonrL/H0fn8OzcpOKRaZCR5eQc2qaN4A/t8",
	"H6X53MjySIKLmC8vXK7fPGLXRf1xcLSgw9h/Y+w98muuDsOfAAAA//9bavjhlgUAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}