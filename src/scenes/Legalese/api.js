import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload } from 'shared/utils';

export async function CreateCertification(certificationRequest) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.CreateSignedCertificationPayload;
  const response = await client.apis.certification.createSignedCertification(
    formatPayload(certificationRequest, payloadDef),
  );
  checkResponse(response, 'failed to create issue due to server error');
}
