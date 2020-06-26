# terraform-provider-aws-synthetics
Terraform provider for AWS Synthetics Canary

#### This Terraform provider was created from PR implementation at https://github.com/terraform-providers/terraform-provider-aws/pull/13140

Once this functionality would be marged to terraform aws provider this repository would be Archived.

## Example

### main.tf
```terraform
provider "synthetics" {
  version = "v0.0.1"
  # assume_role {
  #   role_arn = var.assume_role_arn
  # }
  region = var.region
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
  execution_role_arn   = aws_iam_policy.synthetic.arn
  artifact_s3_location = "s3://${aws_s3_bucket.synthetic_artifacts.id}/canary/"
  zip_file             = data.archive_file.synthetic.output_path
  handler              = "synthetic.handler"
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

```bash
mkdir -p terraform.d/plugins/linux_amd64
wget https://github.com/MaksymBilenko/terraform-provider-aws-synthetics/releases/download/v0.0.1/linux_amd64-terraform-provider-aws-synthetics_v0.0.1 -O terraform.d/plugins/linux_amd64/terraform-provider-synthetics_v0.0.1
terraform init
```