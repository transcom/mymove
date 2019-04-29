import React from 'react';
import { shallow } from 'enzyme';
import { PanelField } from 'shared/EditablePanel';
import { AddressElementDisplay } from '.';

describe('Address component test', () => {
  describe('when address has required data', () => {
    const address = {
      street_address_1: '560 State Street',
      city: 'New York City',
      state: 'NY',
      postal_code: '11217',
    };
    const wrapper = shallow(<AddressElementDisplay address={address} title="Primary" />)
      .find(PanelField)
      .dive();
    const fields = wrapper.find('.field-value').props().children;
    it("should render the address' label", () => {
      const Label = wrapper.find('.field-title').text();
      expect(Label).toBe('Primary');
    });
    it('should render the address itself', () => {
      const addressInfo = fields[0].props.children[0];
      expect(addressInfo).toBe(address.street_address_1);
    });
    it('should have a line break in between addresses', () => {
      const lineBreak = fields[0].props.children[1];
      expect(lineBreak.type).toBe('br');
    });
    it('should not render addresses not present', () => {
      const street_address_2 = fields[1];
      expect(street_address_2).toBeFalsy();
    });
    it('should render the city', () => {
      const city = fields[2];
      expect(city).toBe(address.city);
    });
    it('should have a comma in between the address', () => {
      const comma = fields[3];
      expect(comma).toBe(', ');
    });
    it('should render the state', () => {
      const state = fields[4];
      expect(state).toBe(address.state);
    });
    it('should render a space between state and postal code', () => {
      const space = fields[5];
      expect(space).toBe(' ');
    });
    it('should render the postal code', () => {
      const postalCode = fields[6];
      expect(postalCode).toBe(address.postal_code);
    });
  });
  describe('when address component is provided optional additional addresses', () => {
    it('should render street_address_2', () => {
      const address = {
        street_address_1: '560 State Street',
        street_address_2: 'PO Box 123',
        city: 'New York City',
        state: 'NY',
        postal_code: '11217',
      };
      const wrapper = shallow(<AddressElementDisplay address={address} title="primary" />)
        .find(PanelField)
        .dive();

      const address_2 = wrapper
        .find('.field-value')
        .children()
        .at(1)
        .props().children[0];

      expect(address_2).toBe(address.street_address_2);
    });
  });
  describe('when address component is provided only city, state, postal_code', () => {
    const address = {
      city: 'New York City',
      state: 'NY',
      postal_code: '11217',
    };
    const wrapper = shallow(<AddressElementDisplay address={address} title="Primary" />)
      .find(PanelField)
      .dive();
    const fields = wrapper.find('.field-value').props().children;
    it('should not render street address if not provided', () => {
      expect(fields[0]).toBeFalsy();
    });
    it('should not render line breaks', () => {
      expect(fields[1]).toBeFalsy();
    });
    it('should render city', () => {
      expect(fields[2]).toBe(address.city);
    });
    it('should render state', () => {
      expect(fields[4]).toBe(address.state);
    });
    it('should render postal_code', () => {
      expect(fields[6]).toBe(address.postal_code);
    });
  });
});
