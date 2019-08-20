import React from 'react';
import { mount } from 'enzyme';

import SitAction from './SitAction';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

let wrapper;
const onClickHandler = jest.fn();

describe('SIT action', () => {
  describe('delete action, has icon', () => {
    beforeEach(() => {
      onClickHandler.mockClear();
      wrapper = mount(<SitAction action="Delete" onClick={onClickHandler} icon={faTimes} />);
    });

    it('renders without crashing', () => {
      expect(wrapper.exists('.sit-action')).toBe(true);
    });

    it('has anchor with proper data-cy link', () => {
      expect(wrapper.find('a[data-cy="sit-delete-link"]')).toHaveLength(1);
    });

    it('has svg for icon', () => {
      expect(wrapper.find('a').exists('svg')).toBe(true);
    });

    it('has text for action', () => {
      expect(
        wrapper
          .find('.sit-action')
          .find('a')
          .text(),
      ).toEqual('Delete');
    });

    it('clicking calls handler', () => {
      wrapper.find('a').simulate('click');
      expect(onClickHandler.mock.calls.length).toBe(1);
    });
  });

  describe('place into SIT action, no icon', () => {
    beforeEach(() => {
      onClickHandler.mockClear();
      wrapper = mount(<SitAction action="Place into SIT" onClick={onClickHandler} />);
    });

    it('renders without crashing', () => {
      expect(wrapper.exists('.sit-action')).toBe(true);
    });

    it('has anchor with proper data-cy link', () => {
      expect(wrapper.find('a[data-cy="sit-place-into-sit-link"]')).toHaveLength(1);
    });

    it('has no icon', () => {
      expect(wrapper.find('a').exists('svg')).toBe(false);
    });

    it('has text for action', () => {
      expect(
        wrapper
          .find('.sit-action')
          .find('a')
          .text(),
      ).toEqual('Place into SIT');
    });

    it('clicking calls handler', () => {
      wrapper.find('a').simulate('click');
      expect(onClickHandler.mock.calls.length).toBe(1);
    });
  });
});
