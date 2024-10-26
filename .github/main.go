package main

import (
	"context"
	"fmt"

	"github.com/dagger/dagger/.github/internal/dagger"
)

const (
	daggerVersion                       = "v0.13.6"
	daggerCloudToken                    = "dag_dagger_sBIv6DsjNerWvTqt2bSFeigBUqWxp9bhh3ONSSgeFnw"
	publicRunner                        = "ubuntu-latest"
	repositoryWithAccessToCustomRunners = "dagger/dagger"
)

type CI struct {
	// +private
	Gha *dagger.Gha
}

func New(
	// The dagger repository
	// +optional
	// +defaultPath="/"
	// +ignore=["!.github"]
	repository *dagger.Directory,
) *CI {
	ci := new(CI)

	ci.Gha = dag.Gha(dagger.GhaOpts{
		DaggerVersion: daggerVersion,
		PublicToken:   daggerCloudToken,
		Repository:    repository,
	})

	return ci
}

// Configure Workflows
func (ci *CI) Workflows() *CI {
	return ci.
		onDaggerRunners().
		onDepotRunners().
		onNamespaceRunners()
}

func (ci *CI) onDaggerRunners() *CI {
	workflowOpts := dagger.GhaWithPipelineOpts{
		OnPushBranches:              []string{"main"},
		OnPullRequestOpened:         true,
		OnPullRequestReopened:       true,
		OnPullRequestSynchronize:    true,
		OnPullRequestReadyForReview: true,
		PullRequestConcurrency:      "preempt",
		Permissions:                 []dagger.GhaPermission{dagger.ReadContents},
	}
	return ci.withRunner(NewDaggerRunner(daggerVersion), workflowOpts)
}

func (ci *CI) onDepotRunners() *CI {
	workflowOpts := dagger.GhaWithPipelineOpts{
		OnSchedule:             []string{"6 0 * * *"},
		PullRequestConcurrency: "preempt",
		Permissions:            []dagger.GhaPermission{dagger.ReadContents},
		RunIf:                  fmt.Sprintf("${{ github.repository == '%s' }}", repositoryWithAccessToCustomRunners),
	}
	return ci.withRunner(NewDepotRunner(daggerVersion), workflowOpts)
}

func (ci *CI) onNamespaceRunners() *CI {
	workflowOpts := dagger.GhaWithPipelineOpts{
		OnSchedule:             []string{"6 0 * * *"},
		PullRequestConcurrency: "preempt",
		Permissions:            []dagger.GhaPermission{dagger.ReadContents},
		RunIf:                  fmt.Sprintf("${{ github.repository == '%s' }}", repositoryWithAccessToCustomRunners),
	}
	return ci.withRunner(NewNamespaceRunner(daggerVersion), workflowOpts)
}

func (ci *CI) withRunner(runner Runner, opts dagger.GhaWithPipelineOpts) *CI {
	return ci.
		workflow(runner.Pipeline("docs"), "docs lint", runner.Medium().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-python"), "check --targets=sdk/python", runner.Medium().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-python-dev"), "check --targets=sdk/python", runner.Large().DaggerInDocker().RunsOn(), true, 10, opts).
		workflow(runner.Pipeline("sdk-typescript"), "check --targets=sdk/typescript", runner.Medium().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-typescript-dev"), "check --targets=sdk/typescript", runner.Large().DaggerInDocker().RunsOn(), true, 10, opts).
		workflow(runner.Pipeline("sdk-go"), "check --targets=sdk/go", runner.Medium().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-go-dev"), "check --targets=sdk/go", runner.Large().DaggerInDocker().RunsOn(), true, 10, opts).
		workflow(runner.Pipeline("sdk-java"), "check --targets=sdk/java", runner.Medium().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-java-dev"), "check --targets=sdk/java", runner.Large().DaggerInDocker().RunsOn(), true, 10, opts).
		workflow(runner.Pipeline("sdk-elixir"), "check --targets=sdk/elixir", runner.XLarge().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-elixir-dev"), "check --targets=sdk/elixir", runner.XLarge().DaggerInDocker().RunsOn(), true, 10, opts).
		workflow(runner.Pipeline("sdk-rust"), "check --targets=sdk/rust", runner.XLarge().Cached().RunsOn(), false, 15, opts).
		workflow(runner.Pipeline("sdk-rust-dev"), "check --targets=sdk/rust", runner.XLarge().DaggerInDocker().RunsOn(), true, 15, opts).
		workflow(runner.Pipeline("sdk-php"), "check --targets=sdk/php", runner.Medium().Cached().RunsOn(), false, 10, opts).
		workflow(runner.Pipeline("sdk-php-dev"), "check --targets=sdk/php", runner.Medium().DaggerInDocker().RunsOn(), true, 10, opts)
}

func (ci *CI) workflow(
	name string,
	command string,
	runner []string,
	devEngine bool,
	timeout int,
	opts dagger.GhaWithPipelineOpts,
) *CI {
	opts.Runner = runner
	opts.TimeoutMinutes = timeout
	if devEngine {
		opts.DaggerVersion = "."
	} else {
		opts.DaggerVersion = daggerVersion
	}
	command = fmt.Sprintf("--docker-cfg=file:$HOME/.docker/config.json %s", command)
	ci.Gha = ci.Gha.WithPipeline(name, command, opts)
	return ci
}

// Generate Github Actions pipelines to call our Dagger pipelines
func (ci *CI) Generate() *dagger.Directory {
	return ci.Gha.Config()
}

func (ci *CI) Check(ctx context.Context) error {
	return dag.Dirdiff().AssertEqual(ctx, ci.Gha.Settings().Repository(), ci.Generate(), []string{".github/workflows"})
}
