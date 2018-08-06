import * as helpers from 'shared/ReduxHelpers';
import { GetSpec, GetPublicSpec } from './api';
const resource = 'SWAGGER';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const loadSchema = helpers.generateAsyncActionCreator(resource, GetSpec);

export const loadPublicSchema = helpers.generateAsyncActionCreator(resource, GetPublicSpec);

export default helpers.generateAsyncReducer(resource, v => ({ spec: v }));
