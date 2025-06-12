import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import * as Yup from 'yup';
import userEvent from '@testing-library/user-event';

import CreatePaymentRequestForm from './CreatePaymentRequestForm';

import { MockProviders } from 'testUtils';

describe('CreatePaymentRequestForm', () => {
  // No need to test any other data setting here because all checkboxes are unset to start
  const initialValues = {
    serviceItems: [],
  };

  const createPaymentRequestSchema = Yup.object().shape({
    serviceItems: Yup.array().of(Yup.string()).min(1),
  });

  const twoShipments = [
    {
      id: '1',
      pickupAddress: { streetAddress1: '500 Main Street', city: 'New York', state: 'NY', postalCode: '10001' },
      destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      primeActualWeight: 2000,
    },
    {
      id: '2',
      pickupAddress: { streetAddress1: '33 Bleeker Street', city: 'New York', state: 'NY', postalCode: '10002' },
      destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      primeActualWeight: 2000,
    },
  ];

  const basicAndShipmentsServiceItems = {
    basic: [{ id: '3', reServiceCode: 'MS' }],
    1: [
      { id: '4', reServiceCode: 'DLH' },
      { id: '6', reServiceCode: 'DDFSIT', reServiceName: 'Domestic destination 1st day SIT' },
    ],
    2: [{ id: '5', reServiceCode: 'FSC' }],
    3: [
      { id: '7', reServiceCode: 'IHPK' },
      { id: '8', reServiceCode: 'IHUPK' },
      { id: '8', reServiceCode: 'ISLH' },
      { id: '8', reServiceCode: 'POEFSC' },
    ],
    4: [
      { id: '11', reServiceCode: 'PODFSC' },
      { id: '12', reServiceCode: 'IUBPK' },
      { id: '13', reServiceCode: 'IUBUPK' },
      { id: '14', reServiceCode: 'UBP' },
    ],
  };

  it('renders the form', async () => {
    render(
      <MockProviders>
        <CreatePaymentRequestForm
          initialValues={initialValues}
          createPaymentRequestSchema={createPaymentRequestSchema}
          mtoShipments={twoShipments}
          groupedServiceItems={basicAndShipmentsServiceItems}
          onSubmit={jest.fn()}
          handleSelectAll={jest.fn()}
          handleValidateDate={jest.fn()}
        />
      </MockProviders>,
    );

    // 1 move service item and 1 on each shipment
    expect(screen.getAllByRole('checkbox', { name: 'Add to payment request' }).length).toEqual(4);
    // 1 select all checkbox for each shipment
    expect(screen.getAllByLabelText('Add all service items').length).toEqual(2);
    const submitBtn = screen.getByLabelText('Submit Payment Request');
    expect(submitBtn).toBeDisabled();
  });

  it('enables the submit button when at least one service item is checked', async () => {
    render(
      <MockProviders>
        <CreatePaymentRequestForm
          initialValues={initialValues}
          createPaymentRequestSchema={createPaymentRequestSchema}
          mtoShipments={twoShipments}
          groupedServiceItems={basicAndShipmentsServiceItems}
          onSubmit={jest.fn()}
          handleSelectAll={jest.fn()}
          handleValidateDate={jest.fn()}
        />
      </MockProviders>,
    );

    await userEvent.click(screen.getAllByRole('checkbox', { name: 'Add to payment request' })[0]);

    await waitFor(() => {
      const submitBtn = screen.getByLabelText('Submit Payment Request');
      expect(submitBtn).toBeEnabled();
    });
  });

  it('displays the validation error when no service items are selected', async () => {
    render(
      <MockProviders>
        <CreatePaymentRequestForm
          initialValues={initialValues}
          createPaymentRequestSchema={createPaymentRequestSchema}
          mtoShipments={twoShipments}
          groupedServiceItems={basicAndShipmentsServiceItems}
          onSubmit={jest.fn()}
          handleSelectAll={jest.fn()}
          handleValidateDate={jest.fn()}
        />
      </MockProviders>,
    );

    const basicServiceItemInput = screen.getAllByRole('checkbox', { name: 'Add to payment request' })[0];
    await userEvent.click(basicServiceItemInput);

    await userEvent.click(basicServiceItemInput);

    await waitFor(() => {
      expect(screen.getByText('At least 1 service item must be added when creating a payment request'));
      const submitBtn = screen.getByLabelText('Submit Payment Request');
      expect(submitBtn).toBeDisabled();
    });
  });

  it('selects all service items of the shipment when select all is clicked', async () => {
    const handleSelectAll = (shipmentID, values, setValues) => {
      setValues({ serviceItems: ['4'] });
    };

    render(
      <MockProviders>
        <CreatePaymentRequestForm
          initialValues={initialValues}
          createPaymentRequestSchema={createPaymentRequestSchema}
          mtoShipments={twoShipments}
          groupedServiceItems={basicAndShipmentsServiceItems}
          onSubmit={jest.fn()}
          handleSelectAll={handleSelectAll}
          handleValidateDate={jest.fn()}
        />
      </MockProviders>,
    );

    const shipmentSelectAllInput = screen.getAllByRole('checkbox', { name: 'Add all service items' })[0];
    await userEvent.click(shipmentSelectAllInput);

    await waitFor(() => {
      expect(screen.getByRole('checkbox', { name: 'Add to payment request', checked: true })).toBeInTheDocument();
      const submitBtn = screen.getByLabelText('Submit Payment Request');
      expect(submitBtn).toBeEnabled();
    });
  });

  it('deselects all service items of the shipment when select all is unchecked', async () => {
    const handleSelectAll = (shipmentID, values, setValues, event) => {
      if (!event.target.checked) {
        setValues({ serviceItems: [] });
      } else {
        setValues({ serviceItems: ['4'] });
      }
    };

    render(
      <MockProviders>
        <CreatePaymentRequestForm
          initialValues={initialValues}
          createPaymentRequestSchema={createPaymentRequestSchema}
          mtoShipments={twoShipments}
          groupedServiceItems={basicAndShipmentsServiceItems}
          onSubmit={jest.fn()}
          handleSelectAll={handleSelectAll}
          handleValidateDate={jest.fn()}
        />
      </MockProviders>,
    );

    const shipmentSelectAllInput = screen.getAllByRole('checkbox', { name: 'Add all service items' })[0];
    await userEvent.click(shipmentSelectAllInput);
    await userEvent.click(shipmentSelectAllInput);

    await waitFor(() => {
      const shipmentServiceItemInput = screen.queryAllByRole('checkbox', {
        name: 'Add to payment request',
        checked: true,
      });
      // all checkboxes are back to being unchecked
      expect(shipmentServiceItemInput).toHaveLength(0);
      expect(screen.getByText('At least 1 service item must be added when creating a payment request'));
      const submitBtn = screen.getByLabelText('Submit Payment Request');
      expect(submitBtn).toBeDisabled();
    });
  });

  it('renders the weight billed text input box', async () => {
    render(
      <MockProviders>
        <CreatePaymentRequestForm
          initialValues={initialValues}
          createPaymentRequestSchema={createPaymentRequestSchema}
          mtoShipments={twoShipments}
          groupedServiceItems={basicAndShipmentsServiceItems}
          onSubmit={jest.fn()}
          handleSelectAll={jest.fn()}
          handleValidateDate={jest.fn()}
        />
      </MockProviders>,
    );

    expect(
      screen.getAllByRole('textbox', { name: 'Weight Billed (if different from shipment weight)' }).length,
    ).toBeGreaterThan(0);
  });
});
