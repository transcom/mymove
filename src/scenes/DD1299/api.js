import { getClient, checkResponse } from 'shared/api';

export async function CreateForm1299(formData) {
  const client = await getClient();
  const response = await client.apis.form1299s.createForm1299({
    createForm1299Payload: formData,
  });
  checkResponse(response, 'failed to create form 1299 due to server error');
  //todo: return uuid?
}
