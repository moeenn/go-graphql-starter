package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

type GraphqlRequest struct {
	Query string `json:"query"`
}

type GraphqlOperation struct {
	Name string
	Type string
}

func parseGraphqlAST(graphqlAst *ast.Document) []GraphqlOperation {
	operations := []GraphqlOperation{}
	for _, def := range graphqlAst.Definitions {
		if op, ok := def.(*ast.OperationDefinition); ok {
			//nolint:exhaustruct
			operation := GraphqlOperation{
				Type: op.Operation,
			}

			for _, sel := range op.SelectionSet.Selections {
				if field, ok := sel.(*ast.Field); ok {
					name := field.Name.Value
					operation.Name = name
				}
			}

			operations = append(operations, operation)
		}
	}

	return operations
}

func readRequestBody(r *http.Request) (*GraphqlRequest, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("request method is not POST")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	var graphqlRequest GraphqlRequest
	if err := json.Unmarshal(bodyBytes, &graphqlRequest); err != nil {
		return nil, fmt.Errorf("failed to read request body as json: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return &graphqlRequest, nil
}

// the returned bool value of true means the request should be forwarded to the
// graphql relay.Handler instead.
func parseGraphqlRequestOperations(r *http.Request) ([]GraphqlOperation, bool, error) {
	requestBody, err := readRequestBody(r)
	if err != nil {
		return nil, false, err
	}

	//nolint:exhaustruct
	source := source.NewSource(&source.Source{
		Body: []byte(requestBody.Query),
	})

	//nolint:exhaustruct
	graphqlAst, err := parser.Parse(parser.ParseParams{
		Source: source,
	})

	// if the graphql request body is invalid, we let the graphql relay.Handler
	// report back the errors.
	if err != nil {
		return nil, true, err
	}

	operations := parseGraphqlAST(graphqlAst)
	return operations, false, nil
}

func allOperationsWhitelisted(operations []GraphqlOperation, whitelist []string) bool {
	for _, op := range operations {
		if !slices.Contains(whitelist, op.Name) {
			return false
		}
	}

	return true
}
