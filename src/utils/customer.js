/* eslint-disable import/prefer-default-export */
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';

export const findNextServiceMemberStep = (profileState) => {
  switch (profileState) {
    case profileStates.EMPTY_PROFILE:
      return customerRoutes.CONUS_OCONUS;
    case profileStates.DOD_INFO_COMPLETE:
      return customerRoutes.NAME;
    case profileStates.NAME_COMPLETE:
      return customerRoutes.CONTACT_INFO;
    case profileStates.CONTACT_INFO_COMPLETE:
      return customerRoutes.CURRENT_DUTY_STATION;
    case profileStates.DUTY_STATION_COMPLETE:
      return customerRoutes.CURRENT_ADDRESS;
    case profileStates.ADDRESS_COMPLETE:
      return customerRoutes.BACKUP_ADDRESS;
    case profileStates.BACKUP_ADDRESS_COMPLETE:
      return customerRoutes.BACKUP_CONTACTS;
    default:
      return '/';
  }
};
