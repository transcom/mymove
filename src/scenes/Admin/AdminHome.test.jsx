import React from 'react';
import { mount } from 'enzyme';
import restProvider from 'ra-data-simple-rest';
import AdminHome from 'scenes/Admin/AdminHome';

const dataProvider = restProvider('http://admin/v1/...');

describe('AdminHome tests', () => {
  describe('AdminHome component', () => {
    let wrapper;
    wrapper = mount(<AdminHome dataProvider={dataProvider} />);

    it('renders without crashing', () => {
      expect(wrapper.find('.admin-system-wrapper').length).toEqual(1);
    });
  });
});
