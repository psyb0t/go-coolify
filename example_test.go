package coolify_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-coolify"
)

// This is the README usage example. It needs an instance and credentials.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
	auth := coolify.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("COOLIFY_API_TOKEN"),
			)

			return nil
		},
	)

	client, err := coolify.NewClientWithResponses(
		os.Getenv("COOLIFY_BASE_URL"),
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ListProjectsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
