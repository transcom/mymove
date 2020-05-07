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
});

let wrapper;
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

// not sure how much granularity there should be here
// describe('Shows correct status and icon below tab heading', () => {
// // click into a new move, Basics tab displays faClock and 'Submitted'
//   // office user process the Orders in Document viewer
//   // then clicks 'Approve Basics'
//   // and refreshes the page
//     // the Basics tab displays faCheck and 'Approved'

// // click into the PPM tab, tab should display faClock and one of the following:
// // ‘Move pending’ ‘Payment requested’ ‘In review’
//   // office user processes weight tickets
//   // then clicks 'Approve PPM'
//   // and refreshes the page
//     // PPM tab displays faClock and
// }
