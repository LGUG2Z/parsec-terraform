package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type TfVars struct {
	InstanceType   string `json:"instance_type"`
	Region         string `json:"region"`
	ServerKey      string `json:"server_key"`
	SpotInstanceID string `json:"spot_instance_id"`
	SpotPrice      string `json:"spot_price"`
	SubnetID       string `json:"subnet_id"`
	VpcID          string `json:"vpc_id"`
}

type TfOutputs struct {
	InstanceType struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"instance_type"`
	Region struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"region"`
	ServerKey struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"server_key"`
	SpotInstanceID struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"spot_instance_id"`
	SpotPrice struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"spot_price"`
	SubnetID struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"subnet_id"`
	VpcID struct {
		Sensitive bool   `json:"sensitive"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"vpc_id"`
}

// TODO: Marshal this out as a toml file and write to terraform.tfvars
func (v *TfVars) Write() error {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", installPath, CurrentSession)

	return ioutil.WriteFile(filePath, bytes, 0644)
}

func (v *TfVars) Calculate(ec2Client *ec2.EC2, region, serverKey, instanceType string) error {
	vpcID, err := getVpcID(ec2Client)
	if err != nil {
		return err
	}

	spotPrice, err := getSpotPrice(ec2Client, instanceType)
	if err != nil {
		return err
	}

	spotBid := calculateUserBid(*spotPrice.SpotPrice, bid)
	availabilityZone := *spotPrice.AvailabilityZone

	subnetID, err := getSubnetID(ec2Client, availabilityZone)
	if err != nil {
		return err
	}

	v.InstanceType = instanceType
	v.Region = region
	v.ServerKey = serverKey
	v.SpotPrice = spotBid
	v.SubnetID = subnetID
	v.VpcID = vpcID

	return nil
}

func (v *TfOutputs) Read() error {
	o := tfCmd([]string{TfCmdOutput, TfFlagJSON})
	output, err := executeReturn(o)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(output, &v); err != nil {
		return err
	}

	return nil
}