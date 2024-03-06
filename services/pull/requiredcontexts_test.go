package pull

import (
	"testing"

	model "code.gitea.io/gitea/models/git"
	"code.gitea.io/gitea/modules/structs"
)

func TestMergeRequiredContextsCommitStatus_AllRequiredContextsMatched(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("Build 1"),
		createSuccess("Build 2"),
	)

	requiredContexts := createRequiredContexts(
		"Build 1",
		"Build 2",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusSuccess, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_AllRequiredContextsMatchedWithPattern(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("Build 1"),
		createSuccess("Build 2"),
	)

	requiredContexts := createRequiredContexts(
		"Build*",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusSuccess, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_MissingContext(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("Build 1"),
		createSuccess("Build 2"),
	)

	requiredContexts := createRequiredContexts(
		"Build 1",
		"Build 2",
		"Build 3",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusPending, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_MissingContextWithPattern(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("ci/head"),
		createSuccess("ci/pr"),
	)

	requiredContexts := createRequiredContexts(
		"ci/*",
		"lint",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusPending, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_OneCheckFailed(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("Build 1"),
		createFailure("Build 2"),
	)

	requiredContexts := createRequiredContexts(
		"Build 1",
		"Build 2",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusFailure, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_OneCheckFailedWithWildcard(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("ci/head"),
		createFailure("ci/pr"),
	)

	requiredContexts := createRequiredContexts(
		"ci/*",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusFailure, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_OneCheckPendingWithWildcard(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("ci/head"),
		createPending("ci/pr"),
	)

	requiredContexts := createRequiredContexts(
		"ci/*",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusPending, mergedStatus)
}

func TestMergeRequiredContextsCommitStatus_OneCheckPendingWithCatchAllWildcard(t *testing.T) {
	givenStatuses := createStatuses(
		createSuccess("ci/head"),
		createPending("ci/pr"),
	)

	requiredContexts := createRequiredContexts(
		"*",
	)

	mergedStatus := MergeRequiredContextsCommitStatus(givenStatuses, requiredContexts)

	assertStatus(t, structs.CommitStatusPending, mergedStatus)
}

func assertStatus(t *testing.T, expectedStatus, actualStatus structs.CommitStatusState) {
	if actualStatus != expectedStatus {
		t.Errorf("Expected merged status <%s> but got <%s>", expectedStatus, actualStatus)
	}
}

func createStatuses(statuses ...*model.CommitStatus) []*model.CommitStatus {
	return statuses
}

func createRequiredContexts(contexts ...string) []string {
	return contexts
}

func createSuccess(context string) *model.CommitStatus {
	return &model.CommitStatus{
		Context: context,
		State:   structs.CommitStatusSuccess,
	}
}

func createPending(context string) *model.CommitStatus {
	return &model.CommitStatus{
		Context: context,
		State:   structs.CommitStatusPending,
	}
}

func createFailure(context string) *model.CommitStatus {
	return &model.CommitStatus{
		Context: context,
		State:   structs.CommitStatusFailure,
	}
}
