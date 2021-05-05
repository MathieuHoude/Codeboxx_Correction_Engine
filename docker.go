package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

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
	var curr, _ = os.Getwd()
	err := os.Chdir(filepath.Join(curr, imageName))
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("/tmp/%s_build.tar", time.Now().Format("20060102150405999999"))
	tar := new(archivex.TarFile)
	tar.Create(fileName)
	tar.AddAll(".", true)
	tar.Close()

	dockerBuildContext, _ := os.Open(fileName)
	defer dockerBuildContext.Close()

	keys := getKeys()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
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

	defer buildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, buildResponse.Body)

	err = os.Chdir(curr)
	if err != nil {
		log.Fatal(err)
	}

	deleteFile(fileName)
}

func docker(correctionRequest CorrectionRequest) {
	updateJobStatus(correctionRequest.JobID, "Building")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	pullDockerImage(ctx, cli, correctionRequest.DockerImageName)
	dockerBuild(cli, correctionRequest.DockerImageName)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: correctionRequest.DockerImageName + ":local",
		// Image: imageName + ":test",
		Cmd: []string{"/usr/bin/correction-script.sh"},
		Tty: false,
		Env: []string{
			"JOBID=" + fmt.Sprint(correctionRequest.JobID),
			"DELIVERABLEID=" + fmt.Sprint(correctionRequest.DeliverableID),
			"UNIXDELIVERABLEDEADLINE=" + fmt.Sprint(correctionRequest.UnixDeliverableDeadline),
			"REPOSITORYURL=" + correctionRequest.RepositoryURL,
		},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	updateJobStatus(correctionRequest.JobID, "Correcting")
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
	content, err := ioutil.ReadAll(out)

	if err != nil {
		log.Println("Error in ReadALL", err)
	}

	x := string(content)
	println(x)
}
