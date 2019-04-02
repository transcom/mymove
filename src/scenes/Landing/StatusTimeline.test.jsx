import React from 'react';
import { shallow, mount } from 'enzyme';
import { StatusBlock, PPMStatusTimeline, ShipmentStatusTimeline, ProfileStatusTimeline } from './StatusTimeline';

describe('StatusTimeline', () => {
  describe('PPMStatusTimeline', () => {
    test('renders timeline', () => {
      const ppm = {};
      const wrapper = mount(<PPMStatusTimeline ppm={ppm} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(4);
    });

    test('renders timeline for submitted ppm', () => {
      const ppm = { status: 'SUBMITTED' };
      const wrapper = mount(<PPMStatusTimeline ppm={ppm} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(1);
      expect(completed.prop('code')).toEqual('SUBMITTED');

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('SUBMITTED');
    });

    test('renders timeline for an in-progress ppm', () => {
      const ppm = { status: 'APPROVED', original_move_date: '2019-03-20' };
      const wrapper = mount(<PPMStatusTimeline ppm={ppm} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(3);
      expect(completed.map(b => b.prop('code'))).toEqual(['SUBMITTED', 'PPM_APPROVED', 'IN_PROGRESS']);

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('IN_PROGRESS');
    });
  });

  describe('ShipmentStatusTimeline', () => {
    test('renders timeline', () => {
      const shipment = {};
      const wrapper = mount(<ShipmentStatusTimeline shipment={shipment} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(5);
    });

    test('renders timeline for scheduled hhg', () => {
      const shipment = { status: 'SCHEDULED' };
      const wrapper = mount(<ShipmentStatusTimeline shipment={shipment} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(1);
      expect(completed.prop('code')).toEqual('SCHEDULED');

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('SCHEDULED');
    });

    test('renders timeline for packed hhg', () => {
      const shipment = { status: 'PACKED', actual_pack_date: '2019-03-20', today: '2019-03-20' };
      const wrapper = mount(<ShipmentStatusTimeline shipment={shipment} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(2);
      expect(completed.map(b => b.prop('code'))).toEqual(['SCHEDULED', 'PACKED']);

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('PACKED');
    });
  });

  describe('ProfileStatusTimeline', () => {
    test('renders timeline', () => {
      const profile = {};
      const wrapper = mount(<ProfileStatusTimeline profile={profile} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(4);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(2);
      expect(completed.map(b => b.prop('code'))).toEqual(['PROFILE', 'ORDERS']);

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('ORDERS');
    });
  });
});

describe('StatusBlock', () => {
  test('complete but not current status block', () => {
    const wrapper = shallow(
      <StatusBlock name="Approved" completed={true} current={false} code="PPM_APPROVED" key="PPM_APPROVED" />,
    );

    expect(wrapper.hasClass('ppm_approved')).toEqual(true);
    expect(wrapper.hasClass('status_completed')).toEqual(true);
    expect(wrapper.hasClass('status_current')).toEqual(false);
  });

  test('complete and current status block', () => {
    const wrapper = shallow(
      <StatusBlock name="In Progress" completed={true} current={true} code="IN_PROGRESS" key="IN_PROGRESS" />,
    );

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
