import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import {
  formatAccountingCode,
  formatAddress,
  formatAgent,
  formatCustomerDestination,
  formatPaymentRequestAddressString,
  formatPaymentRequestReviewAddressString,
  getShipmentModificationType,
  retrieveSAC,
  retrieveTAC,
} from './shipmentDisplay';

import { LOA_TYPE } from 'shared/constants';
import { shipmentStatuses, shipmentModificationTypes } from 'constants/shipments';

describe('shipmentDisplay utils', () => {
  describe('formatAccountingCode', () => {
    it('formats the accounting code with the type in parentheses', () => {
      const formattedAccountingCode = formatAccountingCode(1234, 'SAC');
      expect(formattedAccountingCode).toEqual('1234 (SAC)');
    });
  });

  describe('retrieveTAC', () => {
    it.each([
      [LOA_TYPE.HHG, { tac: '1234', ntsTac: '456' }, '1234'],
      [LOA_TYPE.NTS, { tac: '1234', ntsTac: '456' }, '456'],
      ['default', { tac: '1234', ntsTac: '456' }, '1234'],
    ])('returns the correct tac when provided the tac type of %s', (tacType, ordersLOA, expectedTAC) => {
      const retrievedTAC = retrieveTAC(tacType, ordersLOA, expectedTAC);
      expect(retrievedTAC).toEqual(expectedTAC);
    });
  });

  describe('retrieveSAC', () => {
    it.each([
      [LOA_TYPE.HHG, { sac: '1234', ntsSac: '456' }, '1234'],
      [LOA_TYPE.NTS, { sac: '1234', ntsSac: '456' }, '456'],
      ['default', { sac: '1234', ntsTac: '456' }, '1234'],
    ])('returns the correct SAC when provided the tac type of %s', (sacType, ordersLOA, expectedSAC) => {
      const retrievedSAC = retrieveSAC(sacType, ordersLOA, expectedSAC);
      expect(retrievedSAC).toEqual(expectedSAC);
    });
  });

  describe('formatAddress', () => {
    describe('all address parts provided', () => {
      const shipmentAddress = {
        streetAddress1: '555 Main Street',
        streetAddress2: 'P.O. Box 9882',
        city: 'Celebration',
        state: 'FL',
        postalCode: '34747',
      };
      const component = mount(formatAddress(shipmentAddress));
      it('includes full address with comma seperator', () => {
        expect(component.at(0).text()).toEqual('555 Main Street');
        // Must use the character code for nbsp
        expect(component.at(1).text()).toEqual(',\xa0');
        expect(component.at(2).text()).toEqual('P.O. Box 9882');
        expect(component.at(3).text()).toEqual(',\xa0');
        expect(component.at(4).text()).toEqual('Celebration, FL 34747');
      });
    });
    describe('street address is missing', () => {
      const shipmentAddress = {
        city: 'Celebration',
        state: 'FL',
        postalCode: '34747',
      };
      const component = mount(formatAddress(shipmentAddress));
      it('formats as single line', () => {
        expect(component.text()).toEqual('Celebration, FL 34747');
      });
    });
    describe('postal code only', () => {
      const shipmentAddress = {
        postalCode: '34747',
      };
      const component = mount(formatAddress(shipmentAddress));
      it('omits city and state', () => {
        expect(component.text()).toEqual('34747');
      });
    });
  });
  describe('formatAgent', () => {
    it('shows entire agent', () => {
      const agent = {
        firstName: 'John',
        lastName: 'Johnson',
        phone: '(405) 555-1234',
        email: 'johnson@example.com',
      };
      render(formatAgent(agent));
      expect(screen.getByText(`${agent.firstName} ${agent.lastName}`)).toBeInTheDocument();
      expect(screen.getByText(`${agent.phone}`)).toBeInTheDocument();
      expect(screen.getByText(`${agent.email}`)).toBeInTheDocument();
    });

    it('shows just first name and last name', () => {
      const agent = {
        firstName: 'Jane',
        lastName: 'Jamison',
      };
      render(formatAgent(agent));
      expect(screen.getByText(`${agent.firstName} ${agent.lastName}`)).toBeInTheDocument();
      expect(screen.queryByText(`${agent.phone}`)).toBeFalsy();
      expect(screen.queryByText(`${agent.email}`)).toBeFalsy();
    });
  });
  describe('formatDestination', () => {
    it('shows entire address', () => {
      const destinationLocation = {
        streetAddress1: '123 Any Street',
        streetAddress2: 'Apt 4',
        city: 'Los Angeles',
        state: 'CA',
        postalCode: '111111',
      };
      const wrapper = mount(formatCustomerDestination(destinationLocation));
      expect(wrapper.at(0).text()).toEqual(destinationLocation.streetAddress1);
      expect(wrapper.at(2).text()).toEqual(destinationLocation.streetAddress2);
      expect(wrapper.at(4).text()).toEqual(destinationLocation.city);
      expect(wrapper.at(6).text()).toEqual(destinationLocation.state);
      expect(wrapper.at(8).text()).toEqual(destinationLocation.postalCode);
    });

    it('shows postalCode if address is not provided', () => {
      expect(formatCustomerDestination(null, '11111')).toBe('11111');
    });
  });
  describe('formatPaymentRequestAddressString', () => {
    const pickupAddress = {
      city: 'Princeton',
      state: 'NJ',
      postalCode: '08540',
    };
    const destinationAddress = { city: 'Boston', state: 'MA', postalCode: '02101' };
    it('shows expected string when both pickupAddress and destinationAddress are present', () => {
      const wrapper = mount(formatPaymentRequestAddressString(pickupAddress, destinationAddress));

      expect(wrapper.at(0).text()).toEqual(pickupAddress.city);
      expect(wrapper.at(2).text()).toEqual(pickupAddress.state);
      expect(wrapper.at(4).text()).toEqual(pickupAddress.postalCode);
      expect(wrapper.at(8).text()).toEqual(destinationAddress.city);
      expect(wrapper.at(10).text()).toEqual(destinationAddress.state);
      expect(wrapper.at(12).text()).toEqual(destinationAddress.postalCode);
    });
    it('shows expected string when both pickupAddress is missing', () => {
      const wrapper = mount(formatPaymentRequestAddressString(undefined, destinationAddress));

      expect(wrapper.at(0).text()).toEqual('TBD ');
      expect(wrapper.at(3).text()).toEqual(destinationAddress.city);
      expect(wrapper.at(5).text()).toEqual(destinationAddress.state);
      expect(wrapper.at(7).text()).toEqual(destinationAddress.postalCode);
    });
    it('shows expected string when both destinationAddress is missing', () => {
      const wrapper = mount(formatPaymentRequestAddressString(pickupAddress, undefined));

      expect(wrapper.at(0).text()).toEqual(pickupAddress.city);
      expect(wrapper.at(2).text()).toEqual(pickupAddress.state);
      expect(wrapper.at(4).text()).toEqual(pickupAddress.postalCode);
      expect(wrapper.at(8).text()).toEqual(`TBD`);
    });
  });

  describe('formatPaymentRequestAddressString', () => {
    const address = {
      streetAddress1: '123 Any Street',
      streetAddress2: 'Apt 4',
      city: 'Los Angeles',
      state: 'CA',
      postalCode: '111111',
    };

    it('shows expected string when an address is present', () => {
      const addressString = formatPaymentRequestReviewAddressString(address);
      expect(addressString).toEqual('Los Angeles, CA 111111');
    });

    it('shows expected string when an address is not present', () => {
      const addressString = formatPaymentRequestReviewAddressString();
      expect(addressString).toEqual('');
    });
  });

  describe('getShipmentModificationType', () => {
    const canceledShipment = {
      status: shipmentStatuses.CANCELED,
      diversion: false,
    };

    const divertedShipment = {
      status: shipmentStatuses.APPROVED,
      diversion: true,
    };

    it('returns canceled when the shipment status is canceled', () => {
      const shipmentType = getShipmentModificationType(canceledShipment);
      expect(shipmentType).toEqual(shipmentModificationTypes.CANCELED);
    });

    it('returns diversion when the shipment has been marked as a diversion', () => {
      const shipmentType = getShipmentModificationType(divertedShipment);
      expect(shipmentType).toEqual(shipmentModificationTypes.DIVERSION);
    });
  });
});
