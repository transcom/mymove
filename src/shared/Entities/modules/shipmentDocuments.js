import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { moveDocuments } from '../schema';
import { denormalize } from 'normalizr';

export const getShipmentDocumentsLabel = 'Shipments.getAllShipmentDocuments';

export function getAllShipmentDocuments(label, shipmentId) {
  return swaggerRequest(getPublicClient, 'move_docs.indexMoveDocuments', { shipmentId }, { label });
}

export const createShipmentDocumentLabel = 'MoveDocs.createShipmentDocuments';

export function createShipmentDocument(label, shipmentId, createGenericMoveDocumentPayload) {
  return swaggerRequest(
    getPublicClient,
    'move_docs.createGenericMoveDocument',
    { shipmentId, createGenericMoveDocumentPayload },
    { label },
  );
}

// Gives an array of objects of the doucments that belong to a specific shipment
export const selectShipmentDocuments = (state, shipmentId) =>
  Object.values(state.entities.moveDocuments).filter(document => document.shipment_id === shipmentId);

const defaultShipmentDocument = {
  document: { uploads: [] },
  notes: '',
  status: '',
  title: '',
  type: '',
};

export const selectShipmentDocument = (state, id) =>
  denormalize([id], moveDocuments, state.entities)[0] || defaultShipmentDocument;
