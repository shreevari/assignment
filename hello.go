package main

import (
	"fmt"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gorilla/mux"	
	"log"
	"net/http"
	"time"
)

type AWSInstance struct {
	Id string 
	Name string
	PowerState string
}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	
	query := r.URL.Query()
	region, ok := query["region"]
	if !ok {
		w.WriteHeader(400)
		return
	}

	instances, err := describeInstances(region[0], nil)

	if err != nil {
		w.WriteHeader(500)
		return
	}
	
	responseJson, err := json.Marshal(map[string]interface{}{
		"instances": instances,
	})

	if err != nil {
		w.WriteHeader(500)
		return
	}
	fmt.Fprintf(w, string(responseJson))
}

func describeInstances(region string, instanceIds []string) ([]AWSInstance, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	svc := ec2.New(sess)

	var awsInstanceIds []*string
	for _, instanceId := range(instanceIds) {
		awsInstanceIds = append(awsInstanceIds, aws.String(instanceId))
	}
	
	input := &ec2.DescribeInstancesInput{
		InstanceIds: awsInstanceIds,
	}

	descriptions, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	
	var instances []AWSInstance
	for _, reservation := range(descriptions.Reservations) {
		for _, instance := range(reservation.Instances) {
			instances = append(instances, AWSInstance{
				Id: *instance.InstanceId,
				Name: *instance.KeyName,
				PowerState: *instance.State.Name,
			})
		}
	}
	return instances, nil
}

func startInstancesHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	var reqBody map[string][]string
	
	json.NewDecoder(r.Body).Decode(&reqBody)
	query := r.URL.Query()
	region, ok := query["region"]
	if !ok {
		w.WriteHeader(400)
		return
	}
	startingInstanceIds, err := startInstances(region[0], reqBody["instanceIds"])
	if err != nil {
		w.WriteHeader(500)
		return
	}

	const TIMEOUT = 10
	success := false
	for i := 0; i < TIMEOUT && !success; i++ {
		time.Sleep(time.Second)
		instances, err := describeInstances(region[0], startingInstanceIds)
		if err != nil {
			break
		}
		
		for _, instance := range(instances) {
			if instance.PowerState != "running" {
				success = false
				break
			}
			success = true
		}
	}
	
	if success {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Success")
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failure")
	}
	
}

func startInstances(region string, instanceIds []string) ([]string, error) {


	var awsInstanceIds []*string
	
	for _, instanceId := range(instanceIds) {
		awsInstanceIds = append(awsInstanceIds, aws.String(instanceId))
	}


	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	svc := ec2.New(sess)

	input := &ec2.StartInstancesInput{
		InstanceIds: awsInstanceIds,
	}

	result, err := svc.StartInstances(input)
	if err != nil {
		return nil, err
	}

	var startingInstanceIds []string
	for _, startingInstance := range(result.StartingInstances) {
		startingInstanceIds = append(startingInstanceIds, *startingInstance.InstanceId)
	}
	return startingInstanceIds, nil
}


func stopInstancesHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	var reqBody map[string][]string
	
	json.NewDecoder(r.Body).Decode(&reqBody)
	query := r.URL.Query()
	region, ok := query["region"]
	if !ok {
		w.WriteHeader(400)
		return
	}
	stoppingInstanceIds, err := stopInstances(region[0], reqBody["instanceIds"])
	if err != nil {
		w.WriteHeader(500)
		return
	}

	const TIMEOUT = 10
	success := false
	for i := 0; i < TIMEOUT && !success; i++ {
		time.Sleep(time.Second)
		instances, err := describeInstances(region[0], stoppingInstanceIds)
		if err != nil {
			break
		}
		
		for _, instance := range(instances) {
			if instance.PowerState != "stopped" {
				success = false
				break
			}
			success = true
		}
	}
	
	if success {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Success")
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failure")
	}
	
}


func stopInstances(region string, instanceIds []string) ([]string, error) {
	var awsInstanceIds []*string
	
	for _, instanceId := range(instanceIds) {
		awsInstanceIds = append(awsInstanceIds, aws.String(instanceId))
	}


	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	svc := ec2.New(sess)

	input := &ec2.StopInstancesInput{
		InstanceIds: awsInstanceIds,
	}

	result, err := svc.StopInstances(input)
	if err != nil {
		return nil, err
	}

	var stoppingInstanceIds []string
	for _, stoppingInstance := range(result.StoppingInstances) {
		stoppingInstanceIds = append(stoppingInstanceIds, *stoppingInstance.InstanceId)
	}
	return stoppingInstanceIds, nil
}


func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Vary", "Origin")
	(*w).Header().Set("Vary", "Access-Control-Request-Method")
	(*w).Header().Set("Vary", "Access-Control-Request-Headers")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
}

func main() {
	fmt.Println("Application started")
	route := mux.NewRouter().StrictSlash(true)
	route.HandleFunc("/instance", dataHandler)
	route.HandleFunc("/instances/start", startInstancesHandler).Methods("POST", "OPTIONS")
	route.HandleFunc("/instances/stop", stopInstancesHandler).Methods("POST", "OPTIONS")
//	route.HandleFunc("/create-instance", createInstanceHandler).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8000", route))
}
