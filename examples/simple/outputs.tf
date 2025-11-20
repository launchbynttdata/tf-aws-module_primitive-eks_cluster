############################
# Outputs - Normalized
############################

output "resource_id" {
  description = "EKS cluster name (primary identifier)"
  value       = module.eks_cluster.id
}

output "resource_name" {
  description = "EKS cluster name"
  value       = module.eks_cluster.id
}

output "cluster_name" {
  description = "EKS cluster name (for test validation)"
  value       = module.eks_cluster.name
}

output "cluster_arn" {
  description = "EKS cluster ARN"
  value       = module.eks_cluster.arn
}

output "cluster_endpoint" {
  description = "EKS cluster API endpoint"
  value       = module.eks_cluster.endpoint
}

output "cluster_version" {
  description = "Kubernetes version"
  value       = module.eks_cluster.version
}

output "cluster_security_group_id" {
  description = "Cluster security group ID"
  value       = module.eks_cluster.cluster_security_group_id
}

output "cluster_role_arn" {
  description = "IAM role ARN used by the EKS cluster"
  value       = module.eks_service_role.role_arn
}

output "cluster_tags" {
  description = "Tags applied to the EKS cluster"
  value       = module.eks_cluster.tags_all
}

output "resource_names_generated" {
  description = "Map of generated resource names for reference"
  value = {
    eks_cluster = module.resource_names["eks"].minimal
    iam_role    = module.resource_names["iam_role"].minimal
    vpc         = module.resource_names["vpc"].standard
  }
}
