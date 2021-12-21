package execution

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutionFactory(t *testing.T) {
	t.Run("TestDebugExecutorIsReturned", func(t *testing.T) {
		testSite := &types.Site{types.ApiSiteCore{TestMode: true}, types.ApiSiteCredentials{}}

		exec, err := NewCommandExecutor(testSite, &types.Config{})

		assert.Nil(t, err)
		assert.IsType(t, &DebugCommandExecutor{}, exec)
	})

	t.Run("TestK8ExecutorIsReturned", func(t *testing.T) {
		testSite := &types.Site{types.ApiSiteCore{TestMode: false, Namespace: "site-test", LabelSelector: "example.test/instance-name=test"}, types.ApiSiteCredentials{}}

		exec, err := NewCommandExecutor(testSite, &types.Config{})

		assert.Nil(t, err)
		assert.IsType(t, &KubernetesCommandExecutor{}, exec)
	})

	t.Run("TestErrorIsReturned", func(t *testing.T) {
		testSite := &types.Site{types.ApiSiteCore{TestMode: false, Namespace: "", LabelSelector: "example.test/instance-name=test"}, types.ApiSiteCredentials{}}

		exec, err := NewCommandExecutor(testSite, &types.Config{})

		assert.Error(t, err)
		assert.Nil(t, exec)
	})
}
