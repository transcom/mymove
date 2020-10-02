import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes, { string, bool, func } from 'prop-types';
import { get } from 'lodash';

import styles from './SelectMoveType.module.scss';

import { updateMove as updateMoveAction } from 'scenes/Moves/ducks';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import { WizardPage } from 'shared/WizardPage';
import SelectableCard from 'components/Customer/SelectableCard';

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
    const hasPpm = true; // temp until data is fetched
    const hasSubmittedMove = true; // temp until data is fetched
    const ppmCardText =
      'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a a Personally Procured Move (PPM), sometimes called a DITY.';
    const hhgCardText =
      'Your things are packed and moved by professionals, paid for by the government. This is a Household Goods move (HHG).';
    const hhgCardTextPostSubmit = 'Talk with your movers directly if you want to add or change shipments.';
    const ppmCardTextAlreadyChosen = `You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty station.`;
    const footerText = (
      <div className={styles.footer}>
        It&apos;s OK if you&apos;re not sure about your choices. Your move counselor will go over all your options and
        can help make changes if necessary.
      </div>
    );
    return (
      <div className={`grid-container ${styles.gridContainer} ${styles.cardsContainer}`}>
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
              <h1 className="sm-heading">How do you want to move your belongings?</h1>
              <SelectableCard
                label={hasPpm ? 'Do it yourself (already chosen)' : 'Do it yourself'}
                onChange={(e) => this.setMoveType(e)}
                value={SHIPMENT_OPTIONS.PPM}
                name="moveType"
                id={SHIPMENT_OPTIONS.PPM}
                cardText={hasPpm ? ppmCardTextAlreadyChosen : ppmCardText}
                checked={moveType === SHIPMENT_OPTIONS.PPM}
              />
              <SelectableCard
                label="Professional movers"
                onChange={(e) => this.setMoveType(e)}
                value={SHIPMENT_OPTIONS.HHG}
                name="moveType"
                id={SHIPMENT_OPTIONS.HHG}
                cardText={hasSubmittedMove ? hhgCardTextPostSubmit : hhgCardText}
                checked={moveType === SHIPMENT_OPTIONS.HHG}
              />
            </WizardPage>
          </div>
          <div className="tablet:grid-col-2" />
        </div>
      </div>
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
