resource "aws_eks_cluster" "this" {
  name     = var.name
  role_arn = var.role_arn
  version  = var.kubernetes_version

  enabled_cluster_log_types = var.enabled_cluster_log_types

  vpc_config {
    subnet_ids              = var.vpc_config.subnet_ids
    security_group_ids      = try(var.vpc_config.security_group_ids, null)
    endpoint_private_access = try(var.vpc_config.endpoint_private_access, null)
    endpoint_public_access  = try(var.vpc_config.endpoint_public_access, null)
    public_access_cidrs     = try(var.vpc_config.public_access_cidrs, null)
  }

  dynamic "kubernetes_network_config" {
    for_each = var.kubernetes_network_config == null ? [] : [var.kubernetes_network_config]
    content {
      ip_family         = try(kubernetes_network_config.value.ip_family, null)
      service_ipv4_cidr = try(kubernetes_network_config.value.service_ipv4_cidr, null)
    }
  }

  dynamic "encryption_config" {
    for_each = var.encryption_config
    content {
      provider {
        key_arn = encryption_config.value.provider_key_arn
      }
      resources = encryption_config.value.resources
    }
  }

  dynamic "access_config" {
    for_each = var.access_config == null ? [] : [var.access_config]
    content {
      authentication_mode                         = try(access_config.value.authentication_mode, null)
      bootstrap_cluster_creator_admin_permissions = try(access_config.value.bootstrap_cluster_creator_admin_permissions, null)
    }
  }

  dynamic "outpost_config" {
    for_each = var.outpost_config == null ? [] : [var.outpost_config]
    content {
      control_plane_instance_type = outpost_config.value.control_plane_instance_type
      outpost_arns                = outpost_config.value.outpost_arns
    }
  }

  bootstrap_self_managed_addons = var.bootstrap_self_managed_addons

  tags = local.tags

  dynamic "timeouts" {
    for_each = var.timeouts == null ? [] : [var.timeouts]
    content {
      create = try(timeouts.value.create, null)
      update = try(timeouts.value.update, null)
      delete = try(timeouts.value.delete, null)
    }
  }
}
