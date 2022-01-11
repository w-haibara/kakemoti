package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

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
			_, _ = runString(t, fmt.Sprintf("make workflow-gen asl=%s", tt.asl))
			ctx := contextobj.New(context.Background())
			ctx = contextobj.Set(ctx, "aaa", 111)
			out, err := StartExecution(ctx, Options{
				Logfile: "",
				Input:   tt.inputFile,
				ASL:     "workflow.json",
				Timeout: 0,
			})
			if err != nil {
				t.Fatal("StartExecution() failed:", err)
			}
			want, err := os.ReadFile(tt.wantFile)
			if err != nil {
				t.Fatal("os.ReadFile(tt.wantFile) failed:", err)
			}
			if !jsonEqual(t, []byte(out), want) {
				t.Fatalf("FATAL\nWANT: [%s]\nGOT : [%s]\n", want, out)
			}
		})
	}
}

func runString(t *testing.T, str ...string) (out1, out2 string) {
	s := make([]string, 0)
	for _, v := range str {
		s = append(s, strings.Split(v, " ")...)
	}
	return run(t, s[0], s[1:])
}

func run(t *testing.T, name string, args []string) (out1, out2 string) {
	cmd := exec.Command(name, args...) // #nosec G204
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal("cmd.StdoutPipe() failed", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal("cmd.StderrPipe() failed", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal("cmd.Start() failed", err)
	}
	o1, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatal("io.ReadAll(stdout) failed", err)
	}
	o2, err := io.ReadAll(stderr)
	if err != nil {
		t.Fatal("io.ReadAll(stderr) failed", err)
	}
	err = cmd.Wait()
	t.Logf("cmd: [%s %v]\n====== stdout ======\n%s\n====== stderr ======\n%s\n",
		name, args, o1, o2)
	if err != nil {
		t.Fatalf("run(t, %s, %v) failed: %v", name, args, err)
	}
	return string(o1), string(o2)
}

func jsonEqual(t *testing.T, b1, b2 []byte) bool {
	var v1, v2 interface{}
	if err := json.Unmarshal(b1, &v1); err != nil {
		t.Fatal("json.Unmarshal(b1, &v1) failed:", err)
	}
	if err := json.Unmarshal(b2, &v2); err != nil {
		t.Fatal("json.Unmarshal(b2, &v2) failed:", err)
	}
	return reflect.DeepEqual(v1, v2)
}
