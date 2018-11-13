import { GetCookieURL } from './api.js';
import * as helpers from 'shared/ReduxHelpers';

const resource = 'SET_DPS_AUTH_COOKIE';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const getCookieURL = helpers.generateAsyncActionCreator(resource, GetCookieURL);

const initialStateMixin = { schema: {} };

export const dpsAuthCookieReducer = helpers.generateAsyncReducer(resource, v => null, initialStateMixin);
