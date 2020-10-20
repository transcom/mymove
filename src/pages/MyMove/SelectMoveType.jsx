/* eslint-disable react/jsx-props-no-spreading */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { string, bool, func, arrayOf, shape, number } from 'prop-types';
import { get } from 'lodash';

import styles from './SelectMoveType.module.scss';

import wizardStyles from 'pages/MyMove/index.module.scss';
import { updateMove as updateMoveAction } from 'scenes/Moves/ducks';
import { SHIPMENT_OPTIONS, MOVE_STATUSES } from 'shared/constants';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import { WizardPage } from 'shared/WizardPage';
import SelectableCard from 'components/Customer/SelectableCard';
import {
  selectMTOShipmentsByMoveId,
  loadMTOShipments as loadMTOShipmentsAction,
} from 'shared/Entities/modules/mtoShipments';
import { MoveTaskOrderShape } from 'types/moveOrder';
import ConnectedStorageInfoModal from 'components/Customer/modals/StorageInfoModal/StorageInfoModal';

export class SelectMoveType extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moveType: props.selectedMoveType,
      showStorageInfoModal: false,
    };
  }

  componentDidMount() {
    const { loadMTOShipments, move } = this.props;
    loadMTOShipments(move.id);
  }

  setMoveType = (e) => {
    this.setState({ moveType: e.target.value });
  };

  toggleStorageModal = () => {
    this.setState((state) => ({
      showStorageInfoModal: !state.showStorageInfoModal,
    }));
  };

  handleSubmit = () => {
    const { match, updateMove } = this.props;
    const { moveType } = this.state;
    return updateMove(match.params.moveId, moveType);
  };

  render() {
    const {
      pageKey,
      pageList,
      match,
      push,
      isPpmSelectable,
      isHhgSelectable,
      isNtsSelectable,
      isNtsrSelectable,
      shipmentNumber,
    } = this.props;
    const { moveType, showStorageInfoModal } = this.state;
    const ppmCardText =
      'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a a Personally Procured Move (PPM), sometimes called a DITY.';
    const hhgCardText =
      'Your things are packed and moved by professionals, paid for by the government. This is a Household Goods move (HHG).';
    const ntsCardText = `Movers pack and ship things to a storage facility, where they stay until a future move. Your orders might not authorize long-term storage — your counselor can verify. This is an NTS (non-temporary storage) shipment.`;
    const ntsrCardText =
      'Movers pick up things you put into NTS during an earlier move and ship them to your new destination. This is an NTS-R (non-temporary storage release) shipment.';
    const ntsDisabledText =
      'You‘ve already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.';
    const ntsrDisabledText =
      'You‘ve already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.';
    const hhgCardTextPostSubmit = 'Talk with your movers directly if you want to add or change shipments.';
    const ppmCardTextAlreadyChosen = `You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty station.`;
    const selectableCardDefaultProps = {
      onChange: (e) => this.setMoveType(e),
      name: 'moveType',
    };
    const ppmEnabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Do it yourself"
        value={SHIPMENT_OPTIONS.PPM}
        id={SHIPMENT_OPTIONS.PPM}
        cardText={ppmCardText}
        checked={moveType === SHIPMENT_OPTIONS.PPM}
        disabled={false}
      />
    );
    const ppmDisabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Do it yourself (already chosen)"
        value={SHIPMENT_OPTIONS.PPM}
        id={SHIPMENT_OPTIONS.PPM}
        cardText={ppmCardTextAlreadyChosen}
        checked={false}
        disabled={!isPpmSelectable}
      />
    );
    const hhgEnabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Professional movers"
        value={SHIPMENT_OPTIONS.HHG}
        id={SHIPMENT_OPTIONS.HHG}
        cardText={hhgCardText}
        checked={moveType === SHIPMENT_OPTIONS.HHG}
        disabled={false}
      />
    );
    const hhgDisabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Professional movers"
        value={SHIPMENT_OPTIONS.HHG}
        id={SHIPMENT_OPTIONS.HHG}
        cardText={hhgCardTextPostSubmit}
        checked={false}
        disabled={!isHhgSelectable}
      />
    );
    const footerText = (
      <div className={styles.footer}>
        It’s OK if you’re not sure about your choices. Your move counselor will go over all your options and can help
        make changes if necessary.
      </div>
    );
    return (
      <div className={`grid-container ${wizardStyles.gridContainer} ${styles.gridContainer}`}>
        <div className="grid-row">
          <div className="tablet:grid-col-2 desktop:grid-col-2" />
          <div className="tablet:grid-col-8 desktop:grid-col-8">
            <WizardPage
              pageKey={pageKey}
              match={match}
              pageList={pageList}
              dirty
              handleSubmit={this.handleSubmit}
              push={push}
              footerText={footerText}
            >
              <h6 className="sm-heading">Shipment {shipmentNumber}</h6>
              <h1 className={`sm-heading ${styles.selectTypeHeader} ${styles.header}`}>
                {shipmentNumber > 1
                  ? 'How do you want this group of things moved?'
                  : 'How do you want to move your belongings?'}
              </h1>
              <h2>Choose 1 shipment at a time.</h2>
              <p>You can add more later</p>
              {isPpmSelectable ? ppmEnabledCard : ppmDisabledCard}
              {isHhgSelectable ? hhgEnabledCard : hhgDisabledCard}
              <h3>Long-term storage</h3>
              <p>These shipments do count against your weight allowance for this move.</p>
              <SelectableCard
                {...selectableCardDefaultProps}
                label="Put things into long-term storage"
                value={SHIPMENT_OPTIONS.NTS}
                id={SHIPMENT_OPTIONS.NTS}
                cardText={isNtsSelectable ? ntsCardText : ntsDisabledText}
                checked={moveType === SHIPMENT_OPTIONS.NTS && isNtsSelectable}
                disabled={!isNtsSelectable}
                onHelpClick={this.toggleStorageModal}
              />
              {/* TODO - update when NTSR option is added to API */}
              <SelectableCard
                {...selectableCardDefaultProps}
                label="Get things out of long-term storage"
                value={SHIPMENT_OPTIONS.NTSR}
                id={SHIPMENT_OPTIONS.NTSR}
                cardText={isNtsSelectable ? ntsrCardText : ntsrDisabledText}
                checked={moveType === SHIPMENT_OPTIONS.NTSR && isNtsrSelectable}
                disabled={!isNtsrSelectable}
                onHelpClick={this.toggleStorageModal}
              />
            </WizardPage>
          </div>
          <div className="tablet:grid-col-2" />
        </div>

        <ConnectedStorageInfoModal isOpen={showStorageInfoModal} closeModal={this.toggleStorageModal} />
      </div>
    );
  }
}

SelectMoveType.propTypes = {
  pageKey: string.isRequired,
  pageList: arrayOf(string).isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  push: func.isRequired,
  updateMove: func.isRequired,
  loadMTOShipments: func.isRequired,
  selectedMoveType: string.isRequired,
  move: MoveTaskOrderShape.isRequired,
  isPpmSelectable: bool.isRequired,
  isHhgSelectable: bool.isRequired,
  isNtsSelectable: bool.isRequired,
  isNtsrSelectable: bool.isRequired,
  shipmentNumber: number.isRequired,
};

function mapStateToProps(state) {
  const move = selectActiveOrLatestMove(state);
  const hasPpm = !!move.personally_procured_moves?.length;
  // TODO: Make dynamic when we have ability to submit nts/ntsr
  const hasNTS = false;
  const hasNTSR = false;
  const ppmCount = hasPpm ? 1 : 0;
  const mtosCount = selectMTOShipmentsByMoveId(state, move.id)?.length || 0;
  const isMoveDraft = move.status === MOVE_STATUSES.DRAFT;
  const props = {
    move,
    selectedMoveType: get(move, 'selected_move_type'),
    isPpmSelectable: !hasPpm,
    isHhgSelectable: isMoveDraft,
    isNtsSelectable: isMoveDraft && !hasNTS,
    isNtsrSelectable: isMoveDraft && !hasNTSR,
    shipmentNumber: 1 + ppmCount + mtosCount,
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateMove: updateMoveAction, loadMTOShipments: loadMTOShipmentsAction }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);
export { mapStateToProps as _mapStateToProps };
