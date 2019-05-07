import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

const createStorageInTransitLabel = 'StorageInTransits.createStorageInTransit';
const getStorageInTransitsLabel = 'StorageInTransits.getStorageInTransitsForShipment';
const updateStorageInTransitLabel = 'StorageInTransits.updateStorageInTransit';
const approveStorageInTransitLabel = 'StorageInTransits.approveStorageInTransit';

export const selectStorageInTransits = (state, shipmentId) => {
  const storageInTransits = Object.values(state.entities.storageInTransits).filter(
    storageInTransit => storageInTransit.shipment_id === shipmentId,
  );

  return storageInTransits;
};

export function createStorageInTransit(shipmentId, storageInTransit, label = createStorageInTransitLabel) {
  return swaggerRequest(
    getPublicClient,
    'storage_in_transits.createStorageInTransit',
    { shipmentId, storageInTransit },
    { label },
  );
}

export const getStorageInTransitsForShipment = (shipmentId, label = getStorageInTransitsLabel) => {
  return swaggerRequest(getPublicClient, 'storage_in_transits.indexStorageInTransits', { shipmentId }, { label });
};

export function updateStorageInTransit(
  shipmentId,
  storageInTransitId,
  storageInTransit,
  label = updateStorageInTransitLabel,
) {
  return swaggerRequest(
    getPublicClient,
    'storage_in_transits.patchStorageInTransit',
    { shipmentId, storageInTransitId, storageInTransit },
    { label },
  );
}

export function approveStorageInTransit(
  shipmentId,
  storageInTransitId,
  storageInTransitApprovalPayload,
  label = approveStorageInTransitLabel,
) {
  return swaggerRequest(
    getPublicClient,
    'storage_in_transits.approveStorageInTransit',
    {
      shipmentId,
      storageInTransitId,
      storageInTransitApprovalPayload,
    },
    { label },
  );
}
