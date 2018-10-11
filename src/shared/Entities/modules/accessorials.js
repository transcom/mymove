import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { accessorials } from '../schema';
import { denormalize } from 'normalizr';

export function getAllAccessorials(label) {
  return swaggerRequest(getPublicClient, 'accessorials.getTariff400ngItems', {
    label,
  });
}

export const selectAccessorials = state =>
  Object.values(state.entities.accessorials);

export const getAccessorialsLabel = 'Accessorials.getAllAccessorials';

export const selectAccessorial = (state, id) =>
  denormalize([id], accessorials, state.entities)[0];
