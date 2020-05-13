import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router-dom';

import MoveInfo from './MoveInfo';
import store from 'shared/store';
import { mount } from 'enzyme/build';
import { shallow } from 'enzyme';
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

describe('ShipmentInfo tests', () => {
  let wrapper;
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

// // click into the PPM tab, tab should display faClock and one of the following:
// // ‘Move pending’ ‘Payment requested’ ‘In review’
//   // office user processes weight tickets
//   // then clicks 'Approve PPM'
//   // and refreshes the page
//     // PPM tab displays faClock and
// possible statuses
// DRAFT don't think this is touched here
// SUBMITTED
// APPROVED
// PAYMENT_REQUESTED // In review
// COMPLETED
// CANCELED
describe('Test PPM tab icon', () => {
  let wrapper;
  const minProps = {
    ppm: { status: '' },
    ppmAdvance: { status: '' },
  };
  describe('Showing PPM status APPROVED and PPM Advance status APPROVED', () => {
    it('Should show red clock icon', () => {
      // move pending
      wrapper = shallow(<MoveInfo {...minProps} ppm={{ status: 'APPROVED' }} ppmAdvance={{ status: 'APPROVED' }} />);
      wrapper.debug();
      expect(wrapper.find({ 'data-cy': 'ppmTabStatus' }).prop('src')).toEqual('faClock');
    });
  });
  // it('Should show red clock icon and `????` when PPM status is APPROVED', () => {
  //   // APPROVED
  // });
  // it('Should show red clock icon and `Payment requested` when PPM status is PAYMENT_REQUESTED', () => {
  //   // PAYMENT_REQUESTED
  // });
  // it('Should show green check icon and `Completed` when PPM status is COMPLETED', () => {
  //   // COMPLETED
  // });
});

// describe Test PPM tab status
// ... and `Move Pending`
