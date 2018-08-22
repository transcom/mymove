import { getClient, checkResponse } from 'shared/api';
import { formatPayload } from 'shared/utils';

export async function CreateIssue(issueBody) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.CreateIssuePayload;
  const response = await client.apis.issues.createIssue({
    createIssuePayload: formatPayload(issueBody, payloadDef),
  });
  checkResponse(response, 'failed to create issue due to server error');
}
