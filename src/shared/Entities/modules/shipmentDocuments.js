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

export const selectShipmentDocuments = state =>
  Object.values(state.entities.moveDocuments);
