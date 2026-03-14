// Code generated from Pkl module `makila.minecraftgo.properties`. DO NOT EDIT.
package config

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Properties struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Properties
func LoadFromPath(ctx context.Context, path string) (ret Properties, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return ret, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Properties
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Properties, error) {
	var ret Properties
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
