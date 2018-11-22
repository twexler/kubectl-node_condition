// Copyright © 2018 Ted Wexler
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
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	nodeConditionExample = `
	# output current node condition data for each node
	%[1]s node-condition
	# output a simplified json object with condition data for each node
	%[1]s node-condition --output json
`
)

// NodeConditionOptions represents configuration for the command
type NodeConditionOptions struct {
	configFlags *genericclioptions.ConfigFlags
	rawConfig   *rest.Config

	genericclioptions.IOStreams
}

// NewNodeConditionOptions provides an instance of NodeConditionOptions with default values
func NewNodeConditionOptions(streams genericclioptions.IOStreams) *NodeConditionOptions {
	return &NodeConditionOptions{
		configFlags: genericclioptions.NewConfigFlags(),
		IOStreams:   streams,
	}
}

// NewCmdNodeCondition provides a wrapper around NodeConditionOptions
func NewCmdNodeCondition(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewNodeConditionOptions(streams)
	cmd := &cobra.Command{
		Use:   "node-condition",
		Short: "A kubectl plugin to output node conditions",
		Long:  fmt.Sprintf(nodeConditionExample, "kubectl"),
		RunE:  o.Execute,
	}
	cmd.Flags().StringP("output", "o", "cli", `format to output node condition data (one of "cli" or "json")`)
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Execute outputs node condition data to the provided streams
func (o *NodeConditionOptions) Execute(cmd *cobra.Command, args []string) error {
	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(o.rawConfig)
	if err != nil {
		return err
	}
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	if output == "cli" {
		o.cliOutput(nodes)
	} else if output == "json" {
		return o.jsonOutput(nodes)
	} else {
		return fmt.Errorf("unknown output type %s", output)
	}
	return nil
}

func (o *NodeConditionOptions) jsonOutput(nodes *corev1.NodeList) error {
	data := make(map[string]map[string]map[string]string, len(nodes.Items))
	for _, node := range nodes.Items {
		data[node.Name] = make(map[string]map[string]string, len(node.Status.Conditions))
		for _, condition := range node.Status.Conditions {
			data[node.Name][condition.Reason] = map[string]string{
				"status":             string(condition.Status),
				"message":            condition.Message,
				"lastTransitionTime": condition.LastTransitionTime.String(),
			}
		}
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprint(o.Out, string(jsonData))
	return nil
}

func (o *NodeConditionOptions) cliOutput(nodes *corev1.NodeList) {
	for _, node := range nodes.Items {
		fmt.Fprintln(o.Out, node.Name)
		i := 0
		for i < len(node.Name) {
			fmt.Fprint(o.Out, "=")
			i++
		}
		fmt.Fprintf(o.Out, "\n\n")
		table := tablewriter.NewWriter(o.Out)
		table.SetHeader([]string{"Reason", "Status", "Message", "Last transition time"})
		table.SetRowLine(true)

		for _, condition := range node.Status.Conditions {
			table.Append([]string{condition.Reason, string(condition.Status), condition.Message, condition.LastTransitionTime.String()})
		}
		table.Render()
		fmt.Fprintf(o.Out, "\n\n")
	}
}
