package gantry

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ad-freiburg/gantry/types"
)

const def = `version: "2.0"
steps:
  a:
    image: alpine
  b:
    image: alpine
    after:
      - a
services:
  c:
    build:
      context: ./dummy
    depends_on:
      - b
`
const env = `steps:
  b:
    ignore: true
services:
  c:
    keep_alive: replace
`

func checkCallsAndCalled(t *testing.T, runner *NoopRunner, key string, calls int, called int) {
	if c := runner.NumCalls(key); c != calls {
		t.Errorf("incorrect NumCalls for '%s', got: '%d', wanted '%d'", key, c, calls)
	}
	if c := runner.NumCalled(key); c != called {
		t.Errorf("incorrect NumCalled for '%s', got: '%d', wanted '%d'", key, c, called)
	}
}

func setupDefAndEnv(def string, env string) (string, string) {
	tmpDef, err := ioutil.TempFile("", "def")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(tmpDef.Name(), []byte(def), 0644)
	if err != nil {
		os.Remove(tmpDef.Name())
		log.Fatal(err)
	}
	tmpEnv, err := ioutil.TempFile("", "env")
	if err != nil {
		os.Remove(tmpDef.Name())
		log.Fatal(err)
	}
	err = ioutil.WriteFile(tmpEnv.Name(), []byte(env), 0644)
	if err != nil {
		os.Remove(tmpDef.Name())
		os.Remove(tmpEnv.Name())
		log.Fatal(err)
	}
	return tmpDef.Name(), tmpEnv.Name()
}

func TestPipelineGetRunnerForMeta(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		stepname string
		runner   *NoopRunner
	}{
		{"a", localRunner},
		{"b", noopRunner},
		{"c", localRunner},
	}

	for _, c := range cases {
		if v := p.GetRunnerForMeta(p.Definition.Steps[c.stepname].Meta); v != c.runner {
			t.Errorf("incorrect runner for '%s', got: '%#v', wanted '%#v'", c.stepname, v, c.runner)
		}
	}
}

func TestPipelineBuildImages(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ImageBuilder(c)", localRunner, 0, 0},
	}

	if err := p.BuildImages(false); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineBuildImagesForced(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ImageBuilder(c)", localRunner, 0, 0},
	}

	if err := p.BuildImages(true); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelinePullImages(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ImageExistenceChecker(a)", localRunner, 1, 1},
		{"ImageExistenceChecker(b)", noopRunner, 1, 1},
		{"ImageExistenceChecker(c)", localRunner, 0, 0},
		{"ImagePuller(a)", localRunner, 0, 0},
		{"ImagePuller(b)", noopRunner, 0, 0},
		{"ImagePuller(c)", localRunner, 0, 0},
	}

	if err := p.PullImages(false); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelinePullImagesForced(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ImageExistenceChecker(a)", localRunner, 1, 1},
		{"ImageExistenceChecker(b)", noopRunner, 1, 1},
		{"ImageExistenceChecker(c)", localRunner, 0, 0},
		{"ImagePuller(a)", localRunner, 1, 1},
		{"ImagePuller(b)", noopRunner, 1, 1},
		{"ImagePuller(c)", localRunner, 0, 0},
	}

	if err := p.PullImages(true); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineKillContainers(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerKiller(a)", localRunner, 1, 1},
		{"ContainerKiller(b)", noopRunner, 1, 1},
		{"ContainerKiller(c)", localRunner, 1, 1},
		{"ContainerRemover(a)", localRunner, 1, 1},
		{"ContainerRemover(b)", noopRunner, 1, 1},
		{"ContainerRemover(c)", localRunner, 1, 1},
	}

	if err := p.KillContainers(false); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineKillContainersPreRun(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerKiller(a)", localRunner, 1, 1},
		{"ContainerKiller(b)", noopRunner, 1, 1},
		{"ContainerKiller(c)", localRunner, 0, 0},
		{"ContainerRemover(a)", localRunner, 1, 1},
		{"ContainerRemover(b)", noopRunner, 1, 1},
		{"ContainerRemover(c)", localRunner, 0, 0},
	}

	if err := p.KillContainers(true); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineRemoveContainers(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerRemover(a)", localRunner, 1, 1},
		{"ContainerRemover(b)", noopRunner, 1, 1},
		{"ContainerRemover(c)", localRunner, 1, 1},
	}

	if err := p.RemoveContainers(false); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineRemoveContainersPreRun(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerRemover(a)", localRunner, 1, 1},
		{"ContainerRemover(b)", noopRunner, 1, 1},
		{"ContainerRemover(c)", localRunner, 0, 0},
	}

	if err := p.RemoveContainers(true); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineCreateNetwork(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	p.Network = Network("test")

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"NetworkCreator(test)", localRunner, 1, 1},
	}

	if err := p.CreateNetwork(); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineRemoveNetwork(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	p.Network = Network("test")

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"NetworkRemover(test)", localRunner, 1, 1},
	}

	if err := p.RemoveNetwork(); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineExecuteSteps(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	p.Network = Network("test")

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerKiller(a)", localRunner, 1, 1},
		{"ContainerRemover(a)", localRunner, 1, 1},
		{"ContainerRunner(a,test)", localRunner, 1, 1},
		{"ContainerKiller(b)", noopRunner, 1, 1},
		{"ContainerRemover(b)", noopRunner, 1, 1},
		{"ContainerRunner(b,test)", noopRunner, 1, 1},
		{"ContainerKiller(c)", localRunner, 1, 1},
		{"ContainerRemover(c)", localRunner, 1, 1},
		{"ContainerRunner(c,test)", localRunner, 1, 1},
	}

	if err := p.ExecuteSteps(); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineRemoveTempDirData(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(`version: "2.0"
#! TEMP_DIR_IF_EMPTY ${TEMP_STORAGE}
steps:
  a:
    volumes:
    - ${TEMP_STORAGE}:/input
`, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	p.Network = Network("test")

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerKiller(TempDirCleanUp)", localRunner, 1, 1},
		{"ContainerRemover(TempDirCleanUp)", localRunner, 2, 2},
		{"ContainerRunner(TempDirCleanUp,test)", localRunner, 1, 1},
	}

	if err := p.RemoveTempDirData(); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineRemoveTempDirDataNoTempDirs(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(`version: "2.0"
steps:
  a:
`, "")
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	p.Network = Network("test")

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerKiller(TempDirCleanUp)", localRunner, 0, 0},
		{"ContainerRemover(TempDirCleanUp)", localRunner, 0, 0},
		{"ContainerRunner(TempDirCleanUp,test)", localRunner, 0, 0},
	}

	if err := p.RemoveTempDirData(); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineLogs(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, env)
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, tmpEnv, types.StringMap{}, types.StringSet{}, types.StringSet{})
	if err != nil {
		t.Errorf("unexpected error creating pipeline: '%#v'", err)
	}
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner

	cases := []struct {
		key    string
		runner *NoopRunner
		calls  int
		called int
	}{
		{"ContainerLogReader(a,false)", localRunner, 1, 1},
		{"ContainerLogReader(b,false)", noopRunner, 1, 1},
		{"ContainerLogReader(c,false)", localRunner, 1, 1},
	}

	if err := p.Logs(false); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	for _, c := range cases {
		checkCallsAndCalled(t, c.runner, c.key, c.calls, c.called)
	}
}

func TestPipelineCheck(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(def, "")
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, "", types.StringMap{}, types.StringSet{}, types.StringSet{})
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	if err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	if err := p.Check(); err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
}

func TestPipelineCheckNoContainerInformation(t *testing.T) {
	tmpDef, tmpEnv := setupDefAndEnv(`version: "2.0"
steps:
  a:
`, "")
	defer os.Remove(tmpDef)
	defer os.Remove(tmpEnv)

	p, err := NewPipeline(tmpDef, "", types.StringMap{}, types.StringSet{}, types.StringSet{})
	localRunner := NewNoopRunner(false)
	p.localRunner = localRunner
	noopRunner := NewNoopRunner(false)
	p.noopRunner = noopRunner
	if err != nil {
		t.Errorf("unexpected error, got: '%#v', wanted 'nil'", err)
	}
	if err := p.Check(); err == nil {
		t.Errorf("expected error, got: nil")
	}
}

func TestPipelineDefinitionCheckVersion(t *testing.T) {
	p := PipelineDefinition{}
	cases := []struct {
		version string
		err     string
	}{
		{
			version: "",
			err:     fmt.Sprintf("not supported compose file format version: got: 1.0 want >= %d.%d", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin),
		},
		{
			version: "foo",
			err:     "invalid compose file format version: foo",
		},
		{
			version: "x.y",
			err:     "invalid compose file format version: x.y",
		},
		{
			version: "0.y",
			err:     "invalid compose file format version: 0.y",
		},
		{
			version: "1.0",
			err:     fmt.Sprintf("not supported compose file format version: got: 1.0 want >= %d.%d", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin),
		},
		{
			version: "1",
			err:     fmt.Sprintf("not supported compose file format version: got: 1.0 want >= %d.%d", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin),
		},
		{
			version: "2.-1",
			err:     fmt.Sprintf("not supported compose file format version: got: 2.-1 want >= %d.%d", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin),
		},
		{
			version: fmt.Sprintf("%d.%d", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin),
			err:     "",
		},
		{
			version: fmt.Sprintf("%d", DockerComposeFileFormatMajorMin+1),
			err:     "",
		},
		{
			version: fmt.Sprintf("%d.%d", DockerComposeFileFormatMajorMin+1, DockerComposeFileFormatMinorMin),
			err:     "",
		},
		{
			version: fmt.Sprintf("%d.%d", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin+1),
			err:     "",
		},
		{
			version: fmt.Sprintf("%d.%d.1", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin+1),
			err:     fmt.Sprintf("invalid compose file format version: %d.%d.1", DockerComposeFileFormatMajorMin, DockerComposeFileFormatMinorMin+1),
		},
	}
	for i, c := range cases {
		p.Version = c.version
		err := p.checkVersion()
		if c.err != "" && err == nil {
			t.Errorf("expected error @%d, got: nil", i)
			continue
		}
		if c.err != "" && c.err != err.Error() {
			t.Errorf("incorrect error @%d, got: '%s', wanted: '%s'", i, err.Error(), c.err)
		}
	}
}
