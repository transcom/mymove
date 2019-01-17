import React from 'react';
import { shallow } from 'enzyme';
import { getDetailComponent } from './DetailsHelper';

let wrapper;
describe('testing getDetailComponent()', () => {
  describe('returns default details component', () => {
    const DetailComponent = getDetailComponent();
    wrapper = shallow(<DetailComponent />);

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });
  });

  describe('returns 105B/E details component', () => {
    let DetailComponent = getDetailComponent('105B', true);
    wrapper = shallow(<DetailComponent />);
    it('renders 105B details without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });

    DetailComponent = getDetailComponent('105E', true);
    wrapper = shallow(<DetailComponent />);
    it('renders 105E details without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });
  });
});
