import { getClient, checkResponse } from 'shared/api';

export async function IndexMoveDocuments(moveId) {
  const client = await getClient();
  const response = await client.apis.moves.indexMoveDocuments({
    moveId,
  });
  checkResponse(response, 'failed to get move documents due to server error');
  return response.body;
}

export async function CreateMoveDocument(
  moveId,
  uploadIds,
  title,
  moveDocumentType,
  status,
  notes,
) {
  const client = await getClient();
  const response = await client.apis.moves.createMoveDocument({
    moveId,
    createMoveDocumentPayload: {
      upload_ids: uploadIds,
      title: title,
      move_document_type: moveDocumentType,
      status: status,
      notes: notes,
    },
  });
  checkResponse(response, 'failed to create move document due to server error');
  return response.body;
}

export async function UpdateMoveDocument(moveId, moveDocumentId, payload) {
  const client = await getClient();
  const response = await client.apis.moves.updateMoveDocument({
    moveId,
    moveDocumentId,
    updateMoveDocument: payload,
  });
  checkResponse(response, 'failed to update move document due to server error');
  return response.body;
}
