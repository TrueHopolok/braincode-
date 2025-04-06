package judge

type listGenerator [][]string

func NewListGenerator(tests [][]string) InputGenerator {
	return listGenerator(tests)
}

func (l listGenerator) GenerateInput() ([][]string, error) {
	return l, nil
}
