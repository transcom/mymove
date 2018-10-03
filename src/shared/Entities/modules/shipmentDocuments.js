import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { moveDocuments } from '../schema';
import { denormalize } from 'normalizr';

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

export const getShipmentDocumentsLabel = 'Shipments.getAllShipmentDocuments';

const defaultShipmentDocument = {
  document: { uploads: [] },
  notes: '',
  status: '',
  title: '',
  type: '',
};

export const selectShipmentDocument = (state, id) =>
  denormalize([id], moveDocuments, state.entities)[0] ||
  defaultShipmentDocument;
