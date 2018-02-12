import { ensureClientIsLoaded, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await ensureClientIsLoaded();
  return client.spec;
}

export async function CreateForm1299(formData) {
  const client = await ensureClientIsLoaded();
  const response = await client.apis.form1299s.createForm1299({
    createForm1299Payload: formData,
  });
  checkResponse(response, 'failed to create form 1299 due to server error');
  //todo: return uuid?
}
