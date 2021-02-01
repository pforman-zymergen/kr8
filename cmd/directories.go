package cmd

import (
	//"fmt"
	//"fmt"
	log "github.com/sirupsen/logrus"
	//"github.com/olekukonko/tablewriter"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ComponentDef struct {
	Path string `json:"path"`
}

func (c *Clusters) AddItem(item Cluster) Clusters {
	c.Cluster = append(c.Cluster, item)
	return *c
}

func getClusters(searchDir string) (Clusters, error) {

	fileList := make([]string, 0)
	//clusterName := make([]string, 0)
	//pathName := make([]string, 0)

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	if e != nil {
		log.Fatal("Error building cluster list: ", e)
	}

	ClusterData := []Cluster{}
	c := Clusters{ClusterData}

	for _, file := range fileList {

		splitFile := strings.Split(file, "/")
		// get the filename
		fileName := splitFile[len(splitFile)-1]

		if fileName == "cluster.jsonnet" {
			entry := Cluster{Name: splitFile[len(splitFile)-2], Path: strings.Join(splitFile[:len(splitFile)-1], "/")}
			c.AddItem(entry)

		}
	}

	return c, nil

}

func getCluster(searchDir string, clusterName string) string {

	var clusterPath string

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		if filepath.Base(dir) == clusterName && file == "cluster.jsonnet" {
			clusterPath = path
			return nil
		} else {
			return err
		}
	})

	if e != nil {
		log.Fatal("Error building cluster list: ", e)
	}

	if clusterPath == "" {
		log.Fatal("Could not find cluster: ", clusterName)
	}

	return clusterPath

}

func getClusterParams(basePath string, targetPath string) []string {

	// a slice to store results
	var results []string
	results = append(results, targetPath)

	// remove the cluster.jsonnet
	splitFile := strings.Split(targetPath, "/")

	// gets the targetdir without the cluster.jsonnet
	targetDir := strings.Join(splitFile[:len(splitFile)-1], "/")

	// walk through the directory hierachy
	for {
		rel, _ := filepath.Rel(basePath, targetDir)

		// check if there's a params.json in the folder
		if _, err := os.Stat(targetDir + "/params.jsonnet"); err == nil {
			results = append(results, targetDir+"/params.jsonnet")
		}

		// stop if we're in the basePath
		if rel == "." {
			break
		}

		// next!
		targetDir += "/.."
	}

	// jsonnet's import order matters, so we need to reverse the slice
	last := len(results) - 1
	for i := 0; i < len(results)/2; i++ {
		results[i], results[last-i] = results[last-i], results[i]
	}

	return results

}

func renderClusterParams(cmd *cobra.Command, clusterName string, componentName string, clusterParams string) string {
	if clusterName == "" && clusterParams == "" {
		log.Fatal("Please specify a --cluster name and/or --clusterparams")
	}

	var params []string
	var componentMap map[string]ComponentDef

	if clusterName != "" {
		clusterPath := getCluster(clusterDir, clusterName)
		params = getClusterParams(clusterDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}

	compParams := renderJsonnet(cmd, params, "", false, "")

	compString := gjson.Get(compParams, "_components")
	err := json.Unmarshal([]byte(compString.String()), &componentMap)
	if err != nil {
		log.Panic("failed to parse component map:", err)
	}
	componentDefaultsMerged := "{"
	for key, value := range componentMap {
		if componentName != "" && key != componentName {
			continue
		}
		path := baseDir + "/" + value.Path + "/params.jsonnet"
		filec, err := ioutil.ReadFile(path)
		if err != nil {
			log.Panic("Error reading "+path+" :", err)
		}
		componentDefaultsMerged = componentDefaultsMerged + fmt.Sprintf("'%s': %s,", key, string(filec))
	}
	componentDefaultsMerged = componentDefaultsMerged + "}"

	compParams = renderJsonnet(cmd, params, "", false, componentDefaultsMerged)

	return compParams
}
