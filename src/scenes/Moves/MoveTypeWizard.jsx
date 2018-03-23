import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { createMove } from './ducks';
import WizardPage from 'shared/WizardPage';
import MoveType from './MoveType';

export class MoveTypeWizardPage extends Component {
  componentDidMount() {
    // this.props.loadMove(this.props.match.params.moveId);
  }
  handleSubmit = () => {
    const { pendingMoveType, createMove } = this.props;
    //todo: we should make sure this move matches the redux state
    const moveId = this.props.match.params.moveId;
    if (pendingMoveType) {
      //don't create a move unless the type is selected
      createMove(moveId, pendingMoveType);
    }
  };
  render() {
    const { pages, pageKey, pendingMoveType } = this.props;
    const moveType = pendingMoveType;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={moveType !== null}
      >
        <MoveType />
      </WizardPage>
    );
  }
}
MoveTypeWizardPage.propTypes = {
  createMove: PropTypes.func.isRequired,
  pendingMoveType: PropTypes.string,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createMove }, dispatch);
}
function mapStateToProps(state) {
  return state.submittedMoves;
}
export default connect(mapStateToProps, mapDispatchToProps)(MoveTypeWizardPage);
