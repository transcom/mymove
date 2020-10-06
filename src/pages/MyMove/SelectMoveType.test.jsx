import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { SelectMoveType } from 'pages/MyMove/SelectMoveType';

describe('SelectMoveType', () => {
  let defaultProps;

  beforeEach(() => {
    defaultProps = {
      pageList: ['page1', 'anotherPage/:foo/:bar'],
      pageKey: 'page1',
      match: { isExact: false, path: '', url: '' },
      updateMove: jest.fn(),
      push: jest.fn(),
      loadMTOShipments: jest.fn(),
      move: { id: 'mockId', status: 'DRAFT' },
      selectedMoveType: SHIPMENT_OPTIONS.PPM,
      isPpmSelectable: true,
      isHhgSelectable: true,
      shipmentNumber: 4,
    };
  });

  it('should render radio buttons with PPM selected', () => {
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} />);
    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).text()).toContain('I’ll move things myself');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');
  });

  it('should render radio buttons with HHG selected', () => {
    defaultProps.selectedMoveType = SHIPMENT_OPTIONS.HHG;
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} />);
    expect(wrapper.find(Radio).length).toBe(2);

    expect(wrapper.find(Radio).at(1).text()).toContain('The government packs for me and moves me');
    // HHG button should be checked on page load
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('checked');
  });

  it('should disable PPM form option if PPM is already submitted', () => {
    defaultProps.isPpmSelectable = false;

    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} />);

    // PPM button should be disabled on page load and should contained updated text
    const actualComponentText = wrapper.text();
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('disabled');
    expect(actualComponentText).toContain('contact the PPPO at your origin duty station');
    expect(actualComponentText).not.toContain('You arrange to move some or all of your belongings');
  });

  it('should disable HHG form option if move is already submitted', () => {
    defaultProps.isHhgSelectable = false;

    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} />);

    // HHG button should be disabled on page load and should contained updated text
    const actualComponentText = wrapper.text();
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('disabled');
    expect(actualComponentText).toContain('Talk with your movers directly');
    expect(actualComponentText).not.toContain('Professional movers take care of the whole shipment');
  });
});
