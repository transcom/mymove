/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow, mount } from 'enzyme';

import ConnectedOffice, { OfficeWrapper } from './index';

import { MockProviders } from 'testUtils';
import { roleTypes } from 'constants/userRoles';

describe('ConnectedOffice', () => {
  const mockOfficeProps = {
    getCurrentUserInfo: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadPublicSchema: jest.fn(),
  };

  describe('component', () => {
    let wrapper;

    beforeEach(() => {
      wrapper = shallow(<OfficeWrapper {...mockOfficeProps} />);
    });

    it('renders without crashing or erroring', () => {
      const officeWrapper = wrapper.find('div');
      expect(officeWrapper).toBeDefined();
      expect(wrapper.find('SomethingWentWrong')).toHaveLength(0);
    });

    it('renders the basic header by default', () => {
      expect(wrapper.find('QueueHeader')).toHaveLength(1);
    });

    it('fetches initial data', () => {
      expect(mockOfficeProps.getCurrentUserInfo).toHaveBeenCalled();
      expect(mockOfficeProps.loadInternalSchema).toHaveBeenCalled();
      expect(mockOfficeProps.loadPublicSchema).toHaveBeenCalled();
    });

    describe('if an error occurs', () => {
      it('renders the fail whale', () => {
        wrapper.setState({ hasError: true });
        expect(wrapper.find('SomethingWentWrong')).toHaveLength(1);
      });
    });

    describe('if the user is logged in with multiple roles', () => {
      const multiRoleState = {
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
      };

      describe('on a page that isn’t the Select Application page', () => {
        it('renders the Select Application link', () => {
          const app = mount(
            <MockProviders initialState={multiRoleState} initialEntries={['/']}>
              <ConnectedOffice {...mockOfficeProps} location={{ pathname: '/' }} />
            </MockProviders>,
          );

          expect(app.containsMatchingElement(<a href="/select-application">Change user role</a>)).toEqual(true);
        });
      });

      describe('on the Select Application page', () => {
        it('does not render the Select Application link', () => {
          const app = mount(
            <MockProviders initialState={multiRoleState} initialEntries={['/select-application']}>
              <ConnectedOffice {...mockOfficeProps} location={{ pathname: '/select-application' }} />
            </MockProviders>,
          );

          expect(app.containsMatchingElement(<a href="/select-application">Change user role</a>)).toEqual(false);
        });
      });
    });
  });

  describe('routing', () => {
    // TODO - expects should look for actual component content instead of the route path
    // Might have to add testing-library for this because something about enzyme + Suspense + routes are not rendering content

    const loggedInState = {
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
    };

    const loggedOutState = {
      user: {
        isLoading: false,
        userInfo: {
          isLoggedIn: false,
        },
      },
    };

    it('handles the SignIn URL', () => {
      const app = mount(
        <MockProviders initialState={loggedOutState} initialEntries={['/sign-in']}>
          <OfficeWrapper {...mockOfficeProps} />
        </MockProviders>,
      );

      const renderedRoute = app.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/sign-in');
    });

    it('handles the root URL', () => {
      const app = mount(
        <MockProviders initialState={loggedInState} initialEntries={['/']}>
          <OfficeWrapper {...mockOfficeProps} />
        </MockProviders>,
      );

      const renderedRoute = app.find('PrivateRoute');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/');
    });

    it('handles the Select Application URL', () => {
      const app = mount(
        <MockProviders initialState={loggedInState} initialEntries={['/select-application']}>
          <OfficeWrapper {...mockOfficeProps} />
        </MockProviders>,
      );

      const renderedRoute = app.find('PrivateRoute');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/select-application');
    });

    describe('PPM routes', () => {
      const loggedInPPMState = {
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
      };

      it('handles a MoveInfo URL', () => {
        const app = mount(
          <MockProviders initialState={loggedInPPMState} initialEntries={['/queues/new/moves/123']}>
            <OfficeWrapper {...mockOfficeProps} location={{ pathname: '/queues/new/moves/123' }} />
          </MockProviders>,
        );

        const renderedRoute = app.find('PrivateRoute');
        expect(renderedRoute).toHaveLength(1);
        expect(renderedRoute.prop('path')).toEqual('/queues/:queueType/moves/:moveId');
      });

      it('handles a Queues URL', () => {
        const app = mount(
          <MockProviders initialState={loggedInPPMState} initialEntries={['/queues/new']}>
            <OfficeWrapper {...mockOfficeProps} location={{ pathname: '/queues/new' }} />
          </MockProviders>,
        );

        const renderedRoute = app.find('PrivateRoute');
        expect(renderedRoute).toHaveLength(1);
        expect(renderedRoute.prop('path')).toEqual('/queues/:queueType');
      });

      it('handles a OrdersInfo URL', () => {
        const app = mount(
          <MockProviders initialState={loggedInPPMState} initialEntries={['/moves/123/orders']}>
            <OfficeWrapper {...mockOfficeProps} userIsLoggedIn location={{ pathname: '/moves/123/orders' }} />
          </MockProviders>,
        );

        const renderedRoute = app.find('PrivateRoute');
        expect(renderedRoute).toHaveLength(1);
        expect(renderedRoute.prop('path')).toEqual('/moves/:moveId/orders');

        // OrdersInfo does NOT render the header
        expect(app.find('QueueHeader')).toHaveLength(0);
      });

      it('handles a DocumentViewer URL', () => {
        const app = mount(
          <MockProviders initialState={loggedInPPMState} initialEntries={['/moves/123/documents/abc']}>
            <OfficeWrapper {...mockOfficeProps} userIsLoggedIn location={{ pathname: '/moves/123/documents/abc' }} />
          </MockProviders>,
        );

        const renderedRoute = app.find('PrivateRoute');
        expect(renderedRoute).toHaveLength(1);
        expect(renderedRoute.prop('path')).toEqual('/moves/:moveId/documents/:moveDocumentId?');

        // DocumentViewer does NOT render the header
        expect(app.find('QueueHeader')).toHaveLength(0);
      });
    });

    describe('TOO routes', () => {
      const loggedInTOOState = {
        user: {
          isLoading: false,
          userInfo: {
            isLoggedIn: true,
            roles: [
              {
                roleType: roleTypes.TOO,
              },
            ],
          },
        },
      };

      describe('without the feature flags set', () => {
        it('does not handle the moves queue URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTOOState} initialEntries={['/moves/queue']}>
              <OfficeWrapper {...mockOfficeProps} location={{ pathname: '/moves/queue' }} />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(0);
        });
      });

      describe('with the feature flags set', () => {
        it('handles the moves queue URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTOOState} initialEntries={['/moves/queue']}>
              <OfficeWrapper
                context={{ flags: { too: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/moves/queue' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/moves/queue');
        });

        it('handles the TXOMoveInfo URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTOOState} initialEntries={['/moves/123']}>
              <OfficeWrapper
                context={{ flags: { too: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/moves/123' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/moves/:moveOrderId');
        });

        it('handles the CustomerDetails URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTOOState} initialEntries={['/too/123/customer/abc']}>
              <OfficeWrapper
                context={{ flags: { too: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/too/123/customer/abc' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/too/:moveOrderId/customer/:customerId');
        });

        it('handles the Verification URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTOOState} initialEntries={['/verification-in-progress']}>
              <OfficeWrapper
                context={{ flags: { too: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/verification-in-progress' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('Route');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/verification-in-progress');
        });
      });
    });

    describe('TIO routes', () => {
      const loggedInTIOState = {
        user: {
          isLoading: false,
          userInfo: {
            isLoggedIn: true,
            roles: [
              {
                roleType: roleTypes.TIO,
              },
            ],
          },
        },
      };

      describe('without the feature flags set', () => {
        it('does not handle the invoicing queue URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTIOState} initialEntries={['/invoicing/queue']}>
              <OfficeWrapper {...mockOfficeProps} location={{ pathname: '/invoicing/queue' }} />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(0);
        });
      });

      describe('with the feature flags set', () => {
        it('handles the invoicing queue URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTIOState} initialEntries={['/invoicing/queue']}>
              <OfficeWrapper
                context={{ flags: { tio: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/invoicing/queue' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/invoicing/queue');
        });

        it('handles the TXOMoveInfo URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTIOState} initialEntries={['/moves/123']}>
              <OfficeWrapper
                context={{ flags: { too: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/moves/123' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/moves/:moveOrderId');
        });

        it('handles the PaymentRequestIndex URL', () => {
          const app = mount(
            <MockProviders initialState={loggedInTIOState} initialEntries={['/payment_requests']}>
              <OfficeWrapper
                context={{ flags: { tio: true } }}
                {...mockOfficeProps}
                location={{ pathname: '/payment_requests' }}
              />
            </MockProviders>,
          );

          const renderedRoute = app.find('PrivateRoute');
          expect(renderedRoute).toHaveLength(1);
          expect(renderedRoute.prop('path')).toEqual('/payment_requests');
        });
      });
    });
  });
});
