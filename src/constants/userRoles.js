// eslint-disable-next-line import/prefer-default-export
export const roleTypes = {
  TOO: 'task_ordering_officer',
  TIO: 'task_invoicing_officer',
  CUSTOMER: 'customer',
  CONTRACTING_OFFICER: 'contracting_officer',
  SERVICES_COUNSELOR: 'services_counselor',
  PRIME_SIMULATOR: 'prime_simulator',
  QAE: 'qae',
  CUSTOMER_SERVICE_REPRESENTATIVE: 'customer_service_representative',
  GSR: 'gsr',
  HQ: 'headquarters',
};

export const adminOfficeRoles = [
  { roleType: 'customer', name: 'Customer', abbv: '' },
  { roleType: 'task_ordering_officer', name: 'Task Ordering Officer', abbv: 'TOO' },
  { roleType: 'task_invoicing_officer', name: 'Task Invoicing Officer', abbv: 'TIO' },
  { roleType: 'contracting_officer', name: 'Contracting Officer', abbv: 'CO' },
  { roleType: 'services_counselor', name: 'Services Counselor', abbv: 'SC' },
  { roleType: 'prime_simulator', name: 'Prime Simulator', abbv: 'SC' },
  { roleType: 'qae', name: 'Quality Assurance Evaluator', abbv: 'QAE' },
  { roleType: 'customer_service_representative', name: 'Customer Service Representative', abbv: 'CSR' },
  { roleType: 'gsr', name: 'Government Surveillance Representative', abbv: 'GSR' },
  { roleType: 'headquarters', name: 'Headquarters', abbv: 'HQ' },
];

export const officeRoles = [
  roleTypes.TOO,
  roleTypes.TIO,
  roleTypes.SERVICES_COUNSELOR,
  roleTypes.PRIME_SIMULATOR,
  roleTypes.QAE,
  roleTypes.HQ,
  roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
  roleTypes.GSR,
  roleTypes.CONTRACTING_OFFICER,
];
