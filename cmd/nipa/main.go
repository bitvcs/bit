package main

import (
	"fmt"
	"os"

	"github.com/nipalab/nipa/internal/client/cli"
	"github.com/nipalab/nipa/internal/client/securestorage"
	"github.com/nipalab/nipa/internal/client/usecase"
)

func main() {

	/*var rootCmd = &cobra.Command{
		Use:   "nipa",
		Short: "nipa is centralized version control system for your project",
	}

	// Subcommand (e.g., `mycli clone <url>`)
	var cloneCmd = &cobra.Command{
		Use:   "clone [url]",
		Short: "Clone a repository into a new directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			fmt.Printf("Cloning repository from: %s\n", url)
		},
	}

	// Register subcommand to root
	rootCmd.AddCommand(cloneCmd)

	// Execute CLI
	rootCmd.Execute()

	conn, err := grpc.NewClient("127.0.0.1:6745", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	slog.Info("connected to server")

	client := pb.NewNipaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.LoginWithUsernamePassword(ctx, &pb.LoginUsernamePasswordRequest{
		Username: "supernipa",
		Password: "supernipa",
	})
	if err != nil {
		panic(err)
	}

	slog.Info("login successful", "access_token", res.AccessToken, "refresh_token", res.RefreshToken, "expires_in", res.ExpiresIn)*/

	secureStorage := securestorage.New()
	prompter := cli.NewPrompter()

	authUsecase := usecase.NewAuth(nil, secureStorage, prompter)
	repoUsecase := usecase.NewRepo(authUsecase)
	registry := &Registry{
		authUsecase: authUsecase,
		repoUsecase: repoUsecase,
	}

	cliClient := cli.NewCli(registry)
	if err := cliClient.Run(); err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
