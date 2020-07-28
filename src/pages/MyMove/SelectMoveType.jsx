import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes, { string, bool, func } from 'prop-types';
import { get } from 'lodash';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { selectActiveOrLatestMove, updateMove as updateMoveAction } from 'shared/Entities/modules/moves';
import { WizardPage } from 'shared/WizardPage';

export class SelectMoveType extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moveType: props.selectedMoveType,
    };
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
    const { pageKey, pageList, match, push } = this.props;
    const { moveType } = this.state;
    return (
      <WizardPage
        pageKey={pageKey}
        match={match}
        pageList={pageList}
        dirty
        // eslint-disable-next-line camelcase
        handleSubmit={this.handleSubmit}
        push={push}
      >
        <div className="usa-grid">
          <div className="grid-row">
            <div className="grid-col">
              <h1 className="sm-heading">How do you want to move your belongings?</h1>
              <Radio
                id={SHIPMENT_OPTIONS.PPM}
                label="I’ll move things myself"
                value={SHIPMENT_OPTIONS.PPM}
                name="moveType"
                onChange={(e) => this.setMoveType(e)}
                checked={moveType === SHIPMENT_OPTIONS.PPM}
              />
              <ul>
                <li>This is a PPM - “personally procured move”</li>
                <li>You arrange to move some or all of your belongings</li>
                <li>The government pays you an incentive based on weight</li>
                <li>DIY or hire your own movers</li>
              </ul>
              <Radio
                id={SHIPMENT_OPTIONS.HHG}
                label="The government packs for me and moves me"
                value={SHIPMENT_OPTIONS.HHG}
                onChange={(e) => this.setMoveType(e)}
                name="moveType"
                checked={moveType === SHIPMENT_OPTIONS.HHG}
              />
              <ul>
                <li>This is an HHG shipment — “household goods”</li>
                <li>The most popular kind of shipment</li>
                <li>Professional movers take care of the whole shipment</li>
                <li>They pack and move it for you</li>
              </ul>
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
  selectedMoveType: string.isRequired,
};

function mapStateToProps(state) {
  const move = selectActiveOrLatestMove(state);
  const props = {
    move: selectActiveOrLatestMove(state),
    selectedMoveType: get(move, 'selected_move_type'),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateMove: updateMoveAction }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);
