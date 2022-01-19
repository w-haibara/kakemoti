package cli

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/w-haibara/kakemoti/contextobj"
)

func TestStartExecution(t *testing.T) {
	tests := []struct {
		name, asl, inputFile, wantFile string
	}{
		{"pass", "pass", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"pass(result)", "pass_result", "_workflow/inputs/input1.json", "_workflow/outputs/output5.json"},
		{"pass(chain)", "pass_chain", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"pass(parameters)", "pass_parameters", "_workflow/inputs/input1.json", "_workflow/outputs/output11.json"},
		{"pass(intrinsic)", "pass_intrinsic", "_workflow/inputs/input4.json", "_workflow/outputs/output12.json"},
		{"wait", "wait", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"succeed", "succeed", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"fail", "fail", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"choice", "choice", "_workflow/inputs/input2.json", "_workflow/outputs/output2.json"},
		{"choice(fallback)", "choice_fallback", "_workflow/inputs/input2.json", "_workflow/outputs/output7.json"},
		{"choice(boolean expr)", "choice_bool", "_workflow/inputs/input5.json", "_workflow/outputs/output8.json"},
		{"choice(data test expr)", "choice_data_test", "_workflow/inputs/input6.json", "_workflow/outputs/output8.json"},
		{"parallel", "parallel", "_workflow/inputs/input2.json", "_workflow/outputs/output3.json"},
		{"task", "task", "_workflow/inputs/input1.json", "_workflow/outputs/output4.json"},
		{"task(filter)", "task_filter", "_workflow/inputs/input3.json", "_workflow/outputs/output9.json"},
		{"task(catch)", "task_catch", "_workflow/inputs/input1.json", "_workflow/inputs/input1.json"},
		{"task(retry)", "task_retry", "_workflow/inputs/input1.json", "_workflow/outputs/output8.json"},
		{"task(ctx)", "task_ctx", "_workflow/inputs/input1.json", "_workflow/outputs/output10.json"},
	}

	if err := os.Chdir("../"); err != nil {
		t.Fatal(`os.Chdir("../"):`, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := contextobj.New(context.Background())
			ctx = contextobj.Set(ctx, "aaa", 111)
			out, err := StartExecution(ctx, Options{
				Logfile: "",
				Input:   tt.inputFile,
				ASL:     "_workflow/asl/" + tt.asl + ".asl.json",
				Timeout: 0,
			})
			if err != nil {
				t.Fatal("StartExecution() failed:", err)
			}
			want, err := os.ReadFile(tt.wantFile)
			if err != nil {
				t.Fatal("os.ReadFile(tt.wantFile) failed:", err)
			}
			if d := jsonEqual(t, []byte(out), want); d != "" {
				t.Fatalf("FATAL\nGOT:\n%s\n\nWANT:\n%s\n\nDIFF:\n%s", out, want, d)
			}
		})
	}
}

func jsonEqual(t *testing.T, b1, b2 []byte) string {
	var v1, v2 interface{}
	if err := json.Unmarshal(b1, &v1); err != nil {
		t.Fatal("json.Unmarshal(b1, &v1) failed:", err)
	}
	if err := json.Unmarshal(b2, &v2); err != nil {
		t.Fatal("json.Unmarshal(b2, &v2) failed:", err)
	}
	return cmp.Diff(v1, v2)
}
