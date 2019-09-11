import { getClient, checkResponse } from 'shared/Swagger/api';
import * as legalese from './legaleseText';
import { formatPayload } from 'shared/utils';

// This function will be an API call one day. For now loads a sample.
export async function GetCertificationText(hasSIT, hasAdvance) {
  const txt = [legalese.ppmStandardLiability];
  if (hasSIT) txt.push(legalese.storageLiability);
  if (hasAdvance) txt.push(legalese.ppmAdvance);
  txt.push(legalese.additionalInformation);
  return txt.join('');
}

export async function CreateCertification(certificationRequest) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.CreateSignedCertificationPayload;
  const response = await client.apis.certification.createSignedCertification(
    formatPayload(certificationRequest, payloadDef),
  );
  checkResponse(response, 'failed to create issue due to server error');
}
