############################
# Variables (primitive API)
############################

variable "name" {
  description = "Cluster name."
  type        = string
}

variable "role_arn" {
  description = "IAM role ARN that EKS uses to manage other AWS services."
  type        = string
}

variable "kubernetes_version" {
  description = "Desired Kubernetes control-plane version (e.g., 1.30). Null lets EKS choose latest default."
  type        = string
  default     = null
}

variable "enabled_cluster_log_types" {
  description = "EKS control-plane log types to enable. Valid: api, audit, authenticator, controllerManager, scheduler."
  type        = list(string)
  default     = []
}

variable "vpc_config" {
  description = <<DESC
VPC configuration for the cluster endpoint and networking.
Required: subnet_ids.
Optional: security_group_ids, endpoint_private_access, endpoint_public_access, public_access_cidrs.
DESC
  type = object({
    subnet_ids              = list(string)
    security_group_ids      = optional(list(string))
    endpoint_private_access = optional(bool)
    endpoint_public_access  = optional(bool)
    public_access_cidrs     = optional(list(string))
  })
}

variable "kubernetes_network_config" {
  description = <<DESC
Kubernetes network settings. ip_family: IPV4 or IPV6.
service_ipv4_cidr is optional (only for IPV4 clusters).
DESC
  type = object({
    ip_family         = optional(string) # "IPV4" | "IPV6"
    service_ipv4_cidr = optional(string)
  })
  default = null
  validation {
    condition = var.kubernetes_network_config == null || (
      try(upper(var.kubernetes_network_config.ip_family) == "IPV4" ||
      upper(var.kubernetes_network_config.ip_family) == "IPV6", true)
    )
    error_message = "kubernetes_network_config.ip_family must be IPV4 or IPV6 if set."
  }
}

variable "encryption_config" {
  description = <<DESC
EKS secret encryption config. List of rules.
Each item: { provider_key_arn = KMS key ARN, resources = list of resource types, typically ["secrets"] }.
DESC
  type = list(object({
    provider_key_arn = string
    resources        = list(string)
  }))
  default = []
}

variable "access_config" {
  description = <<DESC
Cluster access configuration.
authentication_mode: CONFIG_MAP, API_AND_CONFIG_MAP, or API.
bootstrap_cluster_creator_admin_permissions: bool.
DESC
  type = object({
    authentication_mode                         = optional(string)
    bootstrap_cluster_creator_admin_permissions = optional(bool)
  })
  default = null
  validation {
    condition = var.access_config == null || (
      try(contains(["CONFIG_MAP", "API_AND_CONFIG_MAP", "API"],
      upper(var.access_config.authentication_mode)), true)
    )
    error_message = "access_config.authentication_mode must be one of: CONFIG_MAP, API_AND_CONFIG_MAP, API."
  }
}

variable "outpost_config" {
  description = <<DESC
For EKS on Outposts. Typical fields:
- control_plane_instance_type (e.g., m5.large)
- outpost_arns (list of Outpost ARNs)
DESC
  type = object({
    control_plane_instance_type = string
    outpost_arns                = list(string)
  })
  default = null
}

variable "bootstrap_self_managed_addons" {
  description = "Whether to let EKS create and manage default self-managed add-ons (vpc-cni, coredns, kube-proxy) on cluster creation."
  type        = bool
  default     = null
}

variable "tags" {
  description = "Tags to apply to the cluster."
  type        = map(string)
  default     = {}
}

variable "timeouts" {
  description = "Optional timeouts for create/update/delete."
  type = object({
    create = optional(string)
    update = optional(string)
    delete = optional(string)
  })
  default = null
}
