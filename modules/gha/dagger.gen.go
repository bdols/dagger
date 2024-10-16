// Code generated by dagger. DO NOT EDIT.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/shykes/gha/internal/dagger"
	"github.com/shykes/gha/internal/querybuilder"
	"github.com/shykes/gha/internal/telemetry"
)

var dag = dagger.Connect()

func Tracer() trace.Tracer {
	return otel.Tracer("dagger.io/sdk.go")
}

// used for local MarshalJSON implementations
var marshalCtx = context.Background()

// called by main()
func setMarshalContext(ctx context.Context) {
	marshalCtx = ctx
	dagger.SetMarshalContext(ctx)
}

type DaggerObject = querybuilder.GraphQLMarshaller

type ExecError = dagger.ExecError

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

// convertSlice converts a slice of one type to a slice of another type using a
// converter function
func convertSlice[I any, O any](in []I, f func(I) O) []O {
	out := make([]O, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func (r Gha) MarshalJSON() ([]byte, error) {
	var concrete struct {
		Pipelines []*Pipeline
		Settings  Settings
	}
	concrete.Pipelines = r.Pipelines
	concrete.Settings = r.Settings
	return json.Marshal(&concrete)
}

func (r *Gha) UnmarshalJSON(bs []byte) error {
	var concrete struct {
		Pipelines []*Pipeline
		Settings  Settings
	}
	err := json.Unmarshal(bs, &concrete)
	if err != nil {
		return err
	}
	r.Pipelines = concrete.Pipelines
	r.Settings = concrete.Settings
	return nil
}

func (r Settings) MarshalJSON() ([]byte, error) {
	var concrete struct {
		PublicToken            string
		DaggerVersion          string
		NoTraces               bool
		StopEngine             bool
		AsJson                 bool
		Runner                 string
		PullRequestConcurrency string
		Debug                  bool
		FileExtension          string
		Repository             *dagger.Directory
		TimeoutMinutes         int
		Permissions            []Permission
	}
	concrete.PublicToken = r.PublicToken
	concrete.DaggerVersion = r.DaggerVersion
	concrete.NoTraces = r.NoTraces
	concrete.StopEngine = r.StopEngine
	concrete.AsJson = r.AsJson
	concrete.Runner = r.Runner
	concrete.PullRequestConcurrency = r.PullRequestConcurrency
	concrete.Debug = r.Debug
	concrete.FileExtension = r.FileExtension
	concrete.Repository = r.Repository
	concrete.TimeoutMinutes = r.TimeoutMinutes
	concrete.Permissions = r.Permissions
	return json.Marshal(&concrete)
}

func (r *Settings) UnmarshalJSON(bs []byte) error {
	var concrete struct {
		PublicToken            string
		DaggerVersion          string
		NoTraces               bool
		StopEngine             bool
		AsJson                 bool
		Runner                 string
		PullRequestConcurrency string
		Debug                  bool
		FileExtension          string
		Repository             *dagger.Directory
		TimeoutMinutes         int
		Permissions            []Permission
	}
	err := json.Unmarshal(bs, &concrete)
	if err != nil {
		return err
	}
	r.PublicToken = concrete.PublicToken
	r.DaggerVersion = concrete.DaggerVersion
	r.NoTraces = concrete.NoTraces
	r.StopEngine = concrete.StopEngine
	r.AsJson = concrete.AsJson
	r.Runner = concrete.Runner
	r.PullRequestConcurrency = concrete.PullRequestConcurrency
	r.Debug = concrete.Debug
	r.FileExtension = concrete.FileExtension
	r.Repository = concrete.Repository
	r.TimeoutMinutes = concrete.TimeoutMinutes
	r.Permissions = concrete.Permissions
	return nil
}

func main() {
	ctx := context.Background()

	// Direct slog to the new stderr. This is only for dev time debugging, and
	// runtime errors/warnings.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := dispatch(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func dispatch(ctx context.Context) error {
	ctx = telemetry.InitEmbedded(ctx, resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("dagger-go-sdk"),
		// TODO version?
	))
	defer telemetry.Close()

	// A lot of the "work" actually happens when we're marshalling the return
	// value, which entails getting object IDs, which happens in MarshalJSON,
	// which has no ctx argument, so we use this lovely global variable.
	setMarshalContext(ctx)

	fnCall := dag.CurrentFunctionCall()
	parentName, err := fnCall.ParentName(ctx)
	if err != nil {
		return fmt.Errorf("get parent name: %w", err)
	}
	fnName, err := fnCall.Name(ctx)
	if err != nil {
		return fmt.Errorf("get fn name: %w", err)
	}
	parentJson, err := fnCall.Parent(ctx)
	if err != nil {
		return fmt.Errorf("get fn parent: %w", err)
	}
	fnArgs, err := fnCall.InputArgs(ctx)
	if err != nil {
		return fmt.Errorf("get fn args: %w", err)
	}

	inputArgs := map[string][]byte{}
	for _, fnArg := range fnArgs {
		argName, err := fnArg.Name(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg name: %w", err)
		}
		argValue, err := fnArg.Value(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg value: %w", err)
		}
		inputArgs[argName] = []byte(argValue)
	}

	result, err := invoke(ctx, []byte(parentJson), parentName, fnName, inputArgs)
	if err != nil {
		return fmt.Errorf("invoke: %w", err)
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err = fnCall.ReturnValue(ctx, dagger.JSON(resultBytes)); err != nil {
		return fmt.Errorf("store return value: %w", err)
	}
	return nil
}
func invoke(ctx context.Context, parentJSON []byte, parentName string, fnName string, inputArgs map[string][]byte) (_ any, err error) {
	_ = inputArgs
	switch parentName {
	case "Gha":
		switch fnName {
		case "Validate":
			var parent Gha
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			var repo *dagger.Directory
			if inputArgs["repo"] != nil {
				err = json.Unmarshal([]byte(inputArgs["repo"]), &repo)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg repo", err))
				}
			}
			return (*Gha).Validate(&parent, ctx, repo)
		case "Config":
			var parent Gha
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return (*Gha).Config(&parent, ctx), nil
		case "WithPipeline":
			var parent Gha
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			var name string
			if inputArgs["name"] != nil {
				err = json.Unmarshal([]byte(inputArgs["name"]), &name)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg name", err))
				}
			}
			var command string
			if inputArgs["command"] != nil {
				err = json.Unmarshal([]byte(inputArgs["command"]), &command)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg command", err))
				}
			}
			var module string
			if inputArgs["module"] != nil {
				err = json.Unmarshal([]byte(inputArgs["module"]), &module)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg module", err))
				}
			}
			var runner string
			if inputArgs["runner"] != nil {
				err = json.Unmarshal([]byte(inputArgs["runner"]), &runner)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg runner", err))
				}
			}
			var secrets []string
			if inputArgs["secrets"] != nil {
				err = json.Unmarshal([]byte(inputArgs["secrets"]), &secrets)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg secrets", err))
				}
			}
			var sparseCheckout []string
			if inputArgs["sparseCheckout"] != nil {
				err = json.Unmarshal([]byte(inputArgs["sparseCheckout"]), &sparseCheckout)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg sparseCheckout", err))
				}
			}
			var dispatch bool
			if inputArgs["dispatch"] != nil {
				err = json.Unmarshal([]byte(inputArgs["dispatch"]), &dispatch)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg dispatch", err))
				}
			}
			var noDispatch bool
			if inputArgs["noDispatch"] != nil {
				err = json.Unmarshal([]byte(inputArgs["noDispatch"]), &noDispatch)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg noDispatch", err))
				}
			}
			var lfs bool
			if inputArgs["lfs"] != nil {
				err = json.Unmarshal([]byte(inputArgs["lfs"]), &lfs)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg lfs", err))
				}
			}
			var debug bool
			if inputArgs["debug"] != nil {
				err = json.Unmarshal([]byte(inputArgs["debug"]), &debug)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg debug", err))
				}
			}
			var daggerVersion string
			if inputArgs["daggerVersion"] != nil {
				err = json.Unmarshal([]byte(inputArgs["daggerVersion"]), &daggerVersion)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg daggerVersion", err))
				}
			}
			var timeoutMinutes int
			if inputArgs["timeoutMinutes"] != nil {
				err = json.Unmarshal([]byte(inputArgs["timeoutMinutes"]), &timeoutMinutes)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg timeoutMinutes", err))
				}
			}
			var permissions Permissions
			if inputArgs["permissions"] != nil {
				err = json.Unmarshal([]byte(inputArgs["permissions"]), &permissions)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg permissions", err))
				}
			}
			var onIssueComment bool
			if inputArgs["onIssueComment"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onIssueComment"]), &onIssueComment)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onIssueComment", err))
				}
			}
			var onIssueCommentCreated bool
			if inputArgs["onIssueCommentCreated"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onIssueCommentCreated"]), &onIssueCommentCreated)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onIssueCommentCreated", err))
				}
			}
			var onIssueCommentEdited bool
			if inputArgs["onIssueCommentEdited"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onIssueCommentEdited"]), &onIssueCommentEdited)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onIssueCommentEdited", err))
				}
			}
			var onIssueCommentDeleted bool
			if inputArgs["onIssueCommentDeleted"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onIssueCommentDeleted"]), &onIssueCommentDeleted)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onIssueCommentDeleted", err))
				}
			}
			var onPullRequest bool
			if inputArgs["onPullRequest"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequest"]), &onPullRequest)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequest", err))
				}
			}
			var pullRequestConcurrency string
			if inputArgs["pullRequestConcurrency"] != nil {
				err = json.Unmarshal([]byte(inputArgs["pullRequestConcurrency"]), &pullRequestConcurrency)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg pullRequestConcurrency", err))
				}
			}
			var onPullRequestBranches []string
			if inputArgs["onPullRequestBranches"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestBranches"]), &onPullRequestBranches)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestBranches", err))
				}
			}
			var onPullRequestPaths []string
			if inputArgs["onPullRequestPaths"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestPaths"]), &onPullRequestPaths)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestPaths", err))
				}
			}
			var onPullRequestAssigned bool
			if inputArgs["onPullRequestAssigned"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestAssigned"]), &onPullRequestAssigned)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestAssigned", err))
				}
			}
			var onPullRequestUnassigned bool
			if inputArgs["onPullRequestUnassigned"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestUnassigned"]), &onPullRequestUnassigned)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestUnassigned", err))
				}
			}
			var onPullRequestLabeled bool
			if inputArgs["onPullRequestLabeled"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestLabeled"]), &onPullRequestLabeled)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestLabeled", err))
				}
			}
			var onPullRequestUnlabeled bool
			if inputArgs["onPullRequestUnlabeled"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestUnlabeled"]), &onPullRequestUnlabeled)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestUnlabeled", err))
				}
			}
			var onPullRequestOpened bool
			if inputArgs["onPullRequestOpened"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestOpened"]), &onPullRequestOpened)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestOpened", err))
				}
			}
			var onPullRequestEdited bool
			if inputArgs["onPullRequestEdited"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestEdited"]), &onPullRequestEdited)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestEdited", err))
				}
			}
			var onPullRequestClosed bool
			if inputArgs["onPullRequestClosed"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestClosed"]), &onPullRequestClosed)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestClosed", err))
				}
			}
			var onPullRequestReopened bool
			if inputArgs["onPullRequestReopened"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestReopened"]), &onPullRequestReopened)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestReopened", err))
				}
			}
			var onPullRequestSynchronize bool
			if inputArgs["onPullRequestSynchronize"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestSynchronize"]), &onPullRequestSynchronize)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestSynchronize", err))
				}
			}
			var onPullRequestConvertedToDraft bool
			if inputArgs["onPullRequestConverted_to_draft"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestConverted_to_draft"]), &onPullRequestConvertedToDraft)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestConverted_to_draft", err))
				}
			}
			var onPullRequestLocked bool
			if inputArgs["onPullRequestLocked"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestLocked"]), &onPullRequestLocked)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestLocked", err))
				}
			}
			var onPullRequestUnlocked bool
			if inputArgs["onPullRequestUnlocked"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestUnlocked"]), &onPullRequestUnlocked)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestUnlocked", err))
				}
			}
			var onPullRequestEnqueued bool
			if inputArgs["onPullRequestEnqueued"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestEnqueued"]), &onPullRequestEnqueued)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestEnqueued", err))
				}
			}
			var onPullRequestDequeued bool
			if inputArgs["onPullRequestDequeued"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestDequeued"]), &onPullRequestDequeued)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestDequeued", err))
				}
			}
			var onPullRequestMilestoned bool
			if inputArgs["onPullRequestMilestoned"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestMilestoned"]), &onPullRequestMilestoned)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestMilestoned", err))
				}
			}
			var onPullRequestDemilestoned bool
			if inputArgs["onPullRequestDemilestoned"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestDemilestoned"]), &onPullRequestDemilestoned)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestDemilestoned", err))
				}
			}
			var onPullRequestReadyForReview bool
			if inputArgs["onPullRequestReadyForReview"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestReadyForReview"]), &onPullRequestReadyForReview)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestReadyForReview", err))
				}
			}
			var onPullRequestReviewRequested bool
			if inputArgs["onPullRequestReviewRequested"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestReviewRequested"]), &onPullRequestReviewRequested)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestReviewRequested", err))
				}
			}
			var onPullRequestReviewRequestRemoved bool
			if inputArgs["onPullRequestReviewRequestRemoved"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestReviewRequestRemoved"]), &onPullRequestReviewRequestRemoved)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestReviewRequestRemoved", err))
				}
			}
			var onPullRequestAutoMergeEnabled bool
			if inputArgs["onPullRequestAutoMergeEnabled"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestAutoMergeEnabled"]), &onPullRequestAutoMergeEnabled)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestAutoMergeEnabled", err))
				}
			}
			var onPullRequestAutoMergeDisabled bool
			if inputArgs["onPullRequestAutoMergeDisabled"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPullRequestAutoMergeDisabled"]), &onPullRequestAutoMergeDisabled)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPullRequestAutoMergeDisabled", err))
				}
			}
			var onPush bool
			if inputArgs["onPush"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPush"]), &onPush)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPush", err))
				}
			}
			var onPushTags []string
			if inputArgs["onPushTags"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPushTags"]), &onPushTags)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPushTags", err))
				}
			}
			var onPushBranches []string
			if inputArgs["onPushBranches"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onPushBranches"]), &onPushBranches)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onPushBranches", err))
				}
			}
			var onSchedule []string
			if inputArgs["onSchedule"] != nil {
				err = json.Unmarshal([]byte(inputArgs["onSchedule"]), &onSchedule)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg onSchedule", err))
				}
			}
			return (*Gha).WithPipeline(&parent, name, command, module, runner, secrets, sparseCheckout, dispatch, noDispatch, lfs, debug, daggerVersion, timeoutMinutes, permissions, onIssueComment, onIssueCommentCreated, onIssueCommentEdited, onIssueCommentDeleted, onPullRequest, pullRequestConcurrency, onPullRequestBranches, onPullRequestPaths, onPullRequestAssigned, onPullRequestUnassigned, onPullRequestLabeled, onPullRequestUnlabeled, onPullRequestOpened, onPullRequestEdited, onPullRequestClosed, onPullRequestReopened, onPullRequestSynchronize, onPullRequestConvertedToDraft, onPullRequestLocked, onPullRequestUnlocked, onPullRequestEnqueued, onPullRequestDequeued, onPullRequestMilestoned, onPullRequestDemilestoned, onPullRequestReadyForReview, onPullRequestReviewRequested, onPullRequestReviewRequestRemoved, onPullRequestAutoMergeEnabled, onPullRequestAutoMergeDisabled, onPush, onPushTags, onPushBranches, onSchedule), nil
		case "":
			var parent Gha
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			var noTraces bool
			if inputArgs["noTraces"] != nil {
				err = json.Unmarshal([]byte(inputArgs["noTraces"]), &noTraces)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg noTraces", err))
				}
			}
			var publicToken string
			if inputArgs["publicToken"] != nil {
				err = json.Unmarshal([]byte(inputArgs["publicToken"]), &publicToken)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg publicToken", err))
				}
			}
			var daggerVersion string
			if inputArgs["daggerVersion"] != nil {
				err = json.Unmarshal([]byte(inputArgs["daggerVersion"]), &daggerVersion)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg daggerVersion", err))
				}
			}
			var stopEngine bool
			if inputArgs["stopEngine"] != nil {
				err = json.Unmarshal([]byte(inputArgs["stopEngine"]), &stopEngine)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg stopEngine", err))
				}
			}
			var asJson bool
			if inputArgs["asJson"] != nil {
				err = json.Unmarshal([]byte(inputArgs["asJson"]), &asJson)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg asJson", err))
				}
			}
			var runner string
			if inputArgs["runner"] != nil {
				err = json.Unmarshal([]byte(inputArgs["runner"]), &runner)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg runner", err))
				}
			}
			var fileExtension string
			if inputArgs["fileExtension"] != nil {
				err = json.Unmarshal([]byte(inputArgs["fileExtension"]), &fileExtension)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg fileExtension", err))
				}
			}
			var repository *dagger.Directory
			if inputArgs["repository"] != nil {
				err = json.Unmarshal([]byte(inputArgs["repository"]), &repository)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg repository", err))
				}
			}
			var timeoutMinutes int
			if inputArgs["timeoutMinutes"] != nil {
				err = json.Unmarshal([]byte(inputArgs["timeoutMinutes"]), &timeoutMinutes)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg timeoutMinutes", err))
				}
			}
			return New(noTraces, publicToken, daggerVersion, stopEngine, asJson, runner, fileExtension, repository, timeoutMinutes), nil
		default:
			return nil, fmt.Errorf("unknown function %s", fnName)
		}
	case "":
		return dag.Module().
			WithDescription("Manage Github Actions configurations with Dagger\n\nDaggerizing your CI makes your YAML configurations smaller, but they still exist,\nand they're still a pain to maintain by hand.\n\nThis module aims to finish the job, by letting you generate your remaining\nYAML configuration from a Dagger pipeline, written in your favorite language.\n").
			WithObject(
				dag.TypeDef().WithObject("Gha").
					WithFunction(
						dag.Function("Validate",
							dag.TypeDef().WithObject("Gha")).
							WithDescription("Validate a Github Actions configuration (best effort)").
							WithArg("repo", dag.TypeDef().WithObject("Directory"))).
					WithFunction(
						dag.Function("Config",
							dag.TypeDef().WithObject("Directory")).
							WithDescription("Export the configuration to a .github directory")).
					WithFunction(
						dag.Function("WithPipeline",
							dag.TypeDef().WithObject("Gha")).
							WithDescription("Add a pipeline").
							WithArg("name", dag.TypeDef().WithKind(dagger.StringKind), dagger.FunctionWithArgOpts{Description: "Pipeline name"}).
							WithArg("command", dag.TypeDef().WithKind(dagger.StringKind), dagger.FunctionWithArgOpts{Description: "The Dagger command to execute\nExample 'build --source=.'"}).
							WithArg("module", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "The Dagger module to load"}).
							WithArg("runner", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Dispatch jobs to the given runner"}).
							WithArg("secrets", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Github secrets to inject into the pipeline environment.\nFor each secret, an env variable with the same name is created.\nExample: [\"PROD_DEPLOY_TOKEN\", \"PRIVATE_SSH_KEY\"]"}).
							WithArg("sparseCheckout", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Use a sparse git checkout, only including the given paths\nExample: [\"src\", \"tests\", \"Dockerfile\"]"}).
							WithArg("dispatch", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "(DEPRECATED) allow this pipeline to be manually \"dispatched\""}).
							WithArg("noDispatch", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Disable manual \"dispatch\" of this pipeline"}).
							WithArg("lfs", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Enable lfs on git checkout"}).
							WithArg("debug", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline in debug mode"}).
							WithArg("daggerVersion", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Dagger version to run this pipeline"}).
							WithArg("timeoutMinutes", dag.TypeDef().WithKind(dagger.IntegerKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "The maximum number of minutes to run the pipeline before killing the process"}).
							WithArg("permissions", dag.TypeDef().WithListOf(dag.TypeDef().WithEnum("Permission")).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Permissions to grant the pipeline"}).
							WithArg("onIssueComment", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline on any issue comment activity"}).
							WithArg("onIssueCommentCreated", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onIssueCommentEdited", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onIssueCommentDeleted", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequest", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline on any pull request activity"}).
							WithArg("pullRequestConcurrency", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Configure this pipeline's concurrency for each PR.\nThis is triggered when the pipeline is scheduled concurrently on the same PR.\n  - allow: all instances are allowed to run concurrently\n  - queue: new instances are queued, and run sequentially\n  - preempt: new instances run immediately, older ones are canceled\nPossible values: \"allow\", \"preempt\", \"queue\"", DefaultValue: dagger.JSON("\"allow\"")}).
							WithArg("onPullRequestBranches", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true)).
							WithArg("onPullRequestPaths", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true)).
							WithArg("onPullRequestAssigned", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestUnassigned", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestLabeled", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestUnlabeled", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestOpened", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestEdited", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestClosed", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestReopened", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestSynchronize", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestConverted_to_draft", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestLocked", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestUnlocked", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestEnqueued", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestDequeued", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestMilestoned", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestDemilestoned", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestReadyForReview", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestReviewRequested", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestReviewRequestRemoved", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestAutoMergeEnabled", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPullRequestAutoMergeDisabled", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true)).
							WithArg("onPush", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline on any git push"}).
							WithArg("onPushTags", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline on git push to the specified tags"}).
							WithArg("onPushBranches", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline on git push to the specified branches"}).
							WithArg("onSchedule", dag.TypeDef().WithListOf(dag.TypeDef().WithKind(dagger.StringKind)).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Run the pipeline at a schedule time"})).
					WithField("Settings", dag.TypeDef().WithObject("Settings"), dagger.TypeDefWithFieldOpts{Description: "Settings for this Github Actions project"}).
					WithConstructor(
						dag.Function("New",
							dag.TypeDef().WithObject("Gha")).
							WithArg("noTraces", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Disable sending traces to Dagger Cloud"}).
							WithArg("publicToken", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Public Dagger Cloud token, for open-source projects. DO NOT PASS YOUR PRIVATE DAGGER CLOUD TOKEN!\nThis is for a special \"public\" token which can safely be shared publicly.\nTo get one, contact support@dagger.io"}).
							WithArg("daggerVersion", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Dagger version to run in the Github Actions pipelines", DefaultValue: dagger.JSON("\"latest\"")}).
							WithArg("stopEngine", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Explicitly stop the Dagger Engine after completing the pipeline"}).
							WithArg("asJson", dag.TypeDef().WithKind(dagger.BooleanKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Encode all files as JSON (which is also valid YAML)"}).
							WithArg("runner", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Configure a default runner for all workflows\nSee https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/using-self-hosted-runners-in-a-workflow", DefaultValue: dagger.JSON("\"ubuntu-latest\"")}).
							WithArg("fileExtension", dag.TypeDef().WithKind(dagger.StringKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "File extension to use for generated workflow files", DefaultValue: dagger.JSON("\".gen.yml\"")}).
							WithArg("repository", dag.TypeDef().WithObject("Directory").WithOptional(true), dagger.FunctionWithArgOpts{Description: "Existing repository root, to merge existing content", Ignore: []string{"!.github"}}).
							WithArg("timeoutMinutes", dag.TypeDef().WithKind(dagger.IntegerKind).WithOptional(true), dagger.FunctionWithArgOpts{Description: "Default timeout for CI jobs, in minutes"}))).
			WithEnum(
				dag.TypeDef().WithEnum("Permission").
					WithEnumValue("read_contents").
					WithEnumValue("read_issues").
					WithEnumValue("read_actions").
					WithEnumValue("read_packages").
					WithEnumValue("read_deployments").
					WithEnumValue("read_pull_requests").
					WithEnumValue("read_pages").
					WithEnumValue("read_id_token").
					WithEnumValue("read_repository_projects").
					WithEnumValue("read_statuses").
					WithEnumValue("read_metadata").
					WithEnumValue("read_checks").
					WithEnumValue("read_discussions").
					WithEnumValue("write_contents").
					WithEnumValue("write_issues").
					WithEnumValue("write_actions").
					WithEnumValue("write_packages").
					WithEnumValue("write_deployments").
					WithEnumValue("write_pull_requests").
					WithEnumValue("write_pages").
					WithEnumValue("write_id_token").
					WithEnumValue("write_repository_projects").
					WithEnumValue("write_statuses").
					WithEnumValue("write_metadata").
					WithEnumValue("write_checks").
					WithEnumValue("write_discussions")).
			WithObject(
				dag.TypeDef().WithObject("Settings").
					WithField("PublicToken", dag.TypeDef().WithKind(dagger.StringKind)).
					WithField("DaggerVersion", dag.TypeDef().WithKind(dagger.StringKind)).
					WithField("NoTraces", dag.TypeDef().WithKind(dagger.BooleanKind)).
					WithField("StopEngine", dag.TypeDef().WithKind(dagger.BooleanKind)).
					WithField("AsJson", dag.TypeDef().WithKind(dagger.BooleanKind)).
					WithField("Runner", dag.TypeDef().WithKind(dagger.StringKind)).
					WithField("PullRequestConcurrency", dag.TypeDef().WithKind(dagger.StringKind)).
					WithField("Debug", dag.TypeDef().WithKind(dagger.BooleanKind)).
					WithField("FileExtension", dag.TypeDef().WithKind(dagger.StringKind)).
					WithField("Repository", dag.TypeDef().WithObject("Directory")).
					WithField("TimeoutMinutes", dag.TypeDef().WithKind(dagger.IntegerKind)).
					WithField("Permissions", dag.TypeDef().WithListOf(dag.TypeDef().WithEnum("Permission")))), nil
	default:
		return nil, fmt.Errorf("unknown object %s", parentName)
	}
}
