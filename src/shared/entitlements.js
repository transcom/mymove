import { get, isNull, sum, isEmpty } from 'lodash';

import { selectServiceMemberFromLoggedInUser, selectCurrentOrders } from 'store/entities/selectors';

export function selectEntitlements(ordersInfo, hasDependents = false, spouseHasProGear = false) {
  if (!ordersInfo) {
    return {};
  }
  const entitlement = {
    weight: ordersInfo.authorizedWeight,
    proGear: ordersInfo.entitlement.proGear,
    proGearSpouse: spouseHasProGear ? ordersInfo.entitlement.proGearSpouse : 0,
  };

  entitlement.sum = sum([entitlement.weight, entitlement.proGear, entitlement.proGearSpouse]);
  return entitlement;
}

export function loadEntitlementsFromState(state) {
  // Temp fix until redux refactor finished - get orders from either entities or orders.currentOrders
  const orders = selectCurrentOrders(state);
  if (isEmpty(orders)) {
    return {};
  }

  const hasDependents = get(orders, 'has_dependents', null);
  const spouseHasProGear = get(orders, 'spouse_has_pro_gear', null);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const weightAllotment = serviceMember?.weight_allotment || null;

  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(weightAllotment)) {
    return {};
  }

  return selectEntitlements(weightAllotment, hasDependents, spouseHasProGear);
}
