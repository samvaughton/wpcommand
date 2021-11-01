package execution

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutionFactory(t *testing.T) {
	t.Run("TestDebugExecutorIsReturned", func(t *testing.T) {
		testSite := &types.Site{TestMode: true}

		exec, err := NewCommandExecutor(testSite, &types.Config{})

		assert.Nil(t, err)
		assert.IsType(t, &DebugCommandExecutor{}, exec)
	})

	t.Run("TestK8ExecutorIsReturned", func(t *testing.T) {
		testSite := &types.Site{TestMode: false, Namespace: "site-test", LabelSelector: "example.test/instance-name=test"}

		exec, err := NewCommandExecutor(testSite, &types.Config{})

		assert.Nil(t, err)
		assert.IsType(t, &KubernetesCommandExecutor{}, exec)
	})

	t.Run("TestErrorIsReturned", func(t *testing.T) {
		testSite := &types.Site{TestMode: false, Namespace: "", LabelSelector: "example.test/instance-name=test"}

		exec, err := NewCommandExecutor(testSite, &types.Config{})

		assert.Error(t, err)
		assert.Nil(t, exec)
	})
}
