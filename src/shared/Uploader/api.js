import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateUpload(fileUpload, documentId) {
  const client = await getClient();
  const response = await client.apis.uploads.createUpload({
    file: fileUpload,
    documentId,
  });
  checkResponse(response, 'failed to upload file due to server error');
  return response.body;
}

export async function CreateDocument(fileUpload, serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.documents.createDocument({
    documentPayload: {
      name: 'document',
      service_member_id: serviceMemberId,
    },
  });
  checkResponse(response, 'failed to create document due to server error');
  return CreateUpload(fileUpload, response.body.id);
}

export default CreateDocument;
