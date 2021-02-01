#!/bin/bash

if [ -z "$KR8" ]; then
  KR8=kr8
fi

KR8_ARGS="-d data"
CLUSTER=bats

$KR8 $KR8_ARGS cluster list > expected/cluster_list
$KR8 $KR8_ARGS cluster components --cluster $CLUSTER > expected/cluster_components
$KR8 $KR8_ARGS cluster params --cluster $CLUSTER > expected/cluster_params
$KR8 $KR8_ARGS cluster params --cluster $CLUSTER -C comp1 > expected/cluster_params_comp1
$KR8 $KR8_ARGS cluster params --cluster $CLUSTER -P comp2 > expected/cluster_params_comp2
$KR8 $KR8_ARGS cluster params --cluster $CLUSTER -C comp3 > expected/cluster_params_comp3
$KR8 $KR8_ARGS cluster params --cluster $CLUSTER --clusterparams data/misc/cluster_params.jsonnet > expected/cluster_params_file
# This is an odd case, but creates unique output, while -P intentionally fails
$KR8 $KR8_ARGS cluster params --cluster $CLUSTER -C no_component > expected/cluster_params_no_comp
