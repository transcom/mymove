import {
  APPROVAL_REQUEST_TYPES_DESTINATION_REQUESTS_DISPLAY,
  APPROVAL_REQUEST_TYPES_TASK_ORDER_DISPLAY,
} from 'constants/approvalRequestTypes';
import { QUEUE_TYPES } from 'constants/queues';
import { deleteAssignedOfficeUserForMove, updateAssignedOfficeUserForMove } from 'services/ghcApi';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

export function handleQueueAssignment(moveID, officeUserId, queueType) {
  if (officeUserId === DEFAULT_EMPTY_VALUE) deleteAssignedOfficeUserForMove({ moveID, queueType });
  else updateAssignedOfficeUserForMove({ moveID, officeUserId, queueType });
}

export function getQueue(queueName) {
  // Normalize the queueName to lowercase for case-insensitive matching
  const normalizedQueueName = queueName?.toLowerCase()?.trim() || '';

  // Define the mapping of queue names to queue types
  const queueMappings = {
    counseling: QUEUE_TYPES.COUNSELING,
    'ppm-closeout': QUEUE_TYPES.CLOSEOUT,
    'move-queue': QUEUE_TYPES.TASK_ORDER,
    'payment-requests': QUEUE_TYPES.PAYMENT_REQUEST,
    'destination-requests': QUEUE_TYPES.DESTINATION_REQUESTS,
  };

  // Check if the queueName exists in the mappings
  if (!Object.keys(queueMappings).includes(normalizedQueueName)) {
    throw new Error(`Invalid queue name: ${queueName}. Valid options are: ${Object.keys(queueMappings).join(', ')}`);
  }

  // Return the corresponding queue type
  return queueMappings[normalizedQueueName];
}

export const formatApprovalRequestTypes = (queueType, requestTypes) => {
  const requestLabels = {};

  if (getQueue(queueType) === QUEUE_TYPES.TASK_ORDER) {
    requestTypes.forEach((request) => {
      if (APPROVAL_REQUEST_TYPES_TASK_ORDER_DISPLAY[request]) {
        requestLabels[APPROVAL_REQUEST_TYPES_TASK_ORDER_DISPLAY[request]] = true;
      }
    });
  }

  if (getQueue(queueType) === QUEUE_TYPES.DESTINATION_REQUESTS) {
    requestTypes.forEach((request) => {
      if (APPROVAL_REQUEST_TYPES_DESTINATION_REQUESTS_DISPLAY[request]) {
        requestLabels[APPROVAL_REQUEST_TYPES_DESTINATION_REQUESTS_DISPLAY[request]] = true;
      }
    });
  }

  return Object.keys(requestLabels).sort().join(', ');
};
