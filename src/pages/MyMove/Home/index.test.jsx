/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { mount } from 'enzyme';
import moment from 'moment';
import { generatePath } from 'react-router';

import { Home } from './index';

import { MockProviders } from 'testUtils';
import { formatCustomerDate } from 'utils/formatters';
import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';

jest.mock('containers/FlashMessage/FlashMessage', () => {
  const MockFlash = () => <div>Flash message</div>;
  MockFlash.displayName = 'ConnectedFlashMessage';
  return MockFlash;
});

const defaultProps = {
  serviceMember: {
    id: 'testServiceMemberId',
    current_location: {
      transportation_office: {
        name: 'Test Transportation Office Name',
        phone_lines: ['555-555-5555'],
      },
    },
    weight_allotment: {
      total_weight_self: 8000,
      total_weight_self_plus_dependents: 11000,
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
  orders: {
    id: '123',
    new_duty_location: {
      name: 'Test Location',
    },
  },
  updateShipmentList: jest.fn(),
  history: {
    goBack: jest.fn(),
    push: jest.fn(),
  },
  location: {},
  move: {
    id: 'testMoveId',
    status: 'DRAFT',
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const mountHomeWithProviders = (props = {}) => {
  return mount(
    <MockProviders>
      <Home {...defaultProps} {...props} />
    </MockProviders>,
  );
};

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
        { id: '4321', createdAt: moment().add(1, 'days').toISOString(), shipmentType: SHIPMENT_OPTIONS.HHG },
        {
          id: '4322',
          createdAt: moment().add(2, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.PPM,
          ppmShipment: {
            id: '0001',
            hasRequestedAdvance: false,
          },
        },
        { id: '4323', createdAt: moment().add(2, 'days').toISOString(), shipmentType: SHIPMENT_OPTIONS.HHG },
        { id: '4324', createdAt: moment().add(3, 'days').toISOString(), shipmentType: SHIPMENT_OPTIONS.NTS },
        { id: '4325', createdAt: moment().add(4, 'days').toISOString(), shipmentType: SHIPMENT_OPTIONS.NTSR },
        {
          id: '4327',
          createdAt: moment().add(5, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.PPM,
          ppmShipment: {
            id: '0001',
            hasRequestedAdvance: null,
          },
        },
      ],
      orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
      move: { id: 'testMoveId', status: 'DRAFT' },
    };

    it('contains ppm and hhg cards if those shipments exist', async () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

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
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

      const editHHGShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[0].id,
      });

      const hhgShipment = screen.getAllByTestId('shipment-list-item-container')[0];
      const editButton = within(hhgShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      userEvent.click(editButton);

      expect(defaultProps.history.push).toHaveBeenCalledWith(`${editHHGShipmentPath}?shipmentNumber=1`);
    });

    it('handles edit click to edit ppm shipment route', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

      const editPPMShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[1].id,
      });

      const ppmShipment = screen.getAllByTestId('shipment-list-item-container')[1];
      const editButton = within(ppmShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      userEvent.click(editButton);

      expect(defaultProps.history.push).toHaveBeenCalledWith(`${editPPMShipmentPath}?shipmentNumber=1`);
    });

    it('handles edit click to edit nts shipment route', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

      const editNTSShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[3].id,
      });

      const ntsShipment = screen.getAllByTestId('shipment-list-item-container')[3];
      const editButton = within(ntsShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      userEvent.click(editButton);

      expect(defaultProps.history.push).toHaveBeenCalledWith(editNTSShipmentPath);
    });

    it('handles edit click to edit ntsr shipment route', () => {
      render(
        <MockProviders>
          <Home {...defaultProps} {...props} />
        </MockProviders>,
      );

      const editNTSRShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: props.mtoShipments[4].id,
      });

      const ntsrShipment = screen.getAllByTestId('shipment-list-item-container')[4];
      const editButton = within(ntsrShipment).getByRole('button', { name: 'Edit' });
      expect(editButton).toBeInTheDocument();
      userEvent.click(editButton);

      expect(defaultProps.history.push).toHaveBeenCalledWith(editNTSRShipmentPath);
    });
  });

  describe('if the user has complete PPMs', () => {
    const props = {
      mtoShipments: [
        {
          id: '4327',
          createdAt: moment().add(5, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.PPM,
          ppmShipment: {
            id: '0001',
            hasRequestedAdvance: true,
          },
        },
      ],
      orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
      move: { id: 'testMoveId', status: 'DRAFT' },
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
      mtoShipments: [
        {
          id: '4327',
          createdAt: moment().add(5, 'days').toISOString(),
          shipmentType: SHIPMENT_OPTIONS.PPM,
          ppmShipment: {
            id: '0001',
            hasRequestedAdvance: null,
          },
        },
      ],
      orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
      move: { id: 'testMoveId', status: 'DRAFT' },
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
      orders: { testOrder: 'test', new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
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
      orders: { testOrder: 'test', has_dependents: false, new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
    });

    it('renders the correct weight allowance', () => {
      expect(wrapper.text().includes('8,000 lbs.')).toBe(true);
    });
  });

  describe('if the user has orders with dependents', () => {
    const wrapper = mountHomeWithProviders({
      orders: { testOrder: 'test', has_dependents: true, new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
    });

    it('renders the correct weight allowance', () => {
      expect(wrapper.text().includes('11,000 lbs.')).toBe(true);
    });
  });

  describe('if the user has orders and shipments but has not submitted their move', () => {
    const wrapper = mountHomeWithProviders({
      orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
      mtoShipments: [{ id: 'test123', shipmentType: 'HHG' }],
    });

    it('renders the NeedsSubmitMove helper', () => {
      expect(wrapper.find('HelperNeedsSubmitMove').exists()).toBe(true);
    });
  });

  describe('if the user has orders and a ppm but has not submitted their move', () => {
    const wrapper = mountHomeWithProviders({
      orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
      mtoShipments: [{ id: 'test123', shipmentType: 'PPM', ppmShipment: { id: 'ppm', hasRequestedAdvance: false } }],
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
    });

    it('renders the NeedsSubmitMove helper', () => {
      expect(wrapper.find('HelperNeedsSubmitMove').exists()).toBe(true);
    });
  });

  describe('if the user has submitted their move', () => {
    describe('for PPM moves', () => {
      const orders = {
        id: 'testOrder123',
        new_duty_location: {
          name: 'Test Duty Location',
        },
      };
      const uploadedOrderDocuments = [{ id: 'testDocument354', filename: 'testOrder1.pdf' }];
      const move = { id: 'testMoveId', status: 'SUBMITTED' };
      const mtoShipments = [
        { id: 'test123', shipmentType: 'PPM', ppmShipment: { id: 'ppm', hasRequestedAdvance: false } },
      ];

      const wrapper = mountHomeWithProviders({
        orders,
        uploadedOrderDocuments,
        move,
        mtoShipments,
      });
      const props = { ...defaultProps, orders, uploadedOrderDocuments, move, mtoShipments };

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

      it('renders Step 5', () => {
        render(<Home {...props} />);
        expect(screen.getByText('Manage your PPM')).toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for HHG moves (no PPM)', () => {
      const orders = { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } };
      const uploadedOrderDocuments = [{ id: 'testDocument354', filename: 'testOrder1.pdf' }];
      const mtoShipments = [{ id: 'test123', shipmentType: 'HHG' }];
      const move = { id: 'testMoveId', status: 'SUBMITTED' };
      const wrapper = mountHomeWithProviders({
        orders,
        uploadedOrderDocuments,
        mtoShipments,
        move,
      });
      const props = { ...defaultProps, orders, uploadedOrderDocuments, mtoShipments, move };

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

      it('does not render Step 5', () => {
        render(<Home {...props} />);
        expect(screen.queryByText('Manage your PPM')).not.toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for NTS moves (no PPM)', () => {
      const orders = { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } };
      const uploadedOrderDocuments = [{ id: 'testDocument354', filename: 'testOrder1.pdf' }];
      const mtoShipments = [{ id: 'test123', shipmentType: SHIPMENT_OPTIONS.NTS }];
      const move = { id: 'testMoveId', status: 'SUBMITTED' };
      const wrapper = mountHomeWithProviders({
        orders,
        uploadedOrderDocuments,
        mtoShipments,
        move,
      });
      const props = { ...defaultProps, orders, uploadedOrderDocuments, mtoShipments, move };

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

      it('does not render Step 5', () => {
        render(<Home {...props} />);
        expect(screen.queryByText('Manage your PPM')).not.toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for HHG/PPM combo moves', () => {
      const submittedAt = new Date();
      const orders = {
        id: 'testOrder123',
        new_duty_location: {
          name: 'Test Duty Location',
        },
      };
      const uploadedOrderDocuments = [{ id: 'testDocument354', filename: 'testOrder1.pdf' }];
      const move = { id: 'testMoveId', status: 'SUBMITTED', submitted_at: submittedAt };
      const mtoShipments = [
        { id: 'test122', shipmentType: 'HHG' },
        { id: 'test123', shipmentType: 'PPM', ppmShipment: { id: 'ppm', hasRequestedAdvance: false } },
      ];

      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <Home
            {...defaultProps}
            orders={orders}
            uploadedOrderDocuments={uploadedOrderDocuments}
            move={move}
            mtoShipments={mtoShipments}
          />
        </MockProviders>,
      );
      const props = { ...defaultProps, orders, uploadedOrderDocuments, move, mtoShipments };

      it('renders submitted date at step 4', () => {
        expect(wrapper.find('[data-testid="move-submitted-description"]').text()).toBe(
          `Move submitted ${formatCustomerDate(submittedAt)}.Print the legal agreement`,
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

      it('renders Step 5', () => {
        render(<Home {...props} />);
        expect(screen.getByText('Manage your PPM')).toBeInTheDocument();
      });

      it('add shipments button no longer present', () => {
        render(<Home {...props} />);
        expect(screen.queryByRole('button', { name: 'Add another shipment' })).not.toBeInTheDocument();
      });
    });

    describe('for unapproved amended orders', () => {
      const submittedAt = new Date();
      const orders = {
        id: 'testOrder123',
        new_duty_location: {
          name: 'Test Duty Location',
        },
      };
      const uploadedOrderDocuments = [{ id: 'testDocument354', filename: 'testOrder1.pdf' }];
      const uploadedAmendedOrderDocuments = [{ id: 'testDocument987', filename: 'testOrder2.pdf' }];
      const move = { id: 'testMoveId', status: 'APPROVALS REQUESTED', submitted_at: submittedAt };
      const mtoShipments = [
        { id: 'test123', shipmentType: 'PPM', ppmShipment: { id: 'ppm', hasRequestedAdvance: false } },
      ];

      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <Home
            {...defaultProps}
            orders={orders}
            uploadedOrderDocuments={uploadedOrderDocuments}
            uploadedAmendedOrderDocuments={uploadedAmendedOrderDocuments}
            move={move}
            mtoShipments={mtoShipments}
          />
        </MockProviders>,
      );

      it('renders the HelperAmendedOrders helper', () => {
        expect(wrapper.find('HelperAmendedOrders').exists()).toBe(true);
      });
      it('renders the amended orders alert', () => {
        expect(wrapper.find('[data-testid="unapproved-amended-orders-alert"]').exists()).toBe(true);
      });
    });

    describe('for approved amended orders', () => {
      const submittedAt = new Date();
      const orders = {
        id: 'testOrder123',
        new_duty_location: {
          name: 'Test Duty Location',
        },
      };
      const uploadedOrderDocuments = [{ id: 'testDocument354', filename: 'testOrder1.pdf' }];
      const uploadedAmendedOrderDocuments = [{ id: 'testDocument987', filename: 'testOrder2.pdf' }];
      const move = { id: 'testMoveId', status: 'APPROVED', submitted_at: submittedAt };
      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <Home
            {...defaultProps}
            orders={orders}
            uploadedOrderDocuments={uploadedOrderDocuments}
            uploadedAmendedOrderDocuments={uploadedAmendedOrderDocuments}
            move={move}
          />
        </MockProviders>,
      );

      it('does not render the HelperAmendedOrders helper', () => {
        expect(wrapper.find('HelperAmendedOrders').exists()).toBe(false);
      });

      it('does not render the amended orders alert', () => {
        expect(wrapper.find('[data-testid="unapproved-amended-orders-alert"]').exists()).toBe(false);
      });
    });
  });
});
