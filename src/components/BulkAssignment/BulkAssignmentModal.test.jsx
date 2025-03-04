import React from 'react';
import { act, render, screen, waitFor, fireEvent, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { BulkAssignmentModal } from 'components/BulkAssignment/BulkAssignmentModal';
import { QUEUE_TYPES } from 'constants/queues';
import { MockProviders } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

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
      workload: 0,
    },
    {
      firstName: 'test2',
      lastName: 'person',
      officeUserId: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42c',
      workload: 0,
    },
  ],
  bulkAssignmentMoveIDs: [
    'b3baf6ce-f43b-437c-85be-e1145c0ddb96',
    '962ce8d2-03a2-435c-94ca-6b9ef6c226c1',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed3',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed4',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed5',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed6',
  ],
};

jest.mock('services/ghcApi', () => ({
  getBulkAssignmentData: jest.fn().mockImplementation(() => Promise.resolve(bulkAssignmentData)),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('BulkAssignmentModal', () => {
  it('renders the component', async () => {
    render(
      <MockProviders>
        <BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { level: 3, name: 'Bulk Assignment (6)' })).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(
      <MockProviders>
        <BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />
      </MockProviders>,
    );

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('closes the modal when the Cancel button is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls the submit function when Save button is clicked', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} />);

    const saveButton = await screen.findByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('renders the user data', async () => {
    render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    const userTable = await screen.findByRole('table');
    expect(userTable).toBeInTheDocument();
    expect(screen.getByText('User')).toBeInTheDocument();
    expect(screen.getByText('Current Workload')).toBeInTheDocument();
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
    const row3 = await screen.getAllByTestId('assignment')[1];
    expect(row1.value).toEqual('2');
    expect(row2.value).toEqual('2');
    expect(row3.value).toEqual('2');
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
    expect(screen.getByText('Current Workload')).toBeInTheDocument();
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
            'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed4',
            'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed5',
            'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed6',
          ],
          userData: [
            {
              ID: '045c3048-df9a-4d44-88ed-8cd6e2100e08',
              moveAssignments: 1,
            },
            {
              ID: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42a',
              moveAssignments: 0,
            },
            {
              ID: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42c',
              moveAssignments: 0,
            },
          ],
        },
      };

      expect(onSubmit).toHaveBeenCalledWith(payload);
    });
  });

  it('displays bulk re-assignment when switch is toggled', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    // Check the state of the modal before toggling the modal state
    await waitFor(() => {
      expect(screen.getByLabelText('BulkAssignmentModeSwitch')).toBeInTheDocument();
      expect(screen.getByText('User')).toBeInTheDocument();
      expect(screen.getByText('Current Workload')).toBeInTheDocument();
      expect(screen.getByText('Assignment')).toBeInTheDocument();
      expect(screen.getByTestId('selectDeselectAllButton')).toBeInTheDocument();
      expect(screen.getByText('Equal Assign')).toBeInTheDocument();
    });

    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');
    // Click the switch inside act() to ensure React updates state
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // Check the state of the modal after toggling the modal state
    await waitFor(() => {
      expect(screen.getByText('Bulk Re-Assignment', { exact: false })).toBeInTheDocument();
      expect(screen.getByText('User')).toBeInTheDocument();
      expect(screen.getByText('Current Workload')).toBeInTheDocument();
      expect(screen.getByText('Assignment')).toBeInTheDocument();
      expect(screen.getByText('Re-assign User')).toBeInTheDocument();
      expect(screen.queryByTestId('selectDeselectAllButton')).not.toBeInTheDocument();
      expect(screen.queryByTestId('Equal Assign')).not.toBeInTheDocument();
    });
    // Click again to revert back
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // Check the state of the modal after reverting toggling the modal state
    await waitFor(() => {
      expect(screen.getByLabelText('BulkAssignmentModeSwitch')).toBeInTheDocument();
      expect(screen.getByText('User')).toBeInTheDocument();
      expect(screen.getByText('Current Workload')).toBeInTheDocument();
      expect(screen.getByText('Assignment')).toBeInTheDocument();
      expect(screen.getByTestId('selectDeselectAllButton')).toBeInTheDocument();
      expect(screen.getByText('Equal Assign')).toBeInTheDocument();
    });
  });

  it('only allows bulk re-assignment from one user at a time', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    const expectedLabels = [
      `Bulk Re-Assignment (${bulkAssignmentData.availableOfficeUsers[0].workload})`,
      `Bulk Assignment (${bulkAssignmentData.availableOfficeUsers[0].workload})`,
    ];

    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');
    // Click the switch inside act() to ensure React updates state
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // Check the state of the modal after toggling the modal state
    await waitFor(() => {
      expect(screen.getByText('Bulk Re-Assignment (0)', { exact: false })).toBeInTheDocument();
      expect(screen.getByText('User')).toBeInTheDocument();
      expect(screen.getByText('Current Workload')).toBeInTheDocument();
      expect(screen.getByText('Assignment')).toBeInTheDocument();
      expect(screen.getByText('Re-assign User')).toBeInTheDocument();
      expect(screen.queryByTestId('selectDeselectAllButton')).not.toBeInTheDocument();
      expect(screen.queryByTestId('Equal Assign')).not.toBeInTheDocument();
    });
    // Select a user to re-assign from
    const radios = screen.getAllByRole('radio');
    const radioToReAssign = radios[0];
    const assignmentBoxes = screen.getAllByRole('spinbutton');
    const reAssignBox = assignmentBoxes[0];

    await act(async () => {
      await fireEvent.click(radioToReAssign);
    });
    // Verify that assignment box is disabled
    expect(screen.getByText('Bulk Re-Assignment (1)', { exact: false })).toBeInTheDocument();
    expect(radioToReAssign).toBeChecked();
    expect(reAssignBox).toBeDisabled();
    expect(assignmentBoxes[1]).toBeEnabled();
    expect(assignmentBoxes[2]).toBeEnabled();
  });

  it('cannot save if more reassignments are made than available', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');
    // Click the switch inside act() to ensure React updates state
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // Check the state of the modal after toggling the modal state
    await waitFor(() => {
      expect(screen.getByText('Bulk Re-Assignment (0)', { exact: false })).toBeInTheDocument();
      expect(screen.getByText('User')).toBeInTheDocument();
      expect(screen.getByText('Current Workload')).toBeInTheDocument();
      expect(screen.getByText('Assignment')).toBeInTheDocument();
      expect(screen.getByText('Re-assign User')).toBeInTheDocument();
      expect(screen.queryByTestId('selectDeselectAllButton')).not.toBeInTheDocument();
      expect(screen.queryByTestId('Equal Assign')).not.toBeInTheDocument();
    });
    // Select a user to re-assign from
    const radios = screen.getAllByRole('radio');
    const radioToReAssign = radios[0];
    const assignmentBoxes = screen.getAllByRole('spinbutton');
    const reAssignBox = assignmentBoxes[1];

    await act(async () => {
      await fireEvent.click(radioToReAssign);
    });
    expect(screen.getByText('Bulk Re-Assignment (1)', { exact: false })).toBeInTheDocument();
    // Try to re-assign 2 moves
    await act(async () => {
      await userEvent.type(reAssignBox, '2');
    });

    const saveButton = await screen.getByTestId('modalSubmitButton');
    await userEvent.click(saveButton);

    expect(screen.getByText('Cannot assign more moves than are available.')).toBeInTheDocument();
  });
});
