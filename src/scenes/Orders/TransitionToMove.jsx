import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import { updateMove } from '../Moves/ducks';
import ordersComplete from 'shared/images/orders-complete-gray-icon.png';
import moveIcon from 'shared/images/move-icon.png';
import { selectMoveFromServiceMemberId } from 'shared/Entities/modules/moves';
import { fetchLatestOrders, getLatestOrdersLabel } from 'shared/Entities/modules/orders';
import { getRequestStatus } from 'shared/Swagger/selectors';

export class TransitionToMove extends Component {
  componentDidMount() {
    //  TODO fix this - error moveId not string
    // if (!this.props.selectedMoveType) {
    //   // Make sure the move is always set to PPM since we no longer allow HHGs
    //   this.props.updateMove(this.props.moveId, 'PPM');
    // }
    this.props.fetchLatestOrders(this.props.serviceMemberId);
  }

  render() {
    return (
      <div className="usa-grid">
        <div className="grid-row grid-gap">
          <div className="grid-col-3 desktop:grid-col-2 text-right">
            <img className="sm margin-top-3 desktop:margin-top-1" src={ordersComplete} alt="profile-check" />
          </div>
          <div className="grid-col-9 desktop:grid-col-10">
            <h1 className="sm-heading">Thank you for entering info about your orders.</h1>
          </div>
        </div>
        <div className="grid-row grid-gap">
          <div className="grid-col-3 desktop:grid-col-2 text-right">
            <img className="sm margin-top-5 desktop:margin-top-1" src={moveIcon} alt="onto-move-orders" />
          </div>
          <div className="grid-col-9 desktop:grid-col-10">
            <h1 className="sm-heading">Now let’s schedule your move.</h1>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  // const move = get(state, 'moves.currentMove');
  const serviceMemberId = get(state, 'serviceMember.currentServiceMember.id');
  const showOrdersRequest = getRequestStatus(state, getLatestOrdersLabel);
  const move = selectMoveFromServiceMemberId(state, serviceMemberId);

  const props = {
    serviceMemberId: serviceMemberId,
    moveId: get(move, 'id'),
    selectedMoveType: get(move, 'selected_move_type'),
    loadDependenciesHasSuccess: showOrdersRequest.isSuccess,
    loadDependenciesHasError: showOrdersRequest.error,
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ fetchLatestOrders, updateMove }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(TransitionToMove);
