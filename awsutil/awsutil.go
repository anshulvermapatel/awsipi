package awsutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//StopCluster Fucntion to stop the cluster, will take instaceIDs as input and the session.
func StopCluster(session *session.Session, instanceIDs []*string) awserr.Error {
	svc := ec2.New(session)
	//Stopping Instances
	input := &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
		DryRun:      aws.Bool(false),
	}
	_, err := svc.StopInstances(input)
	awsErr, _ := err.(awserr.Error)
	return awsErr
}

//StartCluster Fucntion to stop the cluster, will take instaceIDs as input and the session.
func StartCluster(session *session.Session, instanceIDs []*string) awserr.Error {
	svc := ec2.New(session)
	//Stopping Instances
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
		DryRun:      aws.Bool(false),
	}
	_, err := svc.StartInstances(input)
	awsErr, _ := err.(awserr.Error)
	return awsErr
}

func getProviders() string {
	c1 := exec.Command("oc", "get", "nodes", "-oyaml")
	//c2 := exec.Command("grep", "provider")
        c2 := exec.Command("grep", "providerID: aws")
	r, w := io.Pipe()
	var stderr bytes.Buffer
	c1.Stderr = &stderr
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	r.Close()

	defer func() {
		if len(stderr.Bytes()) > 0 {
			fmt.Println("Unable to execute 'oc get node' with below error")
			fmt.Println(stderr.String())
			os.Exit(1)
		}
	}()
	return b2.String()
}

// Returns the providers which contains the zone which would be same for all the nodes
// and an array of string which will be the instance ids.
func GetZoneAndInstanceId() (string, []*string) {
	b := getProviders()
	b = strings.TrimSpace(b)
	providers := strings.Split(b, "\n")
	zone := getZone(providers[0])
	instanceIds := getInstanceIds(providers)

	return zone, instanceIds
}

// Pass only 1 element of the instances array because all element will contain same zone hence 1 can be used to extract the zone.
func getZone(provider string) string {
	zone := strings.Split(provider, "/")[3]
	zone = zone[:len(zone)-1]
	return zone
}

//This is to parse the InstanceIds from the providers and put them in a string array
func getInstanceIds(providers []string) []*string {
	instanceIds := make([]*string, 0)
	for _, provider := range providers {
		arr := strings.Split(provider, "/")
		instanceIds = append(instanceIds, &arr[4])
	}
	return instanceIds
}

//GetZoneAndInstanceIDFromFile This will return the zone and the InstanceID from the ClusterFile hence there is no need for the user to be logged-in in the cluster.
func GetZoneAndInstanceIDFromFile(clusterName string) (string, []*string, error) {
	// In InstanceIDs we will store all the instanceIDs
	InstanceIDs := make([]*string, 0)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	// Reading the clusterFile to get the zone and Intance IDs previously saved by init().
	clusterFile, err := ioutil.ReadFile(homeDir + "/.awsipi/" + clusterName)
	if err != nil {
		fmt.Println(err)
	}

	// Temporary array to store the different line of the content in the Clusterfile
	arr := strings.Split(strings.TrimSpace(string(clusterFile)), "\n")

	// getting InstanceIDs ready to get returns.
	for i := 1; i < len(arr); i++ {
		InstanceIDs = append(InstanceIDs, &arr[i])
	}

	// arr[0] is the zone name
	return arr[0], InstanceIDs, err
}
