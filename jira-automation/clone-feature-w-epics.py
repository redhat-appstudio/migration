#!/usr/bin/env python
""" 
Description: Clone a feature and its child epics.
Author: rbean
"""

import argparse
import os
import pprint
import sys

import jira


def get_args():
    """
    Parse args from the command-line.
    """
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--feature",
        required=True,
        help="Key of the feature id the epic should be attached to",
    )
    parser.add_argument(
        "--summary",
        required=True,
        help="Summary of the new feature",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        default=False,
        help="Whether or not to take action on JIRA",
    )
    return parser.parse_args()


args = get_args()

url = os.environ.get("JIRA_URL", "https://issues.redhat.com")
token = os.environ.get("JIRA_TOKEN")
if not token:
    print("Set JIRA_TOKEN environment variable to your JIRA personal access token")
    sys.exit(1)

JIRA = jira.client.JIRA(server=url, token_auth=token)

print("Inspecting JIRA API.")
all_fields = JIRA.fields()
jira_name_map = {field["name"]: field["id"] for field in all_fields}
parent_key = jira_name_map["Parent Link"]
epic_name_key = jira_name_map["Epic Name"]

query = f"key={args.feature} and type=Feature"
print("Confirming the Feature exists:")
print("  > " + query)
results = JIRA.search_issues(query)
if not results:
    print(f"Feature not found via query: {query}")
    sys.exit(1)
origin_feature = results[0]


query = f"'Parent Link'={args.feature}"
print("Gathering child epics:")
print("  > " + query)
epics = JIRA.search_issues(query)
if not epics:
    print(f"No child epics found via query: {query}")
    sys.exit(1)

kwargs = {
    "project": origin_feature.fields.project.key,
    "summary": args.summary,
    "issuetype": origin_feature.fields.issuetype.name,
    "description": origin_feature.fields.description,
}
if not args.dry_run:
    new_feature = JIRA.create_issue(**kwargs)
    new_feature_key = new_feature.key
    print(f"Created feature {new_feature}, as a copy of {origin_feature}")
    print(f"https://issues.redhat.com/browse/{new_feature}")
    if origin_feature.fields.labels:
        new_feature.update(
            {
                "labels": origin_feature.fields.labels,
            }
        )
else:
    print(f"Skipped creating clone of {origin_feature}")
    print(f"  Would have created:")
    pprint.pprint(kwargs)
    new_feature_key = "<unknown>"

for origin_epic in epics:
    kwargs = {
        "project": origin_epic.fields.project.key,
        "summary": origin_epic.fields.summary,
        "issuetype": origin_epic.fields.issuetype.name,
        "description": origin_epic.fields.description or "",
        epic_name_key: getattr(origin_epic.fields, epic_name_key),
        parent_key: new_feature_key,
    }
    if not args.dry_run:
        new_epic = JIRA.create_issue(**kwargs)
        print(f"Created epic {new_epic}, as a copy of {origin_epic}")
        if origin_epic.fields.labels:
            new_epic.update(
                {
                    "labels": origin_epic.fields.labels,
                }
            )
    else:
        print(f"Skipped creating clone of {origin_epic}")
        print(f"  Would have created:")
        pprint.pprint(kwargs)

print("Done.")
