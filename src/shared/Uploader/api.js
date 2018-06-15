import { getClient } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}
