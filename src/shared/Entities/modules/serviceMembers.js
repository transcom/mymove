import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
const loadBackupContactsLabel = 'ServiceMember.loadBackupContacts';

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
