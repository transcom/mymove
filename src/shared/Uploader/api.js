import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateUpload(fileUpload, moveID, documentID) {
  const client = await getClient();
  const response = await client.apis.uploads.createUpload({
    file: fileUpload,
    moveId: moveID,
    documentId: documentID,
  });
  checkResponse(response, 'failed to upload file due to server error');
}

export async function CreateDocument(fileUpload, moveID) {
  const client = await getClient();
  const response = await client.apis.documents.createDocument({
    documentPayload: { name: 'document' },
    moveId: moveID,
  });
  checkResponse(response, 'failed to create document due to server error');
  CreateUpload(fileUpload, moveID, response.body.id);
}

export default CreateDocument;
