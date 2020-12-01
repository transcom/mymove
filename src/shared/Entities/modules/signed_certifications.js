import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { filter } from 'lodash';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';

export const createSignedCertificationLabel = 'SignedCertifications.createSignedCertification';
export const getSignedCertificationsLabel = 'SignedCertifications.indexSignedCertifications';

export function createSignedCertification(
  moveId,
  payload /*shape: {personally_procured_move_id, certification_text, signature, date, certification_type}*/,
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
  const signedCertifications = filter(state.entities.signedCertifications, (cert) => {
    return cert.certification_type === SIGNED_CERT_OPTIONS.PPM_PAYMENT && cert.move_id === moveId;
  });
  if (!signedCertifications.length) {
    return {};
  }
  return signedCertifications[0];
}

export function selectSignedCertification(state) {
  const certifications = Object.values(state.entities.signedCertifications);
  return certifications.length ? certifications[0] : {};
}
