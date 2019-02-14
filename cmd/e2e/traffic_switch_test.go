package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrafficSwitch(t *testing.T) {
	t.Parallel()

	stacksetName := "stackset-switch-traffic"
	firstVersion := "v1"
	firstStack := fmt.Sprintf("%s-%s", stacksetName, firstVersion)
	updatedVersion := "v2"
	updatedStack := fmt.Sprintf("%s-%s", stacksetName, updatedVersion)

	err := createStackSet(stacksetName, newStacksetSpec(stacksetName, firstVersion, false, true))
	require.NoError(t, err)
	_, err = waitForStack(t, stacksetName, firstVersion)
	require.NoError(t, err)
	err = updateStackset(stacksetName, newStacksetSpec(stacksetName, updatedVersion, false, true))
	require.NoError(t, err)

	_, err = waitForIngress(t, stacksetName)
	require.NoError(t, err)

	initialWeights := map[string]float64{firstStack: 100}
	err = trafficWeightsUpdated(t, stacksetName, weightKindActual, initialWeights).await()
	require.NoError(t, err)

	// Switch traffic 50/50
	desiredWeights := map[string]float64{firstStack: 50, updatedStack: 50}
	err = setDesiredTrafficWeights(stacksetName, desiredWeights)
	require.NoError(t, err)
	err = trafficWeightsUpdated(t, stacksetName, weightKindActual, desiredWeights).await()
	require.NoError(t, err)

	// Switch traffic 0/100
	newDesiredWeights := map[string]float64{updatedStack: 100}
	err = setDesiredTrafficWeights(stacksetName, newDesiredWeights)
	require.NoError(t, err)
	err = trafficWeightsUpdated(t, stacksetName, weightKindActual, newDesiredWeights).await()
	require.NoError(t, err)
}