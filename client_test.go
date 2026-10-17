package coolify_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-coolify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToken = "EXAMPLE-DO-NOT-USE"

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *coolify.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := coolify.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testToken)

			return nil
		},
	)

	client, err := coolify.NewClientWithResponses(srv.URL+"/api/v1", auth)
	require.NoError(t, err)

	return client
}

func assertBearerAuth(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer "+testToken, r.Header.Get("Authorization"))
}

func TestListProjects(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/projects", r.URL.Path)
		assertBearerAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`[{"id":7,"uuid":"project-1","name":"Demo"}]`))
		assert.NoError(t, err)
	})

	res, err := client.ListProjectsWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, *res.JSON200, 1)

	project := (*res.JSON200)[0]
	assert.Equal(t, 7, *project.Id)
	assert.Equal(t, "project-1", *project.Uuid)
	assert.Equal(t, "Demo", *project.Name)
}

func TestCreateProject(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/projects", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearerAuth(t, r)

		var body coolify.CreateProjectJSONRequestBody

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Demo", *body.Name)
		assert.Equal(t, "Test project", *body.Description)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"uuid":"project-1"}`))
		assert.NoError(t, err)
	})

	name := "Demo"
	description := "Test project"

	res, err := client.CreateProjectWithResponse(
		context.Background(),
		coolify.CreateProjectJSONRequestBody{
			Name:        &name,
			Description: &description,
		},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())
	require.NotNil(t, res.JSON201)
	assert.Equal(t, "project-1", *res.JSON201.Uuid)
}

func TestGetProjectByUuid(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/projects/project-1", r.URL.Path)
		assertBearerAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"uuid":"project-1","name":"Demo"}`))
		assert.NoError(t, err)
	})

	res, err := client.GetProjectByUuidWithResponse(
		context.Background(),
		"project-1",
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, "project-1", *res.JSON200.Uuid)
}

func TestPreviewNullableFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/applications/app-1/previews/3", r.URL.Path)
		assertBearerAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"domains":null,"docker_compose_domains":null,
			"domain_port_overrides":null,"pull_request_id":3
		}`))
		assert.NoError(t, err)
	})

	res, err := client.GetPreviewDeploymentByPullRequestIdWithResponse(
		context.Background(),
		"app-1",
		3,
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Nil(t, res.JSON200.Domains)
	assert.Nil(t, res.JSON200.DockerComposeDomains)
	assert.Nil(t, res.JSON200.DomainPortOverrides)
}

func TestUpdatePreviewForbidden(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/v1/applications/app-1/previews/3", r.URL.Path)
		assertBearerAuth(t, r)
		w.WriteHeader(http.StatusForbidden)
	})

	res, err := client.UpdatePreviewDomainsByPullRequestIdWithResponse(
		context.Background(),
		"app-1",
		3,
		coolify.UpdatePreviewDomainsByPullRequestIdJSONRequestBody{},
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, res.StatusCode())
}

func TestUnauthorized(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertBearerAuth(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`{"message":"Unauthenticated."}`))
		assert.NoError(t, err)
	})

	res, err := client.ListProjectsWithResponse(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Nil(t, res.JSON200)
	require.NotNil(t, res.JSON401)
	assert.Equal(t, "Unauthenticated.", *res.JSON401.Message)
}

func TestMalformedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"unexpected":true}`))
		assert.NoError(t, err)
	})

	_, err := client.ListProjectsWithResponse(context.Background())

	var typeError *json.UnmarshalTypeError

	require.ErrorAs(t, err, &typeError)
}

func TestCancelledContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("cancelled request reached server")
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.ListProjectsWithResponse(ctx)
	require.ErrorIs(t, err, context.Canceled)
}

func TestLogLinesUnion(t *testing.T) {
	t.Parallel()

	var numeric, all, unlimited coolify.LogTailLines

	require.NoError(t, numeric.FromLogTailLines0(25))
	require.NoError(t, all.FromLogTailLines1(coolify.LogTailLines1("all")))
	require.NoError(t, unlimited.FromLogTailLines0(-1))

	testCases := []struct {
		name  string
		lines coolify.LogTailLines
		want  string
	}{
		{name: "numeric", lines: numeric, want: "25"},
		{name: "all", lines: all, want: "all"},
		{name: "unlimited alias", lines: unlimited, want: "-1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t,
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(
						t,
						"/api/v1/applications/app-1/logs",
						r.URL.Path,
					)
					assert.Equal(t, tc.want, r.URL.Query().Get("lines"))
					assertBearerAuth(t, r)

					w.Header().Set("Content-Type", "application/json")
					_, err := w.Write([]byte(`{"logs":"ready"}`))
					assert.NoError(t, err)
				},
			)

			res, err := client.GetApplicationLogsByUuidWithResponse(
				context.Background(),
				"app-1",
				&coolify.GetApplicationLogsByUuidParams{Lines: &tc.lines},
			)
			require.NoError(t, err)
			require.NotNil(t, res.JSON200)
			assert.Equal(t, "ready", *res.JSON200.Logs)
		})
	}
}
