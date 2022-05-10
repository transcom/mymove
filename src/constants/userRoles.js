// eslint-disable-next-line import/prefer-default-export
export const roleTypes = {
  PPM: 'ppm_office_users',
  TOO: 'transportation_ordering_officer',
  TIO: 'transportation_invoicing_officer',
  CUSTOMER: 'customer',
  CONTRACTING_OFFICER: 'contracting_officer',
  SERVICES_COUNSELOR: 'services_counselor',
  PRIME_SIMULATOR: 'prime_simulator',
  QAE_CSR: 'qae_csr',
};

export const adminOfficeRoles = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
  { roleType: 'transportation_invoicing_officer', name: 'Transportation Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'ppm_office_users', name: 'PPM Office Users' },
  { roleType: 'services_counselor', name: 'Services Counselor' },
  { roleType: 'prime_simulator', name: 'Prime Simulator' },
  { roleType: 'qae_csr', name: 'Quality Assurance and Customer Support' },
];

export const officeRoles = [
  roleTypes.PPM,
  roleTypes.TOO,
  roleTypes.TIO,
  roleTypes.SERVICES_COUNSELOR,
  roleTypes.PRIME_SIMULATOR,
  roleTypes.QAE_CSR,
];
