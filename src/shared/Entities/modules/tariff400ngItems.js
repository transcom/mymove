import { filter, keys, sortBy, parseInt } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { tariff400ngItems } from '../schema';
import { denormalize } from 'normalizr';
import { createSelector } from 'reselect';

export const getTariff400ngItemsLabel = 'Tariff400ngItem.getAlltariff400ngItems';

export function getAllTariff400ngItems(requires_pre_approval, label = getTariff400ngItemsLabel) {
  return swaggerRequest(getPublicClient, 'accessorials.getTariff400ngItems', { requires_pre_approval }, { label });
}

export const selectTariff400ngItems = state =>
  denormalize(keys(state.entities.tariff400ngItems), tariff400ngItems, state.entities);

export const selectSortedTariff400ngItems = createSelector([selectTariff400ngItems], items =>
  // Sorts by the numeric part of code, e.g. "256A" -> 256
  sortBy(items, item => parseInt(item.code.match(/[0-9]+/g))),
);

export const selectSortedPreApprovalTariff400ngItems = createSelector([selectSortedTariff400ngItems], items =>
  filter(items, item => item.requires_pre_approval),
);

export const selectTariff400ngItem = (state, id) => denormalize([id], tariff400ngItems, state.entities)[0];
