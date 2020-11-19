package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/jhoonb/archivex"
)

//SSHKeys contains the pointers to the ssh keys
type SSHKeys struct {
	PublicKey  *string
	PrivateKey *string
}

func getDockerAuthenticationString() string {
	authConfig := types.AuthConfig{
		Username: os.Getenv("DOCKERUSERNAME"),
		Password: os.Getenv("DOCKERPASSWORD"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return authStr
}

func getKeys() SSHKeys {
	publicKey := os.Getenv("PUBLICKEY")
	// privateKey := os.Getenv("PRIVATEKEY")

	cmd := exec.Command("whoami")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	username := string(stdout)
	username = strings.TrimSuffix(username, "\n")

	cmd1 := exec.Command("cat", "/home/"+username+"/.ssh/id_rsa")
	stdout1, err := cmd1.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	privateKey := string(stdout1)

	keys := SSHKeys{
		PublicKey:  &publicKey,
		PrivateKey: &privateKey,
	}

	return keys

}

func pullDockerImage(ctx context.Context, cli *client.Client, imageName string) {
	authStr := getDockerAuthenticationString()

	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
}

func dockerBuild(cli *client.Client, imageName string) {
	tar := new(archivex.TarFile)
	tar.Create("/tmp/rubyResidentialControllerGrading.tar")
	tar.AddAll("rubyResidentialControllerGrading", true)
	tar.Close()

	dockerBuildContext, err := os.Open("/tmp/rubyResidentialControllerGrading.tar")
	defer dockerBuildContext.Close()

	keys := getKeys()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "rubyResidentialControllerGrading/Dockerfile",
		Tags:       []string{imageName + ":local"},
		BuildArgs: map[string]*string{
			"ssh_pub_key": keys.PublicKey,
			"ssh_prv_key": keys.PrivateKey,
		},
	}

	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
	}

	//time.Sleep(time.Second * 3)

	defer buildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, buildResponse.Body)

	deleteFile("/tmp/rubyResidentialControllerGrading.tar")
}

func docker(imageName string, githubHandle string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	pullDockerImage(ctx, cli, imageName)

	dockerBuild(cli, imageName)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName + ":local",
		// Image: imageName + ":test",
		Cmd: []string{"/usr/bin/correction-script.sh"},
		Tty: false,
		Env: []string{"GITHUBHANDLE=" + githubHandle},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}
