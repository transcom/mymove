import { getClient, getPublicClient } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function GetPublicSpec() {
  const client = await getPublicClient();
  return client.spec;
}
