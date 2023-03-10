package lib

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type Change struct {
	JobName       string
	NewValue      string
	PreviousValue string
}

type Deployment struct {
	artifact    string
	version     string
	environment string
	role        string
	autoapprove bool
	sess        *session.Session
	changeSet   []Change
}

type DeploymentOptions struct {
	// Application name as it exists in ECS
	Artifact string
	// Desired version of ECS Application
	Version string
	// Environment (ECS Cluster Name) you would like to deploy to
	Environment string
	// Role is the IAM role to use when invoking a deployment.
	Role string
	// Automatically approve changes
	AutoApprove bool
	// AWS session
	Sess *session.Session
}

func NewDeployment(opts *DeploymentOptions) *Deployment {
	d := Deployment{
		artifact:    opts.Artifact,
		version:     opts.Version,
		environment: opts.Environment,
		role:        opts.Role,
		autoapprove: opts.AutoApprove,
		sess:        opts.Sess,
	}

	return &d
}

func (d Deployment) String() string {
	var ss []string
	for _, change := range d.changeSet {
		ss = append(ss, fmt.Sprintf("[%s]: %s -> %s", change.JobName, change.PreviousValue, change.NewValue))
	}

	return strings.Join(ss, `\n`)
}

func (d Deployment) HasChanges() bool {
	return len(d.changeSet) > 0
}

func (d *Deployment) Build() error {
	var svc *glue.Glue
	if d.role != "" {
		creds := stscreds.NewCredentials(d.sess, d.role)
		svc = glue.New(d.sess, &aws.Config{Credentials: creds})
	} else {
		svc = glue.New(d.sess)
	}

	var nextToken *string
	for {
		// get jobs names matching environment and artifact ID
		listJobsOutput, err := svc.ListJobs(&glue.ListJobsInput{
			MaxResults: aws.Int64(1000),
			NextToken:  nextToken,
			Tags: map[string]*string{
				"Environment": &d.environment,
				"ArtifactID":  &d.artifact,
			},
		})
		if err != nil {
			return fmt.Errorf("error running Glue:ListJobs: %s", err.Error())
		}

		// get the full job details
		batchGetJobsOutput, err := svc.BatchGetJobs(&glue.BatchGetJobsInput{
			JobNames: listJobsOutput.JobNames,
		})
		if err != nil {
			return fmt.Errorf("error running Glue:BatchGetJobs: %s", err.Error())
		}

		// add the changes to the changeset
		for _, job := range batchGetJobsOutput.Jobs {
			if ok, newValue := isUpdatedVersion(*job.Command.ScriptLocation, d.version); ok {
				d.changeSet = append(d.changeSet, Change{
					JobName:       *job.Name,
					PreviousValue: *job.Command.ScriptLocation,
					NewValue:      newValue,
				})
			}
		}

		// set NextToken for subsequent request
		if listJobsOutput.NextToken == nil {
			break
		}
		nextToken = listJobsOutput.NextToken
	}

	return nil
}

func (d Deployment) Deploy() error {

	var svc *glue.Glue
	if d.role != "" {
		creds := stscreds.NewCredentials(d.sess, d.role)
		svc = glue.New(d.sess, &aws.Config{Credentials: creds})
	} else {
		svc = glue.New(d.sess)
	}

	for _, change := range d.changeSet {
		getJobOutput, err := svc.GetJob(&glue.GetJobInput{
			JobName: &change.JobName,
		})
		if err != nil {
			return fmt.Errorf("error running Glue:GetJob: %s", err.Error())
		}

		jobUpdate := toGlueJob(getJobOutput.Job)
		jobUpdate.Command.ScriptLocation = &change.NewValue

		_, err = svc.UpdateJob(&glue.UpdateJobInput{
			JobName:   &change.JobName,
			JobUpdate: jobUpdate,
		})
		if err != nil {
			return fmt.Errorf("error running Glue:UpdateJob: %s", err.Error())
		}
	}

	return nil
}

// toGlueJob converts a glue "Job" to a glue "JobUpdate"
func toGlueJob(glueJob *glue.Job) *glue.JobUpdate {
	// jobUpdate = glueJob.(glue.JobUpdate)
	jsonGlueJob, err := json.Marshal(glueJob)
	if err != nil {
		fmt.Printf("Could not marshal glue job: %s\n", err)
		// return fmt.Errorf("Could not json marshal glue jobs: %s", err)
	}
	// TODO: handle errors cause we're not trash developers
	var jobUpdate glue.JobUpdate
	err = json.Unmarshal(jsonGlueJob, &jobUpdate)
	if err != nil {
		fmt.Printf("Could not unmarshal glue job: %s\n", err)
	}

	// handle shitty aws devs stuff
	// AllocatedCapacity is depricated. Removing.
	jobUpdate.AllocatedCapacity = nil
	if jobUpdate.WorkerType != nil {
		jobUpdate.MaxCapacity = nil
	}

	return &jobUpdate
}
