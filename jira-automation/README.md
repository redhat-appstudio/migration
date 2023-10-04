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
$ export JIRA_TOKEN=...
$ python clone-feature.py --feature RHTAP-383 --summary "New Team Enablement" --dry-run
$ deactivate
```
