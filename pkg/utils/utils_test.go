package utils_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Ensono/taskctl/pkg/utils"
)

func TestConvertEnv(t *testing.T) {
	type args struct {
		env map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{args: args{env: map[string]string{"key1": "val1"}}, want: []string{"key1=val1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ConvertEnv(tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{file: filepath.Join(cwd, "utils.go")}, want: true, name: "file exists"},
		{args: args{file: filepath.Join(cwd, "utils_test.go")}, want: true, name: "test file exists"},
		{args: args{file: filepath.Join(cwd, "manifesto.txt")}, want: false, name: "file does not exist"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.FileExists(tt.args.file); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExitError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{err: &exec.ExitError{}}, want: true},
		{args: args{err: fmt.Errorf("%w", &exec.ExitError{})}, want: true},
		{args: args{err: os.ErrNotExist}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.IsExitError(tt.args.err); got != tt.want {
				t.Errorf("IsExitError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "HTTP URL", args: args{s: "http://github.com/"}, want: true},
		{name: "HTTPS URL", args: args{s: "https://github.com/"}, want: true},
		{name: "Windows path", args: args{s: "C:\\Windows"}, want: false},
		{name: "Mailto", args: args{s: "mailto:admin@example.org"}, want: false},
		{name: "Invalid", args: args{s: "::::::::not-parsed"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.IsURL(tt.args.s); got != tt.want {
				t.Errorf("IsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastLine(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name  string
		args  args
		wantL string
	}{
		{args: args{r: strings.NewReader("line1\nline2")}, wantL: "line2"},
		{args: args{r: strings.NewReader("\nline1")}, wantL: "line1"},
		{args: args{r: strings.NewReader("line1\n")}, wantL: "line1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotL := utils.LastLine(tt.args.r); gotL != tt.wantL {
				t.Errorf("LastLine() = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	type args struct {
		m interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantKeys []string
	}{
		{args: args{m: map[string]bool{"a": true, "b": false}}, wantKeys: []string{"a", "b"}},
		{args: args{m: []string{"a", "b"}}, wantKeys: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys := utils.MapKeys(tt.args.m)
			for _, v := range tt.wantKeys {
				var found bool
				for _, vv := range gotKeys {
					if v == vv {
						found = true
						break
					}
				}
				if found == false {
					t.Errorf("MapKeys() = %v, want %v", gotKeys, tt.wantKeys)
				}
			}
		})
	}
}

func TestRenderString(t *testing.T) {
	type args struct {
		tmpl      string
		variables map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{args: args{tmpl: "hello, {{ .Name }}!", variables: map[string]interface{}{"Name": "world"}}, want: "hello, world!"},
		{args: args{tmpl: "hello, {{ .Name | default \"John\" }}!", variables: map[string]interface{}{"Name": ""}}, want: "hello, John!"},
		{args: args{tmpl: "hello, {{ .Name }}!", variables: make(map[string]interface{})}, wantErr: true},
		{args: args{tmpl: "hello, {{ .Name", variables: make(map[string]interface{})}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.RenderString(tt.args.tmpl, tt.args.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RenderString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGetwd(t *testing.T) {
	wd, _ := os.Getwd()
	if wd != utils.MustGetwd() {
		t.Error()
	}

}

func TestMustGetUserHomeDir(t *testing.T) {
	err := os.Setenv("HOME", "/test")
	if err != nil {
		t.Fatal(err)
	}
	hd := utils.MustGetUserHomeDir()
	if hd != "/test" {
		t.Error()
	}

}

// Test envfile

func TestUtils_Envfile(t *testing.T) {

	envfile := utils.NewEnvFile(func(e *utils.Envfile) {
		// e.Delay =
		e.Exclude = []string{}
		e.Include = []string{}
		// e.Path = def.Envfile.Path
		e.Modify = []utils.ModifyEnv{
			{Pattern: "", Operation: "lower"},
		}
		e.Quote = false
	})

	if err := envfile.Validate(); err == nil {
		t.Error("failed to validate")
	}

	if envfile.GeneratedDir != ".taskctl" {
		t.Error("generated dir not set correctly")
	}

}
