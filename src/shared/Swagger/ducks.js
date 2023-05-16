import { getSpec, getPublicSpec } from './api';

import * as helpers from 'shared/ReduxHelpers';

const resourceInternal = 'SWAGGER_INTERNAL';

export const actionsTypesInternal = helpers.generateAsyncActionTypes(resourceInternal);

export const loadInternalSchema = helpers.generateAsyncActionCreator(resourceInternal, getSpec);

export const swaggerReducerInternal = helpers.generateAsyncReducer(resourceInternal, (v) => ({
  spec: v,
}));

const resourcePublic = 'SWAGGER_PUBLIC';

export const actionsTypesPublic = helpers.generateAsyncActionTypes(resourcePublic);

export const loadPublicSchema = helpers.generateAsyncActionCreator(resourcePublic, getPublicSpec);
export const swaggerReducerPublic = helpers.generateAsyncReducer(resourcePublic, (v) => ({
  spec: v,
}));
