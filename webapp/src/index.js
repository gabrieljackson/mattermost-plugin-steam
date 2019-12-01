import Reducer from './reducers';
import {id as pluginId} from './manifest';

import UserAttribute from './components/user_attribute';

class Plugin {
    async initialize(registry) {
        registry.registerReducer(Reducer);

        registry.registerPopoverUserAttributesComponent(UserAttribute);
    }
}

global.window.registerPlugin(pluginId, new Plugin());
