// This file contains constants related to office user assignments in the move history.
export const ASSIGNMENT_IDS = {
  SERVICES_COUNSELOR: 'sc_counseling_assigned_id',
  CLOSEOUT_COUNSELOR: 'sc_closeout_assigned_id',
  TASK_ORDERING_OFFICER: 'too_task_order_assigned_id',
  TASK_INVOICING_OFFICER: 'tio_assigned_id',
  TASK_ORDERING_OFFICER_DESTINATION: 'too_destination_assigned_id',
};

export const ASSIGNMENT_NAMES = {
  SERVICES_COUNSELOR: {
    ASSIGNED: 'assigned_sc_counseling',
    RE_ASSIGNED: 're_assigned_sc_counseling',
  },
  CLOSEOUT_COUNSELOR: {
    ASSIGNED: 'assigned_sc_closeout',
    RE_ASSIGNED: 're_assigned_sc_closeout',
  },
  TASK_ORDERING_OFFICER: {
    ASSIGNED: 'assigned_too_task_order',
    RE_ASSIGNED: 're_assigned_too_task_order',
  },
  DESTINATION_TASK_ORDERING_OFFICER: {
    ASSIGNED: 'assigned_too_destination',
    RE_ASSIGNED: 're_assigned_too_destination',
  },
  TASK_INVOICING_OFFICER: {
    ASSIGNED: 'assigned_tio',
    RE_ASSIGNED: 're_assigned_tio',
  },
};
