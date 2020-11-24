export const UPDATE_SERVICE_MEMBER = 'UPDATE_SERVICE_MEMBER';
export const UPDATE_BACKUP_CONTACT = 'UPDATE_BACKUP_CONTACT';

export const updateServiceMember = (payload) => ({
  type: UPDATE_SERVICE_MEMBER,
  payload,
});

export const updateBackupContact = (payload) => ({
  type: UPDATE_BACKUP_CONTACT,
  payload,
});
