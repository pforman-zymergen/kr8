local parseYaml = std.native('parseYaml');
local config = std.extVar('kr8');  // imports the config from params.jsonnet
local cluster = std.extVar('kr8_cluster');  // imports the config from params.jsonnet

// This component is mainly to test pruning behaviors.

{
  cluster: cluster,
  config: config,

  // With internal pruning of parameters, these result in internal jsonnet errors
  empty_params: {
    null_param: config.p_null_param,
    list_param: [] + config.p_empty_list_param,
    object_param: {} + config.p_empty_object_param,
    string_param: '' + config.p_empty_string_param,
  },
}
