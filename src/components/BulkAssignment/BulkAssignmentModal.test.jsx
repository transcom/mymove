import React from 'react';
import { act, render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { BulkAssignmentModal } from 'components/BulkAssignment/BulkAssignmentModal';
import { QUEUE_TYPES } from 'constants/queues';
import { MockProviders } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { milmoveLogger } from 'utils/milmoveLog';
import { getBulkAssignmentData } from 'services/ghcApi';

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
      workload: 6,
    },
    {
      firstName: 'test2',
      lastName: 'person',
      officeUserId: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42c',
      workload: 4,
    },
  ],
  bulkAssignmentMoveIDs: [
    'b3baf6ce-f43b-437c-85be-e1145c0ddb96',
    '962ce8d2-03a2-435c-94ca-6b9ef6c226c1',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed3',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed4',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed5',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed6',
    'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed7',
  ],
};

const checkAllAssignmentBoxes = (currentScreen, expectedNumber) => {
  const assignmentBoxesPreSwitch = currentScreen.getAllByRole('spinbutton');
  assignmentBoxesPreSwitch.forEach((assignmentBox) => {
    expect(assignmentBox.value).toEqual(`${expectedNumber}`);
  });
};

const checkCheckboxSelectionStatus = (currentScreen, expectedState) => {
  const userSelectionBoxSecondSwitch = currentScreen
    .getAllByRole('checkbox')
    .filter((checkbox) => checkbox.closest('cell'));
  userSelectionBoxSecondSwitch.forEach((checkbox) => {
    if (expectedState) {
      expect(checkbox).toBeChecked();
    } else {
      expect(checkbox).not.toBeChecked();
    }
  });
};

const checkModalElements = (currentScreen, isBulkReAssignmentMode, expectedMoveCount) => {
  const expectedString = isBulkReAssignmentMode
    ? `Bulk Re-Assignment (${expectedMoveCount})`
    : `Bulk Assignment (${expectedMoveCount})`;
  expect(currentScreen.getByText(expectedString, { exact: false })).toBeInTheDocument();
  expect(currentScreen.getByText('User')).toBeInTheDocument();
  expect(currentScreen.getByText('Current Workload')).toBeInTheDocument();
  expect(currentScreen.getByText('Assignment')).toBeInTheDocument();
  if (isBulkReAssignmentMode) {
    // Is in re-assignment mode
    expect(currentScreen.queryByText('Re-assign Workload')).toBeInTheDocument();
    expect(currentScreen.queryByTestId('selectDeselectAllButton')).not.toBeVisible();
    expect(currentScreen.getByText('Equal Assign')).not.toBeVisible();
  } else {
    // is not in re-assignment mode
    expect(currentScreen.queryByText('Re-assign Workload')).not.toBeInTheDocument();
    expect(currentScreen.queryByTestId('selectDeselectAllButton')).toBeVisible();
    expect(currentScreen.getByText('Equal Assign')).toBeVisible();
  }
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

    expect(await screen.findByRole('heading', { level: 3, name: 'Bulk Assignment (7)' })).toBeInTheDocument();
  });

  it('handles API error gracefully', async () => {
    // Override getBulkAssignmentData to throw an error
    getBulkAssignmentData.mockImplementationOnce(() => Promise.reject(new Error('API error')));

    // Spy on the logger to verify error handling
    jest.spyOn(milmoveLogger, 'error').mockImplementation(() => {});

    render(
      <MockProviders>
        <BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(milmoveLogger.error).toHaveBeenCalledWith(
        expect.stringContaining('Error fetching bulk assignment data:'),
        expect.any(Error),
      );
    });

    // Check that the save button is disabled since data loading failed
    const saveButton = await screen.findByTestId('modalSubmitButton');
    expect(saveButton).toBeDisabled();
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

    const cancelButton = await screen.findByRole('button', { name: 'Close' });

    await act(async () => {
      await userEvent.click(cancelButton);
    });

    await waitFor(() => {
      expect(onClose).toHaveBeenCalledTimes(1);
    });
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
    const assignmentBoxes = await screen.getAllByTestId('assignment');

    const row1 = assignmentBoxes[0];
    const row2 = assignmentBoxes[1];
    await waitFor(() => {
      expect(row1.value).toEqual('3');
      expect(row2.value).toEqual('2');
    });
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
            'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed7',
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

  it('only allows bulk re-assignment from one user at a time', async () => {
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
      checkModalElements(screen, true, 0);
    });
    // Select a user to re-assign from
    const radioButtons = screen.getAllByRole('radio');
    radioButtons.forEach((radioButton) => {
      expect(radioButton).not.toBeChecked();
    });
    const assignmentBoxes = screen.getAllByRole('spinbutton');
    assignmentBoxes.forEach((assignmentBox) => {
      expect(assignmentBox).toBeDisabled();
    });
    const radioToReAssign = radioButtons[0];

    await act(async () => {
      await fireEvent.click(radioToReAssign);
    });
    // Verify that assignment box is disabled
    await waitFor(() => {
      checkModalElements(screen, true, 1);
      expect(assignmentBoxes[1]).toBeEnabled();
      expect(assignmentBoxes[2]).toBeEnabled();
    });

    // select another user and verify that row's assignment box only is disabled
    const radioToReAssign2 = radioButtons[2];
    const reAssignBox2 = assignmentBoxes[2];

    await act(async () => {
      await fireEvent.click(radioToReAssign2);
    });

    await waitFor(() => {
      checkModalElements(screen, true, 4);
      expect(radioToReAssign2).toBeChecked();
      expect(reAssignBox2.value).toEqual('0');
      expect(reAssignBox2).toBeDisabled();
      expect(assignmentBoxes[0]).toBeEnabled();
      expect(assignmentBoxes[1]).toBeEnabled();
    });
  });

  it('cannot save if more reassignments are made than available', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    // Check the state of the modal after toggling the modal state
    await waitFor(() => {
      checkModalElements(screen, false, bulkAssignmentData.bulkAssignmentMoveIDs.length);
    });

    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');
    // Click the switch inside act() to ensure React updates state
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // Check the state of the modal after toggling the modal state
    await waitFor(() => {
      checkModalElements(screen, true, 0);
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

  it('does not persist unsaved assignment values while mode switching', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    const assignmentBoxesPreSwitch = screen.getAllByRole('spinbutton');
    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');

    // check initial state
    await waitFor(() => {
      checkAllAssignmentBoxes(screen, 0);
    });
    // type some values
    await waitFor(async () => {
      await userEvent.type(assignmentBoxesPreSwitch[0], '2');
      await userEvent.type(assignmentBoxesPreSwitch[1], '4');
      await userEvent.type(assignmentBoxesPreSwitch[2], '6');
    });

    // first switch to bulk re assignment
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    const assignmentBoxesFirstSwitch = screen.getAllByRole('spinbutton');

    await waitFor(() => {
      checkAllAssignmentBoxes(screen, 0);
    });
    await waitFor(async () => {
      await userEvent.type(assignmentBoxesFirstSwitch[0], '2');
      await userEvent.type(assignmentBoxesFirstSwitch[1], '4');
      await userEvent.type(assignmentBoxesFirstSwitch[2], '6');
    });

    // switch back to bulk assignment
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    await waitFor(() => {
      checkAllAssignmentBoxes(screen, 0);
    });

    // second switch to bulk re assignment
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    await waitFor(() => {
      checkAllAssignmentBoxes(screen, 0);
    });
  });
  it('keeps all checkboxes selected when switching back to BA from Bulk Re-Assignment', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    const userSelectionBoxPreSwitch = screen.getAllByRole('checkbox').filter((checkbox) => checkbox.closest('cell'));
    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');

    // check initial state
    await waitFor(() => {
      checkCheckboxSelectionStatus(screen, true);
    });
    // deselect a few boxes (all enabled by default)
    await waitFor(async () => {
      await userEvent.click(userSelectionBoxPreSwitch[0]);
      await userEvent.click(userSelectionBoxPreSwitch[2]);
    });

    // switch to bulk re assignment
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // should be in bulk re-assignment mode and checkboxes should not be visible
    await waitFor(() => {
      expect(bulkReAssignToggleSwitch).toBeChecked();
      checkCheckboxSelectionStatus(screen, true);
    });

    // switch back to bulk assignment
    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    // back in bulk assignment mode and all checkoxes are selected
    await waitFor(async () => {
      expect(bulkReAssignToggleSwitch).not.toBeChecked();
      checkCheckboxSelectionStatus(screen, true);
    });
  });
  it('equal assign does nothing when no users are selected', async () => {
    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    // Deselect all users using the select/deselect all checkbox
    const selectDeselectAllButton = await screen.getByTestId('selectDeselectAllButton');
    await userEvent.click(selectDeselectAllButton);

    // Click the equal assign button
    const equalAssignButton = await screen.getByTestId('modalEqualAssignButton');
    await userEvent.click(equalAssignButton);

    // Verify that all assignment inputs remain 0
    checkAllAssignmentBoxes(screen, 0);
  });
  it('resets all assignment boxes to 0 when the radio button is changed', async () => {
    // Force bulk re-assignment mode so the radio buttons are rendered.
    isBooleanFlagEnabled.mockResolvedValue(true);
    await act(async () => {
      render(<BulkAssignmentModal onSubmit={onSubmit} onClose={onClose} queueType={QUEUE_TYPES.COUNSELING} />);
    });

    const bulkReAssignToggleSwitch = screen.getByLabelText('BulkAssignmentModeSwitch');

    await act(async () => {
      await fireEvent.click(bulkReAssignToggleSwitch);
    });

    const reAssignUserRadio = screen.getAllByRole('radio');

    // Fill each assignment input with a non-zero value.
    const assignmentInputs = screen.getAllByTestId('assignment');

    await waitFor(async () => {
      expect(assignmentInputs[0]).toHaveValue(0);
      expect(assignmentInputs[1]).toHaveValue(0);
      expect(assignmentInputs[2]).toHaveValue(0);
    });

    await act(async () => {
      await fireEvent.click(reAssignUserRadio[0]);
    });

    await act(async () => {
      await userEvent.type(assignmentInputs[1], '2');
      await userEvent.type(assignmentInputs[2], '3');
    });

    await act(async () => {
      await fireEvent.click(reAssignUserRadio[1]);
    });

    const assignmentInputs2 = screen.getAllByTestId('assignment');

    await waitFor(async () => {
      expect(assignmentInputs2[0]).toHaveValue(0);
      expect(assignmentInputs2[1]).toHaveValue(0);
      expect(assignmentInputs2[2]).toHaveValue(0);
    });

    await act(async () => {
      await userEvent.click(bulkReAssignToggleSwitch);
    });

    await waitFor(async () => {
      expect(assignmentInputs[0]).toHaveValue(0);
      expect(assignmentInputs[1]).toHaveValue(0);
      expect(assignmentInputs[2]).toHaveValue(0);
    });
  });
});
