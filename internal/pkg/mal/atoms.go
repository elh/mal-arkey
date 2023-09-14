package mal

// This is a gross hack to dig myself out of the fact that I'm not using pointers
// and my variables are immutable. I don't have a mechanism to carry around the
// atom values and make mutable.

// Global atom map from atom uuid to current atom value
// This is a internal implementation detail. This works for Env.
var atoms = map[string]Sexpr{}
