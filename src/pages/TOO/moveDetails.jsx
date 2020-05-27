import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import '../../index.scss';
import '../../ghc_index.scss';
import { get } from 'lodash';
import { getMoveOrder, selectMoveOrder } from '../../shared/Entities/modules/moveTaskOrders';
import { loadOrders } from '../../shared/Entities/modules/orders';
import OrdersTable from '../../components/Office/OrdersTable';

class MoveDetails extends Component {
  componentDidMount() {
    // eslint-disable-next-line react/destructuring-assignment,react/prop-types
    const { moveOrderId } = this.props.match.params;
    // eslint-disable-next-line react/prop-types,react/destructuring-assignment
    this.props.getMoveOrder(moveOrderId);
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveOrder } = this.props;
    return (
      <div className="grid-container-desktop-lg" data-cy="too-move-details">
        <h1>Move details</h1>
        <div className="container">
          <OrdersTable
            ordersInfo={{
              // eslint-disable-next-line react/prop-types
              newDutyStation: get(moveOrder.destinationDutyStation, 'name'),
              // eslint-disable-next-line react/prop-types
              currentDutyStation: get(moveOrder.originDutyStation, 'name'),
              // eslint-disable-next-line react/prop-types
              issuedDate: moveOrder.date_issued,
              // eslint-disable-next-line react/prop-types
              reportByDate: moveOrder.report_by_date,
              // eslint-disable-next-line react/prop-types
              departmentIndicator: moveOrder.department_indicator,
              // eslint-disable-next-line react/prop-types
              ordersNumber: moveOrder.order_number,
              // eslint-disable-next-line react/prop-types
              ordersType: moveOrder.order_type,
              // eslint-disable-next-line react/prop-types
              ordersTypeDetail: moveOrder.order_type_detail,
              // eslint-disable-next-line react/prop-types
              tacMDC: moveOrder.tac,
              // eslint-disable-next-line react/prop-types
              sacSDN: moveOrder.sacSDN,
            }}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;

  return {
    moveOrder: selectMoveOrder(state, moveOrderId),
  };
};

const mapDispatchToProps = {
  getMoveOrder,
  loadOrders,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
