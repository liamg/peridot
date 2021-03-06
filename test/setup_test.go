package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var testContainers []*testContainer

const defaultUser = "someone"

func TestMain(m *testing.M) {

	exitCode := func() int {

		cwd, err := os.Getwd()
		if err != nil {
			return 2
		}

		cmd := exec.Command("go", "build", ".")
		cmd.Dir = filepath.Dir(cwd)
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "build failed: %s\n", err)
			return 1
		}

		return m.Run()
	}()

	for _, container := range testContainers {
		container.Destroy()
	}

	os.Exit(exitCode)
}

type testContainer struct {
	id      string
	hostDir string
	client  *client.Client
}

func (t *testContainer) Destroy() {
	_ = t.client.ContainerRemove(context.Background(), t.id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	_ = os.RemoveAll(t.hostDir)
}

func (t *testContainer) AddPeridot() error {
	return t.AddFile("../peridot")
}

func (t *testContainer) AddFile(hostPath string) error {
	data, err := ioutil.ReadFile(hostPath)
	if err != nil {
		return err
	}
	target := filepath.Join(t.hostDir, filepath.Base(hostPath))
	return ioutil.WriteFile(target, data, 0700) //nolint
}

func (t *testContainer) ReadHomeFile(relativePath string) (string, error) {
	_, exit, err := t.RunAsRoot("chmod", "-R", "777", ".")
	if err != nil {
		return "", err
	}
	if exit > 0 {
		return "", fmt.Errorf("chmod failed")
	}
	data, err := ioutil.ReadFile(filepath.Join(t.hostDir, relativePath))
	return string(data), err
}

func (t *testContainer) WriteConfig(content string) error {
	return t.WriteHomeFile(".config/peridot/config.yml", content)
}

func (t *testContainer) WriteHomeFile(relativePath, content string) error {
	//nolint
	if err := os.MkdirAll(filepath.Join(t.hostDir, filepath.Dir(relativePath)), 0777); err != nil {
		return err
	}
	//nolint
	if err := ioutil.WriteFile(filepath.Join(t.hostDir, relativePath), []byte(content), 0777); err != nil {
		return err
	}
	return nil
}

func (t *testContainer) RunAsUser(cmd string, args ...string) (string, int, error) {
	return t.run(false, cmd, args...)
}

func (t *testContainer) RunAsRoot(cmd string, args ...string) (string, int, error) {
	return t.run(true, cmd, args...)
}

func (t *testContainer) run(root bool, cmd string, args ...string) (string, int, error) {

	user := defaultUser
	home := fmt.Sprintf("/home/%s", user)

	if root {
		user = "root"
		home = "/root"
	}

	fmt.Printf("Running command on %s: %s %s\n", t.id[:12], cmd, strings.Join(args, " "))
	idResp, err := t.client.ContainerExecCreate(context.Background(), t.id, types.ExecConfig{
		User:         user,
		WorkingDir:   home,
		Cmd:          append([]string{cmd}, args...),
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
	})
	if err != nil {
		return "", 0, err
	}

	resp, err := t.client.ContainerExecAttach(context.Background(), idResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", 0, err
	}

	stdout := new(bytes.Buffer)
	if _, err := stdcopy.StdCopy(stdout, stdout, resp.Reader); err != nil && !errors.Is(err, io.EOF) {
		return "", 0, fmt.Errorf("copy failed: %w", err)
	}

	fmt.Println("Waiting for command...")
	inspect, err := t.client.ContainerExecInspect(context.Background(), idResp.ID)
	if err != nil {
		return "", 0, err
	}
	if inspect.Running {
		return "", 0, fmt.Errorf("command is still running")
	}

	if inspect.ExitCode > 0 {
		return stdout.String(), inspect.ExitCode, nil
	}
	return stdout.String(), inspect.ExitCode, nil
}

func (t *testContainer) Stop() error {
	timeout := time.Second * 3
	return t.client.ContainerStop(context.Background(), t.id, &timeout)
}

//nolint
func startContainer(image string) (*testContainer, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "peridot-test")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Created temp dir at %s for %s container.\n", tmpDir, image)

	if err := os.Chmod(tmpDir, 0777); err != nil {
		return nil, err
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	var found bool
	for _, existing := range images {
		for _, tag := range existing.RepoTags {
			if tag == image || tag == fmt.Sprintf("%s:latest", image) {
				fmt.Printf("Found existing image for %s - no need to pull.\n", tag)
				found = true
				break
			}
		}
	}
	if !found {
		fmt.Printf("Pulling image '%s'...\n", image)
		reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}
		_, _ = io.Copy(io.Discard, reader)
		time.Sleep(time.Second)
	}

	containerName := fmt.Sprintf("peridot_%s", image)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, existing := range containers {
		for _, name := range existing.Names {
			if len(name) > 0 && name[1:] == containerName {
				fmt.Printf("Removing old %s container: %s...\n", image, existing.ID)
				if err := cli.ContainerRemove(context.Background(), existing.ID, types.ContainerRemoveOptions{
					RemoveVolumes: true,
					Force:         true,
				}); err != nil {
					return nil, err
				}
			}
		}
	}

	fmt.Printf("Creating container for '%s'...\n", image)
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
			Tty:   true,
			Cmd:   strslice.StrSlice([]string{"sh"}),
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: tmpDir,
					Target: fmt.Sprintf("/home/%s", defaultUser),
				},
			},
		}, nil, nil, containerName)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Starting '%s' container...\n", image)
	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	created := &testContainer{
		id:      cont.ID,
		hostDir: tmpDir,
		client:  cli,
	}

	testContainers = append(testContainers, created)

	if err := created.AddPeridot(); err != nil {
		return nil, err
	}

	_, exit, err := created.RunAsRoot("cp", fmt.Sprintf("/home/%s/peridot", defaultUser), "/usr/bin/")
	if err != nil {
		return nil, err
	}
	if exit > 0 {
		return nil, fmt.Errorf("failed to install peridot in container")
	}

	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	_, exit, err = created.RunAsRoot("useradd", "-u", u.Uid, defaultUser)
	if err != nil {
		return nil, err
	}
	if exit > 0 {
		return nil, fmt.Errorf("failed to add user to container")
	}
	_, exit, err = created.RunAsRoot("sh", "-c", fmt.Sprintf(`apt update && apt install -y sudo && usermod -a -G sudo %s`, defaultUser))
	if err != nil {
		return nil, err
	}
	if exit > 0 {
		return nil, fmt.Errorf("failed to add user to sudoers")
	}

	return created, nil
}
