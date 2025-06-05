package judge

import (
	"errors"
	"fmt"

	"github.com/TrueHopolok/braincode-/judge/lua"
	"github.com/TrueHopolok/braincode-/judge/ml"
)

var errManyCheckers = errors.New("checker / solution defined multiple times")

func NewProblem(doc ml.Document) (Problem, error) {
	var gens []InputGenerator

	if doc.Lua != "" {
		gens = append(gens, NewLuaGenerator(doc.Lua))
	}
	if doc.GeneratorBF != "" {
		gen, err := NewBFGenerator(doc.GeneratorBF)
		if err != nil {
			return Problem{}, fmt.Errorf("provided brainfunk input generator is invalid: %w", err)
		}
		gens = append(gens, gen)
	}

	var checker OutputChecker
	if doc.CheckerBF != "" {
		c, err := NewBFChecker(doc.CheckerBF)
		if err != nil {
			return Problem{}, fmt.Errorf("provided brainfunk output checker is invalid: %w", err)
		}
		checker = c
	}
	if doc.SolutionBF != "" {
		if checker != nil {
			return Problem{}, errManyCheckers
		}
		c, err := NewBFSolution(doc.SolutionBF, doc.Instructions, doc.Steps, doc.Memory)
		if err != nil {
			return Problem{}, fmt.Errorf("provided brainfunk solution is invalid: %w", err)
		}
		checker = c
	}
	if doc.Lua != "" {
		if checker != nil {
			return Problem{}, errManyCheckers
		}
		c, err := NewLuaChecker(doc.Lua)
		if err != nil && !errors.Is(err, lua.ErrNotAChecker) {
			return Problem{}, fmt.Errorf("provided lua source is invalid: %w", err)
		} else if err == nil {
			checker = c
		}
	}

	if checker == nil {
		return Problem{}, errors.New("no checker provided")
	}

	return Problem{
		InputGenerator: CombineGenerators(gens...),
		OutputChecker:  checker,
		Instructions:   doc.Instructions,
		Steps:          doc.Steps,
		Memory:         doc.Memory,
	}, nil
}
