/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
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

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getAllByLabelText('Email')[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.getByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getAllByLabelText('Email')[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(screen.getByLabelText('Customer remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('does not render special NTS What to expect section', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      expect(
        await screen.queryByText(
          'The moving company will find a storage facility approved by the government, and will move your belongings there.',
        ),
      ).not.toBeInTheDocument();
    });

    it('uses the current residence address for pickup address when checked', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.getByLabelText('Use my current address'));

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveValue(
        defaultProps.currentResidence.street_address_1,
      );

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue(defaultProps.currentResidence.city);
      expect(screen.getAllByLabelText('State')[0]).toHaveValue(defaultProps.currentResidence.state);
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue(defaultProps.currentResidence.postal_code);
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.getByLabelText('Yes'));

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveAttribute(
        'name',
        'pickup.address.street_address_1',
      );
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveAttribute('name', 'delivery.address.street_address_1');

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveAttribute('name', 'pickup.address.street_address_2');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveAttribute('name', 'delivery.address.street_address_2');

      expect(screen.getAllByLabelText('City')[0]).toHaveAttribute('name', 'pickup.address.city');
      expect(screen.getAllByLabelText('City')[1]).toHaveAttribute('name', 'delivery.address.city');

      expect(screen.getAllByLabelText('State')[0]).toHaveAttribute('name', 'pickup.address.state');
      expect(screen.getAllByLabelText('State')[1]).toHaveAttribute('name', 'delivery.address.state');

      expect(screen.getAllByLabelText('ZIP')[0]).toHaveAttribute('name', 'pickup.address.postal_code');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveAttribute('name', 'delivery.address.postal_code');
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

      expect(await screen.findByLabelText('Requested pickup date')).toHaveValue('01 Mar 2020');
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText('State')[0]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue('78234');
      expect(screen.getByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(screen.getByLabelText('Customer remarks')).toHaveValue('mock customer remarks');
      expect(screen.getByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('First name')).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getByLabelText('Last name')).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getByLabelText('Phone')).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getByLabelText('Email')).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.queryByText('Delivery location')).not.toBeInTheDocument();
      expect(screen.queryByText(/Receiving agent/)).not.toBeInTheDocument();

      expect(screen.getByLabelText('Customer remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders special NTS What to expect section', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      expect(
        await screen.findByText(
          'The moving company will find a storage facility approved by the government, and will move your belongings there.',
        ),
      ).toBeInTheDocument();
    });
  });

  describe('creating a new NTS-R shipment', () => {
    it('renders the NTS-R shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-R')).toHaveClass('usa-tag');

      expect(screen.queryByText('Pickup location')).not.toBeInTheDocument();
      expect(screen.queryByText(/Releasing agent/)).not.toBeInTheDocument();

      expect(screen.getByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('First name')).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getByLabelText('Last name')).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getByLabelText('Phone')).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getByLabelText('Email')).toHaveAttribute('name', 'delivery.agent.email');

      expect(screen.getByLabelText('Customer remarks')).toBeInstanceOf(HTMLTextAreaElement);
      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('does not render special NTS What to expect section', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      expect(
        await screen.queryByText(
          'The moving company will find a storage facility approved by the government, and will move your belongings there.',
        ),
      ).not.toBeInTheDocument();
    });
  });
});
