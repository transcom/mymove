import { getClient, checkResponse } from 'shared/api';
import { legaleseSample } from './legaleseSample';

function timeout(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
// This function will be an API call one day. For now loads a sample.
export async function GetCertificationText() {
  await timeout(100);
  return legaleseSample;
}

export async function GetCertifications(moveId, limit) {
  const client = await getClient();
  const response = await client.apis.certification.indexSignedCertifications({
    moveId,
    limit,
  });
  checkResponse(response, 'failed to find certs due to server error');
  return response.body;
}

export async function CreateCertification(certificationRequest) {
  const client = await getClient();
  const response = await client.apis.certification.createSignedCertification(
    certificationRequest,
  );
  checkResponse(response, 'failed to create issue due to server error');
}
