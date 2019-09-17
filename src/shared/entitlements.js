import { get, isNull, sum } from 'lodash';

export function selectEntitlements(rankEntitlement, hasDependents = false, spouseHasProGear = false) {
  if (!rankEntitlement) {
    return {};
  }
  const totalKey = hasDependents ? 'total_weight_self_plus_dependents' : 'total_weight_self';
  // eslint-disable-next-line security/detect-object-injection
  const entitlement = {
    // eslint-disable-next-line security/detect-object-injection
    weight: rankEntitlement[totalKey],
    pro_gear: rankEntitlement.pro_gear_weight,
    pro_gear_spouse: spouseHasProGear ? rankEntitlement.pro_gear_weight_spouse : 0,
  };
  entitlement.sum = sum([entitlement.weight, entitlement.pro_gear, entitlement.pro_gear_spouse]);
  return entitlement;
}

export function loadEntitlementsFromState(state) {
  const hasDependents = get(state, 'orders.currentOrders.has_dependents', null);
  const spouseHasProGear = get(state, 'orders.currentOrders.spouse_has_pro_gear', null);
  const weightAllotment = get(state, 'serviceMember.currentServiceMember.weight_allotment', null);
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(weightAllotment)) {
    return {};
  }
  return selectEntitlements(weightAllotment, hasDependents, spouseHasProGear);
}
