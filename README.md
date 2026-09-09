# tf-aws-module_primitive-eks_cluster

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

This Terraform module creates and manages an AWS EKS (Elastic Kubernetes Service) cluster. It provides a primitive-level interface to the `aws_eks_cluster` resource with comprehensive configuration options for VPC integration, encryption, access control, logging, and network configuration.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.100 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_eks_cluster.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_access_config"></a> [access\_config](#input\_access\_config) | Cluster access configuration.<br/>authentication\_mode: CONFIG\_MAP, API\_AND\_CONFIG\_MAP, or API.<br/>bootstrap\_cluster\_creator\_admin\_permissions: bool. | <pre>object({<br/>    authentication_mode                         = optional(string)<br/>    bootstrap_cluster_creator_admin_permissions = optional(bool)<br/>  })</pre> | `null` | no |
| <a name="input_bootstrap_self_managed_addons"></a> [bootstrap\_self\_managed\_addons](#input\_bootstrap\_self\_managed\_addons) | Whether to let EKS create and manage default self-managed add-ons (vpc-cni, coredns, kube-proxy) on cluster creation. | `bool` | `null` | no |
| <a name="input_enabled_cluster_log_types"></a> [enabled\_cluster\_log\_types](#input\_enabled\_cluster\_log\_types) | EKS control-plane log types to enable. Valid: api, audit, authenticator, controllerManager, scheduler. | `list(string)` | `[]` | no |
| <a name="input_encryption_config"></a> [encryption\_config](#input\_encryption\_config) | EKS secret encryption config. List of rules.<br/>Each item: { provider\_key\_arn = KMS key ARN, resources = list of resource types, typically ["secrets"] }. | <pre>list(object({<br/>    provider_key_arn = string<br/>    resources        = list(string)<br/>  }))</pre> | `[]` | no |
| <a name="input_kubernetes_network_config"></a> [kubernetes\_network\_config](#input\_kubernetes\_network\_config) | Kubernetes network settings. ip\_family: IPV4 or IPV6.<br/>service\_ipv4\_cidr is optional (only for IPV4 clusters). | <pre>object({<br/>    ip_family         = optional(string) # "IPV4" | "IPV6"<br/>    service_ipv4_cidr = optional(string)<br/>  })</pre> | `null` | no |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | Desired Kubernetes control-plane version (e.g., 1.30). Null lets EKS choose latest default. | `string` | `null` | no |
| <a name="input_name"></a> [name](#input\_name) | Cluster name. | `string` | n/a | yes |
| <a name="input_outpost_config"></a> [outpost\_config](#input\_outpost\_config) | For EKS on Outposts. Typical fields:<br/>- control\_plane\_instance\_type (e.g., m5.large)<br/>- outpost\_arns (list of Outpost ARNs) | <pre>object({<br/>    control_plane_instance_type = string<br/>    outpost_arns                = list(string)<br/>  })</pre> | `null` | no |
| <a name="input_role_arn"></a> [role\_arn](#input\_role\_arn) | IAM role ARN that EKS uses to manage other AWS services. | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to apply to the cluster. | `map(string)` | `{}` | no |
| <a name="input_timeouts"></a> [timeouts](#input\_timeouts) | Optional timeouts for create/update/delete. | <pre>object({<br/>    create = optional(string)<br/>    update = optional(string)<br/>    delete = optional(string)<br/>  })</pre> | `null` | no |
| <a name="input_vpc_config"></a> [vpc\_config](#input\_vpc\_config) | VPC configuration for the cluster endpoint and networking.<br/>Required: subnet\_ids.<br/>Optional: security\_group\_ids, endpoint\_private\_access, endpoint\_public\_access, public\_access\_cidrs. | <pre>object({<br/>    subnet_ids              = list(string)<br/>    security_group_ids      = optional(list(string))<br/>    endpoint_private_access = optional(bool)<br/>    endpoint_public_access  = optional(bool)<br/>    public_access_cidrs     = optional(list(string))<br/>  })</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_arn"></a> [arn](#output\_arn) | Cluster ARN. |
| <a name="output_certificate_authority_data"></a> [certificate\_authority\_data](#output\_certificate\_authority\_data) | Base64-encoded certificate data required to communicate with the cluster. |
| <a name="output_cluster_primary_security_group_id"></a> [cluster\_primary\_security\_group\_id](#output\_cluster\_primary\_security\_group\_id) | Primary security group ID for the cluster. |
| <a name="output_cluster_security_group_id"></a> [cluster\_security\_group\_id](#output\_cluster\_security\_group\_id) | Cluster security group ID created by EKS. |
| <a name="output_endpoint"></a> [endpoint](#output\_endpoint) | Cluster API server endpoint. |
| <a name="output_id"></a> [id](#output\_id) | Cluster name (resource ID). |
| <a name="output_identity_oidc_issuer"></a> [identity\_oidc\_issuer](#output\_identity\_oidc\_issuer) | OIDC issuer URL if OIDC is enabled. |
| <a name="output_name"></a> [name](#output\_name) | Cluster name. |
| <a name="output_platform_version"></a> [platform\_version](#output\_platform\_version) | EKS platform version. |
| <a name="output_status"></a> [status](#output\_status) | Cluster status. |
| <a name="output_tags_all"></a> [tags\_all](#output\_tags\_all) | All tags, including provider defaults. |
| <a name="output_version"></a> [version](#output\_version) | Actual Kubernetes version running on the control plane. |
<!-- END_TF_DOCS -->

## Module Development

### Pre-Requisites

The following commands should be available on your system:

- `asdf` or `mise`
- `make`
- `python3` (for pre-commit)

Additionally, your `git` user and email must be configured. Run the `make configure` command from the root of the repository to ensure that you meet these requirements.

### Pre-Commit hooks

The [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to Terraform and Golang, as well as some common linting tasks. These will be configured for you when you run `make configure`.

### Local Validation

You should validate the changes you make to any module locally, prior to pushing your changes in a branch to GitHub.

1. Ensure that you have run `make configure` successfully.

2. Ensure you are signed into the appropriate cloud provider (e.g. AWS or Azure) for the module under test in your current console session.

3. Run the Terraform and Golang linters with the following command:

```
make lint
```

4. Once you have satisfied the linters, the following command will build example infrastructure in your configured cloud, run the tests, and then tear down the infrastructure it created:

```
make test
```

The pre-commit validations, as well as the `make lint` and `make test` targets, will all be performed in CI. Running these validations locally prior to opening a PR helps ensure a smooth review and merge process.

### Review & Merge Process

Once your change has been tested locally and your branch pushed up, open a new Pull Request for your branch to the default (main) branch of this repository.

The title of your Pull Request will determine the version bump for this change, and the title must be in [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification) format in order to merge. A breaking change will trigger a major version bump, a feature will trigger a minor version bump, and all other types will trigger a patch version bump.

Ensure your CI workflows are passing; seek approval from teammates and address any feedback; seek any explicit approvals required by the CODEOWNERS file. You may merge the PR as soon as all requirements are met, and a new release and tag will be automatically created for you.

### Automatic Updates

The shared configuration and workflow files in this repository are largely managed through the [launch-terraform-skeleton](https://github.com/launchbynttdata/launch-terraform-skeleton) repository. Outside of perhaps the `.gitignore` to account for specific files being generated by certain Terraform modules (e.g. Lambda functions), there should not be much cause to update these files on a per-repo basis, and making changes to them individually is discouraged.

If desired, you can check for and run these updates locally in a branch if you have the `copier` tool installed. Some example commands are included below:

```
# Check for updates, optionally checking prerelease versions
copier check-update [--prereleases]

# Run an update, using default answers if there are any. We use tasks, which requires --trust to be set.
copier update --defaults --trust [--prereleases]

# Recopy from the source, and --overwrite all templated files in the process
copier recopy --defaults --trust --overwrite [--prereleases]
```

Automatic updates will run through a scheduled workflow, and if the post-update tests are successful, the Pull Request created will automatically merge. Conflicts in the update or failures to test may leave a Pull Request outstanding, which needs to be addressed by a Launch Engineer.

## Features

- **Comprehensive EKS Configuration**: Support for all major EKS cluster configuration options
- **VPC Integration**: Configurable VPC, subnet, and security group associations
- **Encryption**: Optional KMS encryption for cluster secrets
- **Access Control**: Configurable API endpoint access (public, private, or both) and cluster authentication mode (`CONFIG_MAP`, `API`, or `API_AND_CONFIG_MAP`)
- **Logging**: Support for all EKS control plane logging types
- **Network Configuration**: Configurable Kubernetes IP family (IPv4 or IPv6) and optional service IPv4 CIDR
- **Self-Managed Add-ons**: Optional bootstrap of default EKS add-ons (`vpc-cni`, `coredns`, `kube-proxy`)
- **Outpost Support**: Configuration options for AWS Outposts
- **Canonical Tagging**: Automatic tagging with `provisioner=Terraform`

## Usage

### Minimal Example

```hcl
module "eks_cluster" {
  source = "git::https://github.com/launchbynttdata/tf-aws-module_primitive-eks_cluster.git?ref=1.0.0"

  name               = "my-eks-cluster"
  role_arn           = "arn:aws:iam::123456789012:role/eks-cluster-role"
  kubernetes_version = "1.31"

  vpc_config = {
    subnet_ids = ["subnet-abc123", "subnet-def456"]
  }

  tags = {
    Environment = "development"
    Team        = "platform"
  }
}
```

### Complete Example

```hcl
module "eks_cluster" {
  source = "git::https://github.com/launchbynttdata/tf-aws-module_primitive-eks_cluster.git?ref=1.0.0"

  name               = "my-production-cluster"
  role_arn           = "arn:aws:iam::123456789012:role/eks-cluster-role"
  kubernetes_version = "1.31"

  vpc_config = {
    subnet_ids              = ["subnet-abc123", "subnet-def456"]
    security_group_ids      = ["sg-12345678"]
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = ["0.0.0.0/0"]
  }

  encryption_config = [
    {
      provider_key_arn = "arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012"
      resources         = ["secrets"]
    }
  ]

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  kubernetes_network_config = {
    service_ipv4_cidr = "172.20.0.0/16"
    ip_family         = "IPV4"
  }

  access_config = {
    authentication_mode                         = "API_AND_CONFIG_MAP"
    bootstrap_cluster_creator_admin_permissions = true
  }

  bootstrap_self_managed_addons = true

  tags = {
    Environment = "production"
    Team        = "platform"
    CostCenter  = "engineering"
  }
}
```

For additional examples, see the [examples](./examples) directory:
- [Simple](./examples/simple) - Balanced configuration for integration testing

## Validation

This module includes validation rules for cluster configuration:

- If `kubernetes_network_config.ip_family` is set, it must be `IPV4` or `IPV6`
- If `access_config.authentication_mode` is set, it must be `CONFIG_MAP`, `API_AND_CONFIG_MAP`, or `API`
