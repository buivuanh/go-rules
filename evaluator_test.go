package rules_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rules "github.com/sonda2208/go-rules"
)

type Evaluation struct {
	params   map[string]interface{}
	expected bool
	isError  bool
}

type TestCase struct {
	rules       string
	evaluations []Evaluation
}

func TestEvaluator(t *testing.T) {
	dt, err := time.Parse(time.RFC3339, "2019-03-28T11:39:43+07:00")
	require.NoError(t, err)

	dur2m, err := time.ParseDuration("2m")
	require.NoError(t, err)

	dur1m30s, err := time.ParseDuration("1m30s")
	require.NoError(t, err)

	dur45s, err := time.ParseDuration("45s")
	require.NoError(t, err)

	tests := []TestCase{
		{
			`{ "comparator": "||", "rules": [ { "comparator": "&&", "rules": [ { "var": "a", "op": "==", "val": 1 }, { "var": "b", "op": "==", "val": 2 } ] }, { "comparator": "&&", "rules": [ { "var": "c", "op": "==", "val": 3 }, { "var": "d", "op": "==", "val": 4 } ] } ] }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": 1,
						"b": 2,
						"c": 0,
						"d": 0,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
						"b": 0,
						"c": 3,
						"d": 4,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
						"b": 2,
						"c": 3,
						"d": 4,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
						"b": 2,
						"c": 0,
						"d": 0,
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
						"b": 0,
						"c": 3,
						"d": 0,
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
						"b": 0,
						"c": 0,
						"d": 0,
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
						"b": 2,
						"c": 0,
						"d": "string",
					},
					false,
					true,
				},
				{
					map[string]interface{}{
						"a": 1,
						"b": 2,
						"c": true,
						"d": 0,
					},
					false,
					true,
				},
			},
		},
		{
			`{ "comparator": "&&", "rules": [ { "var": "a", "op": "==", "val": 1 }, { "var": "b", "op": "==", "val": "2" }, { "var": "c", "op": "==", "val": 3 } ] }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": 1,
						"b": "2",
						"c": 3,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
						"b": "",
						"c": 3,
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
						"b": 2,
						"c": 3,
					},
					false,
					true,
				},
			},
		},
		{
			`{ "comparator": "&&", "rules": [ { "var": "a", "op": ">", "val": 1 }, { "var": "b", "op": "<", "val": 2 } ] }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": 2,
						"b": 1,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
						"b": 1,
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": 2,
						"b": -1,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 2,
						"b": 2,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "comparator": "||", "rules": [ { "var": "a", "op": ">=", "val": 1 }, { "var": "b", "op": "<", "val": 2 } ] }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": 1,
						"b": 1,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
						"b": 1,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
						"b": 2,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "!=", "val": 1 }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": 2,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "in", "val": [1, 2, 3] }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": 1,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 0,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "in", "val": ["1", "2", "3"] }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": "1",
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": "",
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "==", "val": "2019-03-28T11:39:43+07:00" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": time.Now(),
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": dt,
					},
					true,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "!=", "val": "2019-03-28T11:39:43+07:00" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": time.Now(),
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dt,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": ">", "val": "2019-03-28T11:39:43+07:00" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": time.Now(),
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dt,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "<", "val": "2019-03-28T11:39:43+07:00" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": time.Now(),
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": dt.Add(-1 * time.Hour),
					},
					true,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "==", "val": "1m30s" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": dur2m,
					},
					false,
					false,
				},
				{
					map[string]interface{}{
						"a": dur1m30s,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": 1,
					},
					false,
					true,
				},
			},
		},
		{
			`{ "var": "a", "op": "!=", "val": "1m30s" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": dur2m,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dur1m30s,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": ">", "val": "1m30s" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": dur2m,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dur1m30s,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": ">=", "val": "1m30s" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": dur2m,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dur1m30s,
					},
					true,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "<", "val": "1m30s" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": dur45s,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dur2m,
					},
					false,
					false,
				},
			},
		},
		{
			`{ "var": "a", "op": "<=", "val": "1m30s" }`,
			[]Evaluation{
				{
					map[string]interface{}{
						"a": dur1m30s,
					},
					true,
					false,
				},
				{
					map[string]interface{}{
						"a": dur2m,
					},
					false,
					false,
				},
			},
		},
	}

	for _, test := range tests {
		expr, err := rules.ParseFromJSON([]byte(test.rules))
		require.NoError(t, err)
		require.NotNil(t, expr)

		for _, e := range test.evaluations {
			res, err := rules.Evaluate(expr, e.params)
			if e.isError {
				assert.Error(t, err)
			} else {
				if !assert.NoError(t, err) {
					t.Logf("expr: %s\tparams: %v\n", test.rules, e.params)
				} else {
					if !assert.Equal(t, e.expected, res) {
						t.Log(expr.String())
					}
				}
			}
		}
	}
}

func TestExample(t *testing.T) {
	expression := `
	{
		"comparator": "||",
		"rules": [
		  {
			"comparator": "&&",
			"rules": [
			  { "var": "a", "op": "<", "val": "2019-03-28T11:39:43+07:00" },
			  { "var": "b", "op": "in", "val": [1, 2, 3] }
			]
		  },
		  {
			"comparator": "&&",
			"rules": [
			  { "var": "c", "op": "!=", "val": "string" },
			  { "var": "d", "op": ">=", "val": 4 }
			]
		  }
		]
	}
	`

	parameters := map[string]interface{}{
		"a": time.Now(),
		"b": 1,
		"c": "number",
		"d": 5,
	}

	expr, err := rules.ParseFromJSON([]byte(expression))
	require.NoError(t, err)

	res, err := rules.Evaluate(expr, parameters)
	require.NoError(t, err)
	assert.True(t, res)
}
