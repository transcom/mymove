/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingShipmentForm from './ServicesCounselingShipmentForm';

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
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
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

describe('ServicesCounselingShipmentForm component', () => {
  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      await waitFor(() => {
        expect(screen.queryByText('HHG')).toHaveClass('usa-tag');
      });

      expect(screen.queryByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.queryByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryAllByLabelText('First name')[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.queryAllByLabelText('Last name')[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.queryAllByLabelText('Phone')[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.queryAllByLabelText('Email')[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.queryByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryAllByLabelText('First name')[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.queryAllByLabelText('Last name')[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.queryAllByLabelText('Phone')[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.queryAllByLabelText('Email')[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(screen.queryByLabelText('Customer remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.queryByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('does not render special NTS What to expect section', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      await waitFor(() => {
        expect(screen.queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });

    it('uses the current residence address for pickup address when checked', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.queryByLabelText('Use my current address'));

      await waitFor(() => {
        expect(screen.queryAllByLabelText('Address 1')[0]).toHaveValue(defaultProps.currentResidence.street_address_1);
        expect(screen.queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(screen.queryAllByLabelText('City')[0]).toHaveValue(defaultProps.currentResidence.city);
        expect(screen.queryAllByLabelText('State')[0]).toHaveValue(defaultProps.currentResidence.state);
        expect(screen.queryAllByLabelText('ZIP')[0]).toHaveValue(defaultProps.currentResidence.postal_code);
      });
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.queryByLabelText('Yes'));

      await waitFor(() => {
        expect(screen.queryAllByLabelText('Address 1')[0]).toHaveAttribute('name', 'pickup.address.street_address_1');
        expect(screen.queryAllByLabelText('Address 1')[1]).toHaveAttribute('name', 'delivery.address.street_address_1');

        expect(screen.queryAllByLabelText(/Address 2/)[0]).toHaveAttribute('name', 'pickup.address.street_address_2');
        expect(screen.queryAllByLabelText(/Address 2/)[1]).toHaveAttribute('name', 'delivery.address.street_address_2');

        expect(screen.queryAllByLabelText('City')[0]).toHaveAttribute('name', 'pickup.address.city');
        expect(screen.queryAllByLabelText('City')[1]).toHaveAttribute('name', 'delivery.address.city');

        expect(screen.queryAllByLabelText('State')[0]).toHaveAttribute('name', 'pickup.address.state');
        expect(screen.queryAllByLabelText('State')[1]).toHaveAttribute('name', 'delivery.address.state');

        expect(screen.queryAllByLabelText('ZIP')[0]).toHaveAttribute('name', 'pickup.address.postal_code');
        expect(screen.queryAllByLabelText('ZIP')[1]).toHaveAttribute('name', 'delivery.address.postal_code');
      });
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      render(
        <ServicesCounselingShipmentForm
          {...defaultProps}
          isCreatePage={false}
          selectedMoveType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
        />,
      );

      await waitFor(() => {
        expect(screen.queryByLabelText('Requested pickup date')).toHaveValue('01 Mar 2020');
      });
      expect(screen.queryByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.queryAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(screen.queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.queryAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(screen.queryAllByLabelText('State')[0]).toHaveValue('TX');
      expect(screen.queryAllByLabelText('ZIP')[0]).toHaveValue('78234');
      expect(screen.queryByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      expect(screen.queryByLabelText('Yes')).toBeChecked();
      expect(screen.queryAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.queryAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.queryAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(screen.queryAllByLabelText('State')[1]).toHaveValue('WA');
      expect(screen.queryAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(screen.queryByLabelText('Customer remarks')).toHaveValue('mock customer remarks');
      expect(screen.queryByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      await waitFor(() => {
        expect(screen.queryByText('NTS')).toHaveClass('usa-tag');
      });
      expect(screen.queryByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.queryByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryByLabelText('First name')).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.queryByLabelText('Last name')).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.queryByLabelText('Phone')).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.queryByLabelText('Email')).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.queryByText('Delivery location')).not.toBeInTheDocument();
      expect(screen.queryByText(/Receiving agent/)).not.toBeInTheDocument();

      expect(screen.queryByLabelText('Customer remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.queryByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders special NTS What to expect section', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      await waitFor(() => {
        expect(screen.queryByTestId('nts-what-to-expect')).toBeInTheDocument();
      });
    });
  });

  describe('creating a new NTS-R shipment', () => {
    it('renders the NTS-R shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      await waitFor(() => {
        expect(screen.queryByText('NTS-R')).toHaveClass('usa-tag');
      });

      expect(screen.queryByText('Pickup location')).not.toBeInTheDocument();
      expect(screen.queryByText(/Releasing agent/)).not.toBeInTheDocument();

      expect(screen.queryByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.queryByLabelText('First name')).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.queryByLabelText('Last name')).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.queryByLabelText('Phone')).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.queryByLabelText('Email')).toHaveAttribute('name', 'delivery.agent.email');

      expect(screen.queryByLabelText('Customer remarks')).toBeInstanceOf(HTMLTextAreaElement);
      expect(screen.queryByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('does not render special NTS What to expect section', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      await waitFor(() => {
        expect(screen.queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });
  });
});
