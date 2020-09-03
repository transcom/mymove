import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { SelectMoveType } from 'pages/MyMove/SelectMoveType';

describe('SelectMoveType', () => {
  const defaultProps = {
    pageList: ['page1', 'anotherPage/:foo/:bar'],
    pageKey: 'page1',
    match: { isExact: false, path: '', url: '' },
    updateMove: () => {},
    push: () => {},
    // selectedMoveType: SHIPMENT_OPTIONS.PPM,
    selectedMoveType: SHIPMENT_OPTIONS.HHG,
  };
  // PPMs will be supported again in the future
  // it('should render radio buttons with PPM selected', () => {
  //   // eslint-disable-next-line react/jsx-props-no-spreading
  //   const wrapper = mount(<SelectMoveType {...defaultProps} />);
  //   expect(wrapper.find(Radio).length).toBe(2);

  //   // PPM button should be checked on page load
  //   expect(wrapper.find(Radio).at(0).text()).toContain('I’ll move things myself');
  //   expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');
  // });
  it('should render radio buttons with HHG selected', () => {
    defaultProps.selectedMoveType = SHIPMENT_OPTIONS.HHG;
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} />);
    expect(wrapper.find(Radio).length).toBe(2);

    expect(wrapper.find(Radio).at(1).text()).toContain('The government packs for me and moves me');
    // HHG button should be checked on page load
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('checked');
  });
});
