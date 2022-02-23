package cli

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/w-haibara/kakemoti/compiler"
)

func TestWorkflowExec(t *testing.T) {
	tests := []struct {
		name, asl, inputFile, wantFile string
	}{
		{"pass", "pass", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"pass(result)", "pass_result", "_workflow/inputs/input1.json", "_workflow/outputs/output5.json"},
		{"pass(chain)", "pass_chain", "_workflow/inputs/input1.json", "_workflow/outputs/output1.json"},
		{"pass(parameters)", "pass_parameters", "_workflow/inputs/input1.json", "_workflow/outputs/output11.json"},
		{"pass(intrinsic)", "pass_intrinsic", "_workflow/inputs/input4.json", "_workflow/outputs/output12.json"},
		{"pass(ctxobj)", "pass_ctxobj", "_workflow/inputs/input1.json", "_workflow/outputs/output10.json"},
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
		{"task(ctxobj)", "task_ctxobj", "_workflow/inputs/input1.json", "_workflow/outputs/output10.json"},
		{"map", "map", "_workflow/inputs/input7.json", "_workflow/outputs/output13.json"},
		{"map(concurrency)", "map_concurrency", "_workflow/inputs/input7.json", "_workflow/outputs/output13.json"},
		{"map(ctxobj)", "map_ctxobj", "_workflow/inputs/input8.json", "_workflow/outputs/output14.json"},
		{"map(ctxobj2)", "map_ctxobj2", "_workflow/inputs/input8.json", "_workflow/outputs/output15.json"},
	}

	if err := os.Chdir("../"); err != nil {
		t.Error(`os.Chdir("../"):`, err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			coj := new(compiler.CtxObj)
			v, err := coj.SetByString("$.aaa", 111)
			if err != nil {
				t.Error("coj.Set() failed:", err)
			}
			coj = v
			out, err := WorkflowExec(ctx, coj, Options{
				Logfile: "",
				Input:   tt.inputFile,
				ASL:     "_workflow/asl/" + tt.asl + ".asl.json",
				Timeout: 0,
			})
			if err != nil {
				t.Error("WorkflowExec() failed:", err)
				return
			}
			want, err := os.ReadFile(tt.wantFile)
			if err != nil {
				t.Error("os.ReadFile(tt.wantFile) failed:", err)
				return
			}

			var v1, v2 interface{}
			if err := json.Unmarshal(out, &v1); err != nil {
				t.Error("json.Unmarshal(b1, &v1) failed:", err)
				return
			}
			if err := json.Unmarshal(want, &v2); err != nil {
				t.Error("json.Unmarshal(b2, &v2) failed:", err)
				return
			}

			if d := cmp.Diff(v1, v2); d != "" {
				t.Errorf("FATAL\nGOT:\n%s\n\nWANT:\n%s\n\nDIFF:\n%s", out, want, d)
				return
			}
		})
	}
}
