import * as helpers from 'shared/ReduxHelpers';
import { getSpec, getPublicSpec } from './api';

const resource = 'SWAGGER';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const loadSchema = helpers.generateAsyncActionCreator(resource, getSpec);

export const loadPublicSchema = helpers.generateAsyncActionCreator(
  resource,
  getPublicSpec,
);

export const swaggerReducer = helpers.generateAsyncReducer(resource, v => ({
  spec: v,
}));
