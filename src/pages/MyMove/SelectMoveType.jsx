import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { string, bool, func, arrayOf, shape, number } from 'prop-types';
import { get } from 'lodash';

import styles from './SelectMoveType.module.scss';

import { updateMove as updateMoveAction } from 'scenes/Moves/ducks';
import { SHIPMENT_OPTIONS, MOVE_STATUSES } from 'shared/constants';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import { WizardPage } from 'shared/WizardPage';
import SelectableCard from 'components/Customer/SelectableCard';
import {
  selectMTOShipmentsByMoveId,
  loadMTOShipments as loadMTOShipmentsAction,
} from 'shared/Entities/modules/mtoShipments';
import { MoveTaskOrderShape, MTOShipmentShape } from 'types/moveOrder';

export class SelectMoveType extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moveType: props.selectedMoveType,
    };
  }

  componentDidMount() {
    const { loadMTOShipments, move } = this.props;
    loadMTOShipments(move.id);
  }

  setMoveType = (e) => {
    this.setState({ moveType: e.target.value });
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
      move,
      mtoShipments,
      isPpmSelectable,
      isHhgSelectable,
      shipmentNumber,
    } = this.props;
    const { moveType } = this.state;
    const hasPpm = !!move?.personally_procured_moves?.length; // eslint-disable-line camelcase
    const hasSubmittedMove = move?.status !== MOVE_STATUSES.DRAFT;
    const hasShipments = !!mtoShipments.length;
    const hasAnyShipments = hasPpm || hasShipments;
    const ppmCardText =
      'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a a Personally Procured Move (PPM), sometimes called a DITY.';
    const hhgCardText =
      'Your things are packed and moved by professionals, paid for by the government. This is a Household Goods move (HHG).';
    const ntsCardText = `Movers pack and ship things to a storage facility, where they stay until a future move. Your orders might not authorize long-term storage — your counselor can verify. This is an NTS (non-temporary storage) shipment.`;
    const ntsrCardText =
      'Movers pick up things you put into NTS during an earlier move and ship them to your new destination. This is an NTS-R (non-temporary storage release) shipment.';
    const hasNTSCardText =
      "You've already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.";
    const hasNTSRCardText =
      "You've already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.";
    const hhgCardTextPostSubmit = 'Talk with your movers directly if you want to add or change shipments.';
    const ppmCardTextAlreadyChosen = `You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty station.`;
    const selectableCardDefaultProps = {
      onChange: (e) => this.setMoveType(e),
      name: 'moveType',
    };
    // TODO: Make dynamic when we have ability to submit nts/ntsr
    const hasNTS = false;
    const hasNTSR = false;
    const selectPpmHasNoPpm = (
      <SelectableCard
        {...selectableCardDefaultProps} // eslint-disable-line
        label="Do it yourself"
        value={SHIPMENT_OPTIONS.PPM}
        id={SHIPMENT_OPTIONS.PPM}
        cardText={ppmCardText}
        checked={moveType === SHIPMENT_OPTIONS.PPM && isPpmSelectable}
        disabled={!isPpmSelectable}
      />
    );
    const selectPpmHasPpm = (
      <SelectableCard
        {...selectableCardDefaultProps} // eslint-disable-line
        label="Do it yourself (already chosen)"
        value={SHIPMENT_OPTIONS.PPM}
        id={SHIPMENT_OPTIONS.PPM}
        cardText={ppmCardTextAlreadyChosen}
        checked={moveType === SHIPMENT_OPTIONS.PPM && isPpmSelectable}
        disabled={!isPpmSelectable}
      />
    );
    const selectHhgDefault = (
      <SelectableCard
        {...selectableCardDefaultProps} // eslint-disable-line
        label="Professional movers"
        value={SHIPMENT_OPTIONS.HHG}
        id={SHIPMENT_OPTIONS.HHG}
        cardText={hhgCardText}
        checked={moveType === SHIPMENT_OPTIONS.HHG && isHhgSelectable}
        disabled={!isHhgSelectable}
      />
    );
    const selectHhgSubmittedMove = (
      <SelectableCard
        {...selectableCardDefaultProps} // eslint-disable-line
        label="Professional movers"
        value={SHIPMENT_OPTIONS.HHG}
        id={SHIPMENT_OPTIONS.HHG}
        cardText={hhgCardTextPostSubmit}
        checked={moveType === SHIPMENT_OPTIONS.HHG && isHhgSelectable}
        disabled={!isHhgSelectable}
      />
    );
    const footerText = (
      <div className={`${styles.footer} grid-col-12`}>
        It’s OK if you’re not sure about your choices. Your move counselor will go over all your options and can help
        make changes if necessary.
      </div>
    );
    return (
      <WizardPage
        pageKey={pageKey}
        match={match}
        pageList={pageList}
        dirty
        handleSubmit={this.handleSubmit}
        push={push}
        footerText={footerText}
      >
        <h6>Shipment {shipmentNumber}</h6>
        <h1 className={`${styles.selectTypeHeader} ${styles.header}`}>
          {hasAnyShipments ? 'How do you want this group of things moved?' : 'How do you want to move your belongings?'}
        </h1>
        <h2>Choose 1 shipment at a time.</h2>
        <p>You can add more later</p>
        {hasPpm ? selectPpmHasPpm : selectPpmHasNoPpm}
        {hasSubmittedMove ? selectHhgSubmittedMove : selectHhgDefault}
        <h3>Long-term storage</h3>
        <p>These shipments do count against your weight allowance for this move.</p>
        <SelectableCard
          {...selectableCardDefaultProps} // eslint-disable-line
          label="Put things into long-term storage"
          value={SHIPMENT_OPTIONS.NTS}
          id={SHIPMENT_OPTIONS.NTS}
          cardText={hasNTS ? ntsCardText : hasNTSCardText}
          checked={moveType === SHIPMENT_OPTIONS.NTS && isHhgSelectable}
          disabled={hasNTS}
        />
        <SelectableCard
          {...selectableCardDefaultProps} // eslint-disable-line
          label="Get things out of long-term storage"
          value={SHIPMENT_OPTIONS.NTS}
          id={SHIPMENT_OPTIONS.NTS}
          cardText={hasNTSR ? ntsrCardText : hasNTSRCardText}
          checked={moveType === SHIPMENT_OPTIONS.NTS && isHhgSelectable}
          disabled={hasNTSR}
        />
      </WizardPage>
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
  mtoShipments: arrayOf(MTOShipmentShape).isRequired,
  isPpmSelectable: bool.isRequired,
  isHhgSelectable: bool.isRequired,
  shipmentNumber: number.isRequired,
};

function mapStateToProps(state) {
  const move = selectActiveOrLatestMove(state);
  const hasPpm = !!move.personally_procured_moves?.length;
  const ppmCount = hasPpm ? 1 : 0;
  const hhgCount = selectMTOShipmentsByMoveId(state, move.id)?.length || 0;
  const props = {
    move,
    selectedMoveType: get(move, 'selected_move_type'),
    mtoShipments: selectMTOShipmentsByMoveId(state, move.id),
    isPpmSelectable: !hasPpm,
    isHhgSelectable: move.status === 'DRAFT',
    shipmentNumber: 1 + ppmCount + hhgCount,
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateMove: updateMoveAction, loadMTOShipments: loadMTOShipmentsAction }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);
export { mapStateToProps as _mapStateToProps };
