import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { DeleteCustomerSupportRemarkConfirmationModal } from 'components/ConfirmationModals/DeleteCustomerSupportRemarkConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('DeleteCustomerSupportRemarkConfirmationModal', () => {
  const customerSupportRemarkID = '123456';

  it('renders the component', async () => {
    render(
      <DeleteCustomerSupportRemarkConfirmationModal
        onSubmit={onSubmit}
        onClose={onClose}
        customerSupportRemarkID={customerSupportRemarkID}
      />,
    );

    expect(
      await screen.findByRole('heading', { level: 3, name: 'Are you sure you want to delete this remark?' }),
    ).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(
      <DeleteCustomerSupportRemarkConfirmationModal
        onSubmit={onSubmit}
        onClose={onClose}
        customerSupportRemarkID={customerSupportRemarkID}
      />,
    );

    const closeButton = await screen.findByTestId('modalCloseButton');

    userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('closes the modal when the keep button is clicked', async () => {
    render(
      <DeleteCustomerSupportRemarkConfirmationModal
        onSubmit={onSubmit}
        onClose={onClose}
        customerSupportRemarkID={customerSupportRemarkID}
      />,
    );

    const keepButton = await screen.findByRole('button', { name: 'No, keep it' });

    userEvent.click(keepButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls the submit function when delete button is clicked', async () => {
    render(
      <DeleteCustomerSupportRemarkConfirmationModal
        onSubmit={onSubmit}
        onClose={onClose}
        customerSupportRemarkID={customerSupportRemarkID}
      />,
    );

    const deleteButton = await screen.findByRole('button', { name: 'Yes, Delete' });

    userEvent.click(deleteButton);

    expect(onSubmit).toHaveBeenCalledWith(customerSupportRemarkID);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });
});
