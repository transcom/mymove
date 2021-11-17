/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingShipmentForm from './ServicesCounselingShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const mockPush = jest.fn();

const defaultProps = {
  isCreatePage: true,
  match: { isExact: false, path: '', url: '', params: { moveCode: 'move123', shipementId: 'shipment123' } },
  history: {
    push: mockPush,
  },
  submitHandler: jest.fn(),
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  serviceMember: {
    weightAllotment: {
      totalWeightSelf: 5000,
    },
  },
  moveTaskOrderID: 'mock move id',
  mtoShipments: [],
};

const mockMtoShipment = {
  id: 'shipment123',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  mtoAgents: [
    {
      agentType: 'RELEASING_AGENT',
      email: 'jasn@email.com',
      firstName: 'Jason',
      lastName: 'Ash',
      phone: '999-999-9999',
    },
    {
      agentType: 'RECEIVING_AGENT',
      email: 'rbaker@email.com',
      firstName: 'Riley',
      lastName: 'Baker',
      phone: '863-555-9664',
    },
  ],
};

describe('ServicesCounselingShipmentForm component', () => {
  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use current address')).toBeInstanceOf(HTMLInputElement);
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

      expect(screen.getByText('Customer remarks')).toBeTruthy();

      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(
        screen.queryByText(
          'The moving company will find a storage facility approved by the government, and will move your belongings there.',
        ),
      ).not.toBeInTheDocument();
    });

    it('uses the current residence address for pickup address when checked', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.getByLabelText('Use current address'));

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveValue(
        defaultProps.currentResidence.streetAddress1,
      );

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue(defaultProps.currentResidence.city);
      expect(screen.getAllByLabelText('State')[0]).toHaveValue(defaultProps.currentResidence.state);
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue(defaultProps.currentResidence.postalCode);
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      userEvent.click(screen.getByLabelText('Yes'));

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveAttribute(
        'name',
        'pickup.address.streetAddress1',
      );
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveAttribute('name', 'delivery.address.streetAddress1');

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveAttribute('name', 'pickup.address.streetAddress2');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveAttribute('name', 'delivery.address.streetAddress2');

      expect(screen.getAllByLabelText('City')[0]).toHaveAttribute('name', 'pickup.address.city');
      expect(screen.getAllByLabelText('City')[1]).toHaveAttribute('name', 'delivery.address.city');

      expect(screen.getAllByLabelText('State')[0]).toHaveAttribute('name', 'pickup.address.state');
      expect(screen.getAllByLabelText('State')[1]).toHaveAttribute('name', 'delivery.address.state');

      expect(screen.getAllByLabelText('ZIP')[0]).toHaveAttribute('name', 'pickup.address.postalCode');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveAttribute('name', 'delivery.address.postalCode');
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
      expect(screen.getByLabelText('Use current address')).not.toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText('State')[0]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue('78234');
      expect(screen.getAllByLabelText('First name')[0]).toHaveValue('Jason');
      expect(screen.getAllByLabelText('Last name')[0]).toHaveValue('Ash');
      expect(screen.getAllByLabelText('Phone')[0]).toHaveValue('999-999-9999');
      expect(screen.getAllByLabelText('Email')[0]).toHaveValue('jasn@email.com');
      expect(screen.getByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(screen.getAllByLabelText('First name')[1]).toHaveValue('Riley');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveValue('Baker');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveValue('863-555-9664');
      expect(screen.getAllByLabelText('Email')[1]).toHaveValue('rbaker@email.com');
      expect(screen.getByText('Customer remarks')).toBeTruthy();
      expect(screen.getByText('mock customer remarks')).toBeTruthy();
      expect(screen.getByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup location')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use current address')).toBeInstanceOf(HTMLInputElement);
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

      expect(screen.getByText('Customer remarks')).toBeTruthy();

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

  describe('creating a new NTS-release shipment', () => {
    it('renders the NTS-release shipment form', async () => {
      render(<ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

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

      expect(screen.getByText('Customer remarks')).toBeTruthy();
      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(
        screen.queryByText(
          'The moving company will find a storage facility approved by the government, and will move your belongings there.',
        ),
      ).not.toBeInTheDocument();
    });
  });

  describe('filling the form', () => {
    it('shows an error if the submitHandler returns an error', async () => {
      const mockSubmitHandler = jest.fn(() =>
        // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
        // eslint-disable-next-line prefer-promise-reject-errors
        Promise.reject({
          message: 'A server error occurred editing the shipment details',
          response: {
            body: {
              detail: 'A server error occurred editing the shipment details',
            },
          },
        }),
      );

      render(
        <ServicesCounselingShipmentForm
          {...defaultProps}
          selectedMoveType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const saveButton = screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();

      userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(await screen.findByText('A server error occurred editing the shipment details')).toBeInTheDocument();
      expect(defaultProps.history.push).not.toHaveBeenCalled();
    });

    it('saves the update to the counselor remarks when the save button is clicked', async () => {
      const newCounselorRemarks = 'Counselor remarks';

      const expectedPayload = {
        body: {
          customerRemarks: 'mock customer remarks',
          counselorRemarks: newCounselorRemarks,
          destinationAddress: {
            streetAddress1: '441 SW Rio de la Plata Drive',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
            streetAddress2: '',
          },
          pickupAddress: {
            streetAddress1: '812 S 129th St',
            city: 'San Antonio',
            state: 'TX',
            postalCode: '78234',
            streetAddress2: '',
          },
          agents: [
            {
              agentType: 'RELEASING_AGENT',
              email: 'jasn@email.com',
              firstName: 'Jason',
              lastName: 'Ash',
              phone: '999-999-9999',
            },
            {
              agentType: 'RECEIVING_AGENT',
              email: 'rbaker@email.com',
              firstName: 'Riley',
              lastName: 'Baker',
              phone: '863-555-9664',
            },
          ],
          requestedDeliveryDate: '2020-03-30',
          requestedPickupDate: '2020-03-01',
          shipmentType: 'HHG',
        },
        shipmentID: 'shipment123',
        moveTaskOrderID: 'mock move id',
        normalize: false,
      };

      const patchResponse = {
        ...expectedPayload,
        created_at: '2021-02-08T16:48:04.117Z',
        updated_at: '2021-02-11T16:48:04.117Z',
      };

      const mockSubmitHandler = jest.fn(() => Promise.resolve(patchResponse));

      render(
        <ServicesCounselingShipmentForm
          {...defaultProps}
          selectedMoveType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const counselorRemarks = await screen.findByLabelText('Counselor remarks');

      userEvent.clear(counselorRemarks);

      userEvent.type(counselorRemarks, newCounselorRemarks);

      const saveButton = screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();

      userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(expectedPayload);
      });
    });
  });
});
