package preview

import (
	"context"
	"embed"
	"fmt"
	gh "github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"golang.org/x/oauth2"
	v1App "k8s.io/api/apps/v1"
	v1Core "k8s.io/api/core/v1"
	v1Net "k8s.io/api/networking/v1"
	v1Meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"net/http/httputil"
	"strings"
)

/*
 * Provides an interface to GitHub actions and tracks the actions currently in progress
 */

//go:embed template
var k8Files embed.FS

type BuildTracker struct {
	ghClient *gh.Client
	k8Client *rest.Config
	jobChan  chan types.BuildRequest
}

func NewGithubClient(token string) *gh.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)

	return gh.NewClient(tc)
}

func NewBuildTracker(ghClient *gh.Client, k8Client *rest.Config) *BuildTracker {
	return &BuildTracker{
		ghClient: ghClient,
		k8Client: k8Client,
		jobChan:  make(chan types.BuildRequest, 100),
	}
}

func (jt *BuildTracker) RunGithubWorkflowJob(request types.BuildRequest) error {
	inputs := map[string]interface{}{}

	if request.IsPreviewBuild {
		inputs = map[string]interface{}{
			"previewBuildId":  request.Id,
			"dockerRepo":      request.DockerRegistryName,
			"wordpressDomain": request.WordpressDomain,
		}
	}

	resp, err := jt.ghClient.Actions.CreateWorkflowDispatchEventByFileName(context.Background(), request.RepoOwner, request.RepoName, request.Workflow, gh.CreateWorkflowDispatchEventRequest{
		Ref:    request.RepoRef,
		Inputs: inputs,
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		bytes, dumpErr := httputil.DumpResponse(resp.Response, true)

		if dumpErr != nil {
			return dumpErr
		}

		return errors.New(string(bytes))
	}

	return nil
}

type TemplateContext struct {
	ImageName string
	BuildId   string
	SiteId    string
}

func (jt *BuildTracker) DeployPreviewBuild(templateContext TemplateContext) error {
	namespace := readK8NamespaceFromYaml("template/namespace.yaml", templateContext)
	ingress := readK8IngressFromYaml("template/ingress.yaml", templateContext)
	service := readK8ServiceFromYaml("template/service.yaml", templateContext)
	deployment := readK8DeploymentFromYaml("template/deployment.yaml", templateContext)

	coreClient, err := kubernetes.NewForConfig(jt.k8Client)

	if err != nil {
		return err
	}

	ns := fmt.Sprintf("site-preview-%s", templateContext.BuildId)

	_, err = coreClient.CoreV1().Namespaces().Create(context.Background(), namespace, v1Meta.CreateOptions{})

	if err != nil {
		return err
	}

	_, err = coreClient.NetworkingV1().Ingresses(ns).Create(context.Background(), ingress, v1Meta.CreateOptions{})

	if err != nil {
		return err
	}

	_, err = coreClient.CoreV1().Services(ns).Create(context.Background(), service, v1Meta.CreateOptions{})

	if err != nil {
		return err
	}

	_, err = coreClient.AppsV1().Deployments(ns).Create(context.Background(), deployment, v1Meta.CreateOptions{})

	if err != nil {
		return err
	}

	return nil
}

func readK8NamespaceFromYaml(file string, templateContext TemplateContext) *v1Core.Namespace {
	ns := &v1Core.Namespace{}
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(loadFile(file, templateContext), nil, ns)

	if err != nil {
		panic(err)
	}

	return ns
}

func readK8IngressFromYaml(file string, templateContext TemplateContext) *v1Net.Ingress {
	ing := &v1Net.Ingress{}
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(loadFile(file, templateContext), nil, ing)

	if err != nil {
		panic(err)
	}

	return ing
}

func readK8ServiceFromYaml(file string, templateContext TemplateContext) *v1Core.Service {
	service := &v1Core.Service{}
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(loadFile(file, templateContext), nil, service)

	if err != nil {
		panic(err)
	}

	return service
}

func readK8DeploymentFromYaml(file string, templateContext TemplateContext) *v1App.Deployment {
	dep := &v1App.Deployment{}
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(loadFile(file, templateContext), nil, dep)

	if err != nil {
		panic(err)
	}

	return dep
}

func loadFile(file string, templateContext TemplateContext) []byte {
	objData, err := k8Files.ReadFile(file)

	if err != nil {
		panic(err)
	}

	return []byte(replaceVariablesInString(string(objData), templateContext))
}

func replaceVariablesInString(data string, templateContext TemplateContext) string {
	data = strings.Replace(data, "$IMAGE", templateContext.ImageName, -1)
	data = strings.Replace(data, "$BUILD_ID", templateContext.BuildId, -1)
	data = strings.Replace(data, "$SITE_ID", templateContext.SiteId, -1)

	return data
}
