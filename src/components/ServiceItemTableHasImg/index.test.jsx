import React from 'react';
import { shallow } from 'enzyme';

import { SERVICE_ITEM_STATUS } from '../../shared/constants';

import ServiceItemTableHasImg from './index';

const handleUpdateServiceItemStatus = jest.fn();

describe('ServiceItemTableHasImg', () => {
  it('should render a thumbnail image when an image url is passed in', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          text: 'grandfather clock 7ft x 2ft x 3.5ft',
          imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
        },
      },
    ];

    const wrapper = shallow(
      <ServiceItemTableHasImg
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItemStatus}
      />,
    );

    expect(wrapper.find('.siThumbnail').exists()).toBe(true);
  });

  it('should only render detail text when there is no image url passed in', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          text: 'grandfather clock 7ft x 2ft x 3.5ft',
        },
      },
    ];

    const wrapper = shallow(
      <ServiceItemTableHasImg
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItemStatus}
      />,
    );
    expect(wrapper.find('table').exists()).toBe(true);
    expect(wrapper.find('.siThumbnail').exists()).toBe(false);
    expect(wrapper.find('.detail').text()).toBe(serviceItems[0].details.text);
  });

  it('should render properly when the detail text is an object', () => {
    const serviceItems = [
      {
        id: 'abc123',
        submittedAt: '2020-11-20',
        serviceItem: 'Domestic Crating',
        code: 'DCRT',
        details: {
          text: {
            ZIP: '11111',
            Reason: 'This is the reason',
          },
        },
      },
    ];

    const wrapper = shallow(
      <ServiceItemTableHasImg
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItemStatus}
      />,
    );
    expect(wrapper.find('.detail').contains('This is the reason')).toBe(true);
    expect(wrapper.find('.detail').contains('11111')).toBe(true);
  });
});
