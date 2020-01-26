package cfgtest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	cfg "github.com/mbrt/gmailctl/pkg/config/v1alpha3"
	"github.com/mbrt/gmailctl/pkg/parser"
)

func TestParseEval(t *testing.T) {
	expr := or(
		and(
			fn(parser.FunctionList, parser.OperationOr,
				"list1.gm.com", "list2@gm.com"),
			fn1(parser.FunctionTo, "me@gmail.com"),
		),
		fn1(parser.FunctionSubject, "Subject"),
		fn(parser.FunctionFrom, parser.OperationOr, "@google.com", "b"),
		fn1(parser.FunctionHas, "Important message"),
		fn1(parser.FunctionHas, "foo@bar.com"),
	)
	eval, err := NewEvaluator(expr)
	if err != nil {
		t.Fatalf("NewEvaluator failed: %v", err)
	}
	fmt.Printf("eval: %+v", eval)

	tests := []struct {
		name        string
		message     cfg.Message
		expectMatch bool
	}{
		{
			name: "subject",
			message: cfg.Message{
				Subject: "contains subject yes",
			},
			expectMatch: true,
		},
		{
			name: "list with @",
			message: cfg.Message{
				Lists: []string{"list1@gm.com"},
				To:    []string{"me@gmail.com"},
			},
			expectMatch: true,
		},
		{
			name: "from google",
			message: cfg.Message{
				From: "someone@google.com",
			},
			expectMatch: true,
		},
		{
			name: "has from",
			message: cfg.Message{
				From: "foo@bar.com",
			},
			expectMatch: true,
		},
		// TODO: Implement body match.
		// {
		// 	name: "has body",
		// 	message: cfg.Message{
		// 		Body: "important message",
		// 	},
		// 	expectMatch: true,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			match := eval.Match(tc.message)
			assert.Equal(t, tc.expectMatch, match)
		})
	}
}
