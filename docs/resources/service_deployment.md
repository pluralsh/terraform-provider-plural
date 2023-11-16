---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "plural_service_deployment Resource - terraform-provider-plural"
subcategory: ""
description: |-
  ServiceDeployment resource
---

# plural_service_deployment (Resource)

ServiceDeployment resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (Attributes) Unique cluster id/handle to deploy this ServiceDeployment (see [below for nested schema](#nestedatt--cluster))
- `name` (String) Human-readable name of this ServiceDeployment.
- `namespace` (String) Namespace to deploy this ServiceDeployment.
- `repository` (Object) Repository information used to pull ServiceDeployment. (see [below for nested schema](#nestedatt--repository))

### Optional

- `configuration` (Attributes List) List of [name, value] secrets used to alter this ServiceDeployment configuration. (see [below for nested schema](#nestedatt--configuration))

### Read-Only

- `id` (String) Internal identifier of this ServiceDeployment.

<a id="nestedatt--cluster"></a>
### Nested Schema for `cluster`

Optional:

- `handle` (String)
- `id` (String)


<a id="nestedatt--repository"></a>
### Nested Schema for `repository`

Required:

- `folder` (String)
- `id` (String)
- `ref` (String)


<a id="nestedatt--configuration"></a>
### Nested Schema for `configuration`

Required:

- `name` (String)
- `value` (String)