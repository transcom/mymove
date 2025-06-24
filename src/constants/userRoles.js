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
  { roleType: roleTypes.CUSTOMER, name: 'Customer', abbv: '' },
  { roleType: roleTypes.TOO, name: 'Task Ordering Officer', abbv: 'TOO' },
  { roleType: roleTypes.TIO, name: 'Task Invoicing Officer', abbv: 'TIO' },
  { roleType: roleTypes.CONTRACTING_OFFICER, name: 'Contracting Officer', abbv: 'KO' },
  { roleType: roleTypes.SERVICES_COUNSELOR, name: 'Services Counselor', abbv: 'SC' },
  { roleType: roleTypes.PRIME_SIMULATOR, name: 'Prime Simulator', abbv: 'PRIME' },
  { roleType: roleTypes.QAE, name: 'Quality Assurance Evaluator', abbv: 'QAE' },
  { roleType: roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE, name: 'Customer Service Representative', abbv: 'CSR' },
  { roleType: roleTypes.GSR, name: 'Government Surveillance Representative', abbv: 'GSR' },
  { roleType: roleTypes.HQ, name: 'Headquarters', abbv: 'HQ' },
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
