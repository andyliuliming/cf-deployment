#! /bin/bash
# ./droplets2docker.sh 1:cache_key 2:checksum_algorithm 3:checksum_value 4:to(dir) 5:from(url) 6:start_cmd
# Test case: ./droplet2docker.sh droplets-02a6f504-4906-411d-b5d2-4e95abf6da86-c34f800d-2061-42c5-89a4-d17b841cef16 sha1 55564e7736ddb029f8d05d4bc633e23b41540bb1 . http://cloud-controller-ng.service.cf.internal:9022/internal/v2/droplets/02a6f504-4906-411d-b5d2-4e95abf6da86/55564e7736ddb029f8d05d4bc633e23b41540bb1/download sh boot.sh

if ! [[ (-x "$(command -v az)") && (-x "$(command -v docker)") ]]; then
    echo "Error: Azure-cli and docker env required"
    exit 1
fi

if ! az account list | grep -q tenantId; then
    echo "Error: please do az login first"
    exit 2
fi

ACR_NAME='testacrname'  # default acr
REG_KEY='registrykey'   # default registrykey
CACHE_KEY=$1    # file dir name
HASH_METHOD=$2
HASH_VALUE=$3
URL=$5
CMD="${@:6}"
DOCKER_IMAGE="${ACR_NAME}.azurecr.io/vcap/${HASH_VALUE}"
# log parameters
echo "Params: $*"
if [ "${CMD:0:2}x" != "shx" ] && [ "${CMD:0:2}x" != "pyx" ];then
    echo "Error: App type not support"
    exit 3
fi
cd $(dirname $(readlink -f $0))
if [ ! -d "droplets" ]; then
    mkdir droplets
fi
cd droplets
if [ -d ${CACHE_KEY} ];then
    exit 0  # already deployed and app not changed
fi
mkdir ${CACHE_KEY}
cd ${CACHE_KEY}
# not used yet, but if download file error, it can be checked.
echo ${HASH_VALUE} > ${HASH_METHOD}
curl -L ${URL} -o droplet.tar.gz
mkdir docker
# XXX target dir $4 ignored !!
tar zxvf droplet.tar.gz -C docker
if [ "${CMD:0:2}x" == "shx" ];then
    sed -e "s/@cmd/${CMD}/g" ../../Dockerfile_python > docker/Dockerfile
elif [ "${CMD:0:2}x" == "pyx" ];then
    sed -e "s/@cmd/${CMD}/g" ../../Dockerfile_python > docker/Dockerfile
else
echo "Error: App type not support"
exit 3
fi

docker build -t ${DOCKER_IMAGE} docker/
docker push ${DOCKER_IMAGE}

sed -e "s/@cid/${HASH_VALUE}/g" -e "s/@registrykey/${REG_KEY}/g" -e "s:@dockerimage:${DOCKER_IMAGE}:g" ../../deploy_template.yml > deploy.yml

kubectl create -f deploy.yml
