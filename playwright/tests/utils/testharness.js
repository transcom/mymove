async function buildDefault(request, action) {
  const r = await request.post(`/testharness/build/${action}`);
  if (!r.ok()) {
    const body = await r.body();
    throw Error(`Error with testharness build for '${action}': ${body}`);
  }
  return r.json();
}

export async function buildDefaultAdminUser(request) {
  return buildDefault(request, 'DefaultAdminUser');
}

export async function buildOfficeUserWithTOOAndTIO(request) {
  return buildDefault(request, 'OfficeUserWithTOOAndTIO');
}

export async function buildDefaultMove(request) {
  return buildDefault(request, 'DefaultMove');
}

export async function buildMoveWithOrders(request) {
  return buildDefault(request, 'MoveWithOrders');
}

export async function buildSpouseProGearMove(request) {
  return buildDefault(request, 'SpouseProGearMove');
}

export async function buildNeedsOrdersUser(request) {
  return buildDefault(request, 'NeedsOrdersUser');
}

export async function buildPPMInProgressMove(request) {
  return buildDefault(request, 'PPMInProgressMove');
}

export async function buildWithShipmentMove(request) {
  return buildDefault(request, 'WithShipmentMove');
}

export async function buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(request) {
  return buildDefault(request, 'HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO');
}

export async function buildHHGMoveWithNTSAndNeedsSC(request) {
  return buildDefault(request, 'HHGMoveWithNTSAndNeedsSC');
}

export async function buildNTSRMoveWithPaymentRequest(request) {
  return buildDefault(request, 'NTSRMoveWithPaymentRequest');
}

export async function buildPrimeSimulatorMoveNeedsShipmentUpdate(request) {
  return buildDefault(request, 'PrimeSimulatorMoveNeedsShipmentUpdate');
}

export async function buildWebhookSubscription(request) {
  return buildDefault(request, 'WebhookSubscription');
}
