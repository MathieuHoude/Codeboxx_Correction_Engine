package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
)

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

func docker(gradingRequest GradingRequest) []DeliverableScore {
	updateJobStatus(gradingRequest.JobID, "Building")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	pullDockerImage(ctx, cli, gradingRequest.DockerImageName)
	dockerBuild(cli, gradingRequest.DockerImageName)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: gradingRequest.DockerImageName + ":local",
		// Image: imageName + ":test",
		Cmd: []string{"/usr/bin/correction-script.sh"},
		Tty: false,
		Env: []string{"GITHUBHANDLE=" + gradingRequest.GithubHandle},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	updateJobStatus(gradingRequest.JobID, "Grading")
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

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	//read the first 8 bytes to ignore the HEADER part from docker container logs
	p := make([]byte, 8)
	out.Read(p)
	content, _ := ioutil.ReadAll(out)

	deliverableScores := buildDeliverableScores(content, gradingRequest)

	return deliverableScores
}
