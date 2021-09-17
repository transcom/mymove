import React from 'react';
import { shallow } from 'enzyme';
import restProvider from 'ra-data-simple-rest';
import Home from './Home';

const dataProvider = restProvider('http://admin/v1/...');

describe('AdminHome tests', () => {
  describe.skip('AdminHome component', () => {
    let wrapper;
    wrapper = shallow(<Home dataProvider={dataProvider} />);

    it('renders without crashing', () => {
      expect(wrapper.find('.admin-system-wrapper').length).toEqual(1);
    });
  });
});
