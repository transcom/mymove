import handleQueueAssignment from './queues';

import { deleteAssignedOfficeUserForMove, updateAssignedOfficeUserForMove } from 'services/ghcApi';
import { roleTypes } from 'constants/userRoles';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

jest.mock('services/ghcApi', () => ({
  deleteAssignedOfficeUserForMove: jest.fn(),
  updateAssignedOfficeUserForMove: jest.fn(),
}));

describe('handleQueueAssignment', () => {
  const moveID = 'PHISH4';
  const roleType = roleTypes.SERVICES_COUNSELOR;

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('calls deleteAssignedOfficeUserForMove when officeUserId is DEFAULT_EMPTY_VALUE', () => {
    const officeUserId = DEFAULT_EMPTY_VALUE;

    handleQueueAssignment(moveID, officeUserId, roleType);

    expect(deleteAssignedOfficeUserForMove).toHaveBeenCalledWith({ moveID, roleType });

    expect(updateAssignedOfficeUserForMove).not.toHaveBeenCalled();
  });

  it('calls updateAssignedOfficeUserForMove when officeUserId is not DEFAULT_EMPTY_VALUE', () => {
    const officeUserId = '3466';

    handleQueueAssignment(moveID, officeUserId, roleType);

    expect(updateAssignedOfficeUserForMove).toHaveBeenCalledWith({ moveID, officeUserId, roleType });

    expect(deleteAssignedOfficeUserForMove).not.toHaveBeenCalled();
  });
});
