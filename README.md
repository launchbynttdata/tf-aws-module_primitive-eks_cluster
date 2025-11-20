# tf-aws-module_primitive-eks_cluster

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

This Terraform module creates and manages an AWS EKS (Elastic Kubernetes Service) cluster. It provides a primitive-level interface to the `aws_eks_cluster` resource with comprehensive configuration options for VPC integration, encryption, access control, logging, and network configuration.

## Pre-Commit hooks

[.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to terraform, golang and common linting tasks. There are no custom hooks added.

`commitlint` hook enforces commit message in certain format. The commit contains the following structural elements, to communicate intent to the consumers of your commit messages:

- **fix**: a commit of the type `fix` patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- **feat**: a commit of the type `feat` introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- **BREAKING CHANGE**: a commit that has a footer `BREAKING CHANGE:`, or appends a `!` after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
footers other than BREAKING CHANGE: <description> may be provided and follow a convention similar to git trailer format.
- **build**: a commit of the type `build` adds changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- **chore**: a commit of the type `chore` adds changes that don't modify src or test files
- **ci**: a commit of the type `ci` adds changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- **docs**: a commit of the type `docs` adds documentation only changes
- **perf**: a commit of the type `perf` adds code change that improves performance
- **refactor**: a commit of the type `refactor` adds code change that neither fixes a bug nor adds a feature
- **revert**: a commit of the type `revert` reverts a previous commit
- **style**: a commit of the type `style` adds code changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **test**: a commit of the type `test` adds missing tests or correcting existing tests

Base configuration used for this project is [commitlint-config-conventional (based on the Angular convention)](https://github.com/conventional-changelog/commitlint/tree/master/@commitlint/config-conventional#type-enum)

If you are a developer using vscode, [this](https://marketplace.visualstudio.com/items?itemName=joshbolduc.commitlint) plugin may be helpful.

`detect-secrets-hook` prevents new secrets from being introduced into the baseline. TODO: INSERT DOC LINK ABOUT HOOKS

In order for `pre-commit` hooks to work properly

- You need to have the pre-commit package manager installed. [Here](https://pre-commit.com/#install) are the installation instructions.
- `pre-commit` would install all the hooks when commit message is added by default except for `commitlint` hook. `commitlint` hook would need to be installed manually using the command below

```
pre-commit install --hook-type commit-msg
```

## To test the resource group module locally

1. For development/enhancements to this module locally, you'll need to install all of its components. This is controlled by the `configure` target in the project's [`Makefile`](./Makefile). Before you can run `configure`, familiarize yourself with the variables in the `Makefile` and ensure they're pointing to the right places.

```
make configure
```

This adds in several files and directories that are ignored by `git`. They expose many new Make targets.

2. _THIS STEP APPLIES ONLY TO MICROSOFT AZURE. IF YOU ARE USING A DIFFERENT PLATFORM PLEASE SKIP THIS STEP._ The first target you care about is `env`. This is the common interface for setting up environment variables. The values of the environment variables will be used to authenticate with cloud provider from local development workstation.

`make configure` command will bring down `azure_env.sh` file on local workstation. Devloper would need to modify this file, replace the environment variable values with relevant values.

These environment variables are used by `terratest` integration suit.

Service principle used for authentication(value of ARM_CLIENT_ID) should have below privileges on resource group within the subscription.

```
"Microsoft.Resources/subscriptions/resourceGroups/write"
"Microsoft.Resources/subscriptions/resourceGroups/read"
"Microsoft.Resources/subscriptions/resourceGroups/delete"
```

Then run this make target to set the environment variables on developer workstation.

```
make env
```

3. The first target you care about is `check`.

**Pre-requisites**
Before running this target it is important to ensure that, developer has created files mentioned below on local workstation under root directory of git repository that contains code for primitives/segments. Note that these files are `azure` specific. If primitive/segment under development uses any other cloud provider than azure, this section may not be relevant.

- A file named `provider.tf` with contents below

```
provider "azurerm" {
  features {}
}
```

- A file named `terraform.tfvars` which contains key value pair of variables used.

Note that since these files are added in `gitignore` they would not be checked in into primitive/segment's git repo.

After creating these files, for running tests associated with the primitive/segment, run

```
make check
```

If `make check` target is successful, developer is good to commit the code to primitive/segment's git repo.

`make check` target

- runs `terraform commands` to `lint`,`validate` and `plan` terraform code.
- runs `conftests`. `conftests` make sure `policy` checks are successful.
- runs `terratest`. This is integration test suit.
- runs `opa` tests
- runs `go lint` tests

## Features

- **Comprehensive EKS Configuration**: Support for all major EKS cluster configuration options
- **VPC Integration**: Configurable VPC, subnet, and security group associations
- **Encryption**: Optional KMS encryption for cluster secrets
- **Access Control**: Configurable API endpoint access (public, private, or both)
- **Logging**: Support for all EKS control plane logging types
- **Network Configuration**: Configurable Kubernetes network settings (service and pod CIDR)
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

  encryption_config = {
    provider = {
      key_arn = "arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012"
    }
    resources = ["secrets"]
  }

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  kubernetes_network_config = {
    service_ipv4_cidr = "172.20.0.0/16"
    ip_family         = "ipv4"
  }

  access_config = {
    authentication_mode                         = "API_AND_CONFIG_MAP"
    bootstrap_cluster_creator_admin_permissions = true
  }

  tags = {
    Environment = "production"
    Team        = "platform"
    CostCenter  = "engineering"
  }
}
```

For additional examples, see the [examples](./examples) directory:
- [Minimal](./examples/minimal) - Basic cluster configuration
- [Complete](./examples/complete) - Comprehensive configuration with all features
- [Simple](./examples/simple) - Balanced configuration for integration testing
- [Private Endpoint](./examples/private-endpoint) - Private-only API endpoint access

## Validation

This module includes validation rules to ensure secure and compliant cluster configurations:

- Kubernetes version must be 1.31 or higher
- At least 2 subnets must be specified in `vpc_config`
- If encryption is enabled, at least one resource type must be specified
- If Outpost configuration is used, all required fields must be provided

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.100 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_eks_cluster.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | Cluster name. | `string` | n/a | yes |
| <a name="input_role_arn"></a> [role\_arn](#input\_role\_arn) | IAM role ARN that EKS uses to manage other AWS services. | `string` | n/a | yes |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | Desired Kubernetes control-plane version (e.g., 1.30). Null lets EKS choose latest default. | `string` | `null` | no |
| <a name="input_enabled_cluster_log_types"></a> [enabled\_cluster\_log\_types](#input\_enabled\_cluster\_log\_types) | EKS control-plane log types to enable. Valid: api, audit, authenticator, controllerManager, scheduler. | `list(string)` | `[]` | no |
| <a name="input_vpc_config"></a> [vpc\_config](#input\_vpc\_config) | VPC configuration for the cluster endpoint and networking.<br/>Required: subnet\_ids.<br/>Optional: security\_group\_ids, endpoint\_private\_access, endpoint\_public\_access, public\_access\_cidrs. | <pre>object({<br/>    subnet_ids              = list(string)<br/>    security_group_ids      = optional(list(string))<br/>    endpoint_private_access = optional(bool)<br/>    endpoint_public_access  = optional(bool)<br/>    public_access_cidrs     = optional(list(string))<br/>  })</pre> | n/a | yes |
| <a name="input_kubernetes_network_config"></a> [kubernetes\_network\_config](#input\_kubernetes\_network\_config) | Kubernetes network settings. ip\_family: IPV4 or IPV6.<br/>service\_ipv4\_cidr is optional (only for IPV4 clusters). | <pre>object({<br/>    ip_family         = optional(string) # "IPV4" | "IPV6"<br/>    service_ipv4_cidr = optional(string)<br/>  })</pre> | `null` | no |
| <a name="input_encryption_config"></a> [encryption\_config](#input\_encryption\_config) | EKS secret encryption config. List of rules.<br/>Each item: { provider\_key\_arn = KMS key ARN, resources = list of resource types, typically ["secrets"] }. | <pre>list(object({<br/>    provider_key_arn = string<br/>    resources        = list(string)<br/>  }))</pre> | `[]` | no |
| <a name="input_access_config"></a> [access\_config](#input\_access\_config) | Cluster access configuration.<br/>authentication\_mode: CONFIG\_MAP, API\_AND\_CONFIG\_MAP, or API.<br/>bootstrap\_cluster\_creator\_admin\_permissions: bool. | <pre>object({<br/>    authentication_mode                         = optional(string)<br/>    bootstrap_cluster_creator_admin_permissions = optional(bool)<br/>  })</pre> | `null` | no |
| <a name="input_outpost_config"></a> [outpost\_config](#input\_outpost\_config) | For EKS on Outposts. Typical fields:<br/>- control\_plane\_instance\_type (e.g., m5.large)<br/>- outpost\_arns (list of Outpost ARNs) | <pre>object({<br/>    control_plane_instance_type = string<br/>    outpost_arns                = list(string)<br/>  })</pre> | `null` | no |
| <a name="input_bootstrap_self_managed_addons"></a> [bootstrap\_self\_managed\_addons](#input\_bootstrap\_self\_managed\_addons) | Whether to let EKS create and manage default self-managed add-ons (vpc-cni, coredns, kube-proxy) on cluster creation. | `bool` | `null` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to apply to the cluster. | `map(string)` | `{}` | no |
| <a name="input_timeouts"></a> [timeouts](#input\_timeouts) | Optional timeouts for create/update/delete. | <pre>object({<br/>    create = optional(string)<br/>    update = optional(string)<br/>    delete = optional(string)<br/>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_name"></a> [name](#output\_name) | Cluster name. |
| <a name="output_id"></a> [id](#output\_id) | Cluster name (resource ID). |
| <a name="output_arn"></a> [arn](#output\_arn) | Cluster ARN. |
| <a name="output_endpoint"></a> [endpoint](#output\_endpoint) | Cluster API server endpoint. |
| <a name="output_certificate_authority_data"></a> [certificate\_authority\_data](#output\_certificate\_authority\_data) | Base64-encoded certificate data required to communicate with the cluster. |
| <a name="output_status"></a> [status](#output\_status) | Cluster status. |
| <a name="output_version"></a> [version](#output\_version) | Actual Kubernetes version running on the control plane. |
| <a name="output_platform_version"></a> [platform\_version](#output\_platform\_version) | EKS platform version. |
| <a name="output_cluster_security_group_id"></a> [cluster\_security\_group\_id](#output\_cluster\_security\_group\_id) | Cluster security group ID created by EKS. |
| <a name="output_cluster_primary_security_group_id"></a> [cluster\_primary\_security\_group\_id](#output\_cluster\_primary\_security\_group\_id) | Primary security group ID for the cluster. |
| <a name="output_identity_oidc_issuer"></a> [identity\_oidc\_issuer](#output\_identity\_oidc\_issuer) | OIDC issuer URL if OIDC is enabled. |
| <a name="output_tags_all"></a> [tags\_all](#output\_tags\_all) | All tags, including provider defaults. |
<!-- END_TF_DOCS -->

## Contributing

Contributions are welcome! Please see our [contributing guidelines](CONTRIBUTING.md) for details.

## License

Apache 2.0 Licensed. See [LICENSE](LICENSE) for full details.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.100 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |

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
| <a name="output_platform_version"></a> [platform\_version](#output\_platform\_version) | EKS platform version. |
| <a name="output_status"></a> [status](#output\_status) | Cluster status. |
| <a name="output_tags_all"></a> [tags\_all](#output\_tags\_all) | All tags, including provider defaults. |
| <a name="output_version"></a> [version](#output\_version) | Actual Kubernetes version running on the control plane. |
<!-- END_TF_DOCS -->
