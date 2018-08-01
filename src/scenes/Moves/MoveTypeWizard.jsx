import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { updateMove, loadMove } from './ducks';
import WizardPage from 'shared/WizardPage';
import MoveType from './MoveType';

export class MoveTypeWizardPage extends Component {
  componentDidMount() {
    if (!this.props.currentMove) {
      this.props.loadMove(this.props.match.params.moveId);
    }
  }
  handleSubmit = () => {
    const { pendingMoveType, updateMove } = this.props;
    //todo: we should make sure this move matches the redux state
    const moveId = this.props.match.params.moveId;
    return updateMove(moveId, pendingMoveType);
  };
  render() {
    const { pages, pageKey, pendingMoveType, currentMove, error } = this.props;
    const moveType =
      pendingMoveType || (currentMove && 'selected_move_type' in currentMove);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={Boolean(moveType)}
        dirty={Boolean(pendingMoveType)}
        error={error}
      >
        <MoveType />
      </WizardPage>
    );
  }
}
MoveTypeWizardPage.propTypes = {
  updateMove: PropTypes.func.isRequired,
  pendingMoveType: PropTypes.string,
  currentMove: PropTypes.shape({
    id: PropTypes.string,
  }),
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateMove, loadMove }, dispatch);
}
function mapStateToProps(state) {
  return state.moves;
}
export default connect(mapStateToProps, mapDispatchToProps)(MoveTypeWizardPage);
