import React from 'react';
import { mount } from 'enzyme';

import PrivateRoute, { userIsAuthorized } from './PrivateRoute';

import { MockProviders } from 'testUtils';
import { roleTypes } from 'constants/userRoles';

describe('userIsAuthorized function', () => {
  it('returns true if no roles are required', () => {
    expect(userIsAuthorized()).toEqual(true);
  });

  it('returns false if the user has no roles', () => {
    expect(userIsAuthorized(undefined, [roleTypes.SERVICES_COUNSELOR])).toEqual(false);
  });

  it('returns true if the user has at least one required role', () => {
    expect(userIsAuthorized([roleTypes.TIO], [roleTypes.TIO, roleTypes.SERVICES_COUNSELOR])).toEqual(true);
  });

  it('returns true if the user has at multiple required roles', () => {
    expect(
      userIsAuthorized([roleTypes.TIO, roleTypes.TOO], [roleTypes.TOO, roleTypes.TIO, roleTypes.SERVICES_COUNSELOR]),
    ).toEqual(true);
  });

  it('returns false if the user does not have a required role', () => {
    expect(userIsAuthorized([roleTypes.TIO], [roleTypes.TOO])).toEqual(false);
  });
});

describe('ConnectedPrivateRoute', () => {
  const MyPrivateComponent = () => <div>My page</div>;

  describe('if the user is still loading', () => {
    it('renders the loading placeholder', () => {
      const wrapper = mount(
        <MockProviders
          initialState={{
            auth: {
              isLoading: true,
            },
          }}
          initialEntries={['/']}
        >
          <PrivateRoute />
        </MockProviders>,
      );
      expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(1);
    });
  });

  describe('if the user has loaded', () => {
    describe('and is not logged in', () => {
      const wrapper = mount(
        <MockProviders
          initialState={{
            auth: {
              isLoading: false,
              isLoggedIn: false,
            },
          }}
          initialEntries={['/']}
        >
          <PrivateRoute path="/" component={MyPrivateComponent} />
        </MockProviders>,
      );

      it('does not render the loading placeholder', () => {
        expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
      });
      it('does not render the requested component', () => {
        expect(wrapper.contains(<div>My page</div>)).toEqual(false);
      });

      it('redirects to the sign in URL', () => {
        const redirect = wrapper.find('Redirect');
        expect(redirect).toHaveLength(1);
        expect(redirect.prop('to')).toEqual('/sign-in');
      });
    });

    describe('and is logged in', () => {
      describe('and is not authorized to view the given route', () => {
        const wrapper = mount(
          <MockProviders
            initialState={{
              auth: {
                isLoading: false,
                isLoggedIn: true,
              },
              entities: {
                user: {
                  userId123: {
                    id: 'userId123',
                    roles: [{ roleType: undefined }],
                  },
                },
              },
            }}
            initialEntries={['/']}
          >
            <PrivateRoute component={MyPrivateComponent} requiredRoles={[roleTypes.TOO]} />
          </MockProviders>,
        );

        it('does not render the loading placeholder', () => {
          expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
        });
        it('does not render the requested component', () => {
          expect(wrapper.contains(<div>My page</div>)).toEqual(false);
        });
        it('redirects to the invalid permissions URL', () => {
          const redirect = wrapper.find('Redirect');
          expect(redirect).toHaveLength(1);
          expect(redirect.prop('to')).toEqual('/invalid-permissions');
        });
      });

      describe('and is authorized to view the given route', () => {
        const wrapper = mount(
          <MockProviders
            initialState={{
              auth: {
                isLoading: false,
                isLoggedIn: true,
              },
              entities: {
                user: {
                  userId123: {
                    id: 'userId123',
                    roles: [{ roleType: roleTypes.TIO }],
                  },
                },
              },
            }}
            initialEntries={['/']}
          >
            <PrivateRoute component={MyPrivateComponent} requiredRoles={[roleTypes.TIO]} />
          </MockProviders>,
        );
        it('does not render the loading placeholder', () => {
          expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
        });
        it('renders the requested component', () => {
          expect(wrapper.contains(<div>My page</div>)).toEqual(true);
        });
      });
    });
  });
});
