import { denormalize } from 'normalizr';

import { shipments } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

export function getAllShipmentDocuments(label, shipmentId) {
  return swaggerRequest(
    getPublicClient,
    'move_docs.indexMoveDocuments',
    { shipmentId },
    { label },
  );
}

export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
