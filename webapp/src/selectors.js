import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const steamProfileInfo = (state, id) => getPluginState(state).steamProfiles[id] || {};
