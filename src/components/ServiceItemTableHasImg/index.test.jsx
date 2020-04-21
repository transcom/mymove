import React from 'react';
import { shallow } from 'enzyme';
import ServiceItemTableHasImg from './index';

describe('ServiceItemTableHasImg', () => {
  it('should render a thumbnail image when an image url is passed in', () => {
    const serviceItems = [
      {
        id: 'abc123',
        dateRequested: '20 Nov 2020',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          text: 'grandfather clock 7ft x 2ft x 3.5ft',
          imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
        },
      },
    ];

    const wrapper = shallow(<ServiceItemTableHasImg serviceItems={serviceItems} />);

    expect(wrapper.find('.si-thumbnail').exists()).toBe(true);
  });

  it('should only render detail text when there is no image url passed in', () => {
    const serviceItems = [
      {
        id: 'abc123',
        dateRequested: '20 Nov 2020',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          text: 'grandfather clock 7ft x 2ft x 3.5ft',
        },
      },
    ];

    const wrapper = shallow(<ServiceItemTableHasImg serviceItems={serviceItems} />);

    expect(wrapper.find('table').exists()).toBe(true);
    expect(wrapper.find('.si-thumbnail').exists()).toBe(false);
    expect(wrapper.find('.si-details').text()).toBe(serviceItems[0].details.text);
  });
});
