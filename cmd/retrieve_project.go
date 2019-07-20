// Copyright © 2018 Jasmin Gacic <jasmin@stackpointcloud.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/packethost/packngo"
	"github.com/spf13/cobra"
)

// retriveProjectCmd represents the retriveProject command
var retriveProjectCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves all available projects or a single project",
	Long: `Example:

Retrieve all projects:
packet project get
  
Retrieve a specific project:
packet project get -i [project_UUID]
packet project get -n [project_name]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if projectID != "" && projectName != "" {
			fmt.Println("must specify only one of project-id and project name")
			return
		}
		if projectID == "" {
			listOpt := &packngo.ListOptions{
				Includes: "members",
			}

			projects, _, err := PacknGo.Projects.List(listOpt)
			if err != nil {
				fmt.Println("Client error:", err)
				return
			}

			var data [][]string
			if projectName == "" {
				data = make([][]string, len(projects))
				for i, p := range projects {
					data[i] = []string{p.ID, p.Name, p.Created}
				}
			} else {
				data = make([][]string, 0)
				for _, p := range projects {
					if p.Name == projectName {
						data = append(data, []string{p.ID, p.Name, p.Created})
						break
					}
				}
				if len(data) == 0 {
					fmt.Println("Error: no project found with name", projectName)
					return
				}
			}

			header := []string{"ID", "Name", "Created"}
			output(projects, header, &data)
		} else {
			p, _, err := PacknGo.Projects.GetExtra(projectID, []string{"members"}, []string{})
			if err != nil {
				fmt.Println("Client error:", err)
				return
			}

			data := make([][]string, 1)

			data[0] = []string{p.ID, p.Name, p.Created}
			header := []string{"ID", "Name", "Created"}
			output(p, header, &data)
		}
	},
}

func init() {
	retriveProjectCmd.Flags().StringVarP(&projectName, "project", "n", "", "Name of the project")
	retriveProjectCmd.Flags().StringVarP(&projectID, "project-id", "i", "", "UUID of the project")
	retriveProjectCmd.PersistentFlags().BoolVarP(&isJSON, "json", "j", false, "JSON output")
	retriveProjectCmd.PersistentFlags().BoolVarP(&isYaml, "yaml", "y", false, "YAML output")
}
