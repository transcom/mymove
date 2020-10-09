import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes, { string, bool, func, number, array } from 'prop-types';
import { get } from 'lodash';
import { Radio } from '@trussworks/react-uswds';

import { updateMove as updateMoveAction } from 'scenes/Moves/ducks';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import { WizardPage } from 'shared/WizardPage';
import {
  loadMTOShipments as loadMTOShipmentsAction,
  selectMTOShipmentsByMoveId,
} from 'shared/Entities/modules/mtoShipments';

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
    const { pageKey, pageList, match, push, isPpmSelectable, isHhgSelectable, shipmentNumber } = this.props;
    const { moveType } = this.state;
    const ppmSelectableContent = (
      <ul>
        <li>This is a PPM - “personally procured move”</li>
        <li>You arrange to move some or all of your belongings</li>
        <li>The government pays you an incentive based on weight</li>
        <li>DIY or hire your own movers</li>
      </ul>
    );
    const ppmUnselectableContent = (
      <ul>
        <li>
          You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to
          that shipment, contact the PPPO at your origin duty station.
        </li>
      </ul>
    );

    const hhgSelectableContent = (
      <ul>
        <li>This is an HHG shipment — “household goods”</li>
        <li>The most popular kind of shipment</li>
        <li>Professional movers take care of the whole shipment</li>
        <li>They pack and move it for you</li>
      </ul>
    );
    const hhgUnselectableContent = (
      <ul>
        <li>Talk with your movers directly if you want to add or change shipments.</li>
      </ul>
    );

    return (
      <WizardPage
        pageKey={pageKey}
        match={match}
        pageList={pageList}
        dirty
        handleSubmit={this.handleSubmit}
        push={push}
      >
        <div className="usa-grid">
          <div className="grid-row">
            <div className="grid-col">
              <h6 className="sm-heading">Shipment {shipmentNumber}</h6>
              <h1 className="sm-heading">How do you want to move your belongings?</h1>
              <Radio
                id={SHIPMENT_OPTIONS.PPM}
                label="I’ll move things myself"
                value={SHIPMENT_OPTIONS.PPM}
                name="moveType"
                onChange={(e) => this.setMoveType(e)}
                checked={moveType === SHIPMENT_OPTIONS.PPM && isPpmSelectable}
                disabled={!isPpmSelectable}
              />
              {isPpmSelectable ? ppmSelectableContent : ppmUnselectableContent}
              <Radio
                id={SHIPMENT_OPTIONS.HHG}
                label="The government packs for me and moves me"
                value={SHIPMENT_OPTIONS.HHG}
                onChange={(e) => this.setMoveType(e)}
                name="moveType"
                checked={moveType === SHIPMENT_OPTIONS.HHG && isHhgSelectable}
                disabled={!isHhgSelectable}
              />
              {isHhgSelectable ? hhgSelectableContent : hhgUnselectableContent}
            </div>
          </div>
        </div>
      </WizardPage>
    );
  }
}

SelectMoveType.propTypes = {
  pageKey: PropTypes.string.isRequired,
  pageList: PropTypes.arrayOf(string).isRequired,
  match: PropTypes.shape({
    isExact: bool.isRequired,
    params: PropTypes.shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  push: func.isRequired,
  updateMove: func.isRequired,
  loadMTOShipments: func.isRequired,
  move: PropTypes.shape({
    id: string.isRequired,
    status: string.isRequired,
    personally_procured_moves: array,
  }).isRequired,
  selectedMoveType: string.isRequired,
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
