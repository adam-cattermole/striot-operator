# Output colors
NORMAL="\033[0;39m"
RED="\e[1;31m"
BLUE="\e[1;34m"

DEFAULT_USER="adamcattermole"
DEFAULT_NAME="striot-operator"


get_ns() {
  if [ -n "$1" ]; then
    v="$1"
  else
    v="default"
  fi
  echo $v
}

install() {
  ns=`get_ns $1`

  log "Install operator in namespace $ns"
  kubectl create -f deploy/service_account.yaml -n $ns
  kubectl create -f deploy/role.yaml -n $ns
  kubectl create -f deploy/role_binding.yaml -n $ns
  kubectl create -f deploy/crds/striot.org_topologies_crd.yaml -n $ns
  kubectl create -f deploy/operator.yaml -n $ns
}

start() {
  ns=`get_ns $1`

  log "Starting Topology in namespace $ns"
  kubectl create -f deploy/crds/striot.org_v1alpha1_topology_cr.yaml -n $ns
}

build() {
  if [ -n "$1" ]; then
    image="$1"
  else
    image="$DEFAULT_USER/$DEFAULT_NAME"
  fi

  log "Building operator image $image"
  operator-sdk build $image
}

stop() {
  ns=`get_ns $1`

  log "Stopping operator in Kubernetes for namespace $ns"
  kubectl delete -f deploy/crds/striot.org_v1alpha1_topology_cr.yaml -n $ns
  kubectl delete -f deploy/operator.yaml -n $ns
  kubectl delete -f deploy/role.yaml -n $ns
  kubectl delete -f deploy/role_binding.yaml -n $ns
  kubectl delete -f deploy/service_account.yaml -n $ns
  kubectl delete -f deploy/crds/striot.org_topologies_crd.yaml -n $ns
}

clean() {
  log "clean"
}


help() {
  echo "-----------------------------------------------------------------------"
  echo "                      Available commands                              -"
  echo "-----------------------------------------------------------------------"
  echo -e "$BLUE"
  echo "   > install - Deploy operator"
  echo "   > start   - Deploy examples"
  echo "   > build   - Build operator"
  echo "   > stop    - Remove the operator from k8s"
  echo "   > clean   - Remove the pipeline containers and all assets"
  echo "   > help    - Display this help"
  echo -e "$NORMAL"
  echo "-----------------------------------------------------------------------"
}


log() {
  echo -e "$BLUE > $1 $NORMAL" | ts '[%d-%m-%Y %H:%M:%.S]'
}

error() {
  echo ""
  echo -e "$RED >>> ERROR - $1$NORMAL" | ts '[%d-%m-%Y %H:%M:%.S]'
}



$*
