package execution

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"strconv"
	"strings"
)

type KubernetesCommandExecutor struct {
	Site     *types.Site
	Pod      v1.Pod
	K8Config *rest.Config
}

func NewKubernetesCommandExecutor(site *types.Site) (*KubernetesCommandExecutor, error) {
	pod, err := getPodBySite(site, config.Config.K8.LabelSelector, config.Config.K8RestConfig)

	if err != nil {
		return nil, err
	}

	return &KubernetesCommandExecutor{
		Site:     site,
		Pod:      *pod,
		K8Config: config.Config.K8RestConfig,
	}, nil
}

func (e *KubernetesCommandExecutor) ExecuteCommand(args []string) (*types.CommandResult, error) {
	command := strings.Join(args, " ")

	fields := log.Fields{
		"event":    "execute-command",
		"command":  command,
		"executor": "kubernetes",
		"pod":      e.Pod.Name,
		"site":     e.Site.Key,
	}

	stdout, stderr, err := executeRemoteCommand(e.K8Config, e.Pod, command, nil)

	if err != nil {
		return nil, err
	}

	if stderr != "" {
		return nil, errors.New(stderr)
	} else {
		log.Debugln(stdout)
		log.WithFields(fields).Infoln("success")
	}

	return &types.CommandResult{Command: command, Output: stdout}, nil
}

type ExecuteSiteCommandOpts struct {
	Stdin *io.Reader
}

// ExecuteRemoteCommand executes a remote shell command on the given pod
// returns the output from stdout and stderr
func executeRemoteCommand(restCfg *rest.Config, pod v1.Pod, command string, stdin *io.Reader) (string, string, error) {
	coreClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return "", "", err
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := coreClient.
		CoreV1().
		RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: []string{"/bin/sh", "-c", command},
			Stdin:   stdin != nil,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())

	if err != nil {
		if buf.String() != "" {
			log.Infoln(buf)
		}

		if errBuf.String() != "" {
			log.Infoln(errBuf)
		}

		log.Errorln(err)

		return "", "", errors.Wrapf(err, "failed executing command (spdy) %s on %v/%v", command, pod.Namespace, pod.Name)
	}

	streamOpts := remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: buf,
		Stderr: errBuf,
		Tty:    true,
	}

	if stdin != nil {
		log.Debugln("Stdin active for current command")

		streamOpts.Stdin = *stdin
	}

	err = exec.Stream(streamOpts)

	if err != nil {
		if buf.String() != "" {
			log.Infoln(buf)
		}

		if errBuf.String() != "" {
			log.Infoln(errBuf)
		}

		log.Errorln(err)

		return "", "", errors.Wrapf(err, "failed executing command (stream) %s on %v/%v", command, pod.Namespace, pod.Name)
	}

	return buf.String(), errBuf.String(), nil
}

func getPodBySite(site *types.Site, baseSelector string, k8Config *rest.Config) (*v1.Pod, error) {
	client, err := kubernetes.NewForConfig(k8Config)
	if err != nil {
		panic(err.Error())
	}

	// combine label selectors from wordpress and the specific one
	combinedSelector := fmt.Sprintf("%s,%s", baseSelector, site.LabelSelector)

	log.Info("selecting ", combinedSelector, " for site ", site.Key)

	pods, err := client.CoreV1().Pods(site.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: combinedSelector,
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error occurred finding the pod for site %s", site.Key))
	}

	if len(pods.Items) != 1 {
		return nil, errors.New("could not find pod with the label selector " + site.LabelSelector + ", items returned: " + strconv.Itoa(len(pods.Items)))
	}

	return &pods.Items[0], nil
}
