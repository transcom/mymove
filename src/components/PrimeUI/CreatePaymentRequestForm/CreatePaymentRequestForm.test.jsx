import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import * as Yup from 'yup';
import userEvent from '@testing-library/user-event';

import CreatePaymentRequestForm from './CreatePaymentRequestForm';

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
    },
    {
      id: '2',
      pickupAddress: { streetAddress1: '33 Bleeker Street', city: 'New York', state: 'NY', postalCode: '10002' },
      destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
    },
  ];

  const basicAndShipmentsServiceItems = {
    basic: [{ id: '3', reServiceCode: 'MS' }],
    1: [{ id: '4', reServiceCode: 'DLH' }],
    2: [{ id: '5', reServiceCode: 'FSC' }],
  };

  it('renders the form', async () => {
    render(
      <CreatePaymentRequestForm
        initialValues={initialValues}
        createPaymentRequestSchema={createPaymentRequestSchema}
        mtoShipments={twoShipments}
        groupedServiceItems={basicAndShipmentsServiceItems}
        onSubmit={jest.fn()}
        handleSelectAll={jest.fn()}
        handleValidateDate={jest.fn()}
      />,
    );

    // 1 move service item and 1 on each shipment
    expect(screen.getAllByRole('checkbox', { name: 'Add to payment request' }).length).toEqual(3);
    // 1 select all checkbox for each shipment
    expect(screen.getAllByLabelText('Add all service items').length).toEqual(2);
    expect(screen.getByRole('button', { type: 'submit' })).toBeDisabled();
  });

  it('enables the submit button when at least one service item is checked', async () => {
    render(
      <CreatePaymentRequestForm
        initialValues={initialValues}
        createPaymentRequestSchema={createPaymentRequestSchema}
        mtoShipments={twoShipments}
        groupedServiceItems={basicAndShipmentsServiceItems}
        onSubmit={jest.fn()}
        handleSelectAll={jest.fn()}
        handleValidateDate={jest.fn()}
      />,
    );

    await userEvent.click(screen.getAllByRole('checkbox', { name: 'Add to payment request' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('button', { type: 'submit' })).toBeEnabled();
    });
  });

  it('displays the validation error when no service items are selected', async () => {
    render(
      <CreatePaymentRequestForm
        initialValues={initialValues}
        createPaymentRequestSchema={createPaymentRequestSchema}
        mtoShipments={twoShipments}
        groupedServiceItems={basicAndShipmentsServiceItems}
        onSubmit={jest.fn()}
        handleSelectAll={jest.fn()}
        handleValidateDate={jest.fn()}
      />,
    );

    const basicServiceItemInput = screen.getAllByRole('checkbox', { name: 'Add to payment request' })[0];
    await userEvent.click(basicServiceItemInput);

    await userEvent.click(basicServiceItemInput);

    await waitFor(() => {
      expect(screen.getByText('At least 1 service item must be added when creating a payment request'));
      expect(screen.getByRole('button', { type: 'submit' })).toBeDisabled();
    });
  });

  it('selects all service items of the shipment when select all is clicked', async () => {
    const handleSelectAll = (shipmentID, values, setValues) => {
      setValues({ serviceItems: ['4'] });
    };

    render(
      <CreatePaymentRequestForm
        initialValues={initialValues}
        createPaymentRequestSchema={createPaymentRequestSchema}
        mtoShipments={twoShipments}
        groupedServiceItems={basicAndShipmentsServiceItems}
        onSubmit={jest.fn()}
        handleSelectAll={handleSelectAll}
        handleValidateDate={jest.fn()}
      />,
    );

    const shipmentSelectAllInput = screen.getAllByRole('checkbox', { name: 'Add all service items' })[0];
    await userEvent.click(shipmentSelectAllInput);

    await waitFor(() => {
      expect(screen.getByRole('checkbox', { name: 'Add to payment request', checked: true })).toBeInTheDocument();
      expect(screen.getByRole('button', { type: 'submit' })).toBeEnabled();
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
      <CreatePaymentRequestForm
        initialValues={initialValues}
        createPaymentRequestSchema={createPaymentRequestSchema}
        mtoShipments={twoShipments}
        groupedServiceItems={basicAndShipmentsServiceItems}
        onSubmit={jest.fn()}
        handleSelectAll={handleSelectAll}
        handleValidateDate={jest.fn()}
      />,
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
      expect(screen.getByRole('button', { type: 'submit' })).toBeDisabled();
    });
  });
});
