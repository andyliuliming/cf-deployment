# Step1. create cfproxy VM
create vm in bosh deploy manifest file and create use cfproxy in the vm
add cfproxy ip to cell hosts

# Step2. extract scripts to cfproxy
please keep the cf2kube dir not changed and set os.env("CF2KUBE") to this cf2kube path
eg: os.setEnv("CF2KUBE", "/home/cfproxy/cf2kube")

# Step3. config go env to run server in cf2kube/server/run.sh
demo bashrc                     # go env needed to run server
(
    # ljh # add go env
    export GOROOT=/usr/local/go
    export GOBIN=$GOROOT/bin
    export PATH=$PATH:$GOBIN
    export GOPATH=$HOME/gopath
    export TMPDIR=$HOME/tmp/        # run go without sudo
)

# Step4. pre-requisitions
install docker-engine
install go rt
install azure-cli
create acr (_create_acr_demo.sh is a demo script)
login docker to acr
install kubectl
config kubectl context
> eg: docker login {reg acr name}.azurecr.io -u {reg acr name} -p {secret}
> eg: az acs kubernetes get-credentials --resource-group={resource group name} --name={acs name}
