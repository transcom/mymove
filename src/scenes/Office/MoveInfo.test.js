import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router-dom';

import MoveInfo from './MoveInfo';
import store from 'shared/store';
import { mount } from 'enzyme/build';
import { ReferrerQueueLink } from './MoveInfo';

const dummyFunc = () => {};
const loadDependenciesHasError = null;
const loadDependenciesHasSuccess = false;
const location = {
  pathname: '',
};
const match = {
  params: { moveID: '123456' },
  url: 'www.nino.com',
  path: '/moveIt/moveIt',
};

const push = jest.fn();
let wrapper;

describe('Loads MoveInfo', () => {
  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(
      <Provider store={store}>
        <Router push={push}>
          <MoveInfo
            loadDependenciesHasError={loadDependenciesHasError}
            loadDependenciesHasSuccess={loadDependenciesHasSuccess}
            location={location}
            match={match}
            loadMoveDependencies={dummyFunc}
          />
        </Router>
      </Provider>,
      div,
    );
  });
  it('shows the Basic and PPM tabs', () => {
    wrapper = mount(
      <Provider store={store}>
        <Router push={push}>
          <MoveInfo
            loadDependenciesHasError={loadDependenciesHasError}
            loadDependenciesHasSuccess={loadDependenciesHasSuccess}
            location={location}
            match={match}
            loadMoveDependencies={dummyFunc}
          />
        </Router>
      </Provider>,
    );
    expect(wrapper.find('[data-testid="basics-tab"]'));
    expect(wrapper.find('[data-testid="ppm-tab"]'));
  });
});

describe('ShipmentInfo tests', () => {
  describe('Shows correct queue to return to', () => {
    it('when a referrer is set in history', () => {
      wrapper = mount(
        <Provider store={store}>
          <Router push={jest.fn()}>
            <ReferrerQueueLink
              history={{ location: { state: { referrerPathname: '/queues/ppm_payment_requested' } } }}
            />
          </Router>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('Payment requested');
    });
    it('when no referrer is set', () => {
      wrapper = mount(
        <Provider store={store}>
          <Router push={jest.fn()}>
            <ReferrerQueueLink history={{ location: {} }} />
          </Router>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('New moves');
    });
  });
});
