/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
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
  currentPpm: {},
  loadMTOShipments: jest.fn(),
  orders: {},
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
    const testProps = {
      currentPpm: { id: '12345', createdAt: moment() },
      mtoShipments: [
        { id: '4321', createdAt: moment().add(1, 'days'), shipmentType: SHIPMENT_OPTIONS.HHG },
        { id: '4322', createdAt: moment().subtract(1, 'days'), shipmentType: SHIPMENT_OPTIONS.HHG },
        { id: '4323', createdAt: moment().add(2, 'days'), shipmentType: SHIPMENT_OPTIONS.NTS },
        { id: '4324', createdAt: moment().add(3, 'days'), shipmentType: SHIPMENT_OPTIONS.NTSR },
      ],
    };

    const wrapper = mountHomeWithProviders(testProps);

    it('contains ppm and hhg cards if those shipments exist', () => {
      expect(wrapper.find('ShipmentListItem').length).toBe(5);
      expect(wrapper.find('ShipmentListItem').at(0).text()).toContain('HHG 1');
      expect(wrapper.find('ShipmentListItem').at(1).text()).toContain('PPM');
      expect(wrapper.find('ShipmentListItem').at(2).text()).toContain('HHG 2');
      expect(wrapper.find('ShipmentListItem').at(3).text()).toContain('NTS');
      expect(wrapper.find('ShipmentListItem').at(4).text()).toContain('NTS-release');
    });

    it('handles edit click to edit hhg shipment route', () => {
      const editHHGShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: testProps.mtoShipments[1].id,
      });

      wrapper.find('ShipmentListItem').at(0).find('button').at(1).simulate('click');

      expect(defaultProps.history.push).toHaveBeenCalledWith(`${editHHGShipmentPath}?shipmentNumber=1`);
    });

    it('handles edit click to edit ppm shipment route', () => {
      const editPPMShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: testProps.currentPpm.id,
      });

      wrapper.find('ShipmentListItem').at(1).find('button').at(1).simulate('click');

      expect(defaultProps.history.push).toHaveBeenCalledWith(`${editPPMShipmentPath}?shipmentNumber=1`);
    });

    it('handles edit click to edit nts shipment route', () => {
      const editNTSShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: testProps.mtoShipments[2].id,
      });

      wrapper.find('ShipmentListItem').at(3).find('button').at(1).simulate('click');

      expect(defaultProps.history.push).toHaveBeenCalledWith(editNTSShipmentPath);
    });

    it('handles edit click to edit ntsr shipment route', () => {
      const editNTSRShipmentPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: defaultProps.move.id,
        mtoShipmentId: testProps.mtoShipments[3].id,
      });

      wrapper.find('ShipmentListItem').at(4).find('button').at(1).simulate('click');

      expect(defaultProps.history.push).toHaveBeenCalledWith(editNTSRShipmentPath);
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

  describe('if the user has orders and a currentPpm but has not submitted their move', () => {
    const wrapper = mountHomeWithProviders({
      orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
      uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
      currentPpm: { id: 'testPpm123' },
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
      const currentPpm = { id: 'mockPpm ' };

      const wrapper = mountHomeWithProviders({
        orders,
        uploadedOrderDocuments,
        move,
        currentPpm,
      });

      it('renders the SubmittedMove helper', () => {
        expect(wrapper.find('HelperSubmittedMove').exists()).toBe(true);
      });

      it('Profile step is editable', () => {
        const profileStep = wrapper.find('Step[step="1"]');
        expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      });

      it('Orders Step allows document uploads', () => {
        const ordersStep = wrapper.find('Step[step="2"]');
        expect(ordersStep.prop('editBtnLabel')).toEqual('Upload documents');
      });

      it('renders the SubmittedPPM helper', () => {
        expect(wrapper.find('HelperSubmittedPPM').exists()).toBe(true);
      });
    });

    describe('for HHG moves (no PPM)', () => {
      const wrapper = mountHomeWithProviders({
        orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
        uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
        mtoShipments: [{ id: 'test123', shipmentType: 'HHG' }],
        move: { id: 'testMoveId', status: 'SUBMITTED' },
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
    });

    describe('for NTS moves (no PPM)', () => {
      const wrapper = mountHomeWithProviders({
        orders: { id: 'testOrder123', new_duty_location: { name: 'Test Duty Location' } },
        uploadedOrderDocuments: [{ id: 'testDocument354', filename: 'testOrder1.pdf' }],
        mtoShipments: [{ id: 'test123', shipmentType: 'NTS' }],
        move: { id: 'testMoveId', status: 'SUBMITTED' },
      });

      it('renders the SubmittedMove helper', () => {
        expect(wrapper.find('HelperSubmittedMove').exists()).toBe(true);
      });

      it('Profile step is editable', () => {
        const profileStep = wrapper.find('Step[step="1"]');
        expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      });

      it('Orders Step is not editable, upload documents offered', () => {
        const ordersStep = wrapper.find('Step[step="2"]');
        expect(ordersStep.prop('editBtnLabel')).toEqual('Upload documents');
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
      const currentPpm = { id: 'mockCombo' };
      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <Home
            {...defaultProps}
            orders={orders}
            uploadedOrderDocuments={uploadedOrderDocuments}
            move={move}
            currentPpm={currentPpm}
          />
        </MockProviders>,
      );

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

      it('renders the SubmittedPPM helper', () => {
        expect(wrapper.find('HelperSubmittedPPM').exists()).toBe(true);
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
      const currentPpm = { id: 'mockCombo' };
      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <Home
            {...defaultProps}
            orders={orders}
            uploadedOrderDocuments={uploadedOrderDocuments}
            uploadedAmendedOrderDocuments={uploadedAmendedOrderDocuments}
            move={move}
            currentPpm={currentPpm}
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
      const currentPpm = { id: 'mockCombo' };
      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <Home
            {...defaultProps}
            orders={orders}
            uploadedOrderDocuments={uploadedOrderDocuments}
            uploadedAmendedOrderDocuments={uploadedAmendedOrderDocuments}
            move={move}
            currentPpm={currentPpm}
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
