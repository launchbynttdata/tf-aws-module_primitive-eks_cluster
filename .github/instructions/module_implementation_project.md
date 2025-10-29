# Terraform Primitive Module Implementation Guide

> A comprehensive guide for AI coding agents implementing Terraform primitive modules following Launch organizational standards and conventions.

## 1. Overview

This guide provides the framework, conventions, and best practices for implementing **Terraform primitive modules** at Launch. A primitive module manages a **single AWS resource** and follows strict conventions for structure, testing, and documentation.

### ⚠️ CRITICAL: Framework Setup Must Be First

**Before writing any code, you MUST:**

> Reference Implementation: [`tf-aws-module_primitive-iam_role`](https://github.com/launchbynttdata/tf-aws-module_primitive-iam_role)

1. Copy three essential files EXACTLY: `Makefile`, `.tool-versions`, `.gitignore`
2. Run `make configure` to bootstrap the environment
3. Run `pre-commit install` to enable quality gates
4. Copy `.github/*` directory EXACTLY (all workflows and configurations)
5. Copy `tests/*` directory EXACTLY (entire test framework structure)

**These are MANDATORY, not optional. Copy them exactly as-is from the reference repository.**

**Only after framework setup is complete should you begin writing module code.**

See [Section 9.1 Initial Setup](#91-initial-setup-critical-do-this-first) for detailed instructions.

### What is a Primitive Module?

- Wraps **one AWS resource** (e.g., `aws_vpc_security_group_ingress_rule`, `aws_iam_role`)
- Lives at the **repository root** (not in a subdirectory)
- Provides sensible defaults and lightweight validation
- Exposes resource functionality for composition in higher-level modules
- Follows the **reference implementation** from [`tf-aws-module_primitive-iam_role`](https://github.com/launchbynttdata/tf-aws-module_primitive-iam_role)

## 2. Module Structure & Conventions

### 2.1 Repository Layout

```
repo-root/
├── main.tf                     # Single resource definition
├── variables.tf                # Input variables
├── outputs.tf                  # Output values
├── locals.tf                   # Local values and validation logic
├── versions.tf                 # Terraform and provider constraints
├── README.md                   # Module documentation
├── .gitignore                  # Standard gitignore
├── .tool-versions              # Tool version specifications
├── Makefile                    # Authoritative build/test tool
├── .github/
│   └── workflows/              # CI/CD workflows
├── examples/
│   ├── complete/               # Comprehensive example
│   ├── minimal/                # Minimum viable configuration
│   ├── simple/                 # Used by integration tests
│   └── [feature-specific]/     # Additional examples as needed
└── tests/
    ├── post_deploy_functional/
    ├── post_deploy_functional_readonly/
    └── testimpl/
```

### 2.2 File Conventions

#### main.tf

- Contains **exactly one resource** block
- May also contain data blocks, if necessary.
- Maps variables directly to resource arguments
- Applies merged tags using `local.tags`
- No hardcoded values

#### variables.tf

- Maps to **all major resource arguments**
- Uses precise validation **or no validation**
- Includes descriptive documentation
- Sets sensible defaults where appropriate

#### outputs.tf

- Exposes all resource attributes useful for composition
- Always includes: resource ID, name/identifier, key attributes
- Provides computed/effective values (e.g., effective source descriptor)

#### locals.tf

- Contains validation logic using **check blocks** (Terraform 1.5+)
- Implements canonical tagging pattern
- Defines computed values for outputs
- Keeps logic readable and maintainable

#### versions.tf

- Uses `~>` constraints with pinned releases
- Specifies minimum feature version required
- Example:

  ```hcl
  terraform {
    required_version = "~> 1.5"  # Requires check blocks
    required_providers {
      aws = {
        source  = "hashicorp/aws"
        version = "~> 5.100"    # Feature-complete version
      }
    }
  }
  ```

## 3. Tagging Strategy

### 3.1 Canonical Pattern

**All primitive modules MUST implement this exact tagging pattern:**

```hcl
locals {
  default_tags = {
    provisioner = "Terraform"
  }
  tags = merge(local.default_tags, var.tags)
}
```

### 3.2 Tag Application

- Apply `local.tags` to the resource
- Always provide `var.tags` input variable (type `map(string)`, default `{}`)
- User tags override defaults via `merge()` order
- Reference: [tf-aws-module_primitive-lambda_layer/locals.tf](https://github.com/launchbynttdata/tf-aws-module_primitive-lambda_layer/blob/main/locals.tf#L13-L19)

## 4. Validation Strategy

### 4.1 Validation Policy

**Validations must be either precise or removed.** No imprecise validations allowed.

#### ✅ IMPLEMENT: Precise Validations

- ID format validation (e.g., `^sg-[a-f0-9]+$`, `^pl-[a-f0-9]+$`)
- Mutual exclusivity checks
- Required field combinations
- Protocol-specific requirements (when 100% accurate)

#### ❌ REMOVE: Imprecise Validations

- Port range values (protocol-dependent: TCP uses 1-65535, ICMP uses -1)
- Port range ordering (provider handles this better)
- Complex conditional logic that varies by context
- Any validation that could produce false positives

### 4.2 Validation Implementation

**Use check blocks (Terraform 1.5+), never preconditions:**

```hcl
# locals.tf
locals {
  # Count sources specified
  source_count = sum([
    var.source_a != null ? 1 : 0,
    var.source_b != null ? 1 : 0,
  ])

  validate_source_count = local.source_count == 1
}

# Validation check block
check "source_validation" {
  assert {
    condition     = local.validate_source_count
    error_message = "Exactly one source must be specified: source_a or source_b."
  }
}
```

**Why check blocks?**

- Non-blocking warnings (don't fail entire apply)
- More flexible than preconditions
- Align with "lightweight validation" principle
- Let AWS provider handle complex validations

### 4.3 Variable Validation Example

```hcl
# ✅ GOOD - Precise validation
variable "security_group_id" {
  description = "Security group ID"
  type        = string

  validation {
    condition     = can(regex("^sg-[a-f0-9]+$", var.security_group_id))
    error_message = "Must be a valid security group ID starting with 'sg-'."
  }
}

# ❌ BAD - Imprecise validation
variable "port" {
  description = "Port number"
  type        = number

  # DON'T DO THIS - allows -1 for TCP when it's only valid for ICMP
  validation {
    condition     = var.port >= -1 && var.port <= 65535
    error_message = "Port must be between -1 and 65535."
  }
}

# ✅ GOOD - No validation (let provider handle it)
variable "port" {
  description = "Port number. Required for tcp/udp. Use -1 for ICMP."
  type        = number
  default     = null
}
```

## 5. Examples

### 5.1 Required Examples

Create **at least 3-4 examples** demonstrating different use cases:

1. **complete** - Comprehensive example with multiple configurations
2. **minimal** - Absolute minimum required inputs
3. **simple** - Basic working example (used by integration tests)
4. **[feature-specific]** - Examples for each major feature/use case
5. **modular** - When available, examples should implement [public primitive modules](https://github.com/orgs/launchbynttdata/repositories?language=&q=tf-aws-module_primitive&sort=&type=public) available on Launch's github.

### 5.2 Example Structure

```
examples/complete/
├── main.tf              # Module usage + supporting resources
├── variables.tf         # Example-specific variables
├── outputs.tf           # Normalized output names
├── terraform.tfvars.sample  # Sample values
├── test.tfvars          # Values for integration tests
└── README.md            # Example-specific documentation
```

**⚠️ IMPORTANT: Do NOT include provider blocks in examples**

Examples should **not** contain `provider` blocks. The test framework will create and configure the provider automatically.

```hcl
# ❌ INCORRECT - Do not include this in examples
provider "aws" {
  region = "us-east-1"
}

# ✅ CORRECT - Examples contain only module usage and resources
module "example" {
  source = "../.."
  # module configuration
}
```

**Why?**

- Test framework manages provider configuration
- Enables testing across multiple regions/accounts
- Prevents provider configuration conflicts
- Maintains consistency across all examples

### 5.3 Output Normalization

**All examples MUST use these standardized output names:**

```hcl
# Primary resource outputs (use consistent names)
output "resource_id" {
  description = "Primary resource identifier"
  value       = module.example.id
}

output "resource_name" {
  description = "Resource name or key identifier"
  value       = module.example.name
}

# Additional example-specific outputs are allowed
output "additional_info" {
  description = "Example-specific output"
  value       = module.example.specific_attribute
}
```

**Why normalize?**

- Enables generic test assertions across all examples
- Tests can validate multiple examples without hardcoded expectations
- Consistent interface for users

### 5.4 Security Best Practices in Examples

**All examples creating VPCs MUST include default security group configuration:**

```hcl
resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
  # ... other configuration
}

# Configure default security group to deny all traffic
resource "aws_default_security_group" "default" {
  vpc_id = aws_vpc.example.id

  # No ingress or egress rules = deny all traffic
  tags = {
    Name = "example-default-sg"
  }
}
```

**Why?**

- Resolves Regula FG_R00089 security warnings
- Follows AWS security best practices
- Forces explicit security group rules

## 6. Testing Framework

### 6.1 Test Structure

```
tests/
├── post_deploy_functional/
│   └── main_test.go           # Integration tests (creates resources)
├── post_deploy_functional_readonly/
│   └── main_test.go           # Read-only tests (no creation)
└── testimpl/
    ├── test_impl.go           # Test implementation
    └── types.go               # Type definitions
```

### 6.2 Generic Test Assertions

**Tests MUST be generic and NOT hardcode expectations:**

```go
// ❌ BAD - Hardcoded expectations
func testResourceProperties(t *testing.T, resourceId string) {
    // DON'T hardcode protocol, port, source type
    assert.Equal(t, "tcp", rule.Protocol)
    assert.Equal(t, int32(22), rule.Port)
    assert.NotNil(t, rule.CidrIpv4)
}

// ✅ GOOD - Generic assertions
func testResourceProperties(t *testing.T, resourceId string) {
    // Validate resource exists and has valid configuration
    assert.NotNil(t, rule.Protocol)
    assert.True(t, hasValidSource(rule))  // Any valid source type
}

func hasValidSource(rule *types.Rule) bool {
    return rule.CidrIpv4 != nil ||
           rule.CidrIpv6 != nil ||
           rule.PrefixListId != nil ||
           rule.SecurityGroupId != nil
}
```

**Why generic tests?**

- Support multiple examples with different configurations
- Avoid test failures when valid configurations differ
- Focus on validating correctness, not specific values

### 6.3 Test Execution

```bash
# Bootstrap (first time only)
make configure
pre-commit install

# Run all quality gates
make check
# Includes:
# - terraform fmt (formatting)
# - terraform validate (syntax)
# - tflint (linting)
# - conftest (policy-as-code)
# - regula (security compliance)
# - integration tests (Go)
```

## 7. CI/CD Integration

### 7.1 GitHub Workflows

**Copy workflows AS-IS from reference repository:**

- all workflow files
- all configuration files (dependabot, release-template, etc.)

**Minimal adjustments only:**

- Update module name references
- Adjust paths if necessary
- Keep all quality gates

### 7.2 Makefile is Authoritative

**ALWAYS use make targets, never call tools directly:**

```bash
# ✅ CORRECT
make check
make tfmodule/test/integration

# ❌ INCORRECT (in docs or CI)
terraform fmt
terraform validate
go test ./tests/...
```

**Why?**

- Consistent tool execution
- Proper configuration and flags
- Framework updates propagate automatically

## 8. Documentation Requirements

### 8.1 README.md Structure

~~~markdown
# [Module Name]

[Brief description]

## Pre-Commit hooks

[Hook installation instructions]

## Examples

```hcl
module "example" {
  source = "path/to/module"
  # Example usage
}
```

## Module Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| ... | ... | ... | ... | ... |

## Module Outputs

| Name | Description |
|------|-------------|
| ... | ... |

## Validation

[Describe validation rules implemented]

## Testing

[Describe test execution]

~~~

### 8.2 Auto-Generated Documentation

**Use terraform-docs for input/output tables:**

```bash
# Configured via .pre-commit-config.yaml
pre-commit run terraform_docs -a
```

## 9. Development Workflow

### 9.1 Initial Setup (CRITICAL: Do This First)

**⚠️ MANDATORY: Copy framework files and bootstrap the environment BEFORE any coding begins.**

**All framework files must be copied EXACTLY as-is from the reference repository. Do not modify, skip, or customize them during setup.**

Reference repository: [`tf-aws-module_primitive-iam_role`](https://github.com/launchbynttdata/tf-aws-module_primitive-iam_role)

#### Step 1: Copy Essential Framework Files (EXACT COPY REQUIRED)

Copy these three critical files **exactly as-is** from the reference repository:

```bash
# 1. Navigate to your new module repository
cd /path/to/new-module-repo

# 2. Copy the three essential framework files EXACTLY
cp /path/to/reference-repo/Makefile .
cp /path/to/reference-repo/.tool-versions .
cp /path/to/reference-repo/.gitignore .
```

**Why these files are essential:**

- **Makefile**: Authoritative build/test tool with all quality gates and commands
- **.tool-versions**: Pins exact tool versions (Terraform, Go, etc.) for consistency
- **.gitignore**: Standard ignore patterns for Terraform projects

**DO NOT modify these files during setup. Copy them exactly.**

#### Step 2: Bootstrap the Environment

After copying the framework files, run `make configure` to bootstrap all remaining files:

```bash
# 3. Bootstrap the environment (generates .pre-commit-config.yaml and other configs)
make configure

# 4. Install pre-commit hooks
pre-commit install
```

**What `make configure` does:**

- Generates `.pre-commit-config.yaml` with all quality gates
- Sets up local development environment
- Ensures all required tools are available
- Prepares the repository for development

#### Step 3: Copy GitHub Workflows (EXACT COPY REQUIRED)

**MANDATORY: Copy the entire `.github/` directory exactly as-is:**

```bash
# 5. Copy GitHub workflows and configurations EXACTLY
cp -r /path/to/reference-repo/.github .
```

**What's included in `.github/`:**

- `workflows/pre-commit.yaml` - Pre-commit quality gates CI
- `workflows/tests.yaml` - Integration test automation
- Other workflow configurations

**DO NOT:**

- Skip any workflow files
- Modify workflows during setup
- Create custom workflows before copying the framework

**You may make minimal adjustments AFTER setup:**

- Update module name references
- Adjust paths if absolutely necessary
- Keep all quality gates intact

#### Step 4: Copy Test Framework (EXACT COPY REQUIRED)

**MANDATORY: Copy the entire `tests/` directory structure exactly as-is:**

```bash
# 6. Copy test framework structure EXACTLY
cp -r /path/to/reference-repo/tests .
```

**What's included in `tests/`:**

- `post_deploy_functional/` - Functional tests after deployment
- `post_deploy_functional_readonly/` - Read-only validation tests
- `testimpl/` - Test implementation utilities
- Go module files (`go.mod`, `go.sum`)
- Test execution scripts

**DO NOT:**

- Skip any test directories
- Modify test structure during setup
- Create custom test patterns before copying the framework

**You MUST adapt tests AFTER copying:**

1. Copy the framework exactly first
2. Update Go module imports to match your module
3. Modify test assertions to validate your specific resource
4. Keep the test structure and patterns intact

#### Step 5: Verify Framework Setup

```bash
# 7. Verify framework setup (may have failures until module is implemented)
make check
```

**Expected at this stage:**
- Tools are available and functional
- Pre-commit hooks are installed
- Make targets execute (may fail validation until code exists)

**⚠️ DO NOT START CODING until Steps 1-4 are complete.**

### 9.2 Implementation Order

**After framework setup is complete, implement in this order:**

1. **Core module files** (versions.tf, variables.tf, main.tf, outputs.tf, locals.tf)
2. **Examples** (minimal, complete, simple, feature-specific)
3. **Test adaptation** (modify copied test framework for your resource)
4. **Documentation** (README.md, update with terraform-docs)
5. **Validation** (run make check, fix issues)

#### Test Framework Adaptation

After copying the test framework exactly in Step 4 of Initial Setup, adapt it for your module:

1. **Update Go module name** in `tests/go.mod`:

   ```go
   module github.com/launchbynttdata/tf-aws-module_primitive-[your-resource]
   ```

2. **Update test imports** to reference your module path

3. **Modify test assertions** in `testimpl/test_impl.go`:
   - Change resource type checks (e.g., `ec2.DescribeSecurityGroupRules`)
   - Update validation logic for your specific resource attributes
   - Keep test structure and patterns intact
   - Make tests generic (don't hardcode expected values)

4. **Update test configurations** in each example's `test.tfvars`

5. **Run tests** to verify functionality:

   ```bash
   make tfmodule/test/integration
   ```

**DO NOT:**

- Change test directory structure
- Remove test categories (functional, readonly)
- Skip test implementation
- Create custom test patterns outside the framework

### 9.3 Quality Gates

Before considering implementation complete:

- ✅ `make configure` succeeds
- ✅ `pre-commit install` succeeds
- ✅ `make check` passes (all quality gates)
- ✅ All examples run successfully (init/plan/apply/destroy)
- ✅ Integration tests pass (all examples)
- ✅ Regula security scan: 0 problems
- ✅ README complete with input/output tables
- ✅ Canonical tagging implemented
- ✅ No imprecise validations

## 10. Common Patterns

### 10.1 Conditional Arguments

```hcl
resource "aws_example" "this" {
  # Required arguments
  required_arg = var.required_arg

  # Conditional arguments (only set if not null)
  optional_arg = var.optional_arg

  # Conditionally include blocks
  dynamic "complex_block" {
    for_each = var.enable_feature ? [1] : []
    content {
      setting = var.feature_setting
    }
  }
}
```

### 10.2 Mutual Exclusivity Validation

```hcl
locals {
  # Count mutually exclusive options
  option_count = sum([
    var.option_a != null ? 1 : 0,
    var.option_b != null ? 1 : 0,
    var.option_c != null ? 1 : 0,
  ])

  validate_options = local.option_count == 1
}

check "option_validation" {
  assert {
    condition     = local.validate_options
    error_message = "Exactly one option must be specified: option_a, option_b, or option_c."
  }
}
```

### 10.3 Computed Output Values

```hcl
locals {
  # Compute effective value for output
  effective_value = (
    var.explicit_value != null ? var.explicit_value :
    var.computed_value != null ? var.computed_value :
    "default_value"
  )
}

output "effective_value" {
  description = "The effective value used"
  value       = local.effective_value
}
```

## 11. Lessons Learned & Best Practices

### 11.1 Framework Files Must Be Copied Exactly

**Lesson:** Skipping, modifying, or customizing framework files during setup causes inconsistencies, broken tooling, and failed CI/CD.

**Solution:** Copy ALL framework files exactly as-is from the reference repository BEFORE any development:

- `Makefile`, `.tool-versions`, `.gitignore` (exact copies)
- `.github/*` entire directory (exact copy)
- `tests/*` entire directory structure (exact copy)

**Only adapt AFTER copying:**

- Tests: Update Go imports and assertions for your resource
- Workflows: Minimal adjustments for module name references only

### 11.2 Validation Precision

**Lesson:** Imprecise validations cause false positives and user frustration.

**Solution:** Only implement validations you can make 100% accurate. When in doubt, remove the validation and let the provider handle it.

**Example:** Port validation varies by protocol (TCP: 1-65535, ICMP: -1, all: -1). Without protocol context, validation will be wrong.

### 11.3 Test Flexibility

**Lesson:** Hardcoded test expectations break when examples use different valid configurations.

**Solution:** Make tests protocol/value-agnostic. Validate structure and validity, not specific values.

**Example:** Test that a source exists, not that it's specifically an IPv4 CIDR.

### 11.4 Output Normalization

**Lesson:** Each example using different output names requires unique test code.

**Solution:** Standardize core output names across all examples. Example-specific outputs are fine as additions.

### 11.5 Security by Default

**Lesson:** Examples with default security groups trigger security scanners.

**Solution:** Always configure default security groups to deny all traffic in examples.

### 11.6 No Provider Blocks in Examples

**Lesson:** Including provider blocks in examples causes conflicts with test framework and limits flexibility.

**Solution:** Never include `provider` blocks in example code. The test framework manages provider configuration automatically.

**Benefits:**
- Test framework controls region/account selection
- Enables multi-region testing
- Prevents provider configuration conflicts
- Maintains consistency across examples

### 11.7 Framework Adherence

**Lesson:** Custom approaches create maintenance burden and inconsistency.

**Solution:** Follow the reference framework exactly. Copy files as-is, make minimal adjustments.

### 11.8 Check Blocks vs Preconditions

**Lesson:** Preconditions fail entire operations, preventing valid use cases.

**Solution:** Use check blocks for validation. They warn without blocking, allowing provider-level validation to take precedence.

## 12. Reference Materials

### 12.1 Reference Repositories

- **Framework Reference:** [`tf-aws-module_primitive-iam_role`](https://github.com/launchbynttdata/tf-aws-module_primitive-iam_role)
- **Tagging Reference:** [`tf-aws-module_primitive-lambda_layer`](https://github.com/launchbynttdata/tf-aws-module_primitive-lambda_layer/blob/main/locals.tf#L13-L19)

### 12.2 Key Documentation

- [Terraform Check Blocks](https://developer.hashicorp.com/terraform/language/checks)
- [AWS Provider Documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Terratest Documentation](https://terratest.gruntwork.io/)

### 12.3 Tools

- **Terraform:** ~> 1.5 (for check blocks)
- **AWS Provider:** ~> 5.x (latest stable)
- **Go:** 1.24+ (for tests)
- **Pre-commit:** Latest
- **TFLint:** Latest
- **terraform-docs:** Latest

## 13. Acceptance Criteria Checklist

Use this checklist to verify implementation completeness:

### Framework Setup (MUST BE FIRST)

- [ ] Copied Makefile EXACTLY from reference repo (no modifications)
- [ ] Copied .tool-versions EXACTLY from reference repo (no modifications)
- [ ] Copied .gitignore EXACTLY from reference repo (no modifications)
- [ ] Ran `make configure` successfully
- [ ] Ran `pre-commit install` successfully
- [ ] Copied .github/ directory EXACTLY from reference repo (entire directory)
- [ ] Verified all GitHub workflow files present (pre-commit.yaml, tests.yaml, etc.)
- [ ] Copied tests/ directory EXACTLY from reference repo (entire directory structure)
- [ ] Verified all test subdirectories present (post_deploy_functional/, post_deploy_functional_readonly/, testimpl/)
- [ ] Adapted tests AFTER copying (updated Go imports, assertions, kept structure intact)

### Core Module

- [ ] All core files present (main.tf, variables.tf, outputs.tf, locals.tf, versions.tf)
- [ ] Single resource definition in main.tf
- [ ] All major resource arguments exposed as variables
- [ ] Canonical tagging pattern implemented
- [ ] Check blocks used for validation (if any)
- [ ] No imprecise validations
- [ ] Version constraints use ~> with pinned releases

### Examples

- [ ] At least 3-4 examples created
- [ ] All examples include minimal, complete, simple
- [ ] Normalized output names across examples
- [ ] Each example has terraform.tfvars.sample
- [ ] Each example has test.tfvars for integration tests
- [ ] NO examples contain provider blocks (test framework manages providers)
- [ ] All examples creating VPCs configure default security group
- [ ] All examples run successfully (init/plan/apply/destroy)

### Testing

- [ ] Test framework adapted from reference repo
- [ ] Integration tests pass for all examples
- [ ] Tests are generic (no hardcoded expectations)
- [ ] Go module dependencies correct for target resource

### Quality Gates

- [ ] `make configure` succeeds
- [ ] `pre-commit install` succeeds
- [ ] `make check` passes (all gates)
- [ ] Regula security scan: 0 problems
- [ ] TFLint: 0 issues
- [ ] Conftest: all tests pass
- [ ] terraform fmt: no changes needed
- [ ] terraform validate: success

### Documentation

- [ ] README.md complete
- [ ] Input/output tables auto-generated
- [ ] Example usage included
- [ ] Validation rules documented
- [ ] Test instructions included

### CI/CD

- [ ] GitHub workflows copied from reference
- [ ] Workflows call make targets (not tools directly)
- [ ] All workflow jobs pass

## 14. Quick Start Template

```bash
# PHASE 1: Framework Setup (MANDATORY - DO THIS FIRST)
# =====================================================

# 1. Create new repository for the primitive module, if not already done
mkdir tf-aws-module_primitive-[resource_type]
cd tf-aws-module_primitive-[resource_type]
git init

# 2. Copy the three essential framework files EXACTLY (no modifications)
cp /path/to/tf-aws-module_primitive-iam_role/Makefile .
cp /path/to/tf-aws-module_primitive-iam_role/.tool-versions .
cp /path/to/tf-aws-module_primitive-iam_role/.gitignore .

# 3. Bootstrap environment (generates .pre-commit-config.yaml and other configs)
make configure

# 4. Install pre-commit hooks
pre-commit install

# 5. Copy GitHub workflows EXACTLY (entire directory, no modifications)
cp -r /path/to/tf-aws-module_primitive-iam_role/.github .

# 6. Copy test framework EXACTLY (entire directory structure, no modifications)
cp -r /path/to/tf-aws-module_primitive-iam_role/tests .

# 7. Verify framework setup (tools should be functional)
make check  # May fail validation until module code exists

# PHASE 2: Core Module Implementation (ONLY AFTER PHASE 1 COMPLETE)
# ==================================================================

# 8. Create core module files with canonical patterns
# See templates below

# 9. Implement examples (minimal, complete, simple, etc.)

# 10. Adapt test framework (AFTER copying exactly):
#    - Update Go module name in tests/go.mod
#    - Update test imports
#    - Modify assertions for your resource type
#    - Keep structure and patterns intact

# 11. Generate documentation
make check  # Now should pass all gates

# 12. Run integration tests
make tfmodule/test/integration
```

Core file templates:

**versions.tf**

```hcl
terraform {
  required_version = "~> 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.100"
    }
  }
}
```

**variables.tf**

```hcl
variable "required_arg" {
  description = "Required argument"
  type        = string
}

variable "tags" {
  description = "Tags to apply"
  type        = map(string)
  default     = {}
}
```

**locals.tf**

```hcl
locals {
  default_tags = {
    provisioner = "Terraform"
  }
  tags = merge(local.default_tags, var.tags)
}
```

**main.tf**

```hcl
resource "aws_example" "this" {
  required_arg = var.required_arg
  tags         = local.tags
}

# outputs.tf
output "id" {
  description = "Resource ID"
  value       = aws_example.this.id
}

# 4. Create examples (minimal, complete, simple)
# 5. Adapt test framework
# 6. Run make configure && pre-commit install
# 7. Run make check
# 8. Iterate until all quality gates pass
```

## 15. Troubleshooting

### Common Issues

**Issue:** Tests fail with "wrong protocol" errors
**Solution:** Remove hardcoded protocol expectations from tests. Make assertions generic.

**Issue:** Regula reports default security group warnings
**Solution:** Add `aws_default_security_group` resource to VPC-creating examples.

**Issue:** Port validation rejects valid ICMP configuration
**Solution:** Remove port range validation. Let AWS provider handle it.

**Issue:** Pre-commit hooks fail
**Solution:** Run `make configure` to install required tools.

**Issue:** Integration tests can't find outputs
**Solution:** Ensure all examples use normalized output names.

---

**Document Version:** 1.0
**Last Updated:** October 28, 2025
**Based On:** tf-aws-module_primitive-vpc_security_group_ingress_rule implementation
