# terraform-provider-aws-synthetics
Terraform provider for AWS Synthetics Canary

#### This Terraform provider was created from PR implementation at https://github.com/terraform-providers/terraform-provider-aws/pull/13140

Once this functionality would be marged to terraform aws provider this repository would be Archived.

## Example

### main.tf
```terraform
terraform {
  required_providers {
    synthetics = {
      source = "MaksymBilenko/synthetics"
      version = ">=0.1"
    }
  }
}

data "archive_file" "synthetic" {
  type        = "zip"
  output_path = "${path.module}/files/synthetic.zip"
  source {
    content  = var.synthetic_script
    filename = "nodejs/node_modules/synthetic.js"
  }
}

resource "synthetics_canary" "terraform-deploy-test" {
  name                 = var.synthetic_name
  runtime_version      = "syn-nodejs-2.1"
  execution_role_arn   = aws_iam_policy.synthetic.arn
  artifact_s3_location = "s3://${aws_s3_bucket.synthetic_artifacts.id}/canary/"
  zip_file             = data.archive_file.synthetic.output_path
  handler              = "synthetic.handler"
  run_config {
    memory_in_mb       = 1024
    timeout_in_seconds = 60
  }
  vpc_config {
    security_group_ids = var.synthetic_vpc_config.security_group_ids
    subnet_ids         = var.synthetic_vpc_config.subnet_ids
  }
  schedule {
    duration_in_seconds = var.synthetic_schedule.duration_in_seconds
    expression          = var.synthetic_schedule.expression
  }
}
```
