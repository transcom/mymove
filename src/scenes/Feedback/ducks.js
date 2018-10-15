import { CreateIssue } from './api.js';
import { getUiSchema } from './uiSchema';
import * as helpers from 'shared/ReduxHelpers';

const resource = 'CREATE_ISSUE';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const createIssue = helpers.generateAsyncActionCreator(resource, CreateIssue);

const initialStateMixin = { schema: {}, uiSchema: getUiSchema() };
export const feedbackReducer = helpers.generateAsyncReducer(resource, v => null, initialStateMixin);
