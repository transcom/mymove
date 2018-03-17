import * as helpers from 'shared/ReduxHelpers';
import { GetSpec } from './api';
const resource = 'SWAGGER';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const loadSchema = helpers.generateAsyncActionCreator(resource, GetSpec);

export default helpers.generateAsyncReducer(resource, v => ({ spec: v }));
