import React from 'react';
import { mount } from 'enzyme';
import AdminWrapper from './index';
import restProvider from 'ra-data-simple-rest';

const dataProvider = restProvider('http://admin/v1/...');

describe('AdminIndex tests', () => {
  describe('AdminIndex home page', () => {
    let wrapper;
    wrapper = mount(<AdminWrapper dataProvider={dataProvider} />);

    it('renders without crashing', () => {
      expect(wrapper.find('.admin-system-wrapper').length).toEqual(1);
    });
  });
});
