import { deleteAssignedOfficeUserForMove, updateAssignedOfficeUserForMove } from 'services/ghcApi';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

const handleQueueAssignment = (moveID, officeUserId, roleType) => {
  if (officeUserId === DEFAULT_EMPTY_VALUE) deleteAssignedOfficeUserForMove({ moveID, roleType });
  else updateAssignedOfficeUserForMove({ moveID, officeUserId, roleType });
};

export default handleQueueAssignment;
