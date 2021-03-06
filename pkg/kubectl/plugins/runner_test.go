/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugins

import (
	"os"
	"testing"

	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"
)

func TestExecRunner(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectedMsg string
		expectedErr string
	}{
		{
			name:        "success",
			command:     "echo test ok",
			expectedMsg: "test ok\n",
		},
		{
			name:        "invalid",
			command:     "false",
			expectedErr: "exit status 1",
		},
		{
			name:        "env",
			command:     "echo $KUBECTL_PLUGINS_TEST",
			expectedMsg: "ok\n",
		},
	}

	os.Setenv("KUBECTL_PLUGINS_TEST", "ok")
	defer os.Unsetenv("KUBECTL_PLUGINS_TEST")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streams, _, outBuf, _ := genericclioptions.NewTestIOStreams()

			plugin := &Plugin{
				Description: Description{
					Name:      tt.name,
					ShortDesc: "Test Runner Plugin",
					Command:   tt.command,
				},
			}

			ctx := RunningContext{
				IOStreams:   streams,
				WorkingDir:  ".",
				EnvProvider: &EmptyEnvProvider{},
			}

			runner := &ExecPluginRunner{}
			err := runner.Run(plugin, ctx)

			if outBuf.String() != tt.expectedMsg {
				t.Errorf("%s: unexpected output: %q", tt.name, outBuf.String())
			}

			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("%s: unexpected err output: %v", tt.name, err)
			}
		})
	}
}
