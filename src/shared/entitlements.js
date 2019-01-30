import { get, has, isNull, sum } from 'lodash';

const defaultStorageInTransitDays = 90;

export function getEntitlements(rank, hasDependents = false, spouseHasProGear = false) {
  if (!has(entitlements, rank)) {
    return null;
  }

  const totalKey = hasDependents ? 'total_weight_self_plus_dependents' : 'total_weight_self';
  // eslint-disable-next-line security/detect-object-injection
  const rankEntitlement = entitlements[rank];
  const entitlement = {
    // eslint-disable-next-line security/detect-object-injection
    weight: rankEntitlement[totalKey],
    pro_gear: rankEntitlement.pro_gear_weight,
    pro_gear_spouse: spouseHasProGear ? rankEntitlement.pro_gear_weight_spouse : 0,
    storage_in_transit: defaultStorageInTransitDays,
  };
  entitlement.sum = sum([entitlement.weight, entitlement.pro_gear, entitlement.pro_gear_spouse]);
  return entitlement;
}

export function loadEntitlementsFromState(state) {
  const hasDependents = get(state, 'orders.currentOrders.has_dependents', null);
  const spouseHasProGear = get(state, 'orders.currentOrders.spouse_has_pro_gear', null);
  const rank = get(state, 'serviceMember.currentServiceMember.rank', null);
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
}

/*
 * These entitlements are pulled from the move.mil source code
 * Source: https://github.com/deptofdefense/move.mil/blob/master/lib/data/entitlements.yml
 */
const entitlements = {
  ACADEMY_CADET_MIDSHIPMAN: {
    total_weight_self: 350,
    total_weight_self_plus_dependents: 350,
    pro_gear_weight: 0,
    pro_gear_weight_spouse: 0,
  },
  AVIATION_CADET: {
    total_weight_self: 7000,
    total_weight_self_plus_dependents: 8000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_1: {
    total_weight_self: 5000,
    total_weight_self_plus_dependents: 8000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_2: {
    total_weight_self: 5000,
    total_weight_self_plus_dependents: 8000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_3: {
    total_weight_self: 5000,
    total_weight_self_plus_dependents: 8000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_4: {
    total_weight_self: 7000,
    total_weight_self_plus_dependents: 8000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_5: {
    total_weight_self: 7000,
    total_weight_self_plus_dependents: 9000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_6: {
    total_weight_self: 8000,
    total_weight_self_plus_dependents: 11000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_7: {
    total_weight_self: 11000,
    total_weight_self_plus_dependents: 13000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_8: {
    total_weight_self: 12000,
    total_weight_self_plus_dependents: 14000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  E_9: {
    total_weight_self: 13000,
    total_weight_self_plus_dependents: 15000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_1_W_1_ACADEMY_GRADUATE: {
    total_weight_self: 10000,
    total_weight_self_plus_dependents: 12000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_2_W_2: {
    total_weight_self: 12500,
    total_weight_self_plus_dependents: 13500,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_3_W_3: {
    total_weight_self: 13000,
    total_weight_self_plus_dependents: 14500,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_4_W_4: {
    total_weight_self: 14000,
    total_weight_self_plus_dependents: 17000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_5_W_5: {
    total_weight_self: 16000,
    total_weight_self_plus_dependents: 17500,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_6: {
    total_weight_self: 18000,
    total_weight_self_plus_dependents: 18000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_7: {
    total_weight_self: 18000,
    total_weight_self_plus_dependents: 18000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_8: {
    total_weight_self: 18000,
    total_weight_self_plus_dependents: 18000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_9: {
    total_weight_self: 18000,
    total_weight_self_plus_dependents: 18000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  O_10: {
    total_weight_self: 18000,
    total_weight_self_plus_dependents: 18000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
  CIVILIAN_EMPLOYEE: {
    total_weight_self: 18000,
    total_weight_self_plus_dependents: 18000,
    pro_gear_weight: 2000,
    pro_gear_weight_spouse: 500,
  },
};
