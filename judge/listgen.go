package judge

type listGenerator [][]string

// NewListGenerator creates a generator that simply returns tests without modification.
func NewListGenerator(tests [][]string) InputGenerator {
	return listGenerator(tests)
}

func (l listGenerator) GenerateInput() ([][]string, error) {
	return l, nil
}

type multiGen []InputGenerator

func (l multiGen) GenerateInput() ([][]string, error) {
	var acc [][]string
	for _, g := range l {
		res, err := g.GenerateInput()
		if err != nil {
			return nil, err
		}
		acc = append(acc, res...)
	}
	return acc, nil
}

func CombineGenerators(gens ...InputGenerator) InputGenerator {
	res := make([]InputGenerator, 0, len(gens))

	for _, gen := range gens {
		if res != nil {
			res = append(res, gen)
		}
	}

	if len(res) == 0 {
		return nilGenerator{}
	}

	if len(res) == 1 {
		return res[0]
	}

	return multiGen(res)
}

type nilGenerator struct{}

func (nilGenerator) GenerateInput() ([][]string, error) { return nil, nil }
