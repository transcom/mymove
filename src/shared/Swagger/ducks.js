import * as helpers from 'shared/ReduxHelpers';
import { GetSpec } from './api';
const resource = 'SWAGGER';

const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const loadSchema = helpers.generateAsyncActionCreator(resource, GetSpec);

export default helpers.generateAsyncReducer(resource, v => v);
