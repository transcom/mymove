// eslint-disable-next-line import/prefer-default-export
export const roleTypes = {
  TOO: 'task_ordering_officer',
  TIO: 'task_invoicing_officer',
  CUSTOMER: 'customer',
  CONTRACTING_OFFICER: 'contracting_officer',
  SERVICES_COUNSELOR: 'services_counselor',
  PRIME_SIMULATOR: 'prime_simulator',
  QAE: 'qae',
  HQ: 'headquarters',
  CUSTOMER_SERVICE_REPRESENTATIVE: 'customer_service_representative',
};

export const adminOfficeRoles = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'task_ordering_officer', name: 'Task Ordering Officer' },
  { roleType: 'task_invoicing_officer', name: 'Task Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'services_counselor', name: 'Services Counselor' },
  { roleType: 'prime_simulator', name: 'Prime Simulator' },
  { roleType: 'qae', name: 'Quality Assurance Evaluator' },
  { roleType: 'headquarters', name: 'Headquarters' },
  { roleType: 'customer_service_representative', name: 'Customer Service Representative' },
];

export const officeRoles = [
  roleTypes.TOO,
  roleTypes.TIO,
  roleTypes.SERVICES_COUNSELOR,
  roleTypes.PRIME_SIMULATOR,
  roleTypes.QAE,
  roleTypes.HQ,
  roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
];
