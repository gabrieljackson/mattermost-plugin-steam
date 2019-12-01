import Client from '../client';
import ActionTypes from '../action_types';

import {steamProfileInfo} from 'selectors';

const STEAM_PROFILE_GET_USER_TIMEOUT_MILLISECONDS = 1000 * 60; // 1 minute

export function getSteamUserData(userID) {
    return async (dispatch, getState) => {
        if (!userID) {
            return {};
        }

        const profile = steamProfileInfo(getState(), userID);
        if (profile && profile.last_try && Date.now() - profile.last_try < STEAM_PROFILE_GET_USER_TIMEOUT_MILLISECONDS) {
            return {};
        }

        let data;
        try {
            data = await Client.getSteamProfle(userID);
        } catch (error) {
            if (error.status === 404) {
                dispatch({
                    type: ActionTypes.RECEIVED_STEAM_PROFILE,
                    userID,
                    data: {last_try: Date.now()},
                });
            }
            return {error};
        }

        dispatch({
            type: ActionTypes.RECEIVED_STEAM_PROFILE,
            userID,
            data,
        });

        return {data};
    };
}

/**
 * Stores`showRHSPlugin` action returned by
 * registerRightHandSidebarComponent in plugin initialization.
 */
export function setShowRHSAction(showRHSPluginAction) {
    return {
        type: ActionTypes.RECEIVED_SHOW_RHS_ACTION,
        showRHSPluginAction,
    };
}
