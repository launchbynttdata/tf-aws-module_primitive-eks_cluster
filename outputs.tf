
############################
# Outputs
############################

output "id" {
  description = "Cluster name (resource ID)."
  value       = aws_eks_cluster.this.id
}

output "arn" {
  description = "Cluster ARN."
  value       = aws_eks_cluster.this.arn
}

output "endpoint" {
  description = "Cluster API server endpoint."
  value       = aws_eks_cluster.this.endpoint
}

output "certificate_authority_data" {
  description = "Base64-encoded certificate data required to communicate with the cluster."
  value       = try(aws_eks_cluster.this.certificate_authority[0].data, null)
}

output "status" {
  description = "Cluster status."
  value       = aws_eks_cluster.this.status
}

output "version" {
  description = "Actual Kubernetes version running on the control plane."
  value       = aws_eks_cluster.this.version
}

output "platform_version" {
  description = "EKS platform version."
  value       = aws_eks_cluster.this.platform_version
}

output "cluster_security_group_id" {
  description = "Cluster security group ID created by EKS."
  value       = aws_eks_cluster.this.vpc_config[0].cluster_security_group_id
}

output "cluster_primary_security_group_id" {
  description = "Primary security group ID for the cluster."
  value       = aws_eks_cluster.this.vpc_config[0].cluster_security_group_id
}

output "identity_oidc_issuer" {
  description = "OIDC issuer URL if OIDC is enabled."
  value       = try(aws_eks_cluster.this.identity[0].oidc[0].issuer, null)
}

output "tags_all" {
  description = "All tags, including provider defaults."
  value       = aws_eks_cluster.this.tags_all
}
