#! /bin/bash
# azure-cli and docker env required
# create acr only first time !

exit 0       # update this to work.

if ! [[ (-x "$(command -v az)") && (-x "$(command -v docker)") ]]; then
    echo "Error: Azure-cli and docker env required"
    exit 1
fi

if ! az account list | grep -q tenantId; then
    az login
fi

regName='testacrname'           # will be used in deploy4droplet.sh
groupName='testgroupname'
email='t-jiahli@microsoft.com'
registrykey='registrykey'       # will be used in deploy4droplet.sh
clusterName='testk8scluster'

az group create --name ${groupName}  --location westus
az acr create -n ${regName} -g ${groupName} -l westus --admin-enabled true --sku Basic
credential=$(az acr credential show --name ${regName} --output table | grep ${regName})
password=$(echo ${credential} | cut -d ' ' -f 2)
# debug credential with print password
# echo 'debug credential.password: '${password}

# login acr with credential
docker login ${regName}.azurecr.io -u ${regName} -p ${password}

# use private registry from Kubernetes cluster
kubectl create secret docker-registry ${registrykey} --docker-server ${regName}.azurecr.io --docker-username ${regName} --docker-password ${password} --docker-email ${email}
