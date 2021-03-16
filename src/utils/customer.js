/* eslint-disable import/prefer-default-export */
import { profileStates } from 'constants/customerStates';
import { generalRoutes, customerRoutes } from 'constants/routes';

export const findNextServiceMemberStep = (profileState) => {
  switch (profileState) {
    case profileStates.EMPTY_PROFILE:
      return customerRoutes.CONUS_OCONUS_PATH;
    case profileStates.DOD_INFO_COMPLETE:
      return customerRoutes.NAME_PATH;
    case profileStates.NAME_COMPLETE:
      return customerRoutes.CONTACT_INFO_PATH;
    case profileStates.CONTACT_INFO_COMPLETE:
      return customerRoutes.CURRENT_DUTY_STATION_PATH;
    case profileStates.DUTY_STATION_COMPLETE:
      return customerRoutes.CURRENT_ADDRESS_PATH;
    case profileStates.ADDRESS_COMPLETE:
      return customerRoutes.BACKUP_ADDRESS_PATH;
    case profileStates.BACKUP_ADDRESS_COMPLETE:
      return customerRoutes.BACKUP_CONTACTS_PATH;
    default:
      return generalRoutes.HOME_PATH;
  }
};
