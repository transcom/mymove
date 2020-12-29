import { createSelector } from 'reselect';

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
  if (['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'].indexOf(ppmForMove?.status) > -1)
    return ppmForMove;
  return null;
};

export const selectCurrentPPM = (state) => {
  const move = selectCurrentMove(state);
  return selectPPMForMove(state, move?.id);
};

export const selectPPMEstimateRange = (state) => {
  return state.entities.ppmEstimateRanges?.undefined || null;
};

export const selectPPMSitEstimate = (state) => {
  return state.entities.ppmSitEstimate?.undefined?.estimate || null;
};
