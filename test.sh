#!/usr/bin/env bash

set -xEeuo pipefail

namespace=${1:?missing namespace}
cluster=${2:?missing cluster}
newImage=${3:?missing image}

image=$(kubectl get cluster -n "${namespace}" "${cluster}" -o jsonpath="{.spec.imageName}")
kubectl cnpg hibernate on -n "${namespace}" "${cluster}"
pvc=$(kubectl get pvc -n "${namespace}" -l "cnpg.io/cluster=${cluster},cnpg.io/pvcRole=PG_DATA,cnpg.io/instanceRole=primary" -o name)
kubectl get -n "${namespace}" "${pvc}" -o json | sed -e "s#${image}#${newImage}#g" | kubectl apply -f -
kubectl cnpg hibernate off -n "${namespace}" "${cluster}"
kubectl cnpg fencing on -n "${namespace}" "${cluster}" "*"
kubectl wait -n "${namespace}" "pod/${pvc#*/}" --for=condition=Initialized
sleep 2
kubectl exec -n "${namespace}" "pod/${pvc#*/}" -- bash -exc '
mkdir /var/lib/postgresql/data/new
chmod 0700 /var/lib/postgresql/data/new
cd /var/lib/postgresql/data/new
initdb .
cp /var/lib/postgresql/data/pgdata/custom.conf .
cat >> postgresql.conf << EOF
# load CloudNativePG custom.conf configuration
include '\''custom.conf'\''
EOF
version=$(cat /var/lib/postgresql/data/pgdata/PG_VERSION)
pg_upgrade --link -b /usr/lib/postgresql/${version}/bin --old-datadir /var/lib/postgresql/data/pgdata --new-datadir /var/lib/postgresql/data/new
rm -f /var/lib/postgresql/data/new/delete_old_cluster.sh
cd /var/lib/postgresql/data/pgdata/
find . -depth ! -path . ! -path ./pg_wal -delete
find pg_wal/ -depth ! -path pg_wal/ -delete
mv /var/lib/postgresql/data/new/pg_wal/*  pg_wal/
rmdir /var/lib/postgresql/data/new/pg_wal
mv /var/lib/postgresql/data/new/* .
rmdir /var/lib/postgresql/data/new/
rm -fr /var/lib/postgresql/data/pgdata/pg_tblspc/*/PG_${version}_*/
'
kubectl cnpg fencing off -n "${namespace}" "${cluster}" "*"
