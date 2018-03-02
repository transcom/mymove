import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateIssue(issueBody) {
  const client = await getClient();
  const response = await client.apis.issues.createIssue({
    createIssuePayload: issueBody,
  });
  checkResponse(response, 'failed to create issue due to server error');
}
