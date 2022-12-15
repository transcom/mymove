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

export async function buildSpouseProGearMove(request) {
  return buildDefault(request, 'SpouseProGearMove');
}

export async function buildNeedsOrdersUser(request) {
  return buildDefault(request, 'NeedsOrdersUser');
}

export async function buildPPMInProgressMove(request) {
  return buildDefault(request, 'PPMInProgressMove');
}
