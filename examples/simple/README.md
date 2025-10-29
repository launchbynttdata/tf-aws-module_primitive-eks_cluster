# Simple EKS Cluster Example

This example demonstrates a balanced EKS cluster configuration. Used by integration tests.

## Launch Modules Used

- **tf-aws-module_primitive-iam_role** - EKS service role
- **tf-aws-module_primitive-iam_role_policy_attachment** - Policy attachments
- **tf-aws-module_primitive-vpc** - VPC
- **tf-aws-module_primitive-subnet** - Subnets

## Features

- EKS cluster with specified Kubernetes version
- Basic control plane logging (api, audit)
- Private and public endpoint access
- Minimal configuration for testing

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Cleanup

```bash
terraform destroy
```

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

| Name | Source | Version |
|------|--------|---------|
| <a name="module_resource_names"></a> [resource\_names](#module\_resource\_names) | terraform.registry.launch.nttdata.com/module_library/resource_name/launch | ~> 2.0 |
| <a name="module_eks_service_role"></a> [eks\_service\_role](#module\_eks\_service\_role) | terraform.registry.launch.nttdata.com/module_primitive/iam_role/aws | ~> 0.1 |
| <a name="module_eks_cluster_policy"></a> [eks\_cluster\_policy](#module\_eks\_cluster\_policy) | terraform.registry.launch.nttdata.com/module_primitive/iam_role_policy_attachment/aws | ~> 0.1 |
| <a name="module_eks_vpc_resource_controller_policy"></a> [eks\_vpc\_resource\_controller\_policy](#module\_eks\_vpc\_resource\_controller\_policy) | terraform.registry.launch.nttdata.com/module_primitive/iam_role_policy_attachment/aws | ~> 0.1 |
| <a name="module_vpc"></a> [vpc](#module\_vpc) | terraform.registry.launch.nttdata.com/module_primitive/vpc/aws | ~> 1.0 |
| <a name="module_subnet_1"></a> [subnet\_1](#module\_subnet\_1) | terraform.registry.launch.nttdata.com/module_primitive/subnet/aws | ~> 1.0 |
| <a name="module_subnet_2"></a> [subnet\_2](#module\_subnet\_2) | terraform.registry.launch.nttdata.com/module_primitive/subnet/aws | ~> 1.0 |
| <a name="module_eks_cluster"></a> [eks\_cluster](#module\_eks\_cluster) | ../.. | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_default_security_group.default](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/default_security_group) | resource |
| [aws_availability_zones.available](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/availability_zones) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_resource_names_map"></a> [resource\_names\_map](#input\_resource\_names\_map) | A map of key to resource\_name that will be used by tf-launch-module\_library-resource\_name to generate resource names | <pre>map(object({<br/>    name       = string<br/>    max_length = optional(number, 60)<br/>  }))</pre> | <pre>{<br/>  "eks": {<br/>    "max_length": 100,<br/>    "name": "eks"<br/>  },<br/>  "iam_role": {<br/>    "max_length": 64,<br/>    "name": "iam"<br/>  },<br/>  "subnet": {<br/>    "max_length": 80,<br/>    "name": "snet"<br/>  },<br/>  "vpc": {<br/>    "max_length": 64,<br/>    "name": "vpc"<br/>  }<br/>}</pre> | no |
| <a name="input_instance_env"></a> [instance\_env](#input\_instance\_env) | Number that represents the instance of the environment | `number` | `0` | no |
| <a name="input_instance_resource"></a> [instance\_resource](#input\_instance\_resource) | Number that represents the instance of the resource | `number` | `0` | no |
| <a name="input_logical_product_family"></a> [logical\_product\_family](#input\_logical\_product\_family) | Logical product family name | `string` | `"launch"` | no |
| <a name="input_logical_product_service"></a> [logical\_product\_service](#input\_logical\_product\_service) | Logical product service name | `string` | `"eks"` | no |
| <a name="input_class_env"></a> [class\_env](#input\_class\_env) | Environment class (e.g., dev, qa, prod) | `string` | `"sandbox"` | no |
| <a name="input_region"></a> [region](#input\_region) | AWS region | `string` | `"us-east-2"` | no |
| <a name="input_assume_role_policy"></a> [assume\_role\_policy](#input\_assume\_role\_policy) | IAM assume role policy statements to include in the trust policy. | <pre>list(object({<br/>    sid     = optional(string)<br/>    effect  = optional(string, "Allow")<br/>    actions = list(string)<br/><br/>    # each statement may have multiple principal blocks<br/>    principals = optional(list(object({<br/>      type        = string<br/>      identifiers = list(string)<br/>    })))<br/><br/>    conditions = optional(list(object({<br/>      test     = string       # e.g., "StringEquals"<br/>      variable = string       # e.g., "aws:PrincipalTag/Team"<br/>      values   = list(string) # e.g., ["DevOps"]<br/>    })))<br/>  }))</pre> | n/a | yes |
| <a name="input_vpc_cidr"></a> [vpc\_cidr](#input\_vpc\_cidr) | CIDR block for VPC | `string` | `"10.0.0.0/16"` | no |
| <a name="input_subnet_1_cidr"></a> [subnet\_1\_cidr](#input\_subnet\_1\_cidr) | CIDR block for subnet 1 | `string` | `"10.0.1.0/24"` | no |
| <a name="input_subnet_2_cidr"></a> [subnet\_2\_cidr](#input\_subnet\_2\_cidr) | CIDR block for subnet 2 | `string` | `"10.0.2.0/24"` | no |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | Kubernetes version for the EKS cluster | `string` | `"1.31"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to apply to all resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_resource_id"></a> [resource\_id](#output\_resource\_id) | EKS cluster name (primary identifier) |
| <a name="output_resource_name"></a> [resource\_name](#output\_resource\_name) | EKS cluster name |
| <a name="output_cluster_name"></a> [cluster\_name](#output\_cluster\_name) | EKS cluster name (for test validation) |
| <a name="output_cluster_arn"></a> [cluster\_arn](#output\_cluster\_arn) | EKS cluster ARN |
| <a name="output_cluster_endpoint"></a> [cluster\_endpoint](#output\_cluster\_endpoint) | EKS cluster API endpoint |
| <a name="output_cluster_version"></a> [cluster\_version](#output\_cluster\_version) | Kubernetes version |
| <a name="output_cluster_security_group_id"></a> [cluster\_security\_group\_id](#output\_cluster\_security\_group\_id) | Cluster security group ID |
| <a name="output_cluster_role_arn"></a> [cluster\_role\_arn](#output\_cluster\_role\_arn) | IAM role ARN used by the EKS cluster |
| <a name="output_cluster_tags"></a> [cluster\_tags](#output\_cluster\_tags) | Tags applied to the EKS cluster |
| <a name="output_resource_names_generated"></a> [resource\_names\_generated](#output\_resource\_names\_generated) | Map of generated resource names for reference |
<!-- END_TF_DOCS -->
