import React from 'react';
import { shallow } from 'enzyme';
import { StatusBlock, StatusTimelineContainer } from './StatusTimeline';

describe('StatusTimeline', () => {
  describe('ppm timeline', () => {
    test('renders timeline', () => {
      const ppm = {};
      const wrapper = shallow(<StatusTimelineContainer ppm={ppm} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(4);
    });

    test('renders timeline for submitted ppm', () => {
      const ppm = { status: 'SUBMITTED' };
      const wrapper = shallow(<StatusTimelineContainer ppm={ppm} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(1);
      expect(completed.prop('code')).toEqual('SUBMITTED');

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('SUBMITTED');
    });
  });

  describe('hhg timeline', () => {
    test('renders timeline', () => {
      const shipment = {};
      const wrapper = shallow(<StatusTimelineContainer shipment={shipment} />);

      expect(wrapper.find(StatusBlock)).toHaveLength(5);
    });

    test('renders timeline for scheduled hhg', () => {
      const shipment = { status: 'SCHEDULED' };
      const wrapper = shallow(<StatusTimelineContainer shipment={shipment} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(1);
      expect(completed.prop('code')).toEqual('SCHEDULED');

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('SCHEDULED');
    });

    test('renders timeline for packed hhg', () => {
      const shipment = { status: 'PACKED', actual_pack_date: '2019-03-20', today: '2019-03-20' };
      const wrapper = shallow(<StatusTimelineContainer shipment={shipment} />);

      const completed = wrapper.findWhere(b => b.prop('completed'));
      expect(completed).toHaveLength(2);
      expect(completed.map(b => b.prop('code'))).toEqual(['SCHEDULED', 'PACKED']);

      const current = wrapper.findWhere(b => b.prop('current'));
      expect(current).toHaveLength(1);
      expect(current.prop('code')).toEqual('PACKED');
    });
  });
});
