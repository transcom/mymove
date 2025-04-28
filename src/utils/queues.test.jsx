import { handleQueueAssignment, getQueue, formatApprovalRequestTypes } from './queues';

import { deleteAssignedOfficeUserForMove, updateAssignedOfficeUserForMove } from 'services/ghcApi';
import { QUEUE_TYPES } from 'constants/queues';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';
import {
  APPROVAL_REQUEST_TYPES_DESTINATION_REQUESTS_DISPLAY,
  APPROVAL_REQUEST_TYPES_TASK_ORDER_DISPLAY,
} from 'constants/approvalRequestTypes';

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

describe('getQueue', () => {
  it('should return the correct queue type for valid queue names', () => {
    expect(getQueue('counseling')).toBe(QUEUE_TYPES.COUNSELING);
    expect(getQueue('ppm-closeout')).toBe(QUEUE_TYPES.CLOSEOUT);
    expect(getQueue('move-queue')).toBe(QUEUE_TYPES.TASK_ORDER);
    expect(getQueue('payment-requests')).toBe(QUEUE_TYPES.PAYMENT_REQUEST);
    expect(getQueue('destination-requests')).toBe(QUEUE_TYPES.DESTINATION_REQUESTS);
  });

  it('should be case-insensitive and trim white space for queue names', () => {
    expect(getQueue(' COUNSELING ')).toBe(QUEUE_TYPES.COUNSELING);
    expect(getQueue('Ppm-Closeout')).toBe(QUEUE_TYPES.CLOSEOUT);
  });

  it('should throw an error for null or undefined queue names and empty queue names', () => {
    expect(() => getQueue('')).toThrow('Invalid queue name:');
    expect(() => getQueue(null)).toThrow('Invalid queue name: null');
    expect(() => getQueue(undefined)).toThrow('Invalid queue name: undefined');
    expect(() => getQueue('invalid-queue')).toThrow('Invalid queue name: invalid-queue');
  });
});

describe('formatApprovalRequestTypes', () => {
  it('properly maps Task Order approvalRequestTypes to their front end display', () => {
    const approvalRequestTypeStrings = Object.keys(APPROVAL_REQUEST_TYPES_TASK_ORDER_DISPLAY);

    approvalRequestTypeStrings.forEach((request) => {
      expect(formatApprovalRequestTypes('move-queue', [request])).toBe(
        APPROVAL_REQUEST_TYPES_TASK_ORDER_DISPLAY[request],
      );
    });
  });
  it('properly maps Destination Request approvalRequestTypes to their front end display', () => {
    const approvalRequestTypeStrings = Object.keys(APPROVAL_REQUEST_TYPES_DESTINATION_REQUESTS_DISPLAY);

    approvalRequestTypeStrings.forEach((request) => {
      expect(formatApprovalRequestTypes('destination-requests', [request])).toBe(
        APPROVAL_REQUEST_TYPES_DESTINATION_REQUESTS_DISPLAY[request],
      );
    });
  });
});
