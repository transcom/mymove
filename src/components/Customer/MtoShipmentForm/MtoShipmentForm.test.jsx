/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
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
      render(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.getByText('Pickup date')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Second pickup location')).toBeInstanceOf(HTMLHeadingElement);
      expect(screen.getByTitle('Yes, I have a second pickup location')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second pickup location')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getAllByLabelText('Email')[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.getByText('Delivery date')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not know my delivery address')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getAllByLabelText('Email')[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(
        screen.getByLabelText('Is there anything special about this shipment that the movers should know?'),
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

    it('renders a second address fieldset when the user has a second pickup address', async () => {
      render(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.getByTitle('Yes, I have a second pickup location'));

      const streetAddress1 = await screen.findAllByLabelText('Address 1');
      expect(streetAddress1[1]).toHaveAttribute('name', 'secondaryPickup.address.street_address_1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[1]).toHaveAttribute('name', 'secondaryPickup.address.street_address_2');

      const city = await screen.findAllByLabelText('City');
      expect(city[1]).toHaveAttribute('name', 'secondaryPickup.address.city');

      const state = await screen.findAllByLabelText('State');
      expect(state[1]).toHaveAttribute('name', 'secondaryPickup.address.state');

      const zip = await screen.findAllByLabelText('ZIP');
      expect(zip[1]).toHaveAttribute('name', 'secondaryPickup.address.postal_code');
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      render(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      const streetAddress1 = await screen.findAllByLabelText('Address 1');
      expect(streetAddress1[0]).toHaveAttribute('name', 'pickup.address.street_address_1');
      expect(streetAddress1[1]).toHaveAttribute('name', 'delivery.address.street_address_1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[0]).toHaveAttribute('name', 'pickup.address.street_address_2');
      expect(streetAddress2[1]).toHaveAttribute('name', 'delivery.address.street_address_2');

      const city = await screen.findAllByLabelText('City');
      expect(city[0]).toHaveAttribute('name', 'pickup.address.city');
      expect(city[1]).toHaveAttribute('name', 'delivery.address.city');

      const state = await screen.findAllByLabelText('State');
      expect(state[0]).toHaveAttribute('name', 'pickup.address.state');
      expect(state[1]).toHaveAttribute('name', 'delivery.address.state');

      const zip = await screen.findAllByLabelText('ZIP');
      expect(zip[0]).toHaveAttribute('name', 'pickup.address.postal_code');
      expect(zip[1]).toHaveAttribute('name', 'delivery.address.postal_code');
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      render(
        <MtoShipmentForm
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
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(
        screen.getByLabelText('Is there anything special about this shipment that the movers should know?'),
      ).toHaveValue('mock remarks');
    });

    it('renders the HHG shipment form with a pre-filled secondary address', async () => {
      render(
        <MtoShipmentForm
          {...defaultProps}
          isCreatePage={false}
          selectedMoveType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={{
            ...mockMtoShipment,
            secondaryPickupAddress: {
              street_address_1: '142 E Barrel Hoop Circle',
              street_address_2: '#4A',
              city: 'Corpus Christi',
              state: 'TX',
              postal_code: '78412',
            },
          }}
        />,
      );

      expect(await screen.findByTitle('Yes, I have a second pickup location')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('142 E Barrel Hoop Circle');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('#4A');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('Corpus Christi');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('78412');
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
