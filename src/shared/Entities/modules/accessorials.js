import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { tariff400ngItems } from '../schema';
import { denormalize } from 'normalizr';

export function getAllAccessorials(label) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.getTariff400ngItems',
    {},
    { label },
  );
}

export const selectAccessorials = state => {
  return Object.values(state.entities.tariff400ngItems);
};

export const getAccessorialsLabel = 'Tariff400ngItem.getAlltariff400ngItems';

export const selectAccessorial = (state, id) =>
  denormalize([id], tariff400ngItems, state.entities)[0];
