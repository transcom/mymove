/* eslint-disable import/prefer-default-export */
import { profileStates } from 'constants/customerStates';

export const findNextServiceMemberStep = (serviceMemberId, profileState) => {
  const profilePathPrefix = `/service-member/${serviceMemberId}`;

  switch (profileState) {
    case profileStates.EMPTY_PROFILE:
      return `${profilePathPrefix}/conus-status`;
    case profileStates.DOD_INFO_COMPLETE:
      return `${profilePathPrefix}/name`;
    case profileStates.NAME_COMPLETE:
      return `${profilePathPrefix}/contact-info`;
    case profileStates.CONTACT_INFO_COMPLETE:
      return `${profilePathPrefix}/duty-station`;
    case profileStates.DUTY_STATION_COMPLETE:
      return `${profilePathPrefix}/residence-address`;
    case profileStates.ADDRESS_COMPLETE:
      return `${profilePathPrefix}/backup-mailing-address`;
    case profileStates.BACKUP_ADDRESS_COMPLETE:
      return `${profilePathPrefix}/backup-contacts`;
    default:
      return '/';
  }
};
