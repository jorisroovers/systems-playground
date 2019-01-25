# Terraform
Simple experiments with Terraform using AWS and GCP.

## Pre-requisities
### AWS
AWS setup is pretty straightforward, just generate an ACCESS_KEY/SECRET_KEY combination, and you're good to go.

### GCP
GCP setup is a bit more complicated.

1. Create a service account for Compute Engine at this URL: https://console.cloud.google.com/apis/credentials/serviceaccountkey
2. Download the secret file and move it to /tmp/gcp-account.json
3. Assign 'Compute Admin' role to the newly created Service account: https://console.cloud.google.com/iam-admin/
4. Make sure that the GCP_PROJECT env variable below to match the project name in the secret file.

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