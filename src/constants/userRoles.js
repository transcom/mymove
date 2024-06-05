import { isBooleanFlagEnabled } from 'utils/featureFlags';

// eslint-disable-next-line import/prefer-default-export
export const roleTypes = {
  TOO: 'transportation_ordering_officer',
  TIO: 'task_invoicing_officer',
  CUSTOMER: 'customer',
  CONTRACTING_OFFICER: 'contracting_officer',
  SERVICES_COUNSELOR: 'services_counselor',
  PRIME_SIMULATOR: 'prime_simulator',
  QAE_CSR: 'qae_csr',
  HQ: 'headquarters',
};

const hqRoleFF = await isBooleanFlagEnabled('headquarters_role');

export const adminOfficeRoles = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
  { roleType: 'task_invoicing_officer', name: 'Task Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'services_counselor', name: 'Services Counselor' },
  { roleType: 'prime_simulator', name: 'Prime Simulator' },
  { roleType: 'qae_csr', name: 'Quality Assurance and Customer Support' },
  hqRoleFF ? { roleType: 'headquarters', name: 'Headquarters' } : null,
];

export const officeRoles = [
  roleTypes.TOO,
  roleTypes.TIO,
  roleTypes.SERVICES_COUNSELOR,
  roleTypes.PRIME_SIMULATOR,
  roleTypes.QAE_CSR,
  hqRoleFF ? roleTypes.HQ : null,
];
