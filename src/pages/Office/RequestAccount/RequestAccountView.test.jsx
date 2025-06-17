import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import RequestAccountView from './RequestAccountView';

const mockRolesWithPrivs = [
  { roleType: 'task_ordering_officer', roleName: 'Task Ordering Officer' },
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

describe('RequestAccountView', () => {
  it('renders the form and role checkboxes', () => {
    render(
      <RequestAccountView
        serverError={null}
        onCancel={jest.fn()}
        onSubmit={jest.fn()}
        initialValues={initialValues}
        rolesWithPrivs={mockRolesWithPrivs}
        privileges={mockPrivileges}
      />,
    );
    // Form header
    expect(screen.getByRole('heading', { name: /Request Office Account/i })).toBeInTheDocument();
    // Role checkboxes
    expect(screen.getByLabelText('Task Ordering Officer')).toBeInTheDocument();
    expect(screen.getByLabelText('Headquarters')).toBeInTheDocument();
  });

  it('calls onCancel when Cancel is clicked', () => {
    const onCancel = jest.fn();
    render(
      <RequestAccountView
        serverError={null}
        onCancel={onCancel}
        onSubmit={jest.fn()}
        initialValues={initialValues}
        rolesWithPrivs={mockRolesWithPrivs}
        privileges={mockPrivileges}
      />,
    );
    fireEvent.click(screen.getByRole('button', { name: /Cancel/i }));
    expect(onCancel).toHaveBeenCalled();
  });
});
