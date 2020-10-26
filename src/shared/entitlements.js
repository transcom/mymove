import { get, isNull, sum, isEmpty } from 'lodash';
import { selectActiveOrLatestOrders } from 'shared/Entities/modules/orders';

export function selectEntitlements(rankEntitlement, hasDependents = false, spouseHasProGear = false) {
  if (!rankEntitlement) {
    return {};
  }
  const totalKey = hasDependents ? 'total_weight_self_plus_dependents' : 'total_weight_self';
  //  security/detect-object-injection
  const entitlement = {
    //  security/detect-object-injection
    weight: rankEntitlement[totalKey],
    pro_gear: rankEntitlement.pro_gear_weight,
    pro_gear_spouse: spouseHasProGear ? rankEntitlement.pro_gear_weight_spouse : 0,
  };
  entitlement.sum = sum([entitlement.weight, entitlement.pro_gear, entitlement.pro_gear_spouse]);
  return entitlement;
}

export function loadEntitlementsFromState(state) {
  // Temp fix until redux refactor finished - get orders from either entities or orders.currentOrders
  let orders = selectActiveOrLatestOrders(state);
  if (isEmpty(orders)) {
    return {};
  }
  const hasDependents = get(orders, 'has_dependents', null);
  const spouseHasProGear = get(orders, 'spouse_has_pro_gear', null);
  const weightAllotment = get(state, 'serviceMember.currentServiceMember.weight_allotment', null);
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(weightAllotment)) {
    return {};
  }
  return selectEntitlements(weightAllotment, hasDependents, spouseHasProGear);
}
