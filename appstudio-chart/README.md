# Stonesoup chart

This repository has a helm chart designed to deploy an application and several components on Stonesoup. It also has a script to automate the deployment of several applications by providing a list of GitHub repos. For more information about the script usage, go to its [README](./script/README.md) file.

## Usage 

Update the values file with the parameters to deploy the application and the number of components needed and execute the following command to deploy it in OpenShift:

```bash
helm template . | oc apply -f
```