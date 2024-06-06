/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, within, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { mount } from 'enzyme';
import moment from 'moment';
import { generatePath, MemoryRouter } from 'react-router-dom';
import { v4 } from 'uuid';

import { Home } from './index';

import MOVE_STATUSES from 'constants/moves';
import { ORDERS_TYPE } from 'constants/orders';
import { customerRoutes } from 'constants/routes';
import { shipmentStatuses, ppmShipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCustomerDate } from 'utils/formatters';
import { MockProviders, renderWithRouterProp } from 'testUtils';
import createUpload from 'utils/test/factories/upload';
import { createBaseWeightTicket, createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import {
  createApprovedPPMShipment,
  createPPMShipmentWithFinalIncentive,
  createSubmittedPPMShipment,
} from 'utils/test/factories/ppmShipment';
import { downloadPPMAOAPacket } from 'services/internalApi';

jest.mock('containers/FlashMessage/FlashMessage', () => {
  const MockFlash = () => <div>Flash message</div>;
  MockFlash.displayName = 'ConnectedFlashMessage';
  return MockFlash;
});

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  downloadPPMAOAPacket: jest.fn(),
}));

const defaultProps = {
  serviceMember: {
    id: v4(),
    current_location: {
      transportation_office: {
        name: 'Test Transportation Office Name',
        phone_lines: ['555-555-5555'],
      },
    },
  },
  showLoggedInUser: jest.fn(),
  createServiceMember: jest.fn(),
  getSignedCertification: jest.fn(),
  mtoShipments: [],
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  loadMTOShipments: jest.fn(),
  updateShipmentList: jest.fn(),
  move: {
    id: v4(),
    status: MOVE_STATUSES.DRAFT,
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const orders = {
  id: v4(),
  new_duty_location: {
    id: v4(),
    name: 'Best Location',
  },
  has_dependents: false,
  moves: [defaultProps.move.id],
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  authorizedWeight: 8000,
};

const ordersUpload = createUpload({ fileName: 'testOrders1.pdf' });

const uploadedOrderDocuments = [ordersUpload];

const mtoPPMShipmentId = v4();
const mtoShipmentCreatedDate = new Date();
const ppmShipmentCreatedDate = new Date();

const incompletePPMShipment = {
  id: mtoPPMShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  status: shipmentStatuses.SUBMITTED,
  moveTaskOrderId: defaultProps.move.id,
  ppmShipment: {
    id: v4(),
    shipmentId: mtoPPMShipmentId,
    status: ppmShipmentStatuses.DRAFT,
    expectedDepartureDate: '2022-08-25',
    pickupPostalCode: '90210',
    destinationPostalCode: '30813',
    createdAt: ppmShipmentCreatedDate.toISOString(),
    updatedAt: ppmShipmentCreatedDate.toISOString(),
    eTag: window.btoa(ppmShipmentCreatedDate.toISOString()),
    pickupAddress: {
      streetAddress1: '1 Test Street',
      streetAddress2: '2 Test Street',
      streetAddress3: '3 Test Street',
      city: 'Pickup Test City',
      state: 'NY',
      postalCode: '10001',
    },
    destinationAddress: {
      streetAddress1: '1 Test Street',
      streetAddress2: '2 Test Street',
      streetAddress3: '3 Test Street',
      city: 'Destination Test City',
      state: 'NY',
      postalCode: '11111',
    },
  },
  createdAt: mtoShipmentCreatedDate.toISOString(),
  updatedAt: mtoShipmentCreatedDate.toISOString(),
  eTag: window.btoa(mtoShipmentCreatedDate.toISOString()),
};

const ppmShipmentUpdatedDate = new Date();

const completeUnSubmittedPPM = {
  ...incompletePPMShipment,
  ppmShipment: {
    ...incompletePPMShipment.ppmShipment,
    sitExpected: false,
    estimatedWeight: 4000,
    hasProGear: false,
    estimatedIncentive: 10000000,
    hasRequestedAdvance: true,
    advanceAmountRequested: 30000,
    updatedAt: ppmShipmentUpdatedDate.toISOString(),
    eTag: window.btoa(ppmShipmentUpdatedDate.toISOString()),
  },
};

const submittedPPMShipment = {
  ...completeUnSubmittedPPM,
  ppmShipment: {
    ...completeUnSubmittedPPM.ppmShipment,
    status: ppmShipmentStatuses.SUBMITTED,
  },
};

const approvedDate = new Date();

const approvedPPMShipment = {
  ...submittedPPMShipment,
  status: shipmentStatuses.APPROVED,
  ppmShipment: {
    ...submittedPPMShipment.ppmShipment,
    status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
    actualMoveDate: null,
    actualPickupPostalCode: null,
    actualDestinationPostalCode: null,
    hasReceivedAdvance: null,
    advanceAmountReceived: null,
    weightTickets: [],
    approvedAt: approvedDate.toISOString(),
    updatedAt: approvedDate.toISOString(),
    eTag: window.btoa(approvedDate.toISOString()),
    pickupAddress: {
      streetAddress1: '1 Test Street',
      streetAddress2: '2 Test Street',
      streetAddress3: '3 Test Street',
      city: 'Pickup Test City',
      state: 'NY',
      postalCode: '10001',
    },
    destinationAddress: {
      streetAddress1: '1 Test Street',
      streetAddress2: '2 Test Street',
      streetAddress3: '3 Test Street',
      city: 'Destination Test City',
      state: 'NY',
      postalCode: '11111',
    },
  },
  updatedAt: approvedDate.toISOString(),
  eTag: window.btoa(approvedDate.toISOString()),
};

const ppmShipmentWithActualShipmentInfo = {
  ...approvedPPMShipment,
  ppmShipment: {
    ...approvedPPMShipment.ppmShipment,
    actualMoveDate: approvedPPMShipment.ppmShipment.expectedDepartureDate,
    actualPickupPostalCode: approvedPPMShipment.ppmShipment.pickupPostalCode,
    actualDestinationPostalCode: approvedPPMShipment.ppmShipment.destinationPostalCode,
    hasReceivedAdvance: approvedPPMShipment.ppmShipment.hasRequestedAdvance,
    advanceAmountReceived: approvedPPMShipment.ppmShipment.advanceAmountRequested,
  },
};

const ppmShipmentWithIncompleteWeightTicket = {
  ...ppmShipmentWithActualShipmentInfo,
  ppmShipment: {
    ...ppmShipmentWithActualShipmentInfo.ppmShipment,
    weightTickets: [
      createBaseWeightTicket(
        { serviceMemberId: defaultProps.serviceMember.id },
        { ppmShipmentId: ppmShipmentWithActualShipmentInfo.id },
      ),
    ],
  },
};

const ppmShipmentWithCompleteWeightTicket = {
  ...ppmShipmentWithIncompleteWeightTicket,
  ppmShipment: {
    ...ppmShipmentWithIncompleteWeightTicket.ppmShipment,
    weightTickets: [
      createCompleteWeightTicket(
        { serviceMemberId: defaultProps.serviceMember.id },
        { ppmShipmentId: ppmShipmentWithActualShipmentInfo.id },
      ),
    ],
  },
};

const approvedAdvancePPMShipment = {
  ...incompletePPMShipment,
  ppmShipment: {
    ...incompletePPMShipment.ppmShipment,
    sitExpected: false,
    estimatedWeight: 4000,
    hasProGear: false,
    estimatedIncentive: 10000000,
    hasRequestedAdvance: true,
    advanceAmountRequested: 30000,
    advanceStatus: 'APPROVED',
    status: ppmShipmentStatuses.SUBMITTED,
    updatedAt: ppmShipmentUpdatedDate.toISOString(),
    eTag: window.btoa(ppmShipmentUpdatedDate.toISOString()),
  },
};

const rejectedAdvancePPMShipment = {
  ...incompletePPMShipment,
  ppmShipment: {
    ...incompletePPMShipment.ppmShipment,
    sitExpected: false,
    estimatedWeight: 4000,
    hasProGear: false,
    estimatedIncentive: 10000000,
    hasRequestedAdvance: true,
    advanceAmountRequested: 30000,
    advanceStatus: 'REJECTED',
    status: ppmShipmentStatuses.SUBMITTED,
    updatedAt: ppmShipmentUpdatedDate.toISOString(),
    eTag: window.btoa(ppmShipmentUpdatedDate.toISOString()),
  },
};

const mountHomeWithProviders = (props = {}) => {
  return mount(
    <MockProviders>
      <Home {...defaultProps} {...props} />
    </MockProviders>,
  );
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('Home component', () => {
  describe('with default props', () => {
    const wrapper = mountHomeWithProviders();

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(4);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('Profile Step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
    });
  });

  describe('contents of Step 3', () => {
    const props = {
      mtoShipments: [
        {
          id: v4(),
          createdAt: moment(completeUnSubmittedPPM).subtract(1, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.HHG,
        },
        completeUnSubmittedPPM,
        {
          id: v4(),
          createdAt: moment(completeUnSubmittedPPM).add(2, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.HHG,
        },
        {
          id: v4(),
          createdAt: moment(completeUnSubmittedPPM).add(3, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.NTS,
        },
        {
          id: v4(),
          createdAt: moment(completeUnSubmittedPPM).add(4, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.NTSR,
        },
        {
          ...completeUnSubmittedPPM,
          id: v4(),
          createdAt: moment(completeUnSubmittedPPM).add(5, 'days').toISOString(),
        },
      ],
      orders,
      uploadedOrderDocuments,
    };

    it('contains ppm and hhg cards if those shipments exist', async () => {
      renderWithRouterProp(<Home {...defaultProps} {...props} />);

      const shipmentListItems = screen.getAllByTestId('shipment-list-item-container');
      expect(shipmentListItems.length).toEqual(6);
      expect(shipmentListItems[0]).toHaveTextContent('HHG 1');
      expect(shipmentListItems[1]).toHaveTextContent('PPM');
      expect(shipmentListItems[2]).toHaveTextContent('HHG 2');
      expect(shipmentListItems[3]).toHaveTextContent('NTS');
      expect(shipmentListItems[4]).toHaveTextContent('NTS-release');
      expect(shipmentListItems[5]).toHaveTextContent('PPM 2');
    });

    it('handles edit click to edit hhg shipment route', async () => {
      renderWithRouterProp(<Home {...defaultProps} {...props} />, { navigate: mockNavigate });

      const editHHGShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[0].id,
      });

      const hhgShipment = screen.getAllByTestId('shipment-list-item-container')[0];
      const editButton = within(hhgShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      await userEvent.click(editButton);

      expect(mockNavigate).toHaveBeenCalledWith(`${editHHGShipmentPath}?shipmentNumber=1`);
    });

    it('handles edit click to edit ppm shipment route', async () => {
      renderWithRouterProp(<Home {...defaultProps} {...props} />, { navigate: mockNavigate });

      const editPPMShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[1].id,
      });

      const ppmShipment = screen.getAllByTestId('shipment-list-item-container')[1];
      const editButton = within(ppmShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      await userEvent.click(editButton);

      expect(mockNavigate).toHaveBeenCalledWith(`${editPPMShipmentPath}?shipmentNumber=1`);
    });

    it('handles edit click to edit nts shipment route', async () => {
      renderWithRouterProp(<Home {...defaultProps} {...props} />, { navigate: mockNavigate });

      const editNTSShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[3].id,
      });

      const ntsShipment = screen.getAllByTestId('shipment-list-item-container')[3];
      const editButton = within(ntsShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      await userEvent.click(editButton);

      expect(mockNavigate).toHaveBeenCalledWith(editNTSShipmentPath);
    });

    it('handles edit click to edit ntsr shipment route', async () => {
      renderWithRouterProp(<Home {...defaultProps} {...props} />, { navigate: mockNavigate });

      const editNTSRShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[4].id,
      });

      const ntsrShipment = screen.getAllByTestId('shipment-list-item-container')[4];
      const editButton = within(ntsrShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      await userEvent.click(editButton);

      expect(mockNavigate).toHaveBeenCalledWith(editNTSRShipmentPath);
    });
  });

  describe('if user has submitted PPM with signed agreement', () => {
    const props = {
      mtoShipments: [createSubmittedPPMShipment()],
      orders,
      uploadedOrderDocuments,
    };

    it('gets the data for the legal agreement', () => {
      const mockGetSignedCertification = jest.fn();
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} getSignedCertification={mockGetSignedCertification} />
        </MockProviders>,
      );
      expect(mockGetSignedCertification).toHaveBeenCalledTimes(1);
    });
  });

  describe('if the user has complete PPMs', () => {
    const props = {
      mtoShipments: [completeUnSubmittedPPM],
      orders,
      uploadedOrderDocuments,
    };

    it('does not display incomplete for a complete PPM', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );
      expect(screen.getAllByTestId('shipment-list-item-container')[0]).not.toHaveTextContent('Incomplete');
    });

    it('does not disable the review and submit button', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );
      expect(screen.getByTestId('review-and-submit-btn')).not.toBeDisabled();
    });
  });

  describe('if the user has incomplete PPMs', () => {
    const props = {
      mtoShipments: [incompletePPMShipment],
      orders,
      uploadedOrderDocuments,
    };

    it('displays incomplete for an incomplete PPM', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

      expect(screen.getAllByTestId('shipment-list-item-container')[0]).toHaveTextContent('Incomplete');
    });

    it('disables the review and submit button', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

      expect(screen.getByTestId('review-and-submit-btn')).toBeDisabled();
    });
  });

  describe('if the user does not have orders', () => {
    const wrapper = mountHomeWithProviders();

    it('renders the NeedsOrders helper', () => {
      expect(wrapper.find('HelperNeedsOrders').exists()).toBe(true);
    });

    it('Orders Step is not editable', () => {
      const ordersStep = wrapper.find('Step[step="2"]');
      expect(ordersStep.prop('editBtnLabel')).toEqual('');
    });
  });

  describe('if the user has orders but not shipments', () => {
    const wrapper = mountHomeWithProviders({
      orders,
      uploadedOrderDocuments,
    });

    it('renders the NeedsShipment helper', () => {
      expect(wrapper.find('HelperNeedsShipment').exists()).toBe(true);
    });

    it('Orders Step is editable', () => {
      const ordersStep = wrapper.find('Step[step="2"]');
      expect(ordersStep.prop('editBtnLabel')).toEqual('Edit');
    });
  });

  describe('if the user has orders with no dependents', () => {
    const wrapper = mountHomeWithProviders({
      orders,
      uploadedOrderDocuments,
    });

    it('renders the correct weight allowance', () => {
      expect(wrapper.text().includes('8,000 lbs.')).toBe(true);
    });
  });

  describe('if the user has orders with dependents', () => {
    const wrapper = mountHomeWithProviders({
      orders: { ...orders, has_dependents: true, authorizedWeight: 11000 },
      uploadedOrderDocuments,
    });

    it('renders the correct weight allowance', () => {
      expect(wrapper.text().includes('11,000 lbs.')).toBe(true);
    });
  });

  describe('if the user has orders and shipments but has not submitted their move', () => {
    const wrapper = mountHomeWithProviders({
      orders,
      uploadedOrderDocuments,
      mtoShipments: [{ id: v4(), shipmentType: SHIPMENT_OPTIONS.HHG }],
    });

    it('renders the NeedsSubmitMove helper', () => {
      expect(wrapper.find('HelperNeedsSubmitMove').exists()).toBe(true);
    });
  });

  describe('if the user has orders and a ppm but has not submitted their move', () => {
    const wrapper = mountHomeWithProviders({
      orders,
      mtoShipments: [completeUnSubmittedPPM],
      uploadedOrderDocuments,
    });

    it('renders the NeedsSubmitMove helper', () => {
      expect(wrapper.find('HelperNeedsSubmitMove').exists()).toBe(true);
    });
  });

  describe('if the user has submitted their move', () => {
    const propUpdates = {
      orders,
      uploadedOrderDocuments,
      move: { ...defaultProps.move, status: MOVE_STATUSES.SUBMITTED, submitted_at: new Date().toISOString() },
    };

    describe('for PPM moves', () => {
      const mtoShipments = [submittedPPMShipment];

      const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });

      const props = { ...defaultProps, ...propUpdates, mtoShipments };

      it('renders the SubmittedMove helper', () => {
        expect(wrapper.find('HelperSubmittedMove').exists()).toBe(true);
      });

      it('Profile step is editable', () => {
        const profileStep = wrapper.find('Step[step="1"]');
        expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      });

      it('Orders Step is not editable, upload documents is offered', () => {
        const ordersStep = wrapper.find('Step[step="2"]');
        expect(ordersStep.prop('editBtnLabel')).toEqual('Upload documents');
      });

      it('renders Manage your PPM Step', () => {
        render(<Home {...props} />);
        expect(screen.getByText('Manage your PPM')).toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for advance request approved PPM', () => {
      it('renders advance request submitted for PPM', () => {
        const mtoShipments = [submittedPPMShipment];
        const props = { ...defaultProps, ...propUpdates, mtoShipments };
        render(<Home {...props} />);
        expect(screen.getByText('Advance request submitted')).toBeInTheDocument();
      });

      it('renders advance request submitted for PPM', () => {
        const mtoShipments = [approvedAdvancePPMShipment];
        const props = { ...defaultProps, ...propUpdates, mtoShipments };
        render(
          <MemoryRouter>
            <Home {...props} />
          </MemoryRouter>,
        );
        expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();
      });

      it('renders advance request reviewed with 1 approved PPM', () => {
        const mtoShipments = [approvedAdvancePPMShipment];
        const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });
        const advanceStep = wrapper.find('Step[step="5"]');
        expect(advanceStep.prop('completedHeaderText')).toEqual('Advance request reviewed');

        const props = { ...defaultProps, ...propUpdates, mtoShipments };
        render(
          <MemoryRouter>
            <Home {...props} />
          </MemoryRouter>,
        );
        expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();
      });

      it('renders advance request reviewed for approved advance for PPM with HHG', () => {
        const mtoShipments = [{ id: v4(), shipmentType: SHIPMENT_OPTIONS.HHG }, approvedAdvancePPMShipment];
        const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });
        const advanceStep = wrapper.find('Step[step="5"]');
        expect(advanceStep.prop('completedHeaderText')).toEqual('Advance request reviewed');

        const props = { ...defaultProps, ...propUpdates, mtoShipments };
        render(
          <MemoryRouter>
            <Home {...props} />
          </MemoryRouter>,
        );
        expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();
      });

      it('renders advance request reviewed with 1 approved and 1 rejected advance', () => {
        const mtoShipments = [approvedAdvancePPMShipment, rejectedAdvancePPMShipment];
        const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });
        const advanceStep = wrapper.find('Step[step="5"]');

        expect(advanceStep.prop('completedHeaderText')).toEqual('Advance request reviewed');

        const props = { ...defaultProps, ...propUpdates, mtoShipments };
        render(
          <MemoryRouter>
            <Home {...props} />
          </MemoryRouter>,
        );
        expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();
        expect(screen.getByText('Advance request denied')).toBeInTheDocument();
      });

      it('renders advance request denied for PPM', () => {
        const mtoShipments = [rejectedAdvancePPMShipment];
        const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });
        const advanceStep = wrapper.find('Step[step="5"]');

        expect(advanceStep.prop('completedHeaderText')).toEqual('Advance request denied');
      });

      it('Download AOA Packet PPM - Error', async () => {
        downloadPPMAOAPacket.mockRejectedValue({
          response: { body: { title: 'Error title', detail: 'Error detail' } },
        });

        const mtoShipments = [approvedAdvancePPMShipment];
        const props = { ...defaultProps, ...propUpdates, mtoShipments };
        render(
          <MemoryRouter>
            <Home {...props} />
          </MemoryRouter>,
        );
        expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();

        const downloadAOAButton = screen.getByText('Download AOA Paperwork (PDF)');
        expect(downloadAOAButton).toBeInTheDocument();
        await userEvent.click(downloadAOAButton);

        await waitFor(() => {
          expect(
            screen.getByText(/Something went wrong downloading PPM paperwork./, { exact: false }),
          ).toBeInTheDocument();
          expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
        });
      });

      it('Download AOA Packet PPM - Success', async () => {
        const mockResponse = {
          ok: true,
          headers: {
            'content-disposition': 'filename="test.pdf"',
          },
          status: 200,
          data: null,
        };
        downloadPPMAOAPacket.mockImplementation(() => Promise.resolve(mockResponse));

        const mtoShipments = [approvedAdvancePPMShipment];
        const props = { ...defaultProps, ...propUpdates, mtoShipments };

        render(
          <MemoryRouter>
            <Home {...props} />
          </MemoryRouter>,
        );

        expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();

        const downloadAOAButton = screen.getByText('Download AOA Paperwork (PDF)');
        expect(downloadAOAButton).toBeInTheDocument();

        await userEvent.click(downloadAOAButton);

        await waitFor(() => {
          expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
        });
      });
    });

    describe('for HHG moves (no PPM)', () => {
      const mtoShipments = [{ id: v4(), shipmentType: SHIPMENT_OPTIONS.HHG }];

      const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });

      const props = { ...defaultProps, ...propUpdates, mtoShipments };

      it('renders the SubmittedMove helper', () => {
        expect(wrapper.find('HelperSubmittedMove').exists()).toBe(true);
      });

      it('Profile step is editable', () => {
        const profileStep = wrapper.find('Step[step="1"]');
        expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      });

      it('Orders Step is not editable, upload documents is offered', () => {
        const ordersStep = wrapper.find('Step[step="2"]');
        expect(ordersStep.prop('editBtnLabel')).toEqual('Upload documents');
      });

      it('does not render Manage your PPM Step', () => {
        render(<Home {...props} />);
        expect(screen.queryByText('Manage your PPM')).not.toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for NTS moves (no PPM)', () => {
      const mtoShipments = [{ id: v4(), shipmentType: SHIPMENT_OPTIONS.NTS }];

      const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });

      const props = { ...defaultProps, ...propUpdates, mtoShipments };

      it('renders the SubmittedMove helper', () => {
        expect(wrapper.find('HelperSubmittedMove').exists()).toBe(true);
      });

      it('Profile step is editable', () => {
        const profileStep = wrapper.find('Step[step="1"]');
        expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      });

      it('Orders Step is not editable, upload documents is offered', () => {
        const ordersStep = wrapper.find('Step[step="2"]');
        expect(ordersStep.prop('editBtnLabel')).toEqual('Upload documents');
      });

      it('does not render Manage your PPM Step', () => {
        render(<Home {...props} />);
        expect(screen.queryByText('Manage your PPM')).not.toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for HHG/PPM combo moves', () => {
      const mtoShipments = [{ id: v4(), shipmentType: SHIPMENT_OPTIONS.HHG }, submittedPPMShipment];

      const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments });

      const props = { ...defaultProps, ...propUpdates, mtoShipments };

      it('renders submitted date at step 4', () => {
        expect(wrapper.find('[data-testid="move-submitted-description"]').text()).toBe(
          `Move submitted ${formatCustomerDate(propUpdates.move.submitted_at)}.Print the legal agreement`,
        );
      });

      it('renders secondary button when step 4 is completed', () => {
        expect(wrapper.find('[data-testid="review-and-submit-btn"]').at(1).hasClass('usa-button--secondary')).toBe(
          true,
        );
      });

      it('renders the SubmittedMove helper', () => {
        expect(wrapper.find('HelperSubmittedMove').exists()).toBe(true);
      });

      it('Profile step is editable', () => {
        const profileStep = wrapper.find('Step[step="1"]');
        expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      });

      it('Orders Step is not editable, upload documents is offered', () => {
        const ordersStep = wrapper.find('Step[step="2"]');
        expect(ordersStep.prop('editBtnLabel')).toEqual('Upload documents');
      });

      it('renders Manage your PPM Step', () => {
        render(<Home {...props} />);
        expect(screen.getByText('Manage your PPM')).toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    const amendedOrdersUploadCreateDate = moment(ordersUpload.created_at).add(1, 'days');

    const uploadedAmendedOrderDocuments = [
      createUpload({ fileName: 'testOrder2.pdf', createdAtDate: amendedOrdersUploadCreateDate }),
    ];

    describe('for unapproved amended orders', () => {
      const move = { ...propUpdates.move, status: MOVE_STATUSES.APPROVALS_REQUESTED };

      const mtoShipments = [submittedPPMShipment];

      const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments, move, uploadedAmendedOrderDocuments });

      it('renders the HelperAmendedOrders helper', () => {
        expect(wrapper.find('HelperAmendedOrders').exists()).toBe(true);
      });
      it('renders the amended orders alert', () => {
        expect(wrapper.find('[data-testid="unapproved-amended-orders-alert"]').exists()).toBe(true);
      });
    });

    describe('for approved amended orders', () => {
      const move = { ...propUpdates.move, status: MOVE_STATUSES.APPROVED };

      const mtoShipments = [submittedPPMShipment];

      const wrapper = mountHomeWithProviders({ ...propUpdates, mtoShipments, move, uploadedAmendedOrderDocuments });

      it('does not render the HelperAmendedOrders helper', () => {
        expect(wrapper.find('HelperAmendedOrders').exists()).toBe(false);
      });

      it('does not render the amended orders alert', () => {
        expect(wrapper.find('[data-testid="unapproved-amended-orders-alert"]').exists()).toBe(false);
      });
    });
  });

  describe('if the user has submitted a move with a ppm shipment that has been approved', () => {
    const props = {
      ...defaultProps,
      move: { ...defaultProps.move, status: MOVE_STATUSES.APPROVED, submitted_at: new Date().toISOString() },
      orders,
      uploadedOrderDocuments,
    };

    it('it will render the correct helper text', () => {
      const propsForApprovedShipment = {
        ...props,
        mtoShipments: [createApprovedPPMShipment()],
      };

      render(
        <MockProviders>
          <Home {...propsForApprovedShipment} />
        </MockProviders>,
      );
      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Your move is in progress');
    });

    describe('then when the Upload PPM Documents button is clicked', () => {
      const ppmShipmentWithMultipleIncompleteWeightTickets = {
        ...ppmShipmentWithIncompleteWeightTicket,
        ppmShipment: {
          ...ppmShipmentWithIncompleteWeightTicket.ppmShipment,
          weightTickets: [
            ...ppmShipmentWithIncompleteWeightTicket.ppmShipment.weightTickets,
            createBaseWeightTicket(
              { serviceMemberId: defaultProps.serviceMember.id },
              { ppmShipmentId: ppmShipmentWithIncompleteWeightTicket.id },
            ),
          ],
        },
      };

      const ppmShipmentWithMultipleWeightTickets = {
        ...ppmShipmentWithCompleteWeightTicket,
        ppmShipment: {
          ...ppmShipmentWithCompleteWeightTicket.ppmShipment,
          weightTickets: [
            ...ppmShipmentWithCompleteWeightTicket.ppmShipment.weightTickets,
            createBaseWeightTicket(
              { serviceMemberId: defaultProps.serviceMember.id },
              { ppmShipmentId: ppmShipmentWithActualShipmentInfo.id },
            ),
          ],
        },
      };

      it.each([
        [
          'About Your PPM page if no actual shipment info has been input',
          [approvedPPMShipment],
          generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
            moveId: props.move.id,
            mtoShipmentId: approvedPPMShipment.id,
          }),
        ],
        [
          'Weight Ticket page if weight ticket info is missing',
          [ppmShipmentWithActualShipmentInfo],
          generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
            moveId: props.move.id,
            mtoShipmentId: ppmShipmentWithActualShipmentInfo.id,
          }),
        ],
        [
          'Weight Ticket page if weight ticket info is incomplete',
          [ppmShipmentWithIncompleteWeightTicket],
          generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
            moveId: props.move.id,
            mtoShipmentId: ppmShipmentWithIncompleteWeightTicket.id,
            weightTicketId: ppmShipmentWithIncompleteWeightTicket.ppmShipment.weightTickets[0].id,
          }),
        ],
        [
          'Weight Ticket page for the first weight ticket if there are multiple but none are complete',
          [ppmShipmentWithMultipleIncompleteWeightTickets],
          generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
            moveId: props.move.id,
            mtoShipmentId: ppmShipmentWithMultipleIncompleteWeightTickets.id,
            weightTicketId: ppmShipmentWithMultipleIncompleteWeightTickets.ppmShipment.weightTickets[0].id,
          }),
        ],
        [
          'Review page if weight ticket info is complete',
          [ppmShipmentWithCompleteWeightTicket],
          generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
            moveId: props.move.id,
            mtoShipmentId: ppmShipmentWithCompleteWeightTicket.id,
          }),
        ],
        [
          'Review page if at least one weight ticket is completely filled out',
          [ppmShipmentWithMultipleWeightTickets],
          generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
            moveId: props.move.id,
            mtoShipmentId: ppmShipmentWithMultipleWeightTickets.id,
          }),
        ],
      ])('will route the user to the %s', async (scenarioDescription, mtoShipments, expectedRoute) => {
        renderWithRouterProp(<Home {...props} mtoShipments={mtoShipments} />, { navigate: mockNavigate });

        await userEvent.click(screen.getByRole('button', { name: 'Upload PPM Documents' }));

        expect(mockNavigate).toHaveBeenCalledTimes(1);
        expect(mockNavigate).toHaveBeenCalledWith(expectedRoute);
      });
    });
  });

  describe('if PPM closeout is complete', () => {
    const props = {
      ...defaultProps,
      move: { ...defaultProps.move, status: MOVE_STATUSES.APPROVED, submitted_at: new Date().toISOString() },
      orders,
      uploadedOrderDocuments,
    };

    it('will render the correct helper text', () => {
      const propsForCloseoutCompleteShipment = {
        ...props,
        mtoShipments: [
          createPPMShipmentWithFinalIncentive({
            ppmShipment: {
              status: ppmShipmentStatuses.NEEDS_CLOSEOUT,
              pickupAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Pickup Test City',
                state: 'NY',
                postalCode: '10001',
              },
              destinationAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Destination Test City',
                state: 'NY',
                postalCode: '11111',
              },
            },
          }),
        ],
      };

      render(
        <MockProviders>
          <Home {...propsForCloseoutCompleteShipment} />
        </MockProviders>,
      );
      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(
        'Someone will review all of your PPM documentation',
      );
    });
  });
});
