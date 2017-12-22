// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"k8s-crd-controller/pkg/controller"
	"log"

	"github.com/spf13/cobra"

	clientver "k8s-crd-controller/pkg/client/clientset/versioned"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
)

var Kubeconfig string

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create controller",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := clientcmd.BuildConfigFromFlags("", Kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
		kubeClientset,err := kubernetes.NewForConfig(config)
		if err!=nil {
			log.Fatal(err)
		}
		clientset, err := clientver.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		//controller
		c := controller.NewController(*clientset, *kubeClientset)

		stop := make(chan struct{})
		defer close(stop)
		go c.Run(1, stop)

		//wait forever
		select {}
	},
}
