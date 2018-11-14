import { getClient, checkResponse } from 'shared/Swagger/api';
import * as legalese from './legaleseText';
import { formatPayload } from 'shared/utils';

// This function will be an API call one day. For now loads a sample.
export async function GetCertificationText(hasSIT, hasAdvance, moveType) {
  let txt;
  if (moveType === 'PPM') {
    txt = [legalese.ppmStandardLiability];
  } else if (moveType === 'HHG') {
    txt = [legalese.hhgStandardLiability];
  }

  if (hasSIT) txt.push(legalese.storageLiability);
  if (hasAdvance) txt.push(legalese.ppmAdvance);
  if (moveType === 'PPM') txt.push(legalese.additionalInformation);
  return txt.join('');
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
  const payloadDef = client.spec.definitions.CreateSignedCertificationPayload;
  const response = await client.apis.certification.createSignedCertification(
    formatPayload(certificationRequest, payloadDef),
  );
  checkResponse(response, 'failed to create issue due to server error');
}
