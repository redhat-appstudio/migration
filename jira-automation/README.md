# Installation

```bash
$ python3 -m venv virtualenv
$ source virtualenv/bin/activate
$ pip install jira
$ deactivate
```

# Usage

```bash
$ source virtualenv/bin/activate
# obtain your personal Jira access token here:
# https://issues.redhat.com/secure/ViewProfile.jspa?selectedTab=com.atlassian.pats.pats-plugin:jira-user-personal-access-tokens
$ export JIRA_TOKEN=...
$ python clone-feature-w-epics.py --feature RHTAP-383 --summary "New Team Enablement" --dry-run
$ deactivate
```
