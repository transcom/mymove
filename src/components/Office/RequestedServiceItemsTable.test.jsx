import React from 'react';
import { shallow } from 'enzyme';
import RequestedServiceItemsTable from './RequestedServiceItemsTable';

describe('RequestedServiceItemsTable', () => {
  it('show the correct number of service items in the table', () => {
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

    let wrapper = shallow(<RequestedServiceItemsTable serviceItems={serviceItems} />);

    expect(wrapper.text().includes('1')).toBe(true);

    serviceItems.push({
      id: 'abc1234',
      dateRequested: '20 Nov 2020',
      serviceItem: 'Domestic origin SIT',
      code: 'DOMSIT',
      details: {
        text: 'Another service item',
        imgURL: '',
      },
    });

    wrapper = shallow(<RequestedServiceItemsTable serviceItems={serviceItems} />);
    expect(wrapper.text().includes('2')).toBe(true);
  });
});
