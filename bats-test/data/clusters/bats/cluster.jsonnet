{
  _cluster+: {
    cluster_name: 'bats',
    cluster_type: 'fake',
  },
  _components+: {
    // basics
    comp1: { path: 'components/comp1' },
    // parameter inheritance
    comp2: { path: 'components/comp2' },
    // parameter and cluster pruning
    comp3: { path: 'components/comp3' },
  },
  // Don't do this, but it works.
  comp2+: {
    p_cluster_jsonnet: 'cluster_jsonnet',
  },
}
