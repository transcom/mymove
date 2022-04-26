import React from 'react';
import { mount } from 'enzyme';

import ConnectedCustomerPrivateRoute from './CustomerPrivateRoute';

import { MockProviders } from 'testUtils';

describe('ConnectedCustomerPrivateRoute', () => {
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
          <ConnectedCustomerPrivateRoute />
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
          <ConnectedCustomerPrivateRoute path="/" component={MyPrivateComponent} />
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
        expect(redirect.prop('to')).toEqual({
          hash: undefined,
          pathname: '/sign-in',
          search: undefined,
        });
      });
    });

    describe('and is logged in', () => {
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
                  service_member: 'serviceMember456',
                },
              },
              serviceMembers: {
                serviceMember456: {
                  id: 'serviceMember456',
                },
              },
            },
          }}
          initialEntries={['/']}
        >
          <ConnectedCustomerPrivateRoute path="/" component={MyPrivateComponent} />
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
