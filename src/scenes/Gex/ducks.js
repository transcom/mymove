import { SendGexRequest } from './api.js';
import * as helpers from 'shared/ReduxHelpers';

const resource = 'SEND_GEX_REQUEST';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const sendGexRequest = helpers.generateAsyncActionCreator(
  resource,
  SendGexRequest,
);

const initialStateMixin = { schema: {} };
export const gexReducer = helpers.generateAsyncReducer(
  resource,
  v => null,
  initialStateMixin,
);
