# Terraform
Simple experiments with Terraform using AWS and GCP.

## Pre-requisities


##  Using Terraform

```shell
terraform init # only required once

# Set env variables
export AWS_ACCESS_KEY="foo"
export AWS_SECRET_KEY="bar"
export GCP_PROJECT="hur"

# See changes:
terraform plan -var "aws_access_key=$AWS_ACCESS_KEY" -var "aws_secret_key=$AWS_SECRET_KEY" -var "gcp_project=$GCP_PROJECT"

# Same thing
export TF_VAR_aws_access_key="$AWS_ACCESS_KEY"
export TF_VAR_aws_secret_key="$AWS_SECRET_KEY"
export TF_VAR_gcp_project="$GCP_PROJECT"

terraform plan

# Apply changes
terraform apply
# Apply just a specific part
terraform apply -target="google_compute_instance.default"

# Show output variables
terraform output
terraform output instance-ip # Specific output

# Cleanup
terraform destroy
```

/tmp/gcp-account.json

Creating service account:
https://console.cloud.google.com/apis/credentials/serviceaccountkey

Assign roles:
https://console.cloud.google.com/iam-admin/


'Compute Admin' Role