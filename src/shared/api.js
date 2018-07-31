import Swagger from 'swagger-client';
let client = null;

export async function getClient() {
  if (!client) {
    client = await Swagger('/internal/swagger.yaml');
  }
  return client;
}

export function checkResponse(response, errorMessage) {
  if (!response.ok) {
    throw new Error(`${errorMessage}: ${response.url}: ${response.statusText}`);
  }
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

export async function DeleteUpload(uploadId) {
  const client = await getClient();
  const response = await client.apis.uploads.deleteUpload({
    uploadId,
  });
  checkResponse(response, 'failed to delete file due to server error');
  return response.body;
}

export async function DeleteUploads(uploadIds) {
  const client = await getClient();
  const response = await client.apis.uploads.deleteUploads({
    uploadIds,
  });
  checkResponse(response, 'failed to delete files due to server error');
  return response.body;
}

export async function CreateDocument(name, serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.documents.createDocument({
    documentPayload: {
      name: name,
      service_member_id: serviceMemberId,
    },
  });
  checkResponse(response, 'failed to create document due to server error');
  return response.body;
}
