/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ConnectedDodInfo from './DodInfo';

import { MockProviders } from 'testUtils';

const testServiceMember = {};

const testProps = {
  pages: [],
  pageKey: '',
  match: {
    path: '',
  },
};

describe('DoD Info page', () => {
  const wrapper = mount(
    <MockProviders
      initialState={{
        serviceMember: {
          currentServiceMember: testServiceMember,
        },
      }}
    >
      <ConnectedDodInfo {...testProps} />
    </MockProviders>,
  );
  it('renders without errors', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('renders the correct text content', () => {
    expect(wrapper.contains(<h1>Create your profile</h1>)).toBe(true);
    expect(wrapper.contains(<p>Before we can schedule your move, we need to know a little more about you.</p>)).toBe(
      true,
    );
  });

  it('renders the DoD info form fields', () => {
    expect(wrapper.find('SwaggerField[fieldName="affiliation"]').length).toBe(1);
    expect(wrapper.find('SwaggerField[fieldName="edipi"]').length).toBe(1);
    expect(wrapper.find('SwaggerField[fieldName="rank"]').length).toBe(1);
    expect(wrapper.find('input[name="social_security_number"]').length).toBe(1);
  });
});
