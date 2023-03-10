# glue-deploy

An opinionated deployment process for AWS Glue.

[![Build Status](https://drone.justmiles.io/api/badges/justmiles/glue-deploy/status.svg)](https://drone.justmiles.io/justmiles/glue-deploy)

## Deployment Process

This tool targets Glue jobs based on Tags. The goal of this deployment
process is to update the command's [script location](https://docs.aws.amazon.com/glue/latest/dg/aws-glue-api-jobs-job.html#aws-glue-api-jobs-job-JobCommand)
field of a Glue job without making any other changes to the job. This
does _not_ update the additional libraries argument, just the main
script location.

To target a job for the deployment, the job requires the following two
tags and they must match your deployment options:

- `tag:ArtifactID` - the versioned artifact you are deploying
- `tag:Environment` - the logical environment of said tag

## Prequisites

This deployment process assumes specific tags have been added to your Glue jobs.
Please read through the deployment process above.

## Usage

Download the build for your machine, unzip, and add to your `$PATH`.
Run `glue-deploy --help` to view available commands

Example:

```bash
glue-deploy ship --package myapp --environment qa --version latest
```

TODO:

- provide IAM requirements
- consider updating or setting the `--additional-python-modules` argument
- support color output
- s3 validation
- persist version in SSM
