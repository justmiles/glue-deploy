package lib

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

func Test_toGlueJob(t *testing.T) {

	var now = time.Now()
	type args struct {
		glueJob *glue.Job
	}
	tests := []struct {
		name string
		args args
		want *glue.JobUpdate
	}{
		{
			name: "happy path one",
			args: args{
				glueJob: &glue.Job{
					Name:           aws.String("my job name"),
					Role:           aws.String("arn:aws:iam::123456789000:role/glue"),
					CreatedOn:      &now,
					LastModifiedOn: &now,
					ExecutionProperty: &glue.ExecutionProperty{
						MaxConcurrentRuns: aws.Int64(1),
					},
					Command: &glue.JobCommand{
						Name:           aws.String("glueetl"),
						ScriptLocation: aws.String("s3://mys3bucket/glue/scripts/util/v1.36.0/somejob.py"),
						PythonVersion:  aws.String("3"),
					},
					DefaultArguments: map[string]*string{
						"Class":     aws.String("GlueApp"),
						"--TempDir": aws.String("s3://aws-glue-assets-123456789000-us-east-1/temporary/"),
					},
					MaxRetries:        aws.Int64(0),
					AllocatedCapacity: aws.Int64(2),
					Timeout:           aws.Int64(2880),
					MaxCapacity:       aws.Float64(2.0),
					WorkerType:        aws.String("G.1X"),
					NumberOfWorkers:   aws.Int64(2),
					GlueVersion:       aws.String("4.0"),
				},
			},
			want: &glue.JobUpdate{
				Role: aws.String("arn:aws:iam::123456789000:role/glue"),
				ExecutionProperty: &glue.ExecutionProperty{
					MaxConcurrentRuns: aws.Int64(1),
				},
				Command: &glue.JobCommand{
					Name:           aws.String("glueetl"),
					ScriptLocation: aws.String("s3://mys3bucket/glue/scripts/util/v1.36.0/somejob.py"),
					PythonVersion:  aws.String("3"),
				},
				DefaultArguments: map[string]*string{
					"Class":     aws.String("GlueApp"),
					"--TempDir": aws.String("s3://aws-glue-assets-123456789000-us-east-1/temporary/"),
				},
				MaxRetries:        aws.Int64(0),
				AllocatedCapacity: nil,
				Timeout:           aws.Int64(2880),
				MaxCapacity:       nil,
				WorkerType:        aws.String("G.1X"),
				NumberOfWorkers:   aws.Int64(2),
				GlueVersion:       aws.String("4.0"),
			},
		},

		{
			name: "happy path two",
			args: args{
				glueJob: &glue.Job{
					Name:           aws.String("my job name"),
					Role:           aws.String("arn:aws:iam::123456789000:role/glue"),
					CreatedOn:      &now,
					LastModifiedOn: &now,
					ExecutionProperty: &glue.ExecutionProperty{
						MaxConcurrentRuns: aws.Int64(1),
					},
					Command: &glue.JobCommand{
						Name:           aws.String("glueetl"),
						ScriptLocation: aws.String("s3://mys3bucket/glue/scripts/util/v1.36.0/somejob.py"),
						PythonVersion:  aws.String("3"),
					},
					DefaultArguments: map[string]*string{
						"Class":     aws.String("GlueApp"),
						"--TempDir": aws.String("s3://aws-glue-assets-123456789000-us-east-1/temporary/"),
					},
					MaxRetries:        aws.Int64(0),
					AllocatedCapacity: aws.Int64(2),
					Timeout:           aws.Int64(2880),
					MaxCapacity:       aws.Float64(2.0),
					WorkerType:        aws.String("G.1X"),
					NumberOfWorkers:   aws.Int64(2),
					GlueVersion:       aws.String("4.0"),
				},
			},
			want: &glue.JobUpdate{
				Role: aws.String("arn:aws:iam::123456789000:role/glue"),
				ExecutionProperty: &glue.ExecutionProperty{
					MaxConcurrentRuns: aws.Int64(1),
				},
				Command: &glue.JobCommand{
					Name:           aws.String("glueetl"),
					ScriptLocation: aws.String("s3://mys3bucket/glue/scripts/util/v1.36.0/somejob.py"),
					PythonVersion:  aws.String("3"),
				},
				DefaultArguments: map[string]*string{
					"Class":     aws.String("GlueApp"),
					"--TempDir": aws.String("s3://aws-glue-assets-123456789000-us-east-1/temporary/"),
				},
				MaxRetries:        aws.Int64(0),
				AllocatedCapacity: nil,
				Timeout:           aws.Int64(2880),
				MaxCapacity:       nil,
				WorkerType:        aws.String("G.1X"),
				NumberOfWorkers:   aws.Int64(2),
				GlueVersion:       aws.String("4.0"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toGlueJob(tt.args.glueJob); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toGlueJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeployment_Deploy(t *testing.T) {
	type fields struct {
		artifact    string
		version     string
		environment string
		role        string
		autoapprove bool
		sess        *session.Session
		changeSet   []Change
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Deployment{
				artifact:    tt.fields.artifact,
				version:     tt.fields.version,
				environment: tt.fields.environment,
				role:        tt.fields.role,
				autoapprove: tt.fields.autoapprove,
				sess:        tt.fields.sess,
				changeSet:   tt.fields.changeSet,
			}
			if err := d.Deploy(); (err != nil) != tt.wantErr {
				t.Errorf("Deployment.Deploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
