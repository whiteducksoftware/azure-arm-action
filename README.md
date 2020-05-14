# GitHub Action for Azure Resource Manager (ARM) deployment

A GitHub Action to deploy ARM templates.

![build and publish](https://github.com/whiteducksoftware/azure-arm-action/workflows/build-release/badge.svg)

![white duck logo](img/wd-githubaction-arm.png?raw=true)

## Dependencies

* [Checkout](https://github.com/actions/checkout) To checks-out your repository so the workflow can access any specified ARM template.

## Inputs
* `creds` **Required**   
    [Create Service Principal for Authentication](#Create-Service-Principal-for-Authentication)    

* `resourceGroupName` **Required**   
    Provide the name of a resource group.

* `templateLocation` **Required**  
    Specify the path to the Azure Resource Manager template.  
(See [assets/json/template.json](assets/json/template.json))

* `deploymentMode`   
    Incremental (only add resources to resource group) or Complete (remove extra resources from resource group). Default: `Incremental`.
  
* `deploymentName`  
    Specifies the name of the resource group deployment to create.

* `parametersLocation`   
    Specify the path to the Azure Resource Manager parameters file.  
    (See [assets/json/serviceprincipal.json](assets/json/parameters.json))

## Outputs

* `deploymentName` Specifies the complete deployment name which has been generated

## Usage

```yml
- uses: whiteducksoftware/azure-arm-action@v2.2
  with:
    creds: ${{ secrets.AZURE_CREDENTIALS }}
    resourceGroupName: <YourResourceGroup>
    templateLocation: <path/to/azuredeploy.json>
    deploymentName: <Deployment base name>
```

## Example

```yml
on: [push]
name: AzureLoginSample

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - uses: whiteducksoftware/azure-arm-action@v2.2
      with:
        creds: ${{ secrets.AZURE_CREDENTIALS }}
        resourceGroupName: <YourResourceGroup>
        templateLocation: <path/to/azuredeploy.json>
        parametersLocation: <path/to/parameters.json>
        deploymentName: <Deployment base name>
```
