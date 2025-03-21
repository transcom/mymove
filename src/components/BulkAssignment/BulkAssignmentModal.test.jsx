import React from 'react';
import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { BulkAssignmentModal } from 'components/BulkAssignment/BulkAssignmentModal';
import { QUEUE_TYPES } from 'constants/queues';
import { MockProviders } from 'testUtils';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

const bulkAssignmentData = {
  availableOfficeUsers: [
    {
      firstName: 'sc',
      lastName: 'user',
      officeUserId: '045c3048-df9a-4d44-88ed-8cd6e2100e08',
      workload: 1,
    },
    {
      firstName: 'test1',
      lastName: 'person',
      officeUserId: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42a',
    },
  ],
  bulkAssignmentMoveIDs: [
    'b3baf6ce-f43b-437c-85be-e1145c0ddb96',
    '962ce8d2-03a2-435c-94ca-6b9ef6c226c1',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed3',
  ],
};

jest.mock('services/ghcApi', () => ({
  getBulkAssignmentData: jest.fn().mockImplementation(() => Promise.resolve(bulkAssignmentData)),
}));

describe('BulkAssignmentModal', () => {
  it('renders the component', async () => {
    render(
      <MockProviders>
        <BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { level: 3, name: 'Bulk Assignment (3)' })).toBeInTheDocument();
  });

  it('shows cancel confirmation modal when close icon is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);
    await screen.findByRole('button', { name: 'Cancel' });

    await act(async () => {
      expect(await screen.getByText('person, test1')).toBeInTheDocument();
      const assignment = await screen.getAllByTestId('assignment')[0];
      await userEvent.type(assignment, '1');
    });

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(screen.getByTestId('cancelModalYes')).toBeInTheDocument();
  });

  it('does not show cancel confirmation if form is unchanged and cancel is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('shows cancel confirmation modal when the Cancel button is click if form has changed', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await act(async () => {
      expect(await screen.getByText('user, sc')).toBeInTheDocument();
      const assignment = await screen.getAllByTestId('assignment')[0];
      await userEvent.type(assignment, '1');
    });

    await userEvent.click(cancelButton);
    expect(screen.getByTestId('cancelModalYes')).toBeInTheDocument();
  });

  it('disables the save button if form is unchanged', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);
    const saveButton = await screen.findByTestId('modalSubmitButton');
    expect(saveButton).toBeDisabled();
    await userEvent.click(saveButton);
    expect(onSubmit).toHaveBeenCalledTimes(0);
  });

  it('calls the submit function when Save button is clicked if form is changed', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);

    await screen.findByRole('button', { name: 'Cancel' });

    await act(async () => {
      expect(await screen.getByText('person, test1')).toBeInTheDocument();
      const assignment = await screen.getAllByTestId('assignment')[0];
      await userEvent.type(assignment, '1');
    });
    const saveButton = await screen.findByTestId('modalSubmitButton');
    await userEvent.click(saveButton);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('renders the user data', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    const userTable = await screen.findByRole('table');
    expect(userTable).toBeInTheDocument();
    expect(screen.getByText('User')).toBeInTheDocument();
    expect(screen.getByText('Workload')).toBeInTheDocument();
    expect(screen.getByText('Assignment')).toBeInTheDocument();
    await act(async () => {
      expect(await screen.getByText('user, sc')).toBeInTheDocument();
    });
    expect(screen.getAllByTestId('bulkAssignmentUserWorkload')[0]).toHaveTextContent('1');
  });

  it('equal assign button splits assignment as equally as possible', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    await screen.findByRole('table');
    const equalAssignButton = await screen.getByTestId('modalEqualAssignButton');
    await userEvent.click(equalAssignButton);
    const row1 = await screen.getAllByTestId('assignment')[0];
    const row2 = await screen.getAllByTestId('assignment')[1];
    expect(row1.value).toEqual('2');
    expect(row2.value).toEqual('1');
  });

  it('select/deselect all checkbox works', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    await screen.findByRole('table');
    const selectDeselectAllButton = await screen.getByTestId('selectDeselectAllButton');
    const row1 = await screen.getAllByTestId('bulkAssignmentUserCheckbox')[0];
    const row2 = await screen.getAllByTestId('bulkAssignmentUserCheckbox')[1];

    expect(row1.checked).toEqual(true);
    expect(row2.checked).toEqual(true);
    expect(selectDeselectAllButton).toBeChecked();

    await userEvent.click(selectDeselectAllButton);
    expect(selectDeselectAllButton).not.toBeChecked();
    expect(row1.checked).toEqual(false);
    expect(row2.checked).toEqual(false);

    await userEvent.click(selectDeselectAllButton);
    expect(selectDeselectAllButton).toBeChecked();
    expect(row1.checked).toEqual(true);
    expect(row2.checked).toEqual(true);
  });

  it('submits the bulk assignment data', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    const userTable = await screen.findByRole('table');
    expect(userTable).toBeInTheDocument();
    expect(screen.getByText('User')).toBeInTheDocument();
    expect(screen.getByText('Workload')).toBeInTheDocument();
    expect(screen.getByText('Assignment')).toBeInTheDocument();
    await act(async () => {
      expect(await screen.getByText('user, sc')).toBeInTheDocument();
      const assignment = await screen.getAllByTestId('assignment')[0];
      await userEvent.type(assignment, '1');
    });
    expect(screen.getAllByTestId('bulkAssignmentUserWorkload')[0]).toHaveTextContent('1');

    const saveButton = await screen.getByTestId('modalSubmitButton');
    await userEvent.click(saveButton);
    await waitFor(() => {
      const payload = {
        bulkAssignmentSavePayload: {
          moveData: [
            'b3baf6ce-f43b-437c-85be-e1145c0ddb96',
            '962ce8d2-03a2-435c-94ca-6b9ef6c226c1',
            'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed3',
          ],
          userData: [
            {
              ID: '045c3048-df9a-4d44-88ed-8cd6e2100e08',
              moveAssignments: 1,
            },
          ],
        },
      };

      expect(onSubmit).toHaveBeenCalledWith(payload);
    });
  });

  it('closes the modal when the close is confirmed', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);

    await screen.findByRole('button', { name: 'Cancel' });

    await act(async () => {
      expect(await screen.getByText('person, test1')).toBeInTheDocument();
      const assignment = await screen.getAllByTestId('assignment')[0];
      await userEvent.type(assignment, '1');
    });
    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    const confirmButton = await screen.findByTestId('cancelModalYes');
    await userEvent.click(confirmButton);

    expect(onClose).toHaveBeenCalledTimes(2);
  });

  it('close confirmation goes away when clicking no', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);

    await screen.findByRole('button', { name: 'Cancel' });

    await act(async () => {
      expect(await screen.getByText('person, test1')).toBeInTheDocument();
      const assignment = await screen.getAllByTestId('assignment')[0];
      await userEvent.type(assignment, '1');
    });
    const closeButton = await screen.findByTestId('modalCloseButton');
    await userEvent.click(closeButton);

    const cancelModalNo = await screen.findByTestId('cancelModalNo');
    await userEvent.click(cancelModalNo);

    const confirmButton = await screen.queryByTestId('cancelModalYes');
    expect(confirmButton).not.toBeInTheDocument();
  });
});
