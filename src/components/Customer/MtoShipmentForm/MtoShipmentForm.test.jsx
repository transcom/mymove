/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MtoShipmentForm from './MtoShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const defaultProps = {
  isCreatePage: true,
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  match: { isExact: false, path: '', url: '', params: { moveId: '' } },
  history: {
    goBack: jest.fn(),
    push: jest.fn(),
  },
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
    street_address_2: '',
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self: 5000,
    },
  },
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  pickupAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
};

describe('MtoShipmentForm component', () => {
  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      const { queryByText, queryByLabelText, queryAllByLabelText } = render(
        <MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />,
      );

      await waitFor(() => {
        expect(queryByText('HHG')).toHaveClass('usa-tag');
      });

      expect(queryByText('Pickup date')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(queryByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(queryAllByLabelText('First name')[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(queryAllByLabelText('Last name')[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(queryAllByLabelText('Phone')[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(queryAllByLabelText('Email')[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(queryByText('Delivery date')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(queryAllByLabelText('First name')[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(queryAllByLabelText('Last name')[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(queryAllByLabelText('Phone')[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(queryAllByLabelText('Email')[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(
        queryByLabelText('Is there anything special about this shipment that the movers should know?'),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('does not render special NTS What to expect section', async () => {
      const { queryByTestId } = render(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });

    it('uses the current residence address for pickup address when checked', async () => {
      const { queryByLabelText, queryAllByLabelText } = render(
        <MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />,
      );

      userEvent.click(queryByLabelText('Use my current address'));

      await waitFor(() => {
        expect(queryAllByLabelText('Address 1')[0]).toHaveValue(defaultProps.currentResidence.street_address_1);
        expect(queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(queryAllByLabelText('City')[0]).toHaveValue(defaultProps.currentResidence.city);
        expect(queryAllByLabelText('State')[0]).toHaveValue(defaultProps.currentResidence.state);
        expect(queryAllByLabelText('ZIP')[0]).toHaveValue(defaultProps.currentResidence.postal_code);
      });
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      const { queryByLabelText, queryAllByLabelText } = render(
        <MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />,
      );

      userEvent.click(queryByLabelText('Yes'));

      await waitFor(() => {
        expect(queryAllByLabelText('Address 1')[0]).toHaveAttribute('name', 'pickup.address.street_address_1');
        expect(queryAllByLabelText('Address 1')[1]).toHaveAttribute('name', 'delivery.address.street_address_1');

        expect(queryAllByLabelText(/Address 2/)[0]).toHaveAttribute('name', 'pickup.address.street_address_2');
        expect(queryAllByLabelText(/Address 2/)[1]).toHaveAttribute('name', 'delivery.address.street_address_2');

        expect(queryAllByLabelText('City')[0]).toHaveAttribute('name', 'pickup.address.city');
        expect(queryAllByLabelText('City')[1]).toHaveAttribute('name', 'delivery.address.city');

        expect(queryAllByLabelText('State')[0]).toHaveAttribute('name', 'pickup.address.state');
        expect(queryAllByLabelText('State')[1]).toHaveAttribute('name', 'delivery.address.state');

        expect(queryAllByLabelText('ZIP')[0]).toHaveAttribute('name', 'pickup.address.postal_code');
        expect(queryAllByLabelText('ZIP')[1]).toHaveAttribute('name', 'delivery.address.postal_code');
      });
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      const { queryByLabelText, queryAllByLabelText } = render(
        <MtoShipmentForm
          {...defaultProps}
          isCreatePage={false}
          selectedMoveType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
        />,
      );

      await waitFor(() => {
        expect(queryByLabelText('Requested pickup date')).toHaveValue('01 Mar 2020');
      });
      expect(queryByLabelText('Use my current address')).not.toBeChecked();
      expect(queryAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(queryAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(queryAllByLabelText('State')[0]).toHaveValue('TX');
      expect(queryAllByLabelText('ZIP')[0]).toHaveValue('78234');
      expect(queryByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      expect(queryByLabelText('Yes')).toBeChecked();
      expect(queryAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(queryAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(queryAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(queryAllByLabelText('State')[1]).toHaveValue('WA');
      expect(queryAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(
        queryByLabelText('Is there anything special about this shipment that the movers should know?'),
      ).toHaveValue('mock remarks');
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      const { queryByText, queryByLabelText } = render(
        <MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />,
      );

      await waitFor(() => {
        expect(queryByText('NTS')).toHaveClass('usa-tag');
      });
      expect(queryByText('Pickup date')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(queryByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('First name')).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(queryByLabelText('Last name')).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(queryByLabelText('Phone')).toHaveAttribute('name', 'pickup.agent.phone');
      expect(queryByLabelText('Email')).toHaveAttribute('name', 'pickup.agent.email');

      expect(queryByText('Delivery date')).not.toBeInTheDocument();
      expect(queryByText('Delivery location')).not.toBeInTheDocument();
      expect(queryByText(/Receiving agent/)).not.toBeInTheDocument();

      expect(
        queryByLabelText('Is there anything special about this shipment that the movers should know?'),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders special NTS What to expect section', async () => {
      const { queryByTestId } = render(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).toBeInTheDocument();
      });
    });
  });

  describe('creating a new NTS-R shipment', () => {
    it('renders the NTS-R shipment form', async () => {
      const { queryByText, queryByLabelText } = render(
        <MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />,
      );

      await waitFor(() => {
        expect(queryByText('NTS-R')).toHaveClass('usa-tag');
      });

      expect(queryByText('Pickup date')).not.toBeInTheDocument();
      expect(queryByText('Pickup location')).not.toBeInTheDocument();
      expect(queryByText(/Releasing agent/)).not.toBeInTheDocument();

      expect(queryByText('Delivery date')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(queryByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(queryByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(queryByLabelText('First name')).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(queryByLabelText('Last name')).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(queryByLabelText('Phone')).toHaveAttribute('name', 'delivery.agent.phone');
      expect(queryByLabelText('Email')).toHaveAttribute('name', 'delivery.agent.email');

      expect(
        queryByLabelText('Is there anything special about this shipment that the movers should know?'),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('does not render special NTS What to expect section', async () => {
      const { queryByTestId } = render(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });
  });
});
