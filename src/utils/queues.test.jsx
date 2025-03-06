import { handleQueueAssignment } from './queues';

import { deleteAssignedOfficeUserForMove, updateAssignedOfficeUserForMove } from 'services/ghcApi';
import { QUEUE_TYPES } from 'constants/queues';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

jest.mock('services/ghcApi', () => ({
  deleteAssignedOfficeUserForMove: jest.fn(),
  updateAssignedOfficeUserForMove: jest.fn(),
}));

describe('handleQueueAssignment', () => {
  const moveID = 'PHISH4';
  const queueType = QUEUE_TYPES.COUNSELING;

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('calls deleteAssignedOfficeUserForMove when officeUserId is DEFAULT_EMPTY_VALUE', () => {
    const officeUserId = DEFAULT_EMPTY_VALUE;

    handleQueueAssignment(moveID, officeUserId, queueType);

    expect(deleteAssignedOfficeUserForMove).toHaveBeenCalledWith({ moveID, queueType });

    expect(updateAssignedOfficeUserForMove).not.toHaveBeenCalled();
  });

  it('calls updateAssignedOfficeUserForMove when officeUserId is not DEFAULT_EMPTY_VALUE', () => {
    const officeUserId = '3466';

    handleQueueAssignment(moveID, officeUserId, queueType);

    expect(updateAssignedOfficeUserForMove).toHaveBeenCalledWith({ moveID, officeUserId, queueType });

    expect(deleteAssignedOfficeUserForMove).not.toHaveBeenCalled();
  });
});
