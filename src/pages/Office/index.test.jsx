/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow, mount } from 'enzyme';

import { OfficeWrapper } from './index';

import { MockProviders } from 'testUtils';
import { roleTypes } from 'constants/userRoles';

describe('OfficeWrapper tests', () => {
  let wrapper;

  const mockOfficeProps = {
    getCurrentUserInfo: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadPublicSchema: jest.fn(),
  };

  beforeEach(() => {
    wrapper = shallow(<OfficeWrapper {...mockOfficeProps} />);
  });

  it('renders without crashing or erroring', () => {
    const officeWrapper = wrapper.find('div');
    expect(officeWrapper).toBeDefined();
    expect(wrapper.find('SomethingWentWrong')).toHaveLength(0);
  });

  describe('if an error occurs', () => {
    it('renders the fail whale', () => {
      wrapper.setState({ hasError: true });
      expect(wrapper.find('SomethingWentWrong')).toHaveLength(1);
    });
  });

  describe('routing', () => {
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

    it('handles the root URL', () => {
      const app = mount(
        <MockProviders initialState={loggedInState} initialEntries={['/']}>
          <OfficeWrapper {...mockOfficeProps} />
        </MockProviders>,
      );

      const renderedRoute = app.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/');
    });

    it('handles a MoveInfo URL', () => {
      const app = mount(
        <MockProviders initialState={loggedInState} initialEntries={['/queues/new/moves/123']}>
          <OfficeWrapper {...mockOfficeProps} />
        </MockProviders>,
      );

      const renderedRoute = app.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/queues/:queueType/moves/:moveId');
    });
  });
});
