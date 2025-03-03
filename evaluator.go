package rules

import (
	"errors"
	"fmt"
	"time"
)

func Evaluate(expr Expr, params map[string]interface{}) (bool, error) {
	if expr == nil {
		return false, errors.New("invalid input expression")
	}

	res, err := evaluate(expr, params)
	if err != nil {
		return false, err
	}

	l, ok := res.(*BoolLiteral)
	if !ok {
		return false, errors.New("invalid result expression")
	}

	return l.Val, nil
}

func evaluate(expr Expr, params map[string]interface{}) (Expr, error) {
	if expr == nil {
		return nil, errors.New("invalid expression")
	}

	switch e := expr.(type) {
	case *ParentExpr:
		return evaluate(e.Expr, params)
	case *BinaryExpr:
		l, err := evaluate(e.LHS, params)
		if err != nil {
			return nil, err
		}

		r, err := evaluate(e.RHS, params)
		if err != nil {
			return nil, err
		}

		return compute(e.Op, l, r)
	case *Ident:
		varName := e.Val
		varVal, ok := params[varName]
		if !ok {
			return nil, fmt.Errorf("param %s not found", varName)
		}

		return toLiteral(varVal)
	}

	return expr, nil
}

func compute(op Op, lhs, rhs Expr) (*BoolLiteral, error) {
	switch op {
	case AND:
		return computeAND(lhs, rhs)
	case OR:
		return computeOR(lhs, rhs)
	case EQ:
		return computeEQ(lhs, rhs)
	case NEQ:
		return computeNEQ(lhs, rhs)
	case LT:
		return computeLT(lhs, rhs)
	case LTE:
		return computeLTE(lhs, rhs)
	case GT:
		return computeGT(lhs, rhs)
	case GTE:
		return computeGTE(lhs, rhs)
	case IN:
		return computeIN(lhs, rhs)
	}

	return nil, errors.New("invalid operator")
}

func computeOR(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *BoolLiteral:
		r, ok := rhs.(*BoolLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val || r.Val)}, nil
		}
	}

	return &BoolLiteral{Val: true}, nil
}

func computeAND(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *BoolLiteral:
		r, ok := rhs.(*BoolLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val && r.Val)}, nil
		}
	}

	return &BoolLiteral{Val: true}, nil
}

func computeEQ(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val == r.Val)}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val == r.Val)}, nil
		}
	case *BoolLiteral:
		r, ok := rhs.(*BoolLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val == r.Val)}, nil
		}
	case *TimeLiteral:
		tr, ok := rhs.(*TimeLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val == tr.Val)}, nil
		}

		sr, ok := rhs.(*StringLiteral)
		if ok {
			dt, err := time.Parse(time.RFC3339, sr.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: (l.Val == dt)}, nil
		}
	case *DurationLiteral:
		rv, ok := rhs.(*StringLiteral)
		if ok {
			v, err := time.ParseDuration(rv.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val == v}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}

func computeNEQ(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val != r.Val)}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val != r.Val)}, nil
		}
	case *BoolLiteral:
		r, ok := rhs.(*BoolLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val != r.Val)}, nil
		}
	case *TimeLiteral:
		tr, ok := rhs.(*TimeLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val != tr.Val)}, nil
		}

		sr, ok := rhs.(*StringLiteral)
		if ok {
			dt, err := time.Parse(time.RFC3339, sr.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: (l.Val != dt)}, nil
		}
	case *DurationLiteral:
		rv, ok := rhs.(*StringLiteral)
		if ok {
			v, err := time.ParseDuration(rv.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val != v}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}

func computeLT(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val < r.Val)}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val < r.Val)}, nil
		}
	case *TimeLiteral:
		tr, ok := rhs.(*TimeLiteral)
		if ok {
			return &BoolLiteral{Val: l.Val.Before(tr.Val)}, nil
		}

		sr, ok := rhs.(*StringLiteral)
		if ok {
			dt, err := time.Parse(time.RFC3339, sr.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val.Before(dt)}, nil
		}
	case *DurationLiteral:
		rv, ok := rhs.(*StringLiteral)
		if ok {
			v, err := time.ParseDuration(rv.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val < v}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}

func computeLTE(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val <= r.Val)}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val <= r.Val)}, nil
		}
	case *TimeLiteral:
		tr, ok := rhs.(*TimeLiteral)
		if ok {
			return &BoolLiteral{Val: l.Val.Before(tr.Val)}, nil
		}

		sr, ok := rhs.(*StringLiteral)
		if ok {
			dt, err := time.Parse(time.RFC3339, sr.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val.Before(dt)}, nil
		}
	case *DurationLiteral:
		rv, ok := rhs.(*StringLiteral)
		if ok {
			v, err := time.ParseDuration(rv.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val <= v}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}

func computeGT(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val > r.Val)}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val > r.Val)}, nil
		}
	case *TimeLiteral:
		tr, ok := rhs.(*TimeLiteral)
		if ok {
			return &BoolLiteral{Val: l.Val.After(tr.Val)}, nil
		}

		sr, ok := rhs.(*StringLiteral)
		if ok {
			dt, err := time.Parse(time.RFC3339, sr.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val.After(dt)}, nil
		}
	case *DurationLiteral:
		rv, ok := rhs.(*StringLiteral)
		if ok {
			v, err := time.ParseDuration(rv.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val > v}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}

func computeGTE(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val >= r.Val)}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringLiteral)
		if ok {
			return &BoolLiteral{Val: (l.Val >= r.Val)}, nil
		}
	case *TimeLiteral:
		tr, ok := rhs.(*TimeLiteral)
		if ok {
			return &BoolLiteral{Val: l.Val.After(tr.Val)}, nil
		}

		sr, ok := rhs.(*StringLiteral)
		if ok {
			dt, err := time.Parse(time.RFC3339, sr.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val.After(dt)}, nil
		}
	case *DurationLiteral:
		rv, ok := rhs.(*StringLiteral)
		if ok {
			v, err := time.ParseDuration(rv.Val)
			if err != nil {
				return nil, err
			}

			return &BoolLiteral{Val: l.Val >= v}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}

func computeIN(lhs, rhs Expr) (*BoolLiteral, error) {
	switch l := lhs.(type) {
	case *NumberLiteral:
		r, ok := rhs.(*NumberSliceLiteral)
		if ok {
			res := false
			for _, v := range r.Val {
				if l.Val == v {
					res = true
					break
				}
			}

			return &BoolLiteral{Val: res}, nil
		}
	case *StringLiteral:
		r, ok := rhs.(*StringSliceLiteral)
		if ok {
			res := false
			for _, v := range r.Val {
				if l.Val == v {
					res = true
					break
				}
			}

			return &BoolLiteral{Val: res}, nil
		}
	}

	return nil, fmt.Errorf(`cannot convert "%s" to %s`, rhs.String(), lhs.Type())
}
