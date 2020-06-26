import React from 'react';
import { mount } from 'enzyme';

import { MockProviders } from 'testUtils';
import PrivateRoute, { userIsAuthorized } from './PrivateRoute';
import { roleTypes } from 'constants/userRoles';

describe('userIsAuthorized function', () => {
  it('returns true if no roles are required', () => {
    expect(userIsAuthorized()).toEqual(true);
  });

  it('returns false if the user has no roles', () => {
    expect(userIsAuthorized(undefined, [roleTypes.PPM])).toEqual(false);
  });

  it('returns true if the user has at least one required role', () => {
    expect(userIsAuthorized([roleTypes.TIO], [roleTypes.TIO, roleTypes.PPM])).toEqual(true);
  });

  it('returns true if the user has at multiple required roles', () => {
    expect(userIsAuthorized([roleTypes.TIO, roleTypes.TOO], [roleTypes.TOO, roleTypes.TIO, roleTypes.PPM])).toEqual(
      true,
    );
  });

  it('returns false if the user does not have a required role', () => {
    expect(userIsAuthorized([roleTypes.TIO], [roleTypes.TOO])).toEqual(false);
  });
});

describe('PrivateRouteContainer', () => {
  describe('if the user is still loading', () => {
    it('renders the loading placeholder', () => {
      const wrapper = mount(
        <MockProviders
          initialState={{
            user: {
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
            user: {
              isLoading: false,
              userInfo: {
                isLoggedIn: false,
              },
            },
          }}
          initialEntries={['/']}
        >
          <PrivateRoute render={() => <div>My page</div>} requiredRoles={[roleTypes.TOO]} />
        </MockProviders>,
      );

      it('does not render the loading placeholder', () => {
        expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
      });
      it('does not render the requested component', () => {
        expect(wrapper.contains(<div>My page</div>)).toEqual(false);
      });
      it('displays the Sign In link', () => {
        expect(wrapper.containsMatchingElement(<a href="/auth/login-gov">Sign in</a>)).toEqual(true);
      });
    });

    describe('and is logged in', () => {
      describe('and is not authorized to view the given route', () => {
        const wrapper = mount(
          <MockProviders
            initialState={{
              user: {
                isLoading: false,
                userInfo: {
                  isLoggedIn: true,
                  roles: [
                    {
                      roleType: roleTypes.PPM,
                    },
                  ],
                },
              },
            }}
            initialEntries={['/']}
          >
            <PrivateRoute render={() => <div>My page</div>} requiredRoles={[roleTypes.TOO]} />
          </MockProviders>,
        );

        it('does not render the loading placeholder', () => {
          expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
        });
        it('does not render the requested component', () => {
          expect(wrapper.contains(<div>My page</div>)).toEqual(false);
        });
        it('redirects to the root URL', () => {
          const redirect = wrapper.find('Redirect');
          expect(redirect).toHaveLength(1);
          expect(redirect.prop('to')).toEqual('/');
        });
      });

      describe('and is authorized to view the given route', () => {
        const wrapper = mount(
          <MockProviders
            initialState={{
              user: {
                isLoading: false,
                userInfo: {
                  isLoggedIn: true,
                  roles: [
                    {
                      roleType: roleTypes.PPM,
                    },
                  ],
                },
              },
            }}
            initialEntries={['/']}
          >
            <PrivateRoute render={() => <div>My page</div>} requiredRoles={[roleTypes.PPM]} />
          </MockProviders>,
        );
        it('does not render the loading placeholder', () => {
          expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
        });
        it('renders the requested component', () => {
          expect(wrapper.contains(<div>My page</div>)).toEqual(true);
        });
      });

      describe('and is authorized with multiple roles', () => {
        describe('on a page that isn’t the Select Application page', () => {
          const wrapper = mount(
            <MockProviders
              initialState={{
                user: {
                  isLoading: false,
                  userInfo: {
                    isLoggedIn: true,
                    roles: [
                      {
                        roleType: roleTypes.TOO,
                      },
                      {
                        roleType: roleTypes.TIO,
                      },
                    ],
                  },
                },
              }}
              initialEntries={['/']}
            >
              <PrivateRoute
                render={() => <div>My page</div>}
                requiredRoles={[roleTypes.TOO]}
                path="/my-page"
                location={{ pathname: '/my-page' }}
              />
            </MockProviders>,
          );

          it('does not render the loading placeholder', () => {
            expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
          });
          it('renders the requested component', () => {
            expect(wrapper.contains(<div>My page</div>)).toEqual(true);
          });
          it('renders the Select Application link', () => {
            expect(wrapper.containsMatchingElement(<a href="/select-application">Change user role</a>)).toEqual(true);
          });
        });

        describe('on the Select Application page', () => {
          const wrapper = mount(
            <MockProviders
              initialState={{
                user: {
                  isLoading: false,
                  userInfo: {
                    isLoggedIn: true,
                    roles: [
                      {
                        roleType: roleTypes.TOO,
                      },
                      {
                        roleType: roleTypes.TIO,
                      },
                    ],
                  },
                },
              }}
              initialEntries={['/select-application']}
            >
              <PrivateRoute
                render={() => <div>My page</div>}
                requiredRoles={[roleTypes.TOO]}
                path="/select-application"
                location={{ pathname: '/select-application' }}
              />
            </MockProviders>,
          );
          it('does not render the loading placeholder', () => {
            expect(wrapper.find('[data-name="loading-placeholder"]')).toHaveLength(0);
          });
          it('renders the requested component', () => {
            expect(wrapper.contains(<div>My page</div>)).toEqual(true);
          });
          it('does not render the Select Application link', () => {
            expect(wrapper.containsMatchingElement(<a href="/select-application">Change user role</a>)).toEqual(false);
          });
        });
      });
    });
  });
});
