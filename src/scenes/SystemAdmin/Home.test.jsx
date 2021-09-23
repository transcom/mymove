import React from 'react';
import { shallow } from 'enzyme';
import restProvider from 'ra-data-simple-rest';
import Home from './Home';

const dataProvider = restProvider('http://admin/v1/...');

describe('AdminHome tests', () => {
  /*
   * Currently skipping this test because the component will always emit a warning, due to
   * the interaction between the `basename` property we pass into `createBrowserHistory` and
   * the way that Jest tests run in the react-scripts package.
   *
   * Essentially, using a Create-React-App application won't let us change the basename value
   * on a per-test basis, which causes the warning to be emitted, usually intended for the
   * browser console but shown in test output here. We should migrate this component to the
   * new front end structure and, when we do so, rewrite how the browser history part works.
   */

  describe.skip('AdminHome component', () => {
    let wrapper;
    wrapper = shallow(<Home dataProvider={dataProvider} />);

    it('renders without crashing', () => {
      expect(wrapper.find('.admin-system-wrapper').length).toEqual(1);
    });
  });
});
