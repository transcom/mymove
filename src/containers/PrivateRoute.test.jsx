import React from 'react';
import { mount } from 'enzyme';

import PrivateRoute, { userIsAuthorized } from './PrivateRoute';

import { MockProviders } from 'testUtils';
import { roleTypes } from 'constants/userRoles';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  Navigate: (props) => {
    mockNavigate(props?.to);
    return null;
  },
}));

afterEach(() => {
  jest.resetAllMocks();
});

const tioUserInitialState = {
  auth: {
    isLoading: false,
    isLoggedIn: true,
    activeRole: roleTypes.TIO,
  },
  entities: {
    user: {
      userId123: {
        id: 'userId123',
      },
    },
  },
};

describe('userIsAuthorized function', () => {
  it('returns true if no roles are required', () => {
    expect(userIsAuthorized()).toEqual(true);
  });

  it('returns false if the user has no roles', () => {
    expect(userIsAuthorized(undefined, [roleTypes.SERVICES_COUNSELOR])).toEqual(false);
  });

  it('returns true if the user has at least one required role', () => {
    expect(userIsAuthorized(roleTypes.TIO, [roleTypes.TIO, roleTypes.SERVICES_COUNSELOR])).toEqual(true);
  });

  it('returns false if the user does not have a required role', () => {
    expect(userIsAuthorized(roleTypes.TIO, [roleTypes.TOO])).toEqual(false);
  });
});

const MyPrivateComponent = () => <div>My page</div>;
describe('PrivateRoute', () => {
  it('renders the component if user has the requred role', () => {
    const wrapper = mount(
      <MockProviders initialState={tioUserInitialState}>
        <PrivateRoute requiredRoles={[roleTypes.TIO]}>
          <MyPrivateComponent />
        </PrivateRoute>
      </MockProviders>,
    );
    expect(wrapper.find(MyPrivateComponent)).toHaveLength(1);
  });

  it('renders the component if the user has one of multiple required roles', () => {
    const wrapper = mount(
      <MockProviders initialState={tioUserInitialState}>
        <PrivateRoute requiredRoles={[roleTypes.TIO, roleTypes.TOO]}>
          <MyPrivateComponent />
        </PrivateRoute>
      </MockProviders>,
    );

    expect(wrapper.find(MyPrivateComponent)).toHaveLength(1);
  });

  it('renders the component if no roles are required', () => {
    const wrapper = mount(
      <MockProviders>
        <PrivateRoute>
          <MyPrivateComponent />
        </PrivateRoute>
      </MockProviders>,
    );

    expect(wrapper.find(MyPrivateComponent)).toHaveLength(1);
  });

  it('does not render the compoent if the user does not have a required role', () => {
    const wrapper = mount(
      <MockProviders initialState={tioUserInitialState}>
        <PrivateRoute requiredRoles={[roleTypes.TOO]}>
          <MyPrivateComponent />
        </PrivateRoute>
      </MockProviders>,
    );

    expect(wrapper.find(MyPrivateComponent)).toHaveLength(0);
    expect(mockNavigate).toHaveBeenCalledWith('/invalid-permissions');
  });
});
