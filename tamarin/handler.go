package tamarin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	VARIABLE_INDICATOR = "{}"
	STATIC_INDICATOR   = "{*}"
)

type handler struct {
	verbose                bool
	handleFuncsGET         map[string][]http.HandlerFunc
	handleFuncsPOST        map[string][]http.HandlerFunc
	handleFuncsPATCH       map[string][]http.HandlerFunc
	variableHandlersGET    map[string][]http.HandlerFunc
	variableHandlersPOST   map[string][]http.HandlerFunc
	variableHandlersPATCH  map[string][]http.HandlerFunc
	staticHandlersGET      map[string][]http.HandlerFunc
	staticHandlersPOST     map[string][]http.HandlerFunc
	staticHandlersPATCH    map[string][]http.HandlerFunc
	handleFuncsDELETE      map[string][]http.HandlerFunc
	variableHandlersDELETE map[string][]http.HandlerFunc
	staticHandlersDELETE   map[string][]http.HandlerFunc
}

// NewHandler returns a fresh Handler / Mux.
// The verbose parameter controls log output
func NewHandler(verbose bool) *handler {
	return &handler{
		verbose:                verbose,
		handleFuncsGET:         make(map[string][]http.HandlerFunc),
		handleFuncsPOST:        make(map[string][]http.HandlerFunc),
		variableHandlersGET:    make(map[string][]http.HandlerFunc),
		variableHandlersPOST:   make(map[string][]http.HandlerFunc),
		staticHandlersGET:      make(map[string][]http.HandlerFunc),
		staticHandlersPOST:     make(map[string][]http.HandlerFunc),
		handleFuncsPATCH:       make(map[string][]http.HandlerFunc),
		variableHandlersPATCH:  make(map[string][]http.HandlerFunc),
		staticHandlersPATCH:    make(map[string][]http.HandlerFunc),
		handleFuncsDELETE:      make(map[string][]http.HandlerFunc),
		variableHandlersDELETE: make(map[string][]http.HandlerFunc),
		staticHandlersDELETE:   make(map[string][]http.HandlerFunc),
	}
}

// Post builds a POST endpoint and adds it to the list of sequences
func (h *handler) Post(path string, handlers ...EndpointHandlerFunc) *handler {
	return h.withEndpoint(NewEndpoint(path).WithHandlers(handlers...).WithMethod(http.MethodPost))
}

// Get builds a GET endpoint and adds it to the list of sequences
func (h *handler) Get(path string, handlers ...EndpointHandlerFunc) *handler {
	return h.withEndpoint(NewEndpoint(path).WithHandlers(handlers...).WithMethod(http.MethodGet))
}

// Patch builds a PATCH endpoint and adds it to the list of sequences
func (h *handler) Patch(path string, handlers ...EndpointHandlerFunc) *handler {
	return h.withEndpoint(NewEndpoint(path).WithHandlers(handlers...).WithMethod(http.MethodPatch))
}

// Delete builds a DELETE endpoint and adds it to the list of sequences
func (h *handler) Delete(path string, handlers ...EndpointHandlerFunc) *handler {
	return h.withEndpoint(NewEndpoint(path).WithHandlers(handlers...).WithMethod(http.MethodDelete))
}

// PostF adds a POST handler to the list of sequences
func (h *handler) PostF(path string, handlers ...http.HandlerFunc) *handler {
	if pathIsVariable(path) {
		h.variableHandlersPOST[path] = handlers
	} else if pathIsStatic(path) {
		h.staticHandlersPOST[path] = handlers
	} else {
		h.handleFuncsPOST[path] = handlers
	}
	return h
}

// GetF adds a GET handler to the list of sequences
func (h *handler) GetF(path string, handlers ...http.HandlerFunc) *handler {
	if pathIsVariable(path) {
		h.variableHandlersGET[path] = handlers
	} else if pathIsStatic(path) {
		h.staticHandlersGET[path] = handlers
	} else {
		h.handleFuncsGET[path] = handlers
	}
	return h
}

// PatchF adds a PATCH handler to the list of sequences
func (h *handler) PatchF(path string, handlers ...http.HandlerFunc) *handler {
	if pathIsVariable(path) {
		h.variableHandlersPATCH[path] = handlers
	} else if pathIsStatic(path) {
		h.staticHandlersPATCH[path] = handlers
	} else {
		h.handleFuncsPATCH[path] = handlers
	}
	return h
}

// DeleteF adds a DELETE handler to the list of sequences
func (h *handler) DeleteF(path string, handlers ...http.HandlerFunc) *handler {
	if pathIsVariable(path) {
		h.variableHandlersDELETE[path] = handlers
	} else if pathIsStatic(path) {
		h.staticHandlersDELETE[path] = handlers
	} else {
		h.handleFuncsDELETE[path] = handlers
	}
	return h
}

// withEndpoint adds an Endpoint (HandlerFunc wrapper) to the list of HandlerFuncs to be
// executed for a given path and method
func (s *handler) withEndpoint(e *endpoint) *handler {
	if e == nil {
		return s
	}
	switch e.method {
	case http.MethodGet:
		if pathIsVariable(e.path) {
			s.variableHandlersGET[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStatic(e.path) {
			s.staticHandlersGET[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsGET[e.path] = []http.HandlerFunc{e.Handle}
		}
	case http.MethodPost:
		if pathIsVariable(e.path) {
			s.variableHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStatic(e.path) {
			s.staticHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsPOST[e.path] = []http.HandlerFunc{e.Handle}
		}
	case http.MethodPatch:
		if pathIsVariable(e.path) {
			s.variableHandlersPATCH[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStatic(e.path) {
			s.staticHandlersPATCH[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsPATCH[e.path] = []http.HandlerFunc{e.Handle}
		}
	case http.MethodDelete:
		if pathIsVariable(e.path) {
			s.variableHandlersDELETE[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStatic(e.path) {
			s.staticHandlersDELETE[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsDELETE[e.path] = []http.HandlerFunc{e.Handle}
		}
	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", e.method)
	}

	return s
}

// WithGetEndpoint adds a GET endpoint to the list of endpoints served
func (s *handler) WithGetEndpoint(e *endpoint) *handler {
	if e == nil {
		return s
	}
	e.method = http.MethodGet
	return s.withEndpoint(e)
}

// WithGetEndpoint adds a POST endpoint to the list of endpoints served
func (s *handler) WithPostEndpoint(e *endpoint) *handler {
	if e == nil {
		return s
	}
	e.method = http.MethodPost
	return s.withEndpoint(e)
}

// WithPatchEndpoint adds a PATCH endpoint to the list of endpoints served
func (s *handler) WithPatchEndpoint(e *endpoint) *handler {
	if e == nil {
		return s
	}
	e.method = http.MethodPatch
	return s.withEndpoint(e)
}

// WithDeleteEndpoint adds a DELETE endpoint to the list of endpoints served
func (s *handler) WithDeleteEndpoint(e *endpoint) *handler {
	if e == nil {
		return s
	}
	e.method = http.MethodDelete
	return s.withEndpoint(e)
}

// WithEndpoint adds HandlerFuncs to the list of HandlerFuncs to be
// executed for a given path and method
func (s *handler) WithHandleFuncs(path, httpMethod string, handlerFuncs ...http.HandlerFunc) *handler {
	switch httpMethod {
	case http.MethodGet:
		if pathIsVariable(path) {
			s.variableHandlersGET[path] = handlerFuncs
		} else if pathIsStatic(path) {
			s.staticHandlersGET[path] = handlerFuncs
		} else {
			s.handleFuncsGET[path] = handlerFuncs
		}
	case http.MethodPost:
		if pathIsVariable(path) {
			s.variableHandlersPOST[path] = handlerFuncs
		} else if pathIsStatic(path) {
			s.staticHandlersPOST[path] = handlerFuncs
		} else {
			s.handleFuncsPOST[path] = handlerFuncs
		}
	case http.MethodPatch:
		if pathIsVariable(path) {
			s.variableHandlersPATCH[path] = handlerFuncs
		} else if pathIsStatic(path) {
			s.staticHandlersPATCH[path] = handlerFuncs
		} else {
			s.handleFuncsPATCH[path] = handlerFuncs
		}
	case http.MethodDelete:
		if pathIsVariable(path) {
			s.variableHandlersDELETE[path] = handlerFuncs
		} else if pathIsStatic(path) {
			s.staticHandlersDELETE[path] = handlerFuncs
		} else {
			s.handleFuncsDELETE[path] = handlerFuncs
		}
	case http.MethodOptions:

	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", httpMethod)
	}
	return s
}

// HandlerNames returns the list of all items the Handler is handling
// Useful for startup log output
func (s *handler) HandlerNames() []string {
	names := []string{}
	for key := range s.handleFuncsGET {
		names = append(names, fmt.Sprintf("[%s]                                 -> %s", http.MethodGet, key))
	}
	for key := range s.variableHandlersGET {
		names = append(names, fmt.Sprintf("[%s] [URL contains variable]         -> %s", http.MethodGet, key))
	}
	for key := range s.staticHandlersGET {
		names = append(names, fmt.Sprintf("[%s] [URL refers to static content]  -> %s ", http.MethodGet, key))
	}
	for key := range s.handleFuncsPOST {
		names = append(names, fmt.Sprintf("[%s]                                -> %s", http.MethodPost, key))
	}
	for key := range s.variableHandlersPOST {
		names = append(names, fmt.Sprintf("[%s] [URL contains variable]        -> %s ", http.MethodPost, key))
	}
	for key := range s.staticHandlersPOST {
		names = append(names, fmt.Sprintf("[%s] [URL refers to static content] -> %s ", http.MethodPost, key))
	}
	for key := range s.handleFuncsPATCH {
		names = append(names, fmt.Sprintf("[%s]                                -> %s", http.MethodPatch, key))
	}
	for key := range s.variableHandlersPATCH {
		names = append(names, fmt.Sprintf("[%s] [URL contains variable]        -> %s ", http.MethodPatch, key))
	}
	for key := range s.staticHandlersPATCH {
		names = append(names, fmt.Sprintf("[%s] [URL refers to static content] -> %s ", http.MethodPatch, key))
	}
	for key := range s.handleFuncsDELETE {
		names = append(names, fmt.Sprintf("[%s]                                -> %s", http.MethodDelete, key))
	}
	for key := range s.variableHandlersDELETE {
		names = append(names, fmt.Sprintf("[%s] [URL contains variable]        -> %s ", http.MethodDelete, key))
	}
	for key := range s.staticHandlersDELETE {
		names = append(names, fmt.Sprintf("[%s] [URL refers to static content] -> %s ", http.MethodDelete, key))
	}
	return names
}

// ServeHTTP fulfills the http.Handler interface and is used along with http.ListenAndServe
func (s *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if rw == nil || req == nil || req.URL == nil {
		return
	}

	// TODO: make this customizable and more locked down
	// Temporary CORS passthrough
	if req.Method == http.MethodOptions {
		handleOptions(rw, req)
		return
	}

	// Set open CORS on all other requests
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "*")

	reqPath := req.URL.Path
	if s.verbose {
		log.Printf("Received request for '%s'", reqPath)
	}
	var endpoints []http.HandlerFunc
	var OK bool
	switch req.Method {
	case http.MethodGet:
		endpoints, OK = s.handleFuncsGET[reqPath]
	case http.MethodPost:
		endpoints, OK = s.handleFuncsPOST[reqPath]
	case http.MethodPatch:
		endpoints, OK = s.handleFuncsPATCH[reqPath]
	case http.MethodDelete:
		endpoints, OK = s.handleFuncsDELETE[reqPath]
	}
	if !OK {
		endpoints = s.getVariableHandlerFuncsForPattern(req.URL.Path, req.Method)
		if endpoints == nil {
			endpoints = s.getStaticHandlerFuncsForPattern(req.URL.Path, req.Method)
			if endpoints == nil {
				if s.verbose {
					log.Printf("don't have a handler for %s", reqPath)
				}
				rw.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}
	for _, endpoint := range endpoints {
		endpoint(rw, req)
	}
	if s.verbose {
		log.Printf("Handled request for '%s'", reqPath)
	}
}

func pathIsStatic(path string) bool {
	return strings.Contains(path, STATIC_INDICATOR)
}

func pathIsVariable(path string) bool {
	return strings.Contains(path, VARIABLE_INDICATOR)
}

func (h *handler) getVariableHandlerFuncsForPattern(path, httpMethod string) []http.HandlerFunc {
	var candidateFuncs map[string][]http.HandlerFunc
	switch httpMethod {
	case http.MethodGet:
		candidateFuncs = h.variableHandlersGET
	case http.MethodPost:
		candidateFuncs = h.variableHandlersPOST
	case http.MethodPatch:
		candidateFuncs = h.variableHandlersPATCH
	case http.MethodDelete:
		candidateFuncs = h.variableHandlersDELETE
	default:
		return nil
	}
	for candidatePath, handlers := range candidateFuncs {
		candidatePrefix := variablePrefix(candidatePath)
		if len(path) < len(candidatePrefix) {
			continue
		}
		if strings.EqualFold(candidatePrefix, path[:len(candidatePrefix)]) {
			candidateSplit := strings.Split(candidatePath, "/")
			inputSplit := strings.Split(path, "/")
			if len(candidateSplit) != len(inputSplit) {
				continue
			}
			allMatched := true
			for idx, element := range candidateSplit {
				if element == VARIABLE_INDICATOR {
					continue
				}
				allMatched = allMatched && strings.EqualFold(element, inputSplit[idx])
			}
			if allMatched {
				return handlers
			}
		}
	}
	return nil
}

func (h *handler) getStaticHandlerFuncsForPattern(path, httpMethod string) []http.HandlerFunc {
	var candidateFuncs map[string][]http.HandlerFunc
	switch httpMethod {
	case http.MethodGet:
		candidateFuncs = h.staticHandlersGET
	case http.MethodPost:
		candidateFuncs = h.staticHandlersPOST
	case http.MethodPatch:
		candidateFuncs = h.staticHandlersPATCH
	case http.MethodDelete:
		candidateFuncs = h.staticHandlersDELETE
	default:
		return nil
	}
	for candidatePath, handlers := range candidateFuncs {
		candidatePrefix := staticPrefix(candidatePath)
		if len(path) < len(candidatePrefix) {
			continue
		}
		if strings.EqualFold(candidatePrefix, path[:len(candidatePrefix)]) {
			return handlers
		}
	}
	return nil
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusNoContent)
}

func staticPrefix(input string) string {
	idx := strings.Index(input, STATIC_INDICATOR)
	if idx < 1 {
		return input
	}
	return input[:idx]
}

func variablePrefix(input string) string {
	idx := strings.Index(input, VARIABLE_INDICATOR)
	if idx < 1 {
		return input
	}
	return input[:idx]
}
