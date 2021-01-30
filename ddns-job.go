package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	zoneDomain, domain := obtainInputs()
	log.Printf("Zone Name: %s", *zoneDomain)
	log.Printf("Domain: %s", *domain)

	ip, ipErr := obtainExternalIP()
	failOnError(ipErr)
	log.Printf("Current IP: %s", *ip)

	awsConfig, awsConfigError := config.LoadDefaultConfig(context.TODO())
	failOnAWSError(awsConfigError)

	client := route53.NewFromConfig(awsConfig)
	ctx := context.Background()

	zone, findErr := findZone(ctx, client, zoneDomain)
	failOnError(findErr)
	log.Printf("Hosted Zone ID: %s", *zone.Id)

	upsertErr := updateRecord(ctx, client, zone, domain, ip)
	failOnAWSError(upsertErr)
	log.Printf("Fin")
}

func obtainInputs() (*string, *string) {
	var zoneDomain string
	var domain string

	flag.StringVar(&zoneDomain, "z", "", "The zone name for the hosted zone containing the domain record to be updated.")
	flag.StringVar(&domain, "d", "", "The domain name to have the A record updated to the current external IP.")
	flag.Parse()

	if len(zoneDomain) == 0 || len(domain) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	return &zoneDomain, &domain
}

func obtainExternalIP() (*string, error) {
	var blank = ""
	response, ipErr := http.Get("https://api.ipify.org?format=text")
	if ipErr != nil {
		return &blank, ipErr
	}
	ipRaw, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return &blank, readErr
	}
	closeErr := response.Body.Close()
	if closeErr != nil {
		return &blank, closeErr
	}
	var ip = string(ipRaw)
	return &ip, nil
}

func findZone(ctx context.Context, client *route53.Client, zoneName *string) (*types.HostedZone, error) {

	listHostedZonesInput := &route53.ListHostedZonesInput{}
	zoneList, zoneListError := client.ListHostedZones(ctx, listHostedZonesInput)

	failOnAWSError(zoneListError)

	for _, zone := range zoneList.HostedZones {
		if *zone.Name == *zoneName {
			return &zone, nil
		}
	}

	return &types.HostedZone{}, fmt.Errorf("unable to find zone: %s", *zoneName)
}

func updateRecord(ctx context.Context, client *route53.Client, zone *types.HostedZone, domain *string, ip *string) error {

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeAction("UPSERT"),
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: domain,
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: ip,
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
			Comment: aws.String("Automated update from DDNS Job"),
		},
		HostedZoneId: zone.Id,
	}
	_, err := client.ChangeResourceRecordSets(ctx, params)
	return err
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func failOnAWSError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
