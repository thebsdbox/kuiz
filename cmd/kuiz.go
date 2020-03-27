package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/kuiz/pkg/client"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
	"github.com/thebsdbox/kuiz/pkg/server"
	"github.com/thebsdbox/kuiz/pkg/ui"
)

var kuizCmd = &cobra.Command{
	Use:   "kuiz",
	Short: "",
}

// Release - this struct contains the release information populated when building kube-vip
var Release struct {
	Version string
	Build   string
}

// Configure the level of logging
var logLevel uint32

// The address of the server to bind to
var address string

// The username of the client
var user string

// The path to a quiz file
var path string

func init() {

	// Manage logging
	kuizCmd.PersistentFlags().Uint32Var(&logLevel, "log", 4, "Set the level of logging")
	kuizCmd.PersistentFlags().StringVar(&address, "address", "localhost:8080", "The address and port")

	kuizClient.Flags().StringVarP(&user, "username", "u", "", "username")
	kuizServer.Flags().StringVarP(&path, "quiz", "q", "", "path to a quiz")

	kuizCmd.AddCommand(kuizServer)
	kuizCmd.AddCommand(kuizClient)

}

// Execute - starts the command parsing process
func Execute() {
	if err := kuizCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// q := types.Quiz{
// 		Name: "Todays quiz is about Tea (test question)",
// 		Questions: []types.Question{
// 			types.Question{
// 				Type: types.QUIZMULTI,
// 				MultiChoiceQuestion: &types.MultipleChoiceQuestion{
// 					QuestionText: "Who makes the best \"English Breakfast\" tea?",
// 					Answers: []string{
// 						"Yorkshire Tea",
// 						"PG Tips",
// 						"Typhoo",
// 					},
// 					Answer: "Yorkshire Tea",
// 				},
// 				TimeToReveal: time.Second * 5,
// 			},
// 			types.Question{
// 				Type: types.QUIZMULTI,
// 				MultiChoiceQuestion: &types.MultipleChoiceQuestion{
// 					QuestionText: "Kubeadm does?",
// 					Answers: []string{
// 						"Installs a cluster on bare-metal",
// 						"Installs a cluster in a SaaS",
// 						"Makes tea",
// 						"Bootstraps a cluster",
// 					},
// 					Answer: "Bootstraps a cluster",
// 				},
// 				TimeToReveal: time.Second * 5,
// 			},
// 			types.Question{
// 				Type: types.QUIZYESNO,
// 				YesNoChoiceQuestion: &types.YesNoChoiceQuestion{
// 					QuestionText: "Is \"English Breakfast\" tea the best?",
// 					Answers: []string{
// 						"Yes",
// 						"No",
// 					},
// 					Answer: "Yes",
// 				},
// 				TimeToReveal: time.Second * 5,
// 			},
// 		},
// 	}

var kuizServer = &cobra.Command{
	Use:   "server",
	Short: "Start the kuiz server",
	Run: func(cmd *cobra.Command, args []string) {
		var quiz types.Quiz
		yamlFile, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("Error reading [%s] Err: [%s] ", path, err)
		}
		err = yaml.Unmarshal(yamlFile, &quiz)
		if err != nil {
			log.Printf("Error parsing [%s] Err: [%s] ", path, err)
		}

		// q := types.Quiz{
		// 	Name: "Todays quiz is about Tea (test question)",
		// 	Questions: []types.Question{
		// 		types.Question{
		// 			Type: types.QUIZMULTI,
		// 			MultiChoiceQuestion: &types.MultipleChoiceQuestion{
		// 				QuestionText: "Who makes the best \"English Breakfast\" tea?",
		// 				Answers: []string{
		// 					"Yorkshire Tea",
		// 					"PG Tips",
		// 					"Typhoo",
		// 				},
		// 				Answer: "Yorkshire Tea",
		// 			},
		// 			TimeToReveal: time.Second * 5,
		// 		},
		// 		types.Question{
		// 			Type: types.QUIZMULTI,
		// 			MultiChoiceQuestion: &types.MultipleChoiceQuestion{
		// 				QuestionText: "Kubeadm does?",
		// 				Answers: []string{
		// 					"Installs a cluster on bare-metal",
		// 					"Installs a cluster in a SaaS",
		// 					"Makes tea",
		// 					"Bootstraps a cluster",
		// 				},
		// 				Answer: "Bootstraps a cluster",
		// 			},
		// 			TimeToReveal: time.Second * 5,
		// 		},
		// 		types.Question{
		// 			Type: types.QUIZYESNO,
		// 			YesNoChoiceQuestion: &types.YesNoChoiceQuestion{
		// 				QuestionText: "Is \"English Breakfast\" tea the best?",
		// 				Answers: []string{
		// 					"Yes",
		// 					"No",
		// 				},
		// 				Answer: "Yes",
		// 			},
		// 			TimeToReveal: time.Second * 5,
		// 		},
		// 	},
		// }
		// a, _ := yaml.Marshal(q)
		// fmt.Printf(string(a))
		// return
		v, err := ui.BuildServerView()
		if err != nil {
			panic(err)
		}
		m := server.NewServer(address, &quiz, v)

		go func() {
			err = m.StartServer()
			if err != nil {
				ui.ErrorMessage(err)
				m.StopServer()
			}
		}()
		if err := v.App.Run(); err != nil {
			panic(err)
		}
		m.StopServer()
	},
}

var kuizClient = &cobra.Command{
	Use:   "client",
	Short: "Connect to a kuiz server",
	Run: func(cmd *cobra.Command, args []string) {

		v, err := ui.BuildClientView()
		if err != nil {
			panic(err)
		}

		m := client.NewClient(user, address, v)
		go func() {
			err = m.StartClient()
			if err != nil {
				m.StopClient()
				return
			}
		}()
		if err := v.App.Run(); err != nil {
			panic(err)
		}
		m.StopClient()

	},
}
