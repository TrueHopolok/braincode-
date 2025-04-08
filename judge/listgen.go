package judge

type listGenerator [][]string

// NewListGenerator creates a generator that simply returns tests without modification.
func NewListGenerator(tests [][]string) InputGenerator {
	return listGenerator(tests)
}

func (l listGenerator) GenerateInput() ([][]string, error) {
	return l, nil
}
