import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import OrdersInfoForm from './OrdersInfoForm';

const testProps = {
  onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
  initialValues: { orders_type: '', issue_date: '', report_by_date: '', has_dependents: '', new_duty_station: {} },
  onBack: jest.fn(),
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
  currentStation: {},
};

describe('OrdersInfoForm component', () => {
  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Orders type')).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText('Orders type')).toBeRequired();
      expect(getByLabelText('Orders date')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Orders date')).toBeRequired();
      expect(getByLabelText('Report-by date')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Report-by date')).toBeRequired();
      expect(getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      // TODO - fixed when Yan's PR merged
      // expect(getByLabelText('New duty station')).toBeInstanceOf(HTMLSelectElement);
    });
  });

  it('renders each option for orders type', async () => {
    const { getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    await waitFor(() => {
      const ordersTypeDropdown = getByLabelText('Orders type');
      expect(ordersTypeDropdown).toBeInstanceOf(HTMLSelectElement);

      userEvent.selectOptions(ordersTypeDropdown, 'PERMANENT_CHANGE_OF_STATION');
      expect(ordersTypeDropdown).toHaveValue('PERMANENT_CHANGE_OF_STATION');

      userEvent.selectOptions(ordersTypeDropdown, 'RETIREMENT');
      expect(ordersTypeDropdown).toHaveValue('RETIREMENT');

      userEvent.selectOptions(ordersTypeDropdown, 'SEPARATION');
      expect(ordersTypeDropdown).toHaveValue('SEPARATION');
    });
  });

  it.skip('validates the new duty station against the current duty station', () => {
    //
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, getAllByText } = render(<OrdersInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(3);
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(
      <OrdersInfoForm
        {...testProps}
        initialValues={{ ...testProps.initialValues, new_duty_station: { id: 'testId', name: 'test duty station' } }}
      />,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.selectOptions(getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    userEvent.type(getByLabelText('Orders date'), '08 Nov 2020');
    userEvent.type(getByLabelText('Report-by date'), '26 Nov 2020');
    userEvent.click(getByLabelText('No'));

    // TODO - select duty station option

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          orders_type: 'PERMANENT_CHANGE_OF_STATION',
          has_dependents: 'no',
          issue_date: '08 Nov 2020',
          report_by_date: '26 Nov 2020',
        }),
        expect.anything(),
      );
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<OrdersInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  describe('with initial values', () => {
    const testInitialValues = {
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: 'no',
      new_duty_station: {
        address: {
          city: 'Des Moines',
          country: 'US',
          id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
          postal_code: '50309',
          state: 'IA',
          street_address_1: '987 Other Avenue',
          street_address_2: 'P.O. Box 1234',
          street_address_3: 'c/o Another Person',
        },
        address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        affiliation: 'AIR_FORCE',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
        name: 'Yuma AFB',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
    };

    it('pre-fills the inputs', async () => {
      const { getByLabelText } = render(<OrdersInfoForm {...testProps} initialValues={testInitialValues} />);

      await waitFor(() => {
        expect(getByLabelText('Orders type')).toHaveValue(testInitialValues.orders_type);
        expect(getByLabelText('Orders date')).toHaveValue('08 Nov 2020');
        expect(getByLabelText('Report-by date')).toHaveValue('26 Nov 2020');
        expect(getByLabelText('Yes')).not.toBeChecked();
        expect(getByLabelText('No')).toBeChecked();

        // TODO - fixed when Yan's PR merged
        // expect(getByLabelText('New duty station')).toHaveValue(testInitialValues.new_duty_station.id);
      });
    });
  });
});
