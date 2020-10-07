import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS, MOVE_STATUSES } from 'shared/constants';
import { SelectMoveType } from 'pages/MyMove/SelectMoveType';

describe('SelectMoveType', () => {
  const defaultProps = {
    pageList: ['page1', 'anotherPage/:foo/:bar'],
    pageKey: 'page1',
    match: { isExact: false, path: '', url: '' },
    updateMove: jest.fn(),
    push: jest.fn(),
    loadMTOShipments: jest.fn(),
    move: { id: 'mockId', status: 'DRAFT' },
    selectedMoveType: SHIPMENT_OPTIONS.PPM,
    mtoShipments: {},
    isPpmSelectable: true,
    isHhgSelectable: true,
    shipmentNumber: 4,
  };

  const getWrapper = (props = {}) => {
    return mount(<SelectMoveType {...defaultProps} {...props} />); // eslint-disable-line react/jsx-props-no-spreading
  };

  it('should render radio buttons with PPM selected', () => {
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = getWrapper();
    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).text()).toContain('Do it yourself');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');
  });

  it('should render radio buttons with HHG selected', () => {
    const props = { selectedMoveType: SHIPMENT_OPTIONS.HHG };
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = getWrapper(props);
    expect(wrapper.find(Radio).length).toBe(2);

    expect(wrapper.find(Radio).at(1).text()).toContain('Professional movers');
    // HHG button should be checked on page load
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('checked');
  });

  it('should disable PPM form option if PPM is already submitted', () => {
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} isPpmSelectable={false} />);

    // PPM button should be disabled on page load and should contained updated text
    const actualComponentText = wrapper.text();
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('disabled');
    expect(actualComponentText).toContain('contact the PPPO at your origin duty station');
    expect(actualComponentText).not.toContain('You arrange to move some or all of your belongings');
  });

  it('should disable HHG form option if move is already submitted', () => {
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<SelectMoveType {...defaultProps} isHhgSelectable={false} />);

    // HHG button should be disabled on page load and should contained updated text
    const actualComponentText = wrapper.text();
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('disabled');
    expect(actualComponentText).toContain('Talk with your movers directly');
    expect(actualComponentText).not.toContain('Professional movers take care of the whole shipment');
  });
  describe('when no PPMs or shipments have been created', () => {
    it('should render the correct text', () => {
      const wrapper = getWrapper();
      expect(wrapper.find('h1').text()).toContain('How do you want to move your belongings?');
      expect(wrapper.find('[data-testid="selectableCardText"]').at(0).text()).toContain(
        'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a a Personally Procured Move (PPM), sometimes called a DITY.',
      );
    });
  });
  describe('when a PPM has already been created', () => {
    const props = {
      move: { personally_procured_moves: [{ id: 1 }] },
    };
    it('should render the correct text', () => {
      const wrapper = getWrapper(props);
      expect(wrapper.find('h1').text()).toContain('How do you want this group of things moved?');
      expect(wrapper.find(Radio).at(0).text()).toContain('Do it yourself (already chosen)');
      expect(wrapper.find('[data-testid="selectableCardText"]').at(0).text()).toContain(
        'You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty station.',
      );
    });
  });
  describe('when an HHG has already been created', () => {
    const props = {
      mtoShipments: [{ id: 2 }],
    };
    it('should render the correct text', () => {
      const wrapper = getWrapper(props);
      expect(wrapper.find('h1').text()).toContain('How do you want this group of things moved?');
    });
  });
  describe('when a move has already been submitted', () => {
    const props = {
      move: { status: MOVE_STATUSES.SUBMITTED },
    };
    it('should render the correct text', () => {
      const wrapper = getWrapper(props);
      expect(wrapper.find('[data-testid="selectableCardText"]').at(1).text()).toContain(
        'Talk with your movers directly if you want to add or change shipments.',
      );
    });
  });
});
