import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import { updateMove } from 'shared/Entities/modules/moves';
import ordersComplete from 'shared/images/orders-complete-gray-icon.png';
import moveIcon from 'shared/images/move-icon.png';
import { getCurrentMove, getCurrentMoveID } from '../../shared/UI/ducks';

export class TransitionToMove extends Component {
  componentDidMount() {
    if (!this.props.selectedMoveType) {
      // Make sure the move is always set to PPM since we no longer allow HHGs
      this.props.updateMove(this.props.moveId, 'PPM');
    }
  }

  render() {
    return (
      <div className="usa-grid">
        <div className="lg center">
          <p> Great, we're done with your orders.</p>
          <img className="sm" src={ordersComplete} alt="profile-check" />
        </div>

        <div className="lg center">
          <p>Now, we're ready to schedule your move!</p>
          <img className="sm" src={moveIcon} alt="onto-move-orders" />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const move = getCurrentMove(state);
  const props = {
    moveId: getCurrentMoveID(state),
    selectedMoveType: get(move, 'selected_move_type'),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateMove }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(TransitionToMove);
