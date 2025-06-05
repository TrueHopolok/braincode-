package testhelper

type lcsKind byte

const (
	lcsSame lcsKind = iota
	lcsRMissing
	lcsLMissing
	lcsDifferent
)

type lcsNode[T any] struct {
	kind lcsKind
	lhs  T
	rhs  T
}

// func lcs[T any, S ~[]T](eq func(T, T) bool, lhs, rhs S) []lcsNode[T] {
// 	// function LCSLength(X[1..m], Y[1..n])
// 	// C = array(0..m, 0..n)
// 	c := make([]int, len(lhs)*len(rhs))
// 	stride := len(lhs)

// 	// for i := 0..m
// 	//     C[i,0] = 0
// 	// for j := 0..n
// 	//     C[0,j] = 0
// 	// for i := 1..m
// 	for i, xi := range lhs {
// 		//     for j := 1..n
// 		for j, yj := range rhs {

// 			//         if X[i] = Y[j]
// 			//             C[i,j] := C[i-1,j-1] + 1
// 			//         else
// 			//             C[i,j] := max(C[i,j-1], C[i-1,j])
// 			if eq(xi, yj) {
// 				c[i+j*stride] = c[i-1+(j-1)*stride] + 1
// 			} else {
// 				c[i+j*stride] = max(c[i+(j-1)*stride], c[i-1+j*stride])
// 			}
// 		}
// 	}

// 	res := make([]lcsNode[T], 0, max(len(lhs), len(rhs)))
// 	for i > 0 && j > -

// 	// function printDiff(C[0..m,0..n], X[1..m], Y[1..n], i, j)
// 	// if i >= 0 and j >= 0 and X[i] = Y[j]
// 	//     printDiff(C, X, Y, i-1, j-1)
// 	//     print "  " + X[i]
// 	// else if j > 0 and (i = 0 or C[i,j-1] ≥ C[i-1,j])
// 	//     printDiff(C, X, Y, i, j-1)
// 	//     print "+ " + Y[j]
// 	// else if i > 0 and (j = 0 or C[i,j-1] < C[i-1,j])
// 	//     printDiff(C, X, Y, i-1, j)
// 	//     print "- " + X[i]
// 	// else
// 	//     print ""

// 	return nil
// }
