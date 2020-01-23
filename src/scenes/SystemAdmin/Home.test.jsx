import React from 'react';
import { mount } from 'enzyme';
import restProvider from 'ra-data-simple-rest';
import Home from './Home';
import { HashRouter as Router } from 'react-router-dom';


const dataProvider = restProvider('http://admin/v1/...');

describe('AdminHome tests', () => {
  describe('AdminHome component', () => {
    let wrapper;
    wrapper = mount(<Router><Home dataProvider={dataProvider} /></Router>);

    it('renders without crashing', () => {
      expect(wrapper.find('.admin-system-wrapper').length).toEqual(1);
    });
  });
});
