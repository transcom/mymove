// eslint-disable-next-line import/prefer-default-export
export const roleTypes = {
  TOO: 'transportation_ordering_officer',
  TIO: 'task_invoicing_officer',
  CUSTOMER: 'customer',
  CONTRACTING_OFFICER: 'contracting_officer',
  SERVICES_COUNSELOR: 'services_counselor',
  PRIME_SIMULATOR: 'prime_simulator',
  HQ: 'headquarters',
  CUSTOMER_SERVICE_REPRESENTATIVE: 'customer_service_representative',
  QAE: 'qae',
};

export const adminOfficeRoles = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
  { roleType: 'task_invoicing_officer', name: 'Task Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'services_counselor', name: 'Services Counselor' },
  { roleType: 'prime_simulator', name: 'Prime Simulator' },
  { roleType: 'headquarters', name: 'Headquarters' },
  { roleType: 'customer_service_representative', name: 'Customer Service Representative' },
  { roleType: 'qae', name: 'Quality Assurance Evaluator' },
];

export const officeRoles = [
  roleTypes.TOO,
  roleTypes.TIO,
  roleTypes.SERVICES_COUNSELOR,
  roleTypes.PRIME_SIMULATOR,
  roleTypes.HQ,
  roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
  roleTypes.QAE,
];
