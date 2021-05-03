/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS, MOVE_STATUSES } from 'shared/constants';
import { SelectShipmentType } from 'pages/MyMove/SelectShipmentType';

describe('SelectShipmentType', () => {
  const defaultProps = {
    updateMove: jest.fn(),
    push: jest.fn(),
    loadMTOShipments: jest.fn(),
    move: { id: 'mockId', status: MOVE_STATUSES.DRAFT },
    mtoShipments: [],
  };

  const getWrapper = (props = {}) => {
    return mount(<SelectShipmentType {...defaultProps} {...props} />);
  };

  it('should render radio buttons with no option selected', () => {
    const wrapper = getWrapper();
    expect(wrapper.find(Radio).length).toBe(4);

    // Ppm and HHG text renders
    expect(wrapper.find(Radio).at(0).text()).toContain('Do it yourself');
    expect(wrapper.find(Radio).at(1).text()).toContain('Professional movers');
    // No buttons should not be checked on page load
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').prop('checked')).toBe(false);
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').prop('checked')).toBe(false);
    expect(wrapper.find(Radio).at(2).find('.usa-radio__input').prop('checked')).toBe(false);
    expect(wrapper.find(Radio).at(3).find('.usa-radio__input').prop('checked')).toBe(false);
  });

  describe('modals', () => {
    const wrapper = getWrapper();
    const storageInfoModal = wrapper.find('ConnectedStorageInfoModal');
    const moveInfoModal = wrapper.find('ConnectedMoveInfoModal');

    describe('the storage info modal', () => {
      it('renders', () => {
        expect(storageInfoModal.exists()).toBe(true);
      });

      it('is closed by default', () => {
        expect(wrapper.state('showStorageInfoModal')).toEqual(false);
        expect(storageInfoModal.prop('isOpen')).toEqual(false);
      });

      it('can click the help button in the NTS card', () => {
        const ntsCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.NTS}"]`);
        expect(ntsCard.length).toBe(1);
        ntsCard.find('button[data-testid="helpButton"]').simulate('click');
        expect(wrapper.state('showStorageInfoModal')).toEqual(true);
        expect(wrapper.state('showStorageInfoModal')).toEqual(true);
      });

      it('can close the storage info modal after opening', () => {
        wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');
        expect(wrapper.state('showStorageInfoModal')).toEqual(false);
        expect(storageInfoModal.prop('isOpen')).toEqual(false);
      });
    });

    describe('the move info modal', () => {
      it('renders', () => {
        expect(moveInfoModal.exists()).toBe(true);
      });

      it('is closed by default', () => {
        expect(wrapper.state('showMoveInfoModal')).toEqual(false);
        expect(moveInfoModal.prop('isOpen')).toEqual(false);
      });

      it('can click the help button in the shipment selection card', () => {
        const hhgCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.HHG}"]`);
        expect(hhgCard.length).toBe(1);
        hhgCard.find('button[data-testid="helpButton"]').simulate('click');
        expect(wrapper.state('showMoveInfoModal')).toEqual(true);
        expect(wrapper.state('showMoveInfoModal')).toEqual(true);
      });

      it('can close the move info modal after opening', () => {
        wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');
        expect(wrapper.state('showMoveInfoModal')).toEqual(false);
        expect(moveInfoModal.prop('isOpen')).toEqual(false);
      });
    });
  });

  describe('when no PPMs or shipments have been created', () => {
    it('should render the correct text', () => {
      const wrapper = getWrapper();
      expect(wrapper.find('h1').text()).toContain('How do you want to move your belongings?');
      expect(wrapper.find('[data-testid="selectableCardText"]').at(0).text()).toContain(
        'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a a Personally Procured Move (PPM), sometimes called a DITY.',
      );
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 1');
      expect(wrapper.find('[data-testid="helper-footer"]').text()).toContain('Your move counselor will go');
    });
  });

  describe('when a PPM has already been created', () => {
    const props = {
      move: { personally_procured_moves: [{ id: '1' }] },
    };
    it('should render the correct text', () => {
      const wrapper = getWrapper(props);
      expect(wrapper.find('h1').text()).toContain('How do you want this group of things moved?');
      expect(wrapper.find(Radio).at(0).text()).toContain('Do it yourself (already chosen)');
      expect(wrapper.find('[data-testid="selectableCardText"]').at(0).text()).toContain(
        'You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty station.',
      );
      expect(wrapper.find('[data-testid="selectableCardText"]').at(0).text()).not.toContain(
        'You arrange to move some or all of your belongings',
      );
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 2');
      expect(wrapper.find('[data-testid="helper-footer"]').length).toBe(0);
    });
    it('should disable PPM form option if PPM is already submitted', () => {
      const wrapper = getWrapper(props);
      // PPM button should be disabled on page load
      expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('disabled');
    });
  });

  describe('when some shipments already exist', () => {
    it('should render the correct text', () => {
      const props = {
        mtoShipments: [{ selectedMoveType: SHIPMENT_OPTIONS.HHG, id: '2' }],
      };
      const wrapper = getWrapper(props);
      expect(wrapper.find('h1').text()).toContain('How do you want this group of things moved?');
    });
    it('should render the correct value in the eyebrow for shipment number with 1 existing shipment', () => {
      const props = {
        mtoShipments: [{ selectedMoveType: SHIPMENT_OPTIONS.HHG, id: '2' }],
      };
      const wrapper = getWrapper(props);
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 2');
    });
    it('should render the correct value in the eyebrow for shipment number with 2 existing shipment', () => {
      const props = {
        mtoShipments: [
          { selectedMoveType: SHIPMENT_OPTIONS.HHG, id: '6' },
          { selectedMoveType: SHIPMENT_OPTIONS.NTS, id: '9' },
        ],
      };
      const wrapper = getWrapper(props);
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 3');
    });
    it('should render the correct value in the shipment number with existing HHG and PPM', () => {
      const props = {
        move: { personally_procured_moves: [{ id: '1' }] },
        mtoShipments: [{ selectedMoveType: SHIPMENT_OPTIONS.HHG, id: '2' }],
      };
      const wrapper = getWrapper(props);
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 3');
    });
  });

  describe('when an NTS has already been created', () => {
    const props = {
      mtoShipments: [{ id: '3', shipmentType: SHIPMENT_OPTIONS.NTS }],
      move: { status: MOVE_STATUSES.DRAFT },
    };
    const wrapper = getWrapper(props);

    it('NTS card should render the correct text', () => {
      expect(wrapper.find('[data-testid="selectableCardText"]').at(2).text()).toContain(
        'You’ve already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.',
      );
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'These shipments do count against your weight allowance for this move.',
      );
    });
    it('NTS card should be disabled', () => {
      expect(wrapper.find(Radio).at(2).find('.usa-radio__input').prop('disabled')).toBe(true);
    });
  });

  describe('when an NTSr has already been created', () => {
    const props = {
      mtoShipments: [{ id: '4', shipmentType: SHIPMENT_OPTIONS.NTSR }],
      move: { status: MOVE_STATUSES.DRAFT },
    };
    const wrapper = getWrapper(props);
    it('NTSr card should render the correct text', () => {
      expect(wrapper.find('[data-testid="selectableCardText"]').at(3).text()).toContain(
        'You’ve already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.',
      );
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'These shipments do count against your weight allowance for this move.',
      );
    });
    it('NTSr card should be disabled', () => {
      expect(wrapper.find(Radio).at(3).find('.usa-radio__input').prop('disabled')).toBe(true);
    });
  });
  describe('when an unsubmitted move has both an NTS and an NTSr', () => {
    const props = {
      mtoShipments: [
        { id: '4', shipmentType: SHIPMENT_OPTIONS.NTS },
        { id: '5', shipmentType: SHIPMENT_OPTIONS.NTSR },
      ],
      move: { status: MOVE_STATUSES.DRAFT },
    };
    const wrapper = getWrapper(props);
    it('should render the correct text', () => {
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'Talk to your movers about long-term storage if you need to add it to this move or change a request you made earlier.',
      );
    });
    it('should not show radio cards for NTS or NTSr', () => {
      expect(wrapper.find(Radio).at(2).exists()).toEqual(false);
      expect(wrapper.find(Radio).at(3).exists()).toEqual(false);
    });
  });

  describe('when a move has already been submitted', () => {
    const props = {
      move: {
        status: MOVE_STATUSES.SUBMITTED,
      },
    };
    const wrapper = getWrapper(props);
    it('should render the correct text', () => {
      expect(wrapper.find('[data-testid="selectableCardText"]').at(1).text()).toContain(
        'Talk with your movers directly if you want to add or change shipments.',
      );
      expect(wrapper.find('[data-testid="selectableCardText"]').at(1).text()).not.toContain(
        'Professional movers take care of the whole shipment',
      );
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'Talk to your movers about long-term storage if you need to add it to this move or change a request you made earlier.',
      );
    });
    it('should disable HHG form option', () => {
      // HHG button should be disabled on page load
      expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('disabled');
    });
    it('should not show radio cards for NTS or NTSr', () => {
      expect(wrapper.find(Radio).at(2).exists()).toEqual(false);
      expect(wrapper.find(Radio).at(3).exists()).toEqual(false);
    });
    it('should have selectable PPM if move does not have a PPM, even if the move is already submitted', () => {
      expect(wrapper.find(Radio).at(0).prop('disabled')).toEqual(false);
    });
  });
});
