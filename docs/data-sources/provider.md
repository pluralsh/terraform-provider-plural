---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "plural_provider Data Source - terraform-provider-plural"
subcategory: ""
description: |-
  A representation of a provider you can deploy your clusters to.
---

# plural_provider (Data Source)

A representation of a provider you can deploy your clusters to.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `cloud` (String) The name of the cloud service for this provider.
- `id` (String) Internal identifier of this provider.

### Read-Only

- `editable` (Boolean) Whether this provider is editable.
- `name` (String) Human-readable name of this provider. Globally unique.
- `namespace` (String) The namespace the Cluster API resources are deployed into.
