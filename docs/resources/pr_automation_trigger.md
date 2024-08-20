---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "plural_pr_automation_trigger Resource - terraform-provider-plural"
subcategory: ""
description: |-
  
---

# plural_pr_automation_trigger (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `pr_automation_branch` (String) Branch that should be created against PR Automation base branch.
- `pr_automation_id` (String) ID of the PR Automation that should be triggered.

### Optional

- `context` (Map of String) PR Automation configuration context.
- `repo_slug` (String) Repo slug of the repository PR Automation should be triggered against. If not provided PR Automation repo will be used. Example format for a github repository: <userOrOrg>/<repoName>