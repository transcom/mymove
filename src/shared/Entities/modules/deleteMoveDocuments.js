import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const deleteMoveDocumentLabel = `MoveDocument.deleteMoveDocument`;

export function deleteMoveDocument(moveDocumentId, label = deleteMoveDocumentLabel) {
  return swaggerRequest(getClient, 'move_docs.deleteMoveDocument', { moveDocumentId }, { label });
}
