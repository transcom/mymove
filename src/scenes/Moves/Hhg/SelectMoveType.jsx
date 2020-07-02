import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import ordersComplete from 'shared/images/orders-complete-gray-icon.png';
import moveIcon from 'shared/images/move-icon.png';

export class SelectMoveType extends Component {
  componentDidMount() {
    // if (!this.props.selectedMoveType) {
    //   // Make sure the move is always set to PPM since we no longer allow HHGs
    //   this.props.updateMove(this.props.moveId, 'PPM');
    // }
  }

  render() {
    return (
      <div className="usa-grid">
        <div className="grid-row grid-gap">
          <div className="grid-col-3 desktop:grid-col-2 text-right">
            <img className="sm margin-top-3 desktop:margin-top-1" src={ordersComplete} alt="profile-check" />
          </div>
          <div className="grid-col-9 desktop:grid-col-10">
            <h1 className="sm-heading">SELECT MOVE TYPE</h1>
          </div>
        </div>
        <div className="grid-row grid-gap">
          <div className="grid-col-3 desktop:grid-col-2 text-right">
            <img className="sm margin-top-5 desktop:margin-top-1" src={moveIcon} alt="onto-move-orders" />
          </div>
          <div className="grid-col-9 desktop:grid-col-10">
            <h1 className="sm-heading">IT'S A PPM DON'T CARE WHAT YOU ACTUALLY WANT</h1>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const move = get(state, 'moves.currentMove');
  const props = {
    moveId: get(move, 'id'),
    selectedMoveType: get(move, 'selected_move_type'),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);
