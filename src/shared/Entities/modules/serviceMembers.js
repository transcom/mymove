import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { get } from 'lodash';
const loadBackupContactsLabel = 'ServiceMember.loadBackupContacts';
const updateBackupContactLabel = 'ServiceMember.updateBackupContact';
export const loadServiceMemberLabel = 'ServiceMember.loadServiceMember';
export const updateServiceMemberLabel = 'ServiceMember.updateServiceMember';

export function loadBackupContacts(serviceMemberId) {
  const label = loadBackupContactsLabel;
  const swaggerTag = 'backup_contacts.indexServiceMemberBackupContacts';
  return swaggerRequest(getClient, swaggerTag, { serviceMemberId }, { label });
}

export function updateBackupContact(backupContactId, backupContact) {
  const label = updateBackupContactLabel;
  const swaggerTag = 'backup_contacts.updateServiceMemberBackupContact';
  return swaggerRequest(
    getClient,
    swaggerTag,
    { backupContactId, updateServiceMemberBackupContactPayload: backupContact },
    { label },
  );
}

export function loadServiceMember(serviceMemberId) {
  const label = loadServiceMemberLabel;
  const swaggerTag = 'service_members.showServiceMember';
  return swaggerRequest(getClient, swaggerTag, { serviceMemberId }, { label });
}

export function updateServiceMember(serviceMemberId, serviceMember) {
  const label = updateServiceMemberLabel;
  const swaggerTag = 'service_members.patchServiceMember';
  return swaggerRequest(
    getClient,
    swaggerTag,
    { serviceMemberId, patchServiceMemberPayload: serviceMember },
    { label },
  );
}

export function selectServiceMember(state, serviceMemberId) {
  return get(state, `entities.serviceMembers.${serviceMemberId}`, {});
}

export function selectServiceMemberForOrders(state, ordersId) {
  const orders = get(state, `entities.orders.${ordersId}`);
  if (!orders) {
    return {};
  }
  const serviceMember = get(state, `entities.serviceMembers.${orders.service_member_id}`);
  return serviceMember || {};
}

export function selectServiceMemberForMove(state, moveId) {
  const move = get(state, `entities.moves.${moveId}`);
  if (!move) {
    return {};
  }
  return selectServiceMemberForOrders(state, move.orders_id);
}

export function selectBackupContactForServiceMember(state, serviceMemberId) {
  return Object.values(state.entities.backupContacts).find(backupContact => {
    return backupContact.service_member_id === serviceMemberId;
  });
}
