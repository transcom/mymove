import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AccountingCodesModal from './AccountingCodesModal';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('components/Office/AccountingCodesModal', () => {
  it('renders content with minimal props', () => {
    render(<AccountingCodesModal isOpen shipmentType={SHIPMENT_OPTIONS.NTSR} />);

    expect(screen.getByText('NTS-release')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Edit accounting codes' })).toBeInTheDocument();
    expect(screen.getByText('No TAC code entered.')).toBeInTheDocument();
    expect(screen.getByText('No SAC code entered.')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add or edit codes' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
  });

  it('renders full content', async () => {
    const onClose = jest.fn();
    const onEditCodesClick = jest.fn();

    render(
      <AccountingCodesModal
        isOpen
        shipmentType={SHIPMENT_OPTIONS.NTS}
        TACs={{ HHG: '1234', NTS: '2345' }}
        SACs={{ HHG: 'ABCD', NTS: 'BCDE' }}
        onClose={onClose}
        onEditCodesClick={onEditCodesClick}
        tacType="HHG"
        sacType="NTS"
      />,
    );

    expect(screen.getByLabelText('1234 (HHG)')).toBeChecked();
    expect(screen.getByLabelText('BCDE (NTS)')).toBeChecked();

    userEvent.click(screen.getByRole('button', { name: 'Add or edit codes' }));
    await waitFor(() => expect(onEditCodesClick).toHaveBeenCalledTimes(1));

    userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    userEvent.click(screen.getByTestId('modalCloseButton'));
    await waitFor(() => expect(onClose).toHaveBeenCalledTimes(2));
  });

  it('returns values from form', async () => {
    const onSubmit = jest.fn();

    render(
      <AccountingCodesModal
        isOpen
        shipmentType={SHIPMENT_OPTIONS.NTS}
        TACs={{ HHG: '1234', NTS: '2345' }}
        SACs={{ HHG: 'ABCD', NTS: 'BCDE' }}
        onSubmit={onSubmit}
        tacType="HHG"
        sacType="NTS"
      />,
    );

    userEvent.click(screen.getByRole('button', { name: 'Save' }));
    await waitFor(() => expect(onSubmit).toHaveBeenCalledWith({ tacType: 'HHG', sacType: 'NTS' }));

    userEvent.click(screen.getByLabelText('2345 (NTS)'));
    userEvent.click(screen.getByLabelText('ABCD (HHG)'));
    userEvent.click(screen.getByRole('button', { name: 'Save' }));
    await waitFor(() => expect(onSubmit).toHaveBeenCalledWith({ tacType: 'NTS', sacType: 'HHG' }));
  });
});
