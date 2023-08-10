/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentForm from './ShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ORDERS_TYPE } from 'constants/orders';
import { roleTypes } from 'constants/userRoles';
import { ADDRESS_UPDATE_STATUS, ppmShipmentStatuses } from 'constants/shipments';
import { tooRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { validatePostalCode } from 'utils/validation';

const mockMutateFunction = jest.fn();
jest.mock('@tanstack/react-query', () => ({
  ...jest.requireActual('@tanstack/react-query'),
  useMutation: () => ({ mutate: mockMutateFunction }),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const defaultProps = {
  isCreatePage: true,
  submitHandler: jest.fn(),
  newDutyLocationAddress: {
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
  originDutyLocationAddress: {
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
    agency: '',
  },
  moveTaskOrderID: 'mock move id',
  mtoShipments: [],
  userRole: roleTypes.SERVICES_COUNSELOR,
  orderType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  isForServivcesCounseling: false,
};

const mockMtoShipment = {
  id: 'shipment123',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  hasSecondaryDeliveryAddress: false,
  hasSecondaryPickupAddress: false,
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
  displayDestinationType: true,
  destinationType: 'PLACE_ENTERED_ACTIVE_DUTY',
};

const mockPPMShipment = {
  ...mockMtoShipment,
  ppmShipment: {
    id: 'ppmShipmentID',
    shipmentId: 'shipment123',
    status: ppmShipmentStatuses.NEEDS_ADVANCE_APPROVAL,
    expectedDepartureDate: '2022-04-01',
    pickupPostalCode: '90210',
    destinationPostalCode: '90211',
    sitExpected: false,
    estimatedWeight: 4999,
    hasProGear: false,
    estimatedIncentive: 1234500,
    hasRequestedAdvance: true,
    advanceAmountRequested: 487500,
    advanceStatus: 'APPROVED',
  },
};

const mockDeliveryAddressUpdate = {
  deliveryAddressUpdate: {
    contractorRemarks: 'Test Contractor Remark',
    id: 'c49f7921-5a6e-46b4-bb39-022583574453',
    newAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NTU5MTVa',
      id: '6b57ce91-cabd-4e3b-9f48-ed4627d4878f',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    originalAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NDkyNTha',
      id: '92509013-aafc-4892-a476-2e3b97e6933d',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    shipmentID: '5c84bcf3-92f7-448f-b0e1-e5378b6806df',
    status: 'REQUESTED',
  },
};

const defaultPropsRetirement = {
  ...defaultProps,
  displayDestinationType: true,
  destinationType: 'HOME_OF_RECORD',
  orderType: ORDERS_TYPE.RETIREMENT,
};

const defaultPropsSeparation = {
  ...defaultProps,
  displayDestinationType: true,
  destinationType: 'HOME_OF_SELECTION',
  orderType: ORDERS_TYPE.SEPARATION,
};

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));
const mockRoutingOptions = {
  path: tooRoutes.BASE_SHIPMENT_EDIT_PATH,
  params: { moveCode: 'move123', shipmentId: 'shipment123' },
};

beforeEach(() => {
  jest.clearAllMocks();
});

const renderWithRouter = (ui) => {
  render(<MockProviders {...mockRoutingOptions}>{ui}</MockProviders>);
};

describe('ShipmentForm component', () => {
  describe('when creating a new shipment', () => {
    it('does not show the delete shipment button', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      const deleteButton = screen.queryByRole('button', { name: 'Delete shipment' });
      await waitFor(() => {
        expect(deleteButton).not.toBeInTheDocument();
      });
    });
  });

  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

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

      const deliveryLocationSectionHeadings = screen.getAllByText('Delivery location');
      expect(deliveryLocationSectionHeadings).toHaveLength(2);
      expect(deliveryLocationSectionHeadings[0]).toBeInstanceOf(HTMLParagraphElement);
      expect(deliveryLocationSectionHeadings[1]).toBeInstanceOf(HTMLLegendElement);
      expect(deliveryLocationSectionHeadings[1]).toHaveClass('usa-sr-only');
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[1]).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getAllByLabelText('Email')[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(screen.getByText('Customer remarks')).toBeTruthy();

      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('uses the current residence address for pickup address when checked', async () => {
      const user = userEvent.setup();
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      await user.click(screen.getByLabelText('Use current address'));

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveValue(
        defaultProps.currentResidence.streetAddress1,
      );

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue(defaultProps.currentResidence.city);
      expect(screen.getAllByLabelText('State')[0]).toHaveValue(defaultProps.currentResidence.state);
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue(defaultProps.currentResidence.postalCode);
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      await userEvent.click(screen.getAllByLabelText('Yes')[1]);

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
      renderWithRouter(<ShipmentForm {...defaultPropsRetirement} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      await userEvent.click(screen.getAllByLabelText('Yes')[1]);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.getAllByLabelText('Destination type')[0]).toHaveAttribute('name', 'destinationType');
    });

    it('does not render delivery address type for PCS order type', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);
      await userEvent.click(screen.getAllByLabelText('Yes')[1]);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByLabelText('Destination type')).toBeNull();
    });

    it('renders a delivery address type for separation orders type', async () => {
      renderWithRouter(<ShipmentForm {...defaultPropsSeparation} shipmentType={SHIPMENT_OPTIONS.HHG} />);
      await userEvent.click(screen.getAllByLabelText('Yes')[1]);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.getAllByLabelText('Destination type')[0]).toHaveAttribute('name', 'destinationType');
    });

    it('does not render an Accounting Codes section', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByRole('heading', { name: 'Accounting codes' })).not.toBeInTheDocument();
    });

    it('does not render NTS release-only sections', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByText(/Shipment weight (lbs)/)).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).not.toBeInTheDocument();
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
          displayDestinationType={false}
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
      expect(screen.getAllByLabelText('Yes')[0]).not.toBeChecked();
      expect(screen.getAllByLabelText('Yes')[1]).toBeChecked();
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

      const noDestinationTypeRadioButton = await screen.getAllByLabelText('No')[1];
      await userEvent.click(noDestinationTypeRadioButton);
      expect(screen.getByText('We can use the zip of their new duty location:')).toBeTruthy();
      expect(screen.queryByLabelText('Destination type')).toBeNull();
    });
  });

  describe('editing an already existing HHG shipment for retiree/separatee', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultPropsRetirement}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockShipmentWithDestinationType}
          displayDestinationType
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
      expect(screen.getAllByLabelText('Yes')[0]).not.toBeChecked();
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
      expect(screen.getByLabelText('Destination type')).toHaveValue('PLACE_ENTERED_ACTIVE_DUTY');
      expect(screen.queryByTestId('alert')).not.toBeInTheDocument();

      const noDestinationTypeRadioButton = await screen.getAllByLabelText('No')[1];
      await userEvent.click(noDestinationTypeRadioButton);
      expect(screen.getByText('We can use the zip of their HOR, HOS or PLEAD:')).toBeTruthy();
      expect(screen.getByLabelText('Destination type')).toBeVisible();
    });

    describe('shipment address change request', () => {
      it('displays appropriate alerting when an address change is requested', async () => {
        renderWithRouter(
          <ShipmentForm
            {...defaultPropsRetirement}
            isCreatePage={false}
            shipmentType={SHIPMENT_OPTIONS.HHG}
            mtoShipment={{ ...mockShipmentWithDestinationType, ...mockDeliveryAddressUpdate }}
            displayDestinationType
          />,
        );

        const alerts = await screen.findAllByTestId('alert');
        expect(alerts).toHaveLength(2); // Should have 2 alerts shown due to the address update request
        expect(await alerts[0]).toHaveTextContent('Request needs review. See delivery location to proceed.');
        expect(await alerts[1]).toHaveTextContent(
          'Pending delivery location change request needs review. Review request to proceed.',
        );
      });

      it('opens a closeable modal when Review Request is clicked', async () => {
        const user = userEvent.setup();

        const shipmentType = SHIPMENT_OPTIONS.HHG;

        renderWithRouter(
          <ShipmentForm
            {...defaultPropsRetirement}
            isCreatePage={false}
            shipmentType={shipmentType}
            mtoShipment={{ ...mockShipmentWithDestinationType, ...mockDeliveryAddressUpdate, shipmentType }}
            displayDestinationType
          />,
        );

        const queryForModal = () => screen.queryByTestId('modal');

        const reviewRequestLink = await screen.findByRole('button', { name: 'Review request' });

        // confirm the modal is not already present
        expect(queryForModal()).not.toBeInTheDocument();

        // Open the modal
        await user.click(reviewRequestLink);

        await waitFor(() => expect(queryForModal()).toBeInTheDocument());

        // Close the modal
        const modalCancel = within(queryForModal()).queryByText('Cancel');

        expect(modalCancel).toBeInTheDocument();

        await user.click(modalCancel);

        // Confirm the modal has been closed
        expect(queryForModal()).not.toBeInTheDocument();
      });

      it('allows a shipment address update review to be submitted via the modal', async () => {
        const user = userEvent.setup();

        const shipmentType = SHIPMENT_OPTIONS.HHG;
        const eTag = '8c32882e7793d9da88e0fdfd68672e2ead2f';

        renderWithRouter(
          <ShipmentForm
            {...defaultPropsRetirement}
            isCreatePage={false}
            shipmentType={shipmentType}
            mtoShipment={{ ...mockShipmentWithDestinationType, ...mockDeliveryAddressUpdate, eTag }}
            displayDestinationType
          />,
        );

        const queryForModal = () => screen.queryByTestId('modal');
        const findAlerts = async () => screen.findAllByTestId('alert');

        const reviewRequestLink = await screen.findByRole('button', { name: 'Review request' });

        expect(await findAlerts()).toHaveLength(2);

        // Open the modal
        await user.click(reviewRequestLink);
        const modal = queryForModal();

        expect(modal).toBeInTheDocument();

        // Fill and submit
        const approvalQuestion = within(modal).getByRole('group', { name: 'Approve address change?' });
        const approvalYes = within(approvalQuestion).getByRole('radio', { name: 'Yes' });
        const officeRemarks = within(modal).getByLabelText('Office remarks');
        const save = within(modal).getByRole('button', { name: 'Save' });

        const officeRemarksAnswer = 'Here are my remarks from the office';
        await user.click(approvalYes);
        await user.type(officeRemarks, officeRemarksAnswer);
        await user.click(save);

        // Confirm that the request was triggered
        expect(mockMutateFunction).toHaveBeenCalledTimes(1);
        expect(mockMutateFunction).toHaveBeenCalledWith({
          shipmentID: mockShipmentWithDestinationType.id,
          ifMatchETag: eTag,
          body: {
            status: ADDRESS_UPDATE_STATUS.APPROVED,
            officeRemarks: officeRemarksAnswer,
          },
        });
      });
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTS} />);

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
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          TACs={{ HHG: '1234', NTS: '5678' }}
          shipmentType={SHIPMENT_OPTIONS.NTS}
          mtoShipment={mockMtoShipment}
        />,
      );

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
      expect(screen.getByLabelText('1234 (HHG)')).toBeInTheDocument();
      expect(screen.getByText('No SAC code entered.')).toBeInTheDocument();
    });

    it('does not render NTS release-only sections', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTS} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');
      expect(screen.queryByText(/Shipment weight (lbs)/)).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).not.toBeInTheDocument();
    });
  });

  describe('editing an already existing NTS shipment', () => {
    it('pre-fills the Accounting Codes section', async () => {
      renderWithRouter(
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
          shipmentType={SHIPMENT_OPTIONS.NTS}
        />,
      );

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
      expect(screen.getByLabelText('000012345 (HHG)')).toBeChecked();
    });

    it('sends an empty string when clearing LOA types when updating a shipment', async () => {
      const mockSubmitHandler = jest.fn().mockResolvedValue(null);

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          mtoShipment={{
            ...mockMtoShipment,
            tacType: 'NTS',
            sacType: 'HHG',
          }}
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: '000012345', NTS: '2222' }}
          shipmentType={SHIPMENT_OPTIONS.NTS}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      await userEvent.click(screen.getByTestId('clearSelection-sacType'));
      const saveButton = screen.getByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton); //

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({ tacType: 'NTS', sacType: '' }),
          }),
          expect.objectContaining({
            onError: expect.any(Function),
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('does not send undefined LOA types when creating shipment', async () => {
      const mockSubmitHandler = jest.fn().mockResolvedValue(null);

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          mtoShipment={{
            ...mockMtoShipment,
          }}
          shipmentType={SHIPMENT_OPTIONS.NTS}
          submitHandler={mockSubmitHandler}
          isCreatePage
        />,
      );

      await userEvent.type(screen.getByLabelText('Requested pickup date'), '26 Mar 2022');
      await userEvent.click(screen.getByTestId('useCurrentResidence'));

      const saveButton = screen.getByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton); //

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.not.objectContaining({ tacType: expect.any(String), sacType: expect.any(String) }),
          }),
          expect.objectContaining({
            onError: expect.any(Function),
            onSuccess: expect.any(Function),
          }),
        );
      });
    });
  });

  describe('creating a new NTS-release shipment', () => {
    it('renders the NTS-release shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} />);

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
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
    });

    it('renders the NTS release-only sections', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');
      expect(screen.getByText(/Previously recorded weight \(lbs\)/)).toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).toBeInTheDocument();
    });
  });

  describe('as a TOO', () => {
    it('renders the HHG shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} userRole={roleTypes.TOO} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
      expect(screen.getByLabelText('Requested pickup date')).toBeInTheDocument();
      expect(screen.getByText('Pickup location')).toBeInTheDocument();
      expect(screen.getByLabelText('Requested delivery date')).toBeInTheDocument();
      expect(screen.getByText(/Receiving agent/).parentElement).toBeInTheDocument();
      expect(screen.getByText('Customer remarks')).toBeInTheDocument();
      expect(screen.getByText('Counselor remarks')).toBeInTheDocument();
    });

    it('renders the NTS shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTS} userRole={roleTypes.TOO} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByRole('heading', { level: 2, name: 'Vendor' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Storage facility info' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Storage facility address' })).toBeInTheDocument();
    });

    it('renders the NTS release shipment form', async () => {
      renderWithRouter(
        <ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} userRole={roleTypes.TOO} />,
      );

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.getByRole('heading', { level: 2, name: 'Vendor' })).toBeInTheDocument();
    });
  });

  describe('filling the form', () => {
    it('shows an error if the submitHandler returns an error', async () => {
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError();
      });

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const saveButton = screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();

      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(
        await screen.findByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('shows an error if the submitHandler returns an error when editing a PPM', async () => {
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError();
      });
      validatePostalCode.mockImplementation(() => Promise.resolve(false));

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(
        await screen.findByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('shows an error if the submitHandler returns an error when creating a PPM', async () => {
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError();
      });
      validatePostalCode.mockImplementation(() => Promise.resolve(false));

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage
        />,
      );

      await userEvent.type(screen.getByLabelText('Planned departure date'), '26 Mar 2022');
      await userEvent.type(screen.getByLabelText('Origin ZIP'), '90210');
      await userEvent.type(screen.getByLabelText('Destination ZIP'), '90210');
      await userEvent.type(screen.getByLabelText('Estimated PPM weight'), '1000');

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(
        await screen.findByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('saves the update to the counselor remarks when the save button is clicked', async () => {
      const newCounselorRemarks = 'Counselor remarks';

      const expectedPayload = {
        body: {
          customerRemarks: 'mock customer remarks',
          counselorRemarks: newCounselorRemarks,
          hasSecondaryDeliveryAddress: false,
          hasSecondaryPickupAddress: false,
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
          shipmentType: SHIPMENT_OPTIONS.HHG,
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

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockMtoShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );
      const counselorRemarks = await screen.findByLabelText('Counselor remarks');

      await userEvent.clear(counselorRemarks);

      await userEvent.type(counselorRemarks, newCounselorRemarks);

      const saveButton = screen.getByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(expectedPayload, {
          onSuccess: expect.any(Function),
          onError: expect.any(Function),
        });
      });
    });
  });

  describe('external vendor shipment', () => {
    it('shows the TOO an alert', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.NTSR}
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
      renderWithRouter(
        <ShipmentForm
          // SC is default role from test props
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.NTSR}
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

  describe('creating a new PPM shipment', () => {
    it('displays PPM content', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          isCreatePage
          userRole={roleTypes.SERVICES_COUNSELOR}
          mtoShipment={mockMtoShipment}
        />,
      );

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');
    });
  });

  describe('editing an already existing PPM shipment', () => {
    it('renders the PPM shipment form with pre-filled values', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
        />,
      );

      expect(await screen.getByLabelText('Planned departure date')).toHaveValue('01 Apr 2022');
      await userEvent.click(screen.getByLabelText('Use current ZIP'));
      expect(await screen.getByLabelText('Origin ZIP')).toHaveValue(defaultProps.originDutyLocationAddress.postalCode);
      await userEvent.click(screen.getByLabelText('Use ZIP for new duty location'));

      expect(await screen.getByLabelText('Destination ZIP')).toHaveValue(
        defaultProps.newDutyLocationAddress.postalCode,
      );
      expect(screen.getAllByLabelText('Yes')[0]).not.toBeChecked();
      expect(screen.getAllByLabelText('No')[0]).toBeChecked();
      expect(screen.getByLabelText('Estimated PPM weight')).toHaveValue('4,999');
      expect(screen.getAllByLabelText('Yes')[1]).not.toBeChecked();
      expect(screen.getAllByLabelText('No')[1]).toBeChecked();
    });

    it('renders the PPM shipment form with pre-filled requested values for Advance Page', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(await screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.findByText('Estimated incentive: $12,345').toBeInTheDocument);
      expect(screen.getByLabelText('Amount requested')).toHaveValue('4,875');
      expect((await screen.findByText('Maximum advance: $7,407')).toBeInTheDocument);
      expect(screen.getByLabelText('Approve')).toBeChecked();
      expect(screen.getByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');

      await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'mock counselor remarks',
              ppmShipment: expect.objectContaining({
                hasRequestedAdvance: true,
                advanceAmountRequested: 487500,
                advanceStatus: 'APPROVED',
              }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('validates the Advance Page making counselor remarks required when `Advance Requested?` is changed from Yes to No', async () => {
      const ppmShipmentWithoutRemarks = {
        ...mockPPMShipment,
        counselorRemarks: '',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipmentWithoutRemarks}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('Yes')).toBeChecked();
      // Selecting advance not requested
      await userEvent.click(screen.getByLabelText('No'));
      await waitFor(() => {
        expect(screen.getByLabelText('No')).toBeChecked();
        expect(screen.getByLabelText('Yes')).not.toBeChecked();
      });
      const requiredAlerts = screen.getAllByRole('alert');
      expect(requiredAlerts[0]).toHaveTextContent('Required');

      expect(screen.queryByLabelText('Amount requested')).not.toBeInTheDocument();

      await userEvent.type(screen.getByLabelText('Counselor remarks'), 'retirees are not given advances');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeEnabled();
      });

      await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'retirees are not given advances',
              ppmShipment: expect.objectContaining({ hasRequestedAdvance: false }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('validates the Advance Page making counselor remarks required when advance amount is changed', async () => {
      const ppmShipmentWithoutRemarks = {
        ...mockPPMShipment,
        counselorRemarks: '',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipmentWithoutRemarks}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      const advanceAmountInput = screen.getByLabelText('Amount requested');

      expect(advanceAmountInput).toHaveValue('4,875');
      // Edit a requested advance amount
      await userEvent.clear(advanceAmountInput);
      await userEvent.type(advanceAmountInput, '2,000');
      advanceAmountInput.blur();
      await waitFor(() => {
        expect(advanceAmountInput).toHaveValue('2,000');
      });

      const requiredAlerts = screen.getAllByRole('alert');

      expect(requiredAlerts[0]).toHaveTextContent('Required');
    });

    it('marks amount requested input as min of $1 expected when conditionally displayed', async () => {
      renderWithRouter(
        <ShipmentForm {...defaultProps} isCreatePage={false} isAdvancePage shipmentType={SHIPMENT_OPTIONS.PPM} />,
      );

      const inputHasRequestedAdvance = screen.getByLabelText('Yes');

      await userEvent.click(inputHasRequestedAdvance);

      const advanceAmountRequested = screen.getByLabelText('Amount requested');

      await userEvent.type(advanceAmountRequested, '0');

      expect(advanceAmountRequested).toHaveValue('0');

      await waitFor(() => {
        const requiredAlerts = screen.getAllByRole('alert');
        expect(requiredAlerts[0]).toHaveTextContent('Enter an amount $1 or more.');
      });
    });

    it('sets `Counselor Remarks` as required when an advance request is rejected', async () => {
      const ppmShipmentWithoutRemarks = {
        ...mockPPMShipment,
        counselorRemarks: '',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipmentWithoutRemarks}
        />,
        { wrapper: MockProviders },
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(screen.getByLabelText('Approve')).toBeChecked();
      expect(screen.getByLabelText('Reject')).not.toBeChecked();
      // Rejecting advance request
      await userEvent.click(screen.getByLabelText('Reject'));
      await waitFor(() => {
        expect(screen.getByLabelText('Approve')).not.toBeChecked();
        expect(screen.getByLabelText('Reject')).toBeChecked();
      });
      const requiredAlert = screen.getAllByRole('alert');
      expect(requiredAlert[0]).toHaveTextContent('Required');

      await userEvent.type(
        screen.getByLabelText('Counselor remarks'),
        'I, a service counselor, have rejected your advance request',
      );
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeEnabled();
      });

      await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'I, a service counselor, have rejected your advance request',
              ppmShipment: expect.objectContaining({ advanceStatus: 'REJECTED' }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });
  });

  describe('creating a new PPM shipment', () => {
    it('displays PPM content', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          isCreatePage
          userRole={roleTypes.SERVICES_COUNSELOR}
          mtoShipment={mockMtoShipment}
        />,
      );

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');
    });
  });

  const mockPPMShipmentWithSIT = {
    sitEstimatedCost: 123400,
    sitEstimatedWeight: 2345,
    pickupPostalCode: '12345',
    destinationPostalCode: '54321',
    sitLocation: 'DESTINATION',
    sitEstimatedDepartureDate: '2022-10-29',
    sitEstimatedEntryDate: '2022-08-06',
    sitExpected: true,
  };

  const defaultSITProps = {
    ...defaultProps,
    shipmentType: SHIPMENT_OPTIONS.PPM,
    isAdvancePage: true,
    mtoShipment: {
      ...mockMtoShipment,
      ppmShipment: mockPPMShipmentWithSIT,
    },
    userRole: roleTypes.SERVICES_COUNSELOR,
  };

  describe('as a SC, the SIT details block', () => {
    it('displays when SIT is expected', () => {
      renderWithRouter(<ShipmentForm {...defaultSITProps} />);
      expect(screen.getByRole('heading', { level: 2, name: /Storage in transit \(SIT\)/ })).toBeInTheDocument();
    });
    it('does not display when SIT is not expected', () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultSITProps}
          mtoShipment={{
            ...mockMtoShipment,
            ppmShipment: {
              ...mockPPMShipmentWithSIT,
              sitExpected: false,
            },
          }}
        />,
      );
      expect(screen.queryByRole('heading', { level: 2, name: /Storage in transit \(SIT\)/ })).not.toBeInTheDocument();
    });
    it('does not display for TOO', () => {
      renderWithRouter(<ShipmentForm {...defaultSITProps} userRole={roleTypes.TOO} />);
      expect(screen.queryByRole('heading', { level: 2, name: /Storage in transit \(SIT\)/ })).not.toBeInTheDocument();
    });
  });
});
