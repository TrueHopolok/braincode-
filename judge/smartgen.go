package judge

type smartGenerator struct{ InputGenerator }

func (s *smartGenerator) GenerateInput() ([][]string, error) {
	const maxCachable = 1024 * 128 // 128 KiB

	res, err := s.InputGenerator.GenerateInput()
	if err != nil {
		return nil, err
	}

	if _, ok := s.InputGenerator.(listGenerator); ok {
		return res, err
	}

	size := inputSize(res)

	if size <= maxCachable {
		s.InputGenerator = listGenerator(res)
	}

	return res, err
}

func inputSize(input [][]string) (res int) {
	for _, ss := range input {
		for _, s := range ss {
			res += len(s)
		}
	}
	return res
}
