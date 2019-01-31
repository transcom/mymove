import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
const loadBackupContactsLabel = 'ServiceMember.loadBackupContacts';
const updateBackupContactLabel = 'ServiceMember.updateBackupContact';

export function loadBackupContacts(serviceMemberId) {
  const label = loadBackupContactsLabel;
  const swaggerTag = 'backup_contacts.indexServiceMemberBackupContacts';
  return swaggerRequest(getClient, swaggerTag, { serviceMemberId }, { label });
}

export function selectBackupContactForServiceMember(state, serviceMemberId) {
  return Object.values(state.entities.backupContacts).find(backupContact => {
    return backupContact.service_member_id === serviceMemberId;
  });
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
