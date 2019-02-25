import { denormalize } from 'normalizr';

import { storageInTransit as storageInTransitModel } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { get, filter, keys } from 'lodash';

const createStorageInTransitLabel = 'StorageInTransits.createStorageInTransit';

export const selectStorageInTransits = (state, shipmentId) => {
  let filteredItems = denormalize(
    keys(get(state, 'entities.storageInTransit', {})),
    storageInTransitModel,
    state.entities,
  );

  return filter(filteredItems, item => {
    return item.shipmentId === shipmentId;
  });
};

export function createStorageInTransit(shipmentId, storageInTransit, label = createStorageInTransitLabel) {
  return swaggerRequest(
    getPublicClient,
    'storage_in_transits.createStorageInTransit',
    { shipmentId, storageInTransit },
    { label },
  );
}
