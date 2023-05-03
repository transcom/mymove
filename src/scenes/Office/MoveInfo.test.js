import React from 'react';
import ReactDOM from 'react-dom';

import MoveInfo from './MoveInfo';
import { useLocation } from 'react-router-dom';
import { mount } from 'enzyme/build';
import { ReferrerQueueLink } from './MoveInfo';
import { MockProviders } from 'testUtils';

const dummyFunc = () => {};
const loadDependenciesHasError = null;
const loadDependenciesHasSuccess = false;

const params = { moveID: '123456' };

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: jest.fn().mockReturnValue({ pathname: '/' }),
}));

let wrapper;

beforeEach(() => {
  jest.resetAllMocks();
});

describe('Loads MoveInfo', () => {
  // TODO: fix this tests- currently only rendering the Loader
  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(
      <MockProviders params={params} path="/moveIt/moveIt">
        <MoveInfo
          loadDependenciesHasError={loadDependenciesHasError}
          loadDependenciesHasSuccess={loadDependenciesHasSuccess}
          loadMoveDependencies={dummyFunc}
        />
      </MockProviders>,
      div,
    );
  });
  it.skip('shows the Basic and PPM tabs', () => {
    // TODO: apply loadDependenciesHasError and loadDependenciesHasSuccess values through store (currently renders Loader only)
    wrapper = mount(
      <MockProviders params={params} path="/moveIt/moveIt">
        <MoveInfo loadDependenciesHasError={false} loadDependenciesHasSuccess={true} loadMoveDependencies={dummyFunc} />
      </MockProviders>,
    );
    expect(wrapper.find('[data-testid="basics-tab"]').length).toBe(1);
    expect(wrapper.find('[data-testid="ppm-tab"]').length).toBe(1);
  });
});

describe('ShipmentInfo tests', () => {
  describe('Shows correct queue to return to', () => {
    it('when a referrer is set in history', () => {
      useLocation.mockReturnValue({
        pathname: '/',
        state: { referrerPathname: '/queues/ppm_payment_requested' },
      });

      wrapper = mount(
        <MockProviders mockLocation={{ state: { referrerPathname: '/queues/ppm_payment_requested' } }}>
          <ReferrerQueueLink />
        </MockProviders>,
      );
      expect(wrapper.text()).toEqual('Payment requested');
    });
    it('when no referrer is set', () => {
      useLocation.mockReturnValue({
        pathname: '/',
      });

      wrapper = mount(
        <MockProviders>
          <ReferrerQueueLink />
        </MockProviders>,
      );
      expect(wrapper.text()).toEqual('New moves');
    });
  });
});
