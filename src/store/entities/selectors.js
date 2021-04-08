import { createSelector } from 'reselect';

import { profileStates } from 'constants/customerStates';
import { MOVE_STATUSES, NULL_UUID } from 'shared/constants';

/**
 * Use this file for selecting "slices" of state from Redux and for computed
 * properties given state. Selectors can be memoized for performance.
 * Documentation: https://github.com/reduxjs/reselect
 */

/** User */
export const selectLoggedInUser = (state) => {
  if (state.entities.user) return Object.values(state.entities.user)[0];
  return null;
};

/** Service Member */
export const selectServiceMemberFromLoggedInUser = (state) => {
  const user = selectLoggedInUser(state);
  if (!user || !user.service_member) return null;
  return state.entities.serviceMembers?.[`${user.service_member}`] || null;
};

export const selectCurrentDutyStation = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  return serviceMember?.current_station || null;
};

export const selectServiceMemberProfileState = createSelector(selectServiceMemberFromLoggedInUser, (serviceMember) => {
  if (!serviceMember) return profileStates.EMPTY_PROFILE;

  /* eslint-disable camelcase */
  const {
    rank,
    edipi,
    affiliation,
    first_name,
    last_name,
    telephone,
    personal_email,
    phone_is_preferred,
    email_is_preferred,
    current_station,
    residential_address,
    backup_mailing_address,
    backup_contacts,
  } = serviceMember;

  if (!rank || !edipi || !affiliation) return profileStates.EMPTY_PROFILE;
  if (!first_name || !last_name) return profileStates.DOD_INFO_COMPLETE;
  if (!telephone || !personal_email || !(phone_is_preferred || email_is_preferred)) return profileStates.NAME_COMPLETE;
  if (!current_station || !current_station.id || current_station.id === NULL_UUID)
    return profileStates.CONTACT_INFO_COMPLETE;
  if (!residential_address) return profileStates.DUTY_STATION_COMPLETE;
  if (!backup_mailing_address) return profileStates.ADDRESS_COMPLETE;
  if (!backup_contacts || !backup_contacts.length) return profileStates.BACKUP_ADDRESS_COMPLETE;
  return profileStates.BACKUP_CONTACTS_COMPLETE;
  /* eslint-enable camelcase */
});

// TODO: this is similar to service_member.isProfileComplete and we should figure out how to use just one if possible
export const selectIsProfileComplete = createSelector(
  selectServiceMemberFromLoggedInUser,
  (serviceMember) =>
    !!(
      serviceMember &&
      serviceMember.rank &&
      serviceMember.edipi &&
      serviceMember.affiliation &&
      serviceMember.first_name &&
      serviceMember.last_name &&
      serviceMember.telephone &&
      serviceMember.personal_email &&
      serviceMember.current_station?.id &&
      serviceMember.residential_address?.postal_code &&
      serviceMember.backup_mailing_address?.postal_code &&
      serviceMember.backup_contacts?.length > 0
    ),
);

/** Backup Contacts */
export const selectBackupContacts = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const backupContactIds = serviceMember?.backup_contacts || [];
  return backupContactIds.map((id) => state.entities.backupContacts?.[`${id}`]);
};

/** Orders */
export const selectOrdersForLoggedInUser = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const ordersIds = serviceMember?.orders || [];
  return ordersIds.map((id) => state.entities.orders?.[`${id}`]);
};

export const selectCurrentOrders = (state) => {
  const orders = selectOrdersForLoggedInUser(state);
  const [activeOrders] = orders.filter(
    (o) => ['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED'].indexOf(o?.status) > -1,
  );

  return activeOrders || orders[0] || null;
};

export const selectOrdersById = (state, id) => {
  return state.entities.orders?.[`${id}`] || null;
};

export const selectUploadsForCurrentOrders = (state) => {
  const orders = selectCurrentOrders(state);
  return orders ? orders.uploaded_orders?.uploads : [];
};

/** Moves */
export const selectMovesForLoggedInUser = (state) => {
  const orders = selectOrdersForLoggedInUser(state);
  const moves = orders?.reduce((prev, cur) => {
    return prev.concat(cur.moves?.map((id) => state.entities.moves?.[`${id}`]) || []);
  }, []);

  return moves;
};

export const selectMovesForCurrentOrders = (state) => {
  const activeOrders = selectCurrentOrders(state);
  const moveIds = activeOrders?.moves || [];
  return moveIds.map((id) => state.entities.moves?.[`${id}`]);
};

export const selectCurrentMove = (state) => {
  const moves = selectMovesForCurrentOrders(state);
  const [activeMove] = moves.filter(
    (m) => ['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED'].indexOf(m?.status) > -1,
  );
  return activeMove || moves[0] || null;
};

export const selectMoveIsApproved = createSelector(selectCurrentMove, (move) => move?.status === 'APPROVED');

export const selectMoveIsInDraft = createSelector(selectCurrentMove, (move) => move?.status === MOVE_STATUSES.DRAFT);

export const selectHasCanceledMove = createSelector(selectMovesForLoggedInUser, (moves) =>
  moves.some((m) => m.status === 'CANCELED'),
);

export const selectMoveType = createSelector(selectCurrentMove, (move) => move?.selected_move_type);

/** MTO Shipments */
export const selectMTOShipmentsForCurrentMove = (state) => {
  const currentMove = selectCurrentMove(state);
  return Object.values(state.entities.mtoShipments)?.filter((m) => m.moveTaskOrderID === currentMove?.id);
};

export function selectMTOShipmentById(state, id) {
  return state.entities?.mtoShipments?.[`${id}`] || null;
}

/** PPMs */
export const selectPPMForMove = (state, moveId) => {
  const ppmForMove = Object.values(state.entities.personallyProcuredMoves).find((ppm) => ppm.move_id === moveId);
  if (['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'].indexOf(ppmForMove?.status) > -1) {
    return ppmForMove;
  }
  return null;
};

export const selectCurrentPPM = (state) => {
  const move = selectCurrentMove(state);
  return selectPPMForMove(state, move?.id);
};

export const selectHasCurrentPPM = (state) => {
  return !!selectCurrentPPM(state);
};

export function selectPPMEstimateRange(state) {
  return state.entities?.ppmEstimateRanges?.undefined || null;
}

export function selectPPMSitEstimate(state) {
  return state.entities?.ppmSitEstimate?.undefined?.estimate || null;
}

export function selectReimbursementById(state, reimbursementId) {
  return state.entities?.reimbursements?.[`${reimbursementId}`] || null;
}

export const selectEntitlementsForLoggedInUser = createSelector(
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  (serviceMember, orders) => {
    const entitlement = {
      pro_gear: serviceMember.weight_allotment?.pro_gear_weight,
      pro_gear_spouse: orders.spouse_has_pro_gear ? serviceMember.weight_allotment?.pro_gear_weight_spouse : 0,
    };

    if (orders.has_dependents) {
      entitlement.weight = serviceMember.weight_allotment?.total_weight_self_plus_dependents;
    } else {
      entitlement.weight = serviceMember.weight_allotment?.total_weight_self;
    }

    entitlement.sum = [entitlement.weight, entitlement.pro_gear, entitlement.pro_gear_spouse].reduce(
      (acc, num) => acc + num,
      0,
    );

    return entitlement;
  },
);
