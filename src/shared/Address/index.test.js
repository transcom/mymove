import React from 'react';
import { mount } from 'enzyme';
import { AddressElementDisplay } from '.';

describe('Address component test', () => {
  it('should render with one address', () => {
    const address = {
      street_address_1: '560 State Street',
      city: 'New York City',
      state: 'NY',
      postal_code: '11217',
    };
    const wrapper = mount(
      <AddressElementDisplay address={address} title="Primary" />,
    );
    const Label = wrapper.find('.field-title').props().children;
    const [
      addressSpan,
      lineBreak,
      street_address_2,
      street_address_3,
      city,
      comma,
      state,
      space,
      postal_code,
    ] = wrapper.find('.field-value').props().children;
    expect(Label).toBe('Primary');
    expect(addressSpan).toBe(address.street_address_1);
    expect(lineBreak.type).toBe('br');
    expect(street_address_2).toBeFalsy();
    expect(street_address_3).toBeFalsy();
    expect(city).toBe(address.city);
    expect(comma).toBe(', ');
    expect(state).toBe(address.state);
    expect(space).toBe(' ');
    expect(postal_code).toBe(address.postal_code);
  });
  it('should render with additional addresses', () => {
    const address = {
      street_address_1: '560 State Street',
      street_address_2: 'PO Box 123',
      street_address_3: 'c/o Harry Potter',
      city: 'New York City',
      state: 'NY',
      postal_code: '11217',
    };
    const wrapper = mount(
      <AddressElementDisplay address={address} title="Primary" />,
    );

    const [, , address_2, address_3] = wrapper
      .find('.field-value')
      .props().children;
    expect(address_2.props.children[0]).toBe(address.street_address_2);
    expect(address_3.props.children[0]).toBe(address.street_address_3);
  });
});
