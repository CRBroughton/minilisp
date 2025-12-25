package main

func builtinAdd(args []*Expr) *Expr {
	sum := 0
	for _, arg := range args {
		sum += arg.Num
	}
	return makeNum(sum)
}

func builtinSub(args []*Expr) *Expr {
	if len(args) == 0 {
		return makeNum(0)
	}
	result := args[0].Num
	for i := 1; i < len(args); i++ {
		result -= args[i].Num
	}
	return makeNum(result)
}

func builtinMul(args []*Expr) *Expr {
	result := 1
	for _, args := range args {
		result *= args.Num
	}
	return makeNum(result)
}

func builtinDiv(args []*Expr) *Expr {
	return makeNum(args[0].Num / args[1].Num)
}
