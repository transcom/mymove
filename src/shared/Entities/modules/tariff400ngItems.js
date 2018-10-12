import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { tariff400ngItems } from '../schema';
import { denormalize } from 'normalizr';

export function getAllTariff400ngItems(label) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.getTariff400ngItems',
    {},
    { label },
  );
}

export const selectTariff400ngItems = state => {
  return Object.values(state.entities.tariff400ngItems);
};

export const getTariff400ngItemsLabel =
  'Tariff400ngItem.getAlltariff400ngItems';

export const selectTariff400ngItem = (state, id) =>
  denormalize([id], tariff400ngItems, state.entities)[0];
