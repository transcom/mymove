import React from 'react';

import RequestAccountForm from './RequestAccountForm';

import { MockProviders } from 'testUtils';

const mockRolesWithPrivs = [
  { roleType: 'task_ordering_officer', roleName: 'Task Ordering Officer' },
  { roleType: 'task_invoicing_officer', roleName: 'Task Invoicing Officer' },
  { roleType: 'contracting_officer', roleName: 'Contracting Officer' },
  { roleType: 'services_counselor', roleName: 'Services Counselor' },
  { roleType: 'qae', roleName: 'Quality Assurance Evaluator' },
  { roleType: 'customer_service_representative', roleName: 'Customer Service Representative' },
  { roleType: 'gsr', roleName: 'Government Surveillance Representative' },
  { roleType: 'headquarters', roleName: 'Headquarters' },
];

const mockPrivileges = [];

const initialValues = {
  officeAccountRequestFirstName: '',
  officeAccountRequestMiddleInitial: '',
  officeAccountRequestLastName: '',
  officeAccountRequestEmail: '',
  officeAccountRequestTelephone: '',
  officeAccountRequestEdipi: '',
  officeAccountRequestOtherUniqueId: '',
  officeAccountTransportationOffice: undefined,
};

export default {
  title: 'Office Components/RequestAccountForm',
  component: RequestAccountForm,
  parameters: { layout: 'fullscreen' },
};

export const Blank = () => (
  <MockProviders>
    <RequestAccountForm
      onCancel={() => {}}
      onSubmit={() => {}}
      initialValues={initialValues}
      rolesWithPrivs={mockRolesWithPrivs}
      privileges={mockPrivileges}
    />
  </MockProviders>
);
