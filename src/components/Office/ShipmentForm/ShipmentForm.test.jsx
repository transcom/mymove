/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentForm from './ShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ORDERS_TYPE } from 'constants/orders';
import { roleTypes } from 'constants/userRoles';

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
  userRole: roleTypes.SERVICES_COUNSELOR,
  orderType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
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

const mockShipmentWithDestinationType = {
  ...mockMtoShipment,
  destinationType: 'HOME_OF_SELECTION',
};

const defaultPropsRetirement = {
  ...defaultProps,
  orderType: ORDERS_TYPE.RETIREMENT,
};

const defaultPropsSeparation = {
  ...defaultProps,
  orderType: ORDERS_TYPE.SEPARATION,
};

describe('ShipmentForm component', () => {
  describe('when creating a new shipment', () => {
    it('does not show the delete shipment button', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      const deleteButton = screen.queryByRole('button', { name: 'Delete shipment' });
      await waitFor(() => {
        expect(deleteButton).not.toBeInTheDocument();
      });
    });
  });
  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

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
    });

    it('uses the current residence address for pickup address when checked', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

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
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

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

    it('renders a delivery address type for retirement orders type', async () => {
      render(<ShipmentForm {...defaultPropsRetirement} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);
      userEvent.click(screen.getByLabelText('Yes'));

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.getAllByLabelText('Destination type')[0]).toHaveAttribute('name', 'destinationType');
    });

    it('does not render delivery address type for PCS order type', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);
      userEvent.click(screen.getByLabelText('Yes'));

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByLabelText('Destination type')).toBeNull();
    });

    it('renders a delivery address type for separation orders type', async () => {
      render(<ShipmentForm {...defaultPropsSeparation} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);
      userEvent.click(screen.getByLabelText('Yes'));

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.getAllByLabelText('Destination type')[0]).toHaveAttribute('name', 'destinationType');
    });

    it('does not render an Accounting Codes section', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByRole('heading', { name: 'Accounting codes' })).not.toBeInTheDocument();
    });

    it('does not render NTS release-only sections', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByText(/Shipment weight (lbs)/)).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).not.toBeInTheDocument();
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      render(
        <ShipmentForm
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

  describe('editing an already existing HHG shipment for retiree/separatee', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      render(
        <ShipmentForm
          {...defaultPropsRetirement}
          isCreatePage={false}
          selectedMoveType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockShipmentWithDestinationType}
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
      expect(screen.getByLabelText('Destination type')).toHaveValue('HOME_OF_SELECTION');
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

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

      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
    });

    it('renders an Accounting Codes section', async () => {
      render(
        <ShipmentForm
          {...defaultProps}
          TACs={{ HHG: '1234', NTS: '5678' }}
          selectedMoveType={SHIPMENT_OPTIONS.NTS}
          mtoShipment={mockMtoShipment}
        />,
      );

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
      expect(screen.getByLabelText('1234 (HHG)')).toBeInTheDocument();
      expect(screen.getByText('No SAC code entered.')).toBeInTheDocument();
    });

    it('does not render NTS release-only sections', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');
      expect(screen.queryByText(/Shipment weight (lbs)/)).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).not.toBeInTheDocument();
    });
  });

  describe('editing an already existing NTS shipment', () => {
    it('pre-fills the Accounting Codes section', async () => {
      render(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          mtoShipment={{
            ...mockMtoShipment,
            tacType: 'NTS',
            sacType: 'HHG',
          }}
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: '000012345' }}
          selectedMoveType={SHIPMENT_OPTIONS.NTS}
        />,
      );

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
      expect(screen.getByLabelText('000012345 (HHG)')).toBeChecked();
    });
  });

  describe('creating a new NTS-release shipment', () => {
    it('renders the NTS-release shipment form', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.queryByText('Pickup location')).not.toBeInTheDocument();
      expect(screen.queryByText(/Releasing agent/)).not.toBeInTheDocument();
      expect(screen.queryByLabelText('Yes')).not.toBeInTheDocument();
      expect(screen.queryByLabelText('No')).not.toBeInTheDocument();

      expect(screen.getByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('First name')).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getByLabelText('Last name')).toHaveAttribute('name', 'delivery.agent.lastName');

      expect(screen.getByText('Customer remarks')).toBeTruthy();
      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
    });

    it('renders an Accounting Codes section', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
    });

    it('renders the NTS release-only sections', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');
      expect(screen.getByText(/Previously recorded weight \(lbs\)/)).toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).toBeInTheDocument();
    });
  });

  describe('as a TOO', () => {
    it('renders the NTS shipment form', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} userRole={roleTypes.TOO} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByRole('heading', { level: 2, name: 'Vendor' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Storage facility info' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Storage facility address' })).toBeInTheDocument();
    });

    it('renders the NTS release shipment form', async () => {
      render(<ShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTSR} userRole={roleTypes.TOO} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.getByRole('heading', { level: 2, name: 'Vendor' })).toBeInTheDocument();
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
        <ShipmentForm
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
        <ShipmentForm
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

  describe('external vendor shipment', () => {
    it('shows the TOO an alert', async () => {
      render(
        <ShipmentForm
          {...defaultProps}
          selectedMoveType={SHIPMENT_OPTIONS.NTSR}
          mtoShipment={{ ...mockMtoShipment, usesExternalVendor: true }}
          isCreatePage={false}
          userRole={roleTypes.TOO}
        />,
      );

      expect(
        await screen.findByText(
          'The GHC prime contractor is not handling the shipment. Information will not be automatically shared with the movers handling it.',
        ),
      ).toBeInTheDocument();
    });

    it('does not show the SC an alert', async () => {
      render(
        <ShipmentForm
          // SC is default role from test props
          {...defaultProps}
          selectedMoveType={SHIPMENT_OPTIONS.NTSR}
          mtoShipment={{ ...mockMtoShipment, usesExternalVendor: true }}
          isCreatePage={false}
        />,
      );

      await waitFor(() => {
        expect(
          screen.queryByText(
            'The GHC prime contractor is not handling the shipment. Information will not be automatically shared with the movers handling it.',
          ),
        ).not.toBeInTheDocument();
      });
    });
  });
});
