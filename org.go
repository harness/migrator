package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strconv"
)

func createOrg(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	if len(migrationReq.OrgName) == 0 {
		promptConfirm = true
		migrationReq.OrgName = TextInput("Name of the Org - ")
	}
	if len(migrationReq.OrgIdentifier) == 0 {
		promptConfirm = true
		migrationReq.OrgIdentifier = TextInput("Identifier for the Org - ")
	}

	log.WithFields(log.Fields{
		"Account":       migrationReq.Account,
		"OrgIdentifier": migrationReq.OrgIdentifier,
		"OrgName":       migrationReq.OrgName,
	}).Info("Org creation details")

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with org creation?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := fmt.Sprintf("%s/api/organizations?accountIdentifier=%s", urlMap[migrationReq.Environment][NG], migrationReq.Account)

	log.Info("Creating the org....")

	_, err := Post(url, migrationReq.Auth, OrgBody{
		Org: OrgDetails{
			Identifier:  migrationReq.OrgIdentifier,
			Name:        migrationReq.OrgName,
			Description: "",
		}})

	if err == nil {
		log.Info("Created the org!")
	} else {
		log.Error("There was an error creating the org")
		return err
	}

	return nil
}

func bulkRemoveOrg(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	names := Split(migrationReq.Names, ",")
	identifiers := Split(migrationReq.Identifiers, ",")
	if len(names) == 0 && len(identifiers) == 0 {
		log.Fatal("No names or identifiers for the organisations provided. Aborting")
	}
	if len(names) > 0 && len(identifiers) > 0 {
		log.Fatal("Both names and identifiers for the organisations provided. Aborting")
	}

	n := len(identifiers)
	if len(names) > 0 {
		n = len(names)
	}
	if promptConfirm {
		confirm := ConfirmInput("Are you sure you want to proceed with deletion of " + strconv.Itoa(n) + " organisations?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	if len(names) > 0 {
		organisations := getOrganisations()
		for _, name := range names {
			id := findOrgIdByName(organisations, name)
			if len(id) > 0 {
				identifiers = append(identifiers, id)
			}
		}
		log.Debugf("Valid identifers for the given names are - %s", identifiers)
	}

	for _, identifier := range identifiers {
		deleteOrg(identifier)
	}
	log.Info("Finished operation for all given organisations")
	return nil
}

func deleteOrg(orgId string) {
	url := fmt.Sprintf("%s/api/organizations/%s?accountIdentifier=%s", urlMap[migrationReq.Environment][NG], orgId, migrationReq.Account)

	log.Infof("Deleting the org with identifier %s", orgId)

	_, err := Delete(url, migrationReq.Auth)

	if err == nil {
		log.Infof("Successfully deleted the org - %s", orgId)
	} else {
		log.Errorf("Failed to delete the org - %s", orgId)
	}
}

func getOrganisations() []OrgDetails {
	url := fmt.Sprintf("%s/api/aggregate/organizations?accountIdentifier=%s&pageSize=1000", urlMap[migrationReq.Environment][NG], migrationReq.Account)
	resp, err := Get(url, migrationReq.Auth)
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch organisations", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch organisations", err)
	}
	var orgListBody OrgListBody
	err = json.Unmarshal(byteData, &orgListBody)
	if err != nil {
		log.Fatal("Failed to fetch organisations", err)
	}
	var details []OrgDetails

	for _, o := range orgListBody.Organisations {
		details = append(details, o.Org.Org)
	}
	return details
}

func findOrgIdByName(organisations []OrgDetails, orgName string) string {
	for _, o := range organisations {
		if o.Name == orgName {
			return o.Identifier
		}
	}
	return ""
}
