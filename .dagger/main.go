package main

import (
	"context"
	"dagger/workshop/internal/dagger"
)

type Workshop struct{}

func (m *Workshop) Build(
	// +defaultPath=/
	source *dagger.Directory,
) *dagger.File {
	return dag.Go().
		Build(source, dagger.GoBuildOpts{
			Trimpath: true,
			Platform: dagger.Platform("linux/amd64"),
		})
	// 	return dag.Container().
	// 		From("golang").
	// 		WithWorkdir("/work").
	// 		WithDirectory(".", source).
	// 		WithExec([]string{"go", "build", "-o", "app"}).
	// 		File("/work/app")
}

func (m *Workshop) Publish(
	ctx context.Context,
	// +defaultPath=/
	source *dagger.Directory,
) error {
	ctr := dag.Container().
		From("alpine").
		WithFile("/usr/local/bin/app", m.Build(source))

	_, err := ctr.Publish(ctx, "ttl.sh/workshop-dagger:1h")
	if err != nil {
		return err
	}
	return nil
}

func (m *Workshop) Lint(
	// +defaultPath=/
	source *dagger.Directory,
) *dagger.Container {
	return dag.GolangciLint().
		Run(source)
}

func (m *Workshop) Test(
	// +defaultPath=/
	source *dagger.Directory,
) *dagger.Container {
	return dag.Go().
		WithSource(source).
		Exec([]string{"go", "test", "-v", "./..."})
	//return dag.Container().
	//	From("golang").
	//	WithWorkdir("/work").
	//	WithDirectory(".", source).
	//	WithExec([]string{"go", "test", "-v", "./..."})
}

func (m *Workshop) Container(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("alpine").
		WithFile("/usr/local/bin/app", m.Build(source))
}
