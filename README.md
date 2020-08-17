# GitHub Action for Azure Resource Manager (ARM) deployment

A GitHub Action to deploy ARM templates.

![build and publish](https://github.com/whiteducksoftware/azure-arm-action/workflows/build-release/badge.svg)
[![go report](https://goreportcard.com/badge/github.com/whiteducksoftware/azure-arm-action)](https://goreportcard.com/report/github.com/whiteducksoftware/azure-arm-action)
[![License: MIT](https://img.shields.io/badge/License-MIT-suceess.svg)](LICENSE)

![white duck logo](img/wd-githubaction-arm.png?raw=true)

## Dependencies

* [Checkout](https://github.com/actions/checkout) To checks-out your repository so the workflow can access any specified ARM template.

## Inputs
* `creds` **Required**   
    [Create Service Principal for Authentication](#Create-Service-Principal-for-Authentication)    

* `templateLocation` **Required**  
    Specify the path to the Azure Resource Manager template.  
(See [assets/json/template.json](test/template.json))

* `resourceGroupName`    
    Provide the name of a resource group.    
    If you dont pass a resource group name the template will be deployed at subscription scope

* `deploymentMode`   
    Incremental (only add resources to resource group) or Complete (remove extra resources from resource group). Default: `Incremental`.
  
* `deploymentName`  
    Specifies the name of the resource group deployment to create.

* `parameters`   
    Specify the path to the Azure Resource Manager parameters file or pass them as space delimited Key-Value Pairs.  
    (See [examples/Advanced.md](examples/Advanced.md))

* `overrideParameters`   
    Specify the path to the Azure Resource Manager override parameters file or pass them as space delimited Key-Value Pairs.  
    (See [examples/Advanced.md](examples/Advanced.md))

## Outputs
Every template output will be exported as output. For example the output is called `containerName` then it will be available with `${{ steps.STEP.outputs.containerName }}`    
For more Information see [examples/Advanced.md](examples/Advanced.md).    
Additionally are the following outputs available:
* `deploymentName` Specifies the complete deployment name which has been generated

## Usage

```yml
- uses: whiteducksoftware/azure-arm-action@v3.2
  with:
    creds: ${{ secrets.AZURE_CREDENTIALS }}
    resourceGroupName: <YourResourceGroup>
    templateLocation: <path/to/azuredeploy.json>
    deploymentName: <Deployment base name>
```

## Example

```yml
on: [push]
name: ARMActionSample

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - uses: whiteducksoftware/azure-arm-action@v3.2
      with:
        creds: ${{ secrets.AZURE_CREDENTIALS }}
        resourceGroupName: <YourResourceGroup>
        templateLocation: <path/to/azuredeploy.json>
        parameters: <path/to/parameters.json> OR <KEY=VALUE>
        deploymentName: <Deployment base name>
```
For more advanced workflows see [examples/Advanced.md](examples/Advanced.md).
