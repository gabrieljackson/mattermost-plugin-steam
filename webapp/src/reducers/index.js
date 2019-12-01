import {combineReducers} from 'redux';

import ActionTypes from '../action_types';

function steamProfiles(state = {}, action) {
    switch (action.type) {
    case ActionTypes.RECEIVED_STEAM_PROFILE: {
        const nextState = {...state};
        nextState[action.userID] = action.data;
        return nextState;
    }
    default:
        return state;
    }
}

export default combineReducers({
    steamProfiles,
});
