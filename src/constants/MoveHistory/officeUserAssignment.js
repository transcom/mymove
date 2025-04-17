// This file contains constants related to office user assignments in the move history.
export const ASSIGNMENT_IDS = {
  SERVICE_COUNSELOR: 'sc_assigned_id',
  TASK_ORDERING_OFFICER: 'too_assigned_id',
  TASK_INVOICING_OFFICER: 'tio_assigned_id',
  TASK_ORDERING_OFFICER_DESTINATION: 'too_destination_assigned_id',
};

export const ASSIGNMENT_NAMES = {
  SERVICE_COUNSELOR: {
    ASSIGNED: 'assigned_sc',
    RE_ASSIGNED: 're_assigned_sc',
  },
  SERVICE_COUNSELOR_PPM: {
    ASSIGNED: 'assigned_sc_ppm',
    RE_ASSIGNED: 're_assigned_sc_ppm',
  },
  TASK_ORDERING_OFFICER: {
    ASSIGNED: 'assigned_too',
    RE_ASSIGNED: 're_assigned_too',
  },
  TASK_INVOICING_OFFICER: {
    ASSIGNED: 'assigned_tio',
    RE_ASSIGNED: 're_assigned_tio',
  },
};
