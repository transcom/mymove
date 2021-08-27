import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ApproveRejectModal } from '.';
import userEvent from '@testing-library/user-event';

const approveBtnOnClick = jest.fn();
const rejectBtnOnClick = jest.fn();

describe('AcceptRejectModal component test', () => {
  // Positive tests
  it('renders without crashing with required props', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);
    expect(await screen.getByRole('button', { name: 'Approve' })).toBeInTheDocument();
  });

  it('renders nothing without crashing with required props and showModal prop false', async () => {
    render(
      <ApproveRejectModal
        showModal={false}
        approveBtnOnClick={approveBtnOnClick}
        rejectBtnOnClick={rejectBtnOnClick}
      />,
    );
    expect(await screen.queryByRole('button', { name: 'Approve' })).not.toBeInTheDocument();
  });

  it('handleApproveClick() is called', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);
    const approveButton = await screen.findByRole('button', { name: 'Approve' });
    expect(approveButton).toBeInTheDocument();
    userEvent.click(approveButton);

    await waitFor(() => {
      expect(approveBtnOnClick).toHaveBeenCalled();
    });
  });

  it('handleRejectionClick() is called', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);

    const rejectionToggle = await screen.getByRole('button', { name: 'Reject' });
    expect(rejectionToggle).toBeInTheDocument();
    userEvent.click(rejectionToggle);

    const rejectionReason = await screen.getByLabelText('Rejection reason');
    expect(rejectionReason).toBeInTheDocument();
    userEvent.type(rejectionReason, 'rejected');

    const rejectButton = await screen.getByRole('button', { name: 'Reject' });
    expect(rejectButton).toBeInTheDocument();
    userEvent.click(rejectButton);

    expect(rejectBtnOnClick).toHaveBeenCalled();
  });

  it('reject button is disabled if reject reason is empty', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);

    const rejectionToggle = await screen.getByRole('button', { name: 'Reject' });
    expect(rejectionToggle).toBeInTheDocument();
    userEvent.click(rejectionToggle);

    const rejectionReason = await screen.getByLabelText('Rejection reason');
    expect(rejectionReason).toBeInTheDocument();
    userEvent.clear(rejectionReason);

    expect(await screen.getByRole('button', { name: 'Reject' })).toBeDisabled();
  });

  it('reject button is enabled if reject reason is filled', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);

    const rejectionToggle = await screen.getByRole('button', { name: 'Reject' });
    expect(rejectionToggle).toBeInTheDocument();
    userEvent.click(rejectionToggle);

    const rejectionReason = await screen.getByLabelText('Rejection reason');
    expect(rejectionReason).toBeInTheDocument();
    userEvent.type(rejectionReason, 'rejected');

    expect(await screen.getByRole('button', { name: 'Reject' })).not.toBeDisabled();
  });

  it('clicking cancel button resets the modal', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);

    const rejectionToggle = await screen.getByRole('button', { name: 'Reject' });
    expect(rejectionToggle).toBeInTheDocument();
    userEvent.click(rejectionToggle);

    const cancelButton = await screen.getByRole('button', { name: 'Cancel' });
    userEvent.click(cancelButton);

    expect(await screen.getByRole('button', { name: 'Approve' })).toBeInTheDocument();
    expect(await screen.getByTestId('rejectionToggle')).toBeInTheDocument();
    expect(await screen.queryByLabelText('Rejection reason')).not.toBeInTheDocument();
    expect(await screen.queryByTestId('rejectionButton')).not.toBeInTheDocument();
  });

  it('rejection toggle hides approve button and shows rejection reason', async () => {
    render(<ApproveRejectModal approveBtnOnClick={approveBtnOnClick} rejectBtnOnClick={rejectBtnOnClick} />);

    const rejectionToggle = await screen.getByRole('button', { name: 'Reject' });
    expect(rejectionToggle).toBeInTheDocument();
    userEvent.click(rejectionToggle);

    expect(await screen.queryByLabelText('Approve')).not.toBeInTheDocument();
    expect(rejectionToggle).not.toBeInTheDocument();
    expect(await screen.getByLabelText('Rejection reason')).toBeInTheDocument();
  });

  // Negative tests
  it('tries renders but crashes with no required props', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    render(<ApproveRejectModal />);
    await expect(consoleErrorSpy).toHaveBeenCalled();
  });
});
