import React from 'react';
import { shallow, mount } from 'enzyme';

import { StatusBlock, ProfileStatusTimeline } from './StatusTimeline';
import { PPMStatusTimeline } from './PPMStatusTimeline';

describe('StatusTimeline', () => {
  describe('PPMStatusTimeline', () => {
    test('renders timeline', () => {
      const ppm = {};
      const wrapper = mount(<PPMStatusTimeline ppm={ppm} getSignedCertification={jest.fn()} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(5);
    });

    test('renders timeline for submitted ppm', () => {
      const ppm = { status: 'SUBMITTED' };
      const wrapper = mount(<PPMStatusTimeline ppm={ppm} getSignedCertification={jest.fn()} />);

      const completed = wrapper.findWhere((b) => b.prop('completed'));
      expect(completed).toHaveLength(1);
      expect(completed.prop('code')).toEqual('SUBMITTED');

      const current = wrapper.findWhere((b) => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('SUBMITTED');
    });

    test('renders timeline for an in-progress ppm', () => {
      const ppm = { status: 'APPROVED', original_move_date: '2019-03-20' };
      const wrapper = mount(<PPMStatusTimeline ppm={ppm} getSignedCertification={jest.fn()} />);

      const completed = wrapper.findWhere((b) => b.prop('completed'));
      expect(completed).toHaveLength(3);
      expect(completed.map((b) => b.prop('code'))).toEqual(['SUBMITTED', 'PPM_APPROVED', 'IN_PROGRESS']);

      const current = wrapper.findWhere((b) => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('IN_PROGRESS');
    });
  });

  describe('ProfileStatusTimeline', () => {
    test('renders timeline', () => {
      const profile = {};
      const wrapper = mount(<ProfileStatusTimeline profile={profile} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(4);

      const completed = wrapper.findWhere((b) => b.prop('completed'));
      expect(completed).toHaveLength(2);
      expect(completed.map((b) => b.prop('code'))).toEqual(['PROFILE', 'ORDERS']);

      const current = wrapper.findWhere((b) => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('ORDERS');
    });
  });
});

describe('StatusBlock', () => {
  test('complete but not current status block', () => {
    const wrapper = shallow(
      <StatusBlock name="Approved" completed current={false} code="PPM_APPROVED" key="PPM_APPROVED" />,
    );

    expect(wrapper.hasClass('ppm_approved')).toEqual(true);
    expect(wrapper.hasClass('status_completed')).toEqual(true);
    expect(wrapper.hasClass('status_current')).toEqual(false);
  });

  test('complete and current status block', () => {
    const wrapper = shallow(<StatusBlock name="In Progress" completed current code="IN_PROGRESS" key="IN_PROGRESS" />);

    expect(wrapper.hasClass('in_progress')).toEqual(true);
    expect(wrapper.hasClass('status_completed')).toEqual(true);
    expect(wrapper.hasClass('status_current')).toEqual(true);
  });

  test('incomplete status block', () => {
    const wrapper = shallow(
      <StatusBlock name="Delivered" completed={false} current={false} code="DELIVERED" key="DELIVERED" />,
    );

    expect(wrapper.hasClass('delivered')).toEqual(true);
    expect(wrapper.hasClass('status_completed')).toEqual(false);
    expect(wrapper.hasClass('status_current')).toEqual(false);
  });
});
