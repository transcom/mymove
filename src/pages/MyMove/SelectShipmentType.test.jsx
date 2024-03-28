/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount, shallow } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { isBooleanFlagEnabled } from '../../utils/featureFlags';
import { FEATURE_FLAG_KEYS, SHIPMENT_OPTIONS, MOVE_STATUSES } from '../../shared/constants';

import { SelectShipmentType } from 'pages/MyMove/SelectShipmentType';

jest.mock('../../utils/featureFlags', () => ({
  ...jest.requireActual('../../utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('SelectShipmentType', () => {
  const defaultProps = {
    updateMove: jest.fn(),
    router: { navigate: jest.fn() },

    loadMTOShipments: jest.fn(),
    move: { id: 'mockId', status: MOVE_STATUSES.DRAFT },
    mtoShipments: [],
  };

  const getWrapper = (props = {}) => {
    return mount(<SelectShipmentType {...defaultProps} {...props} />);
  };

  it('should render radio buttons with no option selected', () => {
    const wrapper = getWrapper();
    // set state to true for mount render
    wrapper.setState({ enablePPM: true });
    wrapper.setState({ enableNTS: true });
    wrapper.setState({ enableNTSR: true });
    expect(wrapper.find(Radio).length).toBe(4);

    // Ppm and HHG text renders
    expect(wrapper.find(Radio).at(0).text()).toContain('Movers pack and ship it, paid by the government (HHG)');
    expect(wrapper.find(Radio).at(1).text()).toContain('Move it yourself and get paid for it (PPM)');
    // No buttons should not be checked on page load
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').prop('checked')).toBe(false);
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').prop('checked')).toBe(false);
    expect(wrapper.find(Radio).at(2).find('.usa-radio__input').prop('checked')).toBe(false);
    expect(wrapper.find(Radio).at(3).find('.usa-radio__input').prop('checked')).toBe(false);
  });

  describe('modals', () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    const wrapper = getWrapper();
    // set state to true for mount render for test case
    wrapper.setState({ enablePPM: true });
    wrapper.setState({ enableNTS: true });
    wrapper.setState({ enableNTSR: true });
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

  describe('feature flags for shipment types show/hide', () => {
    it('feature flags for shipment types hide SelectableCard', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));

      const props = {};
      const wrapper = shallow(<SelectShipmentType {...defaultProps} {...props} />);
      await wrapper;
      const hhgCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.HHG}"]`);
      expect(hhgCard.length).toBe(1);
      const ppmCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.PPM}"]`);
      expect(ppmCard.length).toBe(0);
      const ntsCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.NTS}"]`);
      expect(ntsCard.length).toBe(0);
      const ntsrCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.NTSR}"]`);
      expect(ntsrCard.length).toBe(0);

      expect(wrapper.state('enablePPM')).toEqual(false);
      expect(wrapper.state('enableNTS')).toEqual(false);
      expect(wrapper.state('enableNTSR')).toEqual(false);

      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.PPM);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTS);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTSR);
    });

    it('feature flags for shipment types show SelectableCard', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

      const props = {};
      const wrapper = shallow(<SelectShipmentType {...defaultProps} {...props} />);
      await wrapper;
      const hhgCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.HHG}"]`);
      expect(hhgCard.length).toBe(1);
      const ppmCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.PPM}"]`);
      expect(ppmCard.length).toBe(1);
      const ntsCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.NTS}"]`);
      expect(ntsCard.length).toBe(1);
      const ntsrCard = wrapper.find(`SelectableCard[id="${SHIPMENT_OPTIONS.NTSR}"]`);
      expect(ntsrCard.length).toBe(1);

      expect(wrapper.state('enablePPM')).toEqual(true);
      expect(wrapper.state('enableNTS')).toEqual(true);
      expect(wrapper.state('enableNTSR')).toEqual(true);

      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.PPM);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTS);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTSR);
    });
  });

  describe('when no PPMs or shipments have been created', () => {
    it('should render the correct text', () => {
      const wrapper = getWrapper();
      // set state to true for mount render for test case
      wrapper.setState({ enablePPM: true });
      wrapper.setState({ enableNTS: true });
      wrapper.setState({ enableNTSR: true });
      expect(wrapper.find('h1').text()).toContain('How should this shipment move?');
      expect(wrapper.find('.usa-checkbox__label-description').at(1).text()).toContain(
        'You pack and move your personal property or make other arrangements, The government pays you for the weight you move. This is a Personally Procured Move (PPM).',
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
      const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
      // set state to true for mount render for test case
      wrapper.setState({ enablePPM: true });
      wrapper.setState({ enableNTS: true });
      wrapper.setState({ enableNTSR: true });
      expect(wrapper.find(Radio).at(1).text()).toContain('PPM');
      expect(wrapper.find('.usa-checkbox__label-description').at(0).text()).toContain(
        'Talk with your movers directly if you want to add or change shipments.',
      );
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 2');
      expect(wrapper.find('[data-testid="helper-footer"]').length).toBe(0);
    });
    it('should disable PPM form option if PPM is already submitted', () => {
      const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
      // PPM button should be disabled on page load
      expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('disabled');
    });
  });

  describe('when some shipments already exist', () => {
    it('should render the correct value in the eyebrow for shipment number with 1 existing shipment', () => {
      const props = {
        mtoShipments: [{ id: '2' }],
      };
      const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 2');
    });
    it('should render the correct value in the eyebrow for shipment number with 2 existing shipment', () => {
      const props = {
        mtoShipments: [{ id: '6' }, { id: '9' }],
      };
      const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 3');
    });
    it('should render the correct value in the shipment number with existing HHG and PPM', () => {
      const props = {
        move: { personally_procured_moves: [{ id: '1' }] },
        mtoShipments: [{ id: '2' }],
      };
      const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
      expect(wrapper.find('[data-testid="number-eyebrow"]').text()).toContain('Shipment 3');
    });
  });

  describe('when an NTS has already been created', () => {
    const props = {
      mtoShipments: [{ id: '3', shipmentType: SHIPMENT_OPTIONS.NTS }],
      move: { status: MOVE_STATUSES.DRAFT },
    };

    const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
    // set state to true for mount render for test case
    wrapper.setState({ enablePPM: true });
    wrapper.setState({ enableNTS: true });
    wrapper.setState({ enableNTSR: true });
    it('NTS card should render the correct text', () => {
      expect(wrapper.find('.usa-checkbox__label-description').at(2).text()).toContain(
        'You’ve already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.',
      );
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'Your orders might not authorize long-term storage — your counselor can verify.',
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
    const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
    // set state to true for mount render for test case
    wrapper.setState({ enablePPM: true });
    wrapper.setState({ enableNTS: true });
    wrapper.setState({ enableNTSR: true });
    it('NTSr card should render the correct text', () => {
      expect(wrapper.find('.usa-checkbox__label-description').at(3).text()).toContain(
        'You’ve already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.',
      );
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'Your orders might not authorize long-term storage — your counselor can verify.',
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
    const wrapper = mount(<SelectShipmentType {...defaultProps} {...props} />);
    // set state to true for mount render for test case
    wrapper.setState({ enablePPM: true });
    wrapper.setState({ enableNTS: true });
    wrapper.setState({ enableNTSR: true });
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
    // set state to true for mount render for test case
    wrapper.setState({ enablePPM: true });
    wrapper.setState({ enableNTS: true });
    wrapper.setState({ enableNTSR: true });
    it('should render the correct text', () => {
      expect(wrapper.find('.usa-checkbox__label-description').at(0).text()).toContain(
        'Talk with your movers directly if you want to add or change shipments.',
      );
      expect(wrapper.find('.usa-checkbox__label-description').at(0).text()).not.toContain(
        'Professional movers take care of the whole shipment',
      );
      expect(wrapper.find('[data-testid="long-term-storage-heading"] + p').text()).toEqual(
        'Talk to your movers about long-term storage if you need to add it to this move or change a request you made earlier.',
      );
    });
    it('should disable HHG form option', () => {
      // HHG button should be disabled on page load
      expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('disabled');
    });
    it('should not show radio cards for NTS or NTSr', () => {
      expect(wrapper.find(Radio).at(2).exists()).toEqual(false);
      expect(wrapper.find(Radio).at(3).exists()).toEqual(false);
    });
    it('should have selectable PPM if move does not have a PPM, even if the move is already submitted', () => {
      expect(wrapper.find(Radio).at(1).prop('disabled')).toEqual(false);
    });
  });
});
