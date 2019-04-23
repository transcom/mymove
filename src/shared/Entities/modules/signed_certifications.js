import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { filter } from 'lodash';

export const createSignedCertificationLabel = 'SignedCertifications.createSignedCertification';
export const getSignedCertificationsLabel = 'SignedCertifications.indexSignedCertifications';

export function createSignedCertification(
  moveId,
  payload /*shape: {personally_procured_move_id, shipment_id, certification_text, signature, date, certification_type}*/,
  label = createSignedCertificationLabel,
) {
  return swaggerRequest(
    getClient,
    'certification.createSignedCertification',
    {
      moveId: moveId,
      createSignedCertificationPayload: {
        ...payload,
      },
    },
    { label },
  );
}

export function getSignedCertification(moveId, label = getSignedCertificationsLabel) {
  return swaggerRequest(getClient, 'certification.indexSignedCertification', { moveId }, { label });
}

export function selectPaymentRequestCertificationForMove(state, moveId) {
  const signedCertifications = filter(state.entities.signedCertifications, cert => {
    return cert.certification_type === 'PPM_PAYMENT' && cert.move_id === moveId;
  });
  if (!signedCertifications.length) {
    return {};
  }
  return signedCertifications[0];
}
