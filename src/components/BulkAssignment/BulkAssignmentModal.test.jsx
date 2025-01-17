import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { BulkAssignmentModal } from 'components/BulkAssignment/BulkAssignmentModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

const data = {
  availableOfficeUsers: [
    {
      firstName: 'John',
      lastName: 'Snow',
      officeUserId: '123',
      workload: 0,
    },
    {
      firstName: 'Jane',
      lastName: 'Doe',
      officeUserId: '456',
      workload: 1,
    },
    {
      firstName: 'Jimmy',
      lastName: 'Page',
      officeUserId: '789',
      workload: 50,
    },
  ],
  bulkAssignmentMoveIDs: ['1', '2', '3', '4', '5'],
};

describe('BulkAssignmentModal', () => {
  it('renders the component', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} bulkAssignmentData={data} />);

    expect(await screen.findByRole('heading', { level: 3, name: 'Bulk Assignment (5)' })).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} bulkAssignmentData={data} />);

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('closes the modal when the Cancel button is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} bulkAssignmentData={data} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls the submit function when Save button is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} bulkAssignmentData={data} />);

    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('renders the user data', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} bulkAssignmentData={data} />);

    const userTable = await screen.findByRole('table');

    expect(userTable).toBeInTheDocument();
    expect(screen.getByText('User')).toBeInTheDocument();
    expect(screen.getByText('Workload')).toBeInTheDocument();
    expect(screen.getByText('Assignment')).toBeInTheDocument();

    expect(screen.getByText('Snow, John')).toBeInTheDocument();
    expect(screen.getAllByTestId('bulkAssignmentUserWorkload')[0]).toHaveTextContent('0');
  });
});
