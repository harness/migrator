package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const ExpressionPattern = "\\$\\{[\\w-.\"()]+}"

var ExpressionsMap = map[string]string{
	"infra.kubernetes.namespace": "<+infra.namespace>",
	"infra.kubernetes.infraId":   "<+INFRA_KEY>",
	"infra.helm.releaseName":     "<+infra.releaseName>",
	"infra.name":                 "<+infra.name>",

	// Env Expressions
	"env.name":            "<+env.name>",
	"env.description":     "<+env.description>",
	"env.environmentType": "<+env.type>",
	"env.uuid":            "<+env.identifier>",

	// Service Expressions
	"service.name":        "<+service.name>",
	"service.Name":        "<+service.name>",
	"Service.name":        "<+service.name>",
	"service.tag":         "<+service.tags>",
	"service.uuid":        "<+service.identifier>",
	"service.description": "<+service.description>",

	// Artifact Expressions
	"artifact.metadata.image":          "<+artifact.image>",
	"artifact.metadata.tag":            "<+artifact.tag>",
	"artifact.source.dockerconfig":     "<+artifact.imagePullSecret>",
	"artifact.metadata.fileName":       "<+artifact.fileName>",
	"artifact.metadata.format":         "<+artifact.repositoryFormat>",
	"artifact.metadata.getSHA()":       "<+artifact.metadata.SHA>",
	"artifact.metadata.groupId":        "<+artifact.groupId>",
	"artifact.metadata.package":        "<+artifact.metadata.package>",
	"artifact.metadata.region":         "<+artifact.metadata.region>",
	"artifact.metadata.repository":     "<+artifact.repository>",
	"artifact.metadata.repositoryName": "<+artifact.repositoryName>",
	"artifact.metadata.url":            "<+artifact.url>",

	// Rollback Artifact Expressions
	"rollbackArtifact.metadata.image":          "<+rollbackArtifact.image>",
	"rollbackArtifact.metadata.tag":            "<+rollbackArtifact.tag>",
	"rollbackArtifact.source.dockerconfig":     "<+rollbackArtifact.imagePullSecret>",
	"rollbackArtifact.metadata.fileName":       "<+rollbackArtifact.fileName>",
	"rollbackArtifact.metadata.format":         "<+rollbackArtifact.repositoryFormat>",
	"rollbackArtifact.metadata.getSHA()":       "<+rollbackArtifact.metadata.SHA>",
	"rollbackArtifact.metadata.groupId":        "<+rollbackArtifact.groupId>",
	"rollbackArtifact.metadata.package":        "<+rollbackArtifact.metadata.package>",
	"rollbackArtifact.metadata.region":         "<+rollbackArtifact.metadata.region>",
	"rollbackArtifact.metadata.repository":     "<+rollbackArtifact.repository>",
	"rollbackArtifact.metadata.repositoryName": "<+rollbackArtifact.repositoryName>",
	"rollbackArtifact.metadata.url":            "<+rollbackArtifact.url>",

	// Application Expressions
	"app.name":        "<+project.name>",
	"app.description": "<+project.description>",

	// Http Step
	"httpResponseCode": "<+httpResponseCode>",
	"httpResponseBody": "<+httpResponseBody>",
	"httpMethod":       "<+httpMethod>",
	"httpUrl":          "<+httpUrl>",
}

func ReplaceCurrentGenExpressionsWithNextGen(*cli.Context) (err error) {
	extensions := strings.Split(migrationReq.FileExtensions, ",")
	for i, ext := range extensions {
		extensions[i] = "." + ext
	}

	foundExpressionsMap := make(map[string][]string)

	// Fetch all expressions per file
	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && EndsWith(info.Name(), extensions) {
			content, err := ReadFile(path)
			if err != nil {
				return err
			}
			foundExpressions := Set(FindAllExpressions(content))
			if len(foundExpressions) > 0 {
				foundExpressionsMap[path] = foundExpressions
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(foundExpressionsMap) == 0 {
		log.Info("No files found containing Harness expressions!")
		return nil
	}

	// Render a table with summary of expressions found
	data := make(map[string]interface{})
	for path, expList := range foundExpressionsMap {
		data[path] = strings.Join(expList, ", ")
	}
	renderTable("Files containing expressions", data)

	if migrationReq.DryRun {
		log.Info("Dry run is set to true. Skipping expressions replacement for all files")
		return err
	}

	// We are going to do an actual replacement
	for path, expList := range foundExpressionsMap {
		content, err := ReadFile(path)
		if err != nil {
			return err
		}
		str := ReplaceAllExpressions(content, expList)
		err = WriteToFile(path, []byte(str))
		if err != nil {
			return err
		}
		log.Infof("Replaced expressions from %s", path)
	}
	return
}

func FindAllExpressions(str string) []string {
	r := regexp.MustCompile(ExpressionPattern)
	return r.FindAllString(str, -1)
}

func ReplaceAllExpressions(str string, expressions []string) string {
	for _, exp := range expressions {
		temp := exp[2 : len(exp)-1]
		val, ok := ExpressionsMap[temp]
		if ok {
			str = strings.ReplaceAll(str, exp, val)
		}
	}
	return str
}
