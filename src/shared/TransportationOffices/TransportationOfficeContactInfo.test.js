import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import { TransportationOfficeContactInfo } from './TransportationOfficeContactInfo';

describe('TransportationOfficeContactInfo tests', () => {
  let wrapper;
  it('renders without crashing', () => {
    const div = document.createElement('div');
    wrapper = shallow(<TransportationOfficeContactInfo />, div);
    expect(wrapper.find('div').length).toEqual(3);
  });
});
