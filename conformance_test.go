/*
Copyright 2022 The Dapr Authors
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

package main

import (
	"os"
	"testing"

	confBindings "github.com/dapr/components-contrib/tests/conformance/bindings"
	confPubsub "github.com/dapr/components-contrib/tests/conformance/pubsub"
	confState "github.com/dapr/components-contrib/tests/conformance/state"
	"github.com/dapr/components-contrib/tests/conformance/utils"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/dapr/dapr/pkg/components/bindings"
	"github.com/dapr/dapr/pkg/components/pubsub"
	"github.com/dapr/dapr/pkg/components/state"

	"github.com/dapr/kit/logger"

	"github.com/stretchr/testify/require"
)

var l = logger.NewLogger("dapr-conformance-tests")

func TestPluggableConformance(t *testing.T) {
	socket := "/tmp/socket.sock"
	if soc, ok := os.LookupEnv("INPUT_SOCKET"); ok {
		socket = soc
	}

	var componentMetadata map[string]string
	if metadata, ok := os.LookupEnv("INPUT_METADATA"); ok {
		require.NoError(t, yaml.Unmarshal([]byte(metadata), &componentMetadata))
	}

	operations, allOperations := []string{}, true
	if operationsList, ok := os.LookupEnv("INPUT_OPERATIONS"); ok {
		require.NoError(t, yaml.Unmarshal([]byte(operationsList), &operations))
		allOperations = len(operations) == 0
	}

	operationMap := make(map[string]struct{})

	for _, op := range operations {
		operationMap[op] = struct{}{}
	}

	common := utils.CommonConfig{
		AllOperations: allOperations,
		Operations:    operationMap,
	}

	additionalConfig := make(map[string]any)
	if config, ok := os.LookupEnv("INPUT_CONFIG"); ok && len(config) != 0 {
		require.NoError(t, yaml.Unmarshal([]byte(config), &additionalConfig))
	}

	t.Run("pubsub", func(t *testing.T) {
		pubsub := pubsub.NewGRPCPubSub(l, socket)
		testConf := confPubsub.TestConfig{
			CommonConfig: common,
		}

		require.NoError(t, mapstructure.Decode(additionalConfig, &testConf))

		confPubsub.ConformanceTests(t, componentMetadata, pubsub, testConf)
	})

	t.Run("state", func(t *testing.T) {
		stateStore := state.NewGRPCStateStore(l, socket)

		testConf := confState.TestConfig{
			CommonConfig: common,
		}

		require.NoError(t, mapstructure.Decode(additionalConfig, &testConf))

		confState.ConformanceTests(t, componentMetadata, stateStore, testConf)
	})

	t.Run("bindings", func(t *testing.T) {
		inputBinding := bindings.NewGRPCInputBinding(l, socket)
		outputBinding := bindings.NewGRPCOutputBinding(l, socket)

		testConf := confBindings.TestConfig{
			CommonConfig: common,
		}

		require.NoError(t, mapstructure.Decode(additionalConfig, &testConf))
		confBindings.ConformanceTests(t, componentMetadata, inputBinding, outputBinding, testConf)
	})
}
