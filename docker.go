package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

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

func pullDockerImage(ctx context.Context, cli *client.Client, imageName string) {
	authStr := getDockerAuthenticationString()

	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
}

func getKeys() SSHKeys {

	var publicKeyPointer *string
	var privateKeyPointer *string
	publicKey := os.Getenv("PUBLICKEY")
	privateKey := os.Getenv("PRIVATEKEY")
	publicKeyPointer = &publicKey
	privateKeyPointer = &privateKey

	keys := SSHKeys{
		PublicKey:  publicKeyPointer,
		PrivateKey: privateKeyPointer,
	}

	return keys

}

func docker(githubHandle string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// imageName := "mathieuhoude/ruby-residential-controller-grading:1.1"
	// pullDockerImage(ctx, cli, imageName)
	imageName := "ruby-residential-controller-grading"

	tar := new(archivex.TarFile)
	tar.Create("/tmp/rubyResidentialControllerGrading.tar")
	tar.AddAll("contrib", true)
	tar.AddAll("src", true)
	tar.AddAll("test", true)
	tar.AddAll("s2i", true)
	tar.AddAll("help", true)
	tar.AddAll("image", true)
	tar.AddAll("licenses", true)
	tar.Close()

	dockerBuildContext, err := os.Open("./rubyResidentialControllerGrading/")
	defer dockerBuildContext.Close()

	keys := getKeys()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		BuildArgs: map[string]*string{
			"PUBLICKEY":  keys.PublicKey,
			"PRIVATEKEY": keys.PrivateKey,
		},
	}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer buildResponse.Body.Close()

	// test := []string{"PUBLICKEY=" + os.Getenv("PUBLICKEY"), "PRIVATEKEY=" + os.Getenv("PRIVATEKEY"), "GITHUBHANDLE" + githubHandle}
	// test2 := strings.Split(test, "\n")

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"/usr/bin/correction-script.sh"},
		Tty:   false,
		Env:   []string{"GITHUBHANDLE" + githubHandle},
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
