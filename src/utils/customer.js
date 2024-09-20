/* eslint-disable import/prefer-default-export */

import { profileStates } from 'constants/customerStates';
import { generalRoutes, customerRoutes } from 'constants/routes';

export const findNextServiceMemberStep = (profileState) => {
  switch (profileState) {
    case profileStates.VALIDATION_REQUIRED:
      return customerRoutes.VALIDATION_CODE_PATH;
    case profileStates.EMPTY_PROFILE:
      return customerRoutes.CONUS_OCONUS_PATH;
    case profileStates.DOD_INFO_COMPLETE:
      return customerRoutes.NAME_PATH;
    case profileStates.NAME_COMPLETE:
      return customerRoutes.CONTACT_INFO_PATH;
    case profileStates.CONTACT_INFO_COMPLETE:
      return customerRoutes.CURRENT_ADDRESS_PATH;
    case profileStates.ADDRESS_COMPLETE:
      return customerRoutes.BACKUP_ADDRESS_PATH;
    case profileStates.BACKUP_ADDRESS_COMPLETE:
      return customerRoutes.BACKUP_CONTACTS_PATH;
    default:
      return generalRoutes.HOME_PATH;
  }
};

export const generateUniqueDodid = () => {
  const prefix = 'SM';

  // Custom epoch start date (e.g., 2024-01-01), generates something like 1704067200000
  const customEpoch = new Date('2024-01-01').getTime();
  const now = Date.now();

  // Calculate milliseconds since custom epoch, then convert to an 8-digit integer
  const uniqueNumber = Math.floor((now - customEpoch) / 1000); // Dividing by 1000 to reduce to seconds

  // Convert the unique number to a string, ensuring it has 8 digits
  const uniqueStr = uniqueNumber.toString().slice(0, 8).padStart(8, '0');

  return prefix + uniqueStr;
};

export const generateUniqueEmplid = () => {
  const prefix = 'SM';
  const customEpoch = new Date('2024-01-01').getTime();
  const now = Date.now();
  const uniqueNumber = Math.floor((now - customEpoch) / 1000) % 100000; // Modulo 100000 ensures it's 5 digits
  const uniqueStr = uniqueNumber.toString().padStart(5, '0');
  return prefix + uniqueStr;
};
