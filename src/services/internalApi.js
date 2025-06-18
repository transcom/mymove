import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor, responseInterceptor, makeSwaggerRequestRaw } from './swaggerRequest';

import { UPLOAD_SCAN_STATUS } from 'shared/constants';

let internalClient = null;

// setting up the same config from Swagger/api.js
export async function getInternalClient() {
  if (!internalClient) {
    internalClient = await Swagger({
      url: '/internal/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }

  return internalClient;
}

// Attempt at catch-all error handling
export function getResponseError(response, defaultErrorMessage) {
  if (!response) return defaultErrorMessage;

  const detail = response.body?.detail || response.statusText || defaultErrorMessage;
  const invalidFields = response.body?.invalidFields || response.body?.invalid_fields;

  if (invalidFields && typeof invalidFields === 'object') {
    const fieldErrors = Object.entries(invalidFields)
      .map(([field, messages]) => {
        if (Array.isArray(messages)) {
          return `${field}: ${messages.join(', ')}`;
        }
        return `${field}: ${messages}`;
      })
      .join('\n');

    return `${detail}\n${fieldErrors}`;
  }

  return detail;
}

export async function makeInternalRequest(operationPath, params = {}, options = {}) {
  const client = await getInternalClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function makeInternalRequestRaw(operationPath, params = {}) {
  const client = await getInternalClient();
  return makeSwaggerRequestRaw(client, operationPath, params);
}

export async function validateCode(body) {
  return makeInternalRequestRaw('validation_code.validateCode', {
    body: {
      validationCode: body.parameterValue,
    },
  });
}

export async function getLoggedInUser(normalize = true) {
  return makeInternalRequest('users.showLoggedInUser', {}, { normalize });
}

export async function getLoggedInUserQueries(key, normalize = false) {
  return makeInternalRequest('users.showLoggedInUser', {}, { normalize });
}

export async function getBooleanFeatureFlagForUser(key, flagContext) {
  const normalize = false;
  return makeInternalRequest('featureFlags.booleanFeatureFlagForUser', { key, flagContext }, { normalize });
}

export async function getBooleanFeatureFlagUnauthenticated(key, flagContext) {
  const normalize = false;
  return makeInternalRequest('featureFlags.booleanFeatureFlagUnauthenticated', { key, flagContext }, { normalize });
}

export async function getVariantFeatureFlagForUser(key, flagContext) {
  const normalize = false;
  return makeInternalRequest('featureFlags.variantFeatureFlagForUser', { key, flagContext }, { normalize });
}

export async function getMTOShipmentsForMove(moveTaskOrderID, normalize = true) {
  return makeInternalRequest(
    'mtoShipment.listMTOShipments',
    { moveTaskOrderID },
    { normalize, label: 'mtoShipment.listMTOShipments', schemaKey: 'mtoShipments' },
  );
}

/** BELOW API CALLS ARE NOT NORMALIZED BY DEFAULT */

/** SERVICE MEMBERS */
export async function registerUser(payload) {
  return makeInternalRequest('registration.customerRegistration', { registration: payload }, { normalize: false });
}

export async function createServiceMember(serviceMember = {}) {
  return makeInternalRequest(
    'service_members.createServiceMember',
    { createServiceMemberPayload: serviceMember },
    { normalize: false },
  );
}

export async function getServiceMember(serviceMemberId) {
  return makeInternalRequest(
    'service_members.showServiceMember',
    {
      serviceMemberId,
    },
    {
      normalize: false,
    },
  );
}

export async function patchServiceMember(serviceMember) {
  return makeInternalRequest(
    'service_members.patchServiceMember',
    {
      serviceMemberId: serviceMember.id,
      patchServiceMemberPayload: serviceMember,
    },
    {
      normalize: false,
    },
  );
}

/** OKTA PROFILE */
// this will call the backend and patch the Okta profile
export async function getOktaUser() {
  return makeInternalRequest('okta_profile.showOktaInfo');
}

export async function updateOktaUser(oktaUser) {
  return makeInternalRequest(
    'okta_profile.updateOktaInfo',
    {
      updateOktaUserProfileData: oktaUser,
    },
    {
      normalize: false,
    },
  );
}

/** BACKUP CONTACTS */
export async function createBackupContactForServiceMember(serviceMemberId, backupContact) {
  return makeInternalRequest(
    'backup_contacts.createServiceMemberBackupContact',
    {
      serviceMemberId,
      createBackupContactPayload: backupContact,
    },
    {
      normalize: false,
    },
  );
}

export async function patchBackupContact(backupContact) {
  return makeInternalRequest(
    'backup_contacts.updateServiceMemberBackupContact',
    {
      backupContactId: backupContact.id,
      updateServiceMemberBackupContactPayload: backupContact,
    },
    {
      normalize: false,
    },
  );
}

/** ORDERS */
export async function getOrdersForServiceMember(serviceMemberId) {
  return makeInternalRequest(
    'service_members.showServiceMemberOrders',
    {
      serviceMemberId,
    },
    {
      normalize: false,
    },
  );
}

export async function getOrders(ordersId) {
  return makeInternalRequest(
    'orders.showOrders',
    {
      ordersId,
    },
    {
      normalize: false,
    },
  );
}

export async function createOrders(orders) {
  return makeInternalRequest(
    'orders.createOrders',
    {
      createOrders: orders,
    },
    {
      normalize: false,
    },
  );
}

export async function patchOrders(orders) {
  return makeInternalRequest(
    'orders.updateOrders',
    {
      ordersId: orders.id,
      updateOrders: orders,
    },
    {
      normalize: false,
    },
  );
}

/** UPLOADS */
export async function createUpload(file) {
  return makeInternalRequest(
    'uploads.createUpload',
    {
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function createUploadForAmendedOrdersDocument(file, ordersId) {
  return makeInternalRequest(
    'orders.uploadAmendedOrders',
    {
      ordersId,
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function createUploadForAdditionalDocuments(file, moveId) {
  return makeInternalRequest(
    'moves.uploadAdditionalDocuments',
    {
      moveId,
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function createUploadForDocument(file, documentId) {
  return makeInternalRequest(
    'uploads.createUpload',
    {
      documentId,
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function createUploadForPPMDocument(ppmShipmentId, documentId, file, weightReceipt) {
  return makeInternalRequest(
    'ppm.createPPMUpload',
    {
      ppmShipmentId,
      documentId,
      file,
      weightReceipt,
    },
    {
      normalize: false,
    },
  );
}

// Subscribes to the server sent event for the antivirus status
// This is needed because files, while uploaded, will be inaccessible
// until the AV scan has cleared it
export function waitForAvScan(uploadId, { signal } = {}) {
  return new Promise((resolve, reject) => {
    // Catch if the server event aborted before starting the func
    if (signal?.aborted) {
      reject(new DOMException('aborted', 'AbortError'));
      return;
    }

    // Init server sent event
    const es = new EventSource(`/internal/uploads/${uploadId}/status`, {
      withCredentials: true,
    });

    // Catch user cancellation helper
    // this is to close the currently running sse
    function abortListener() {
      es.close();
      signal?.removeEventListener('abort', abortListener);
      reject(new DOMException('aborted', 'AbortError'));
    }

    // Cleanup helper for the sse
    function cleanup() {
      es.close();
      signal?.removeEventListener('abort', abortListener);
    }

    // Add listener for if the user cancels the event
    signal?.addEventListener('abort', abortListener);

    // Handle incremental SSE messages
    // these cases are all provided by the
    // CustomGetUploadStatusResponse from the backend
    // See pkg/handlers/ghcapi/uploads.go
    es.onmessage = ({ data }) => {
      switch (data) {
        case UPLOAD_SCAN_STATUS.PROCESSING:
          break;
        case UPLOAD_SCAN_STATUS.CLEAN:
          cleanup();
          resolve(UPLOAD_SCAN_STATUS.CLEAN);
          break;
        case UPLOAD_SCAN_STATUS.INFECTED:
          cleanup();
          reject(new Error(UPLOAD_SCAN_STATUS.INFECTED));
          break;
        case UPLOAD_SCAN_STATUS.CONNECTION_CLOSED:
          cleanup();
          reject(new Error(UPLOAD_SCAN_STATUS.CONNECTION_CLOSED));
          break;
        default:
          cleanup();
          reject(new Error('Unknown server response from antivirus scan'));
      }
    };

    // Network SSE responded an error, close the sse
    es.onerror = (err) => {
      cleanup();
      reject(err ?? new Error('Server sent event error when listening for antivirus status'));
    };
  });
}

export async function deleteUpload(uploadId, orderId, ppmId) {
  return makeInternalRequest(
    'uploads.deleteUpload',
    {
      uploadId,
      orderId,
      ppmId,
    },
    {
      normalize: false,
    },
  );
}

export async function deleteAdditionalDocumentUpload(uploadId, moveId) {
  return makeInternalRequest(
    'uploads.deleteUpload',
    {
      uploadId,
      moveId,
    },
    {
      normalize: false,
    },
  );
}

/** MOVES */
export async function getAllMoves(serviceMemberId) {
  return makeInternalRequest(
    'moves.getAllMoves',
    {
      serviceMemberId,
    },
    {
      normalize: false,
    },
  );
}

export async function getMove(moveId) {
  return makeInternalRequest(
    'moves.showMove',
    {
      moveId,
    },
    {
      normalize: false,
    },
  );
}

export async function patchMove(moveId, move, ifMatchETag) {
  return makeInternalRequest(
    'moves.patchMove',
    {
      moveId,
      'If-Match': ifMatchETag,
      patchMovePayload: move,
    },
    {
      normalize: false,
    },
  );
}

export async function submitMoveForApproval(moveId, certificate) {
  return makeInternalRequest(
    'moves.submitMoveForApproval',
    {
      moveId,
      submitMoveForApprovalPayload: {
        certificate,
      },
    },
    {
      normalize: false,
    },
  );
}

export async function submitAmendedOrders(moveId) {
  return makeInternalRequest(
    'moves.submitAmendedOrders',
    {
      moveId,
    },
    {
      normalize: false,
    },
  );
}

/** MTO SHIPMENTS */
export async function createMTOShipment(mtoShipment) {
  return makeInternalRequest(
    'mtoShipment.createMTOShipment',
    {
      body: mtoShipment,
    },
    {
      normalize: false,
    },
  );
}

export async function patchMTOShipment(mtoShipmentId, mtoShipment, ifMatchETag) {
  return makeInternalRequest(
    'mtoShipment.updateMTOShipment',
    {
      mtoShipmentId,
      'If-Match': ifMatchETag,
      body: mtoShipment,
    },
    {
      normalize: false,
    },
  );
}

export async function deleteMTOShipment(mtoShipmentId) {
  return makeInternalRequest(
    'mtoShipment.deleteShipment',
    {
      mtoShipmentId,
    },
    {
      normalize: false,
    },
  );
}

export async function createWeightTicket(ppmShipmentId) {
  return makeInternalRequest(
    'ppm.createWeightTicket',
    {
      ppmShipmentId,
    },
    {
      normalize: false,
    },
  );
}

export async function patchWeightTicket(ppmShipmentId, weightTicketId, payload, eTag) {
  return makeInternalRequest(
    'ppm.updateWeightTicket',
    {
      ppmShipmentId,
      weightTicketId,
      'If-Match': eTag,
      updateWeightTicketPayload: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function deleteWeightTicket(ppmShipmentId, weightTicketId) {
  return makeInternalRequest(
    'ppm.deleteWeightTicket',
    {
      ppmShipmentId,
      weightTicketId,
    },
    {
      normalize: false,
    },
  );
}

export async function createProGearWeightTicket(ppmShipmentId) {
  return makeInternalRequest(
    'ppm.createProGearWeightTicket',
    {
      ppmShipmentId,
    },
    {
      normalize: false,
    },
  );
}

export async function patchProGearWeightTicket(ppmShipmentId, proGearWeightTicketId, payload, eTag) {
  return makeInternalRequest(
    'ppm.updateProGearWeightTicket',
    {
      ppmShipmentId,
      proGearWeightTicketId,
      'If-Match': eTag,
      updateProGearWeightTicket: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function deleteProGearWeightTicket(ppmShipmentId, proGearWeightTicketId) {
  return makeInternalRequest(
    'ppm.deleteProGearWeightTicket',
    {
      ppmShipmentId,
      proGearWeightTicketId,
    },
    {
      normalize: false,
    },
  );
}

export async function createMovingExpense(ppmShipmentId) {
  return makeInternalRequest(
    'ppm.createMovingExpense',
    {
      ppmShipmentId,
    },
    {
      normalize: false,
    },
  );
}

export async function patchMovingExpense(ppmShipmentId, movingExpenseId, payload, eTag) {
  return makeInternalRequest(
    'ppm.updateMovingExpense',
    {
      ppmShipmentId,
      movingExpenseId,
      'If-Match': eTag,
      updateMovingExpense: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function deleteMovingExpense(ppmShipmentId, movingExpenseId) {
  return makeInternalRequest(
    'ppm.deleteMovingExpense',
    {
      ppmShipmentId,
      movingExpenseId,
    },
    {
      normalize: false,
    },
  );
}

export async function submitPPMShipmentSignedCertification(ppmShipmentId, payload) {
  return makeInternalRequest(
    'ppm.submitPPMShipmentDocumentation',
    {
      ppmShipmentId,
      savePPMShipmentSignedCertificationPayload: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function searchTransportationOffices(search) {
  return makeInternalRequest('transportation_offices.getTransportationOffices', { search }, { normalize: false });
}

export async function downloadPPMAOAPacket(ppmShipmentId) {
  return makeInternalRequestRaw('ppm.showAOAPacket', { ppmShipmentId });
}

export async function downloadPPMPaymentPacket(ppmShipmentId) {
  return makeInternalRequestRaw('ppm.showPaymentPacket', { ppmShipmentId });
}

export async function searchLocationByZipCityState(search, includePOBoxes) {
  return makeInternalRequest('addresses.getLocationByZipCityState', { search, includePOBoxes }, { normalize: false });
}

export async function dateSelectionIsWeekendHoliday(countryCode, date) {
  return makeInternalRequestRaw(
    'calendar.isDateWeekendHoliday',
    {
      countryCode,
      date,
    },
    { normalize: false },
  );
}

export async function showCounselingOffices(dutyLocationId) {
  return makeInternalRequestRaw('transportation_offices.showCounselingOffices', { dutyLocationId });
}

export async function cancelMove(moveId) {
  return makeInternalRequestRaw('office.cancelMove', { moveId });
}

export async function getPayGradeOptions(affiliation) {
  return makeInternalRequestRaw('orders.getPayGrades', { affiliation });
}
