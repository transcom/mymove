import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateDocument(fileUpload) {
  const client = await getClient();
  const response = await client.apis.documents.createDocument({
    file: fileUpload,
  });
  checkResponse(response, 'failed to upload file due to server error');
}

export default CreateDocument;
