import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import '../../index.scss';
import '../../ghc_index.scss';
// import OrdersTable from '../../components/Office/OrdersTable';
import { get } from 'lodash';
import { getMoveByLocator, selectMoveByLocator } from '../../shared/Entities/modules/moves';
import { loadOrders, selectOrdersForMove } from '../../shared/Entities/modules/orders';
import OrdersTable from '../../components/Office/OrdersTable';

// import OrdersTable from "../../components/Office/OrdersTable";

class MoveDetails extends Component {
  componentDidMount() {
    // eslint-disable-next-line react/destructuring-assignment,react/prop-types
    const { locator } = this.props.match.params;
    // eslint-disable-next-line react/prop-types,react/destructuring-assignment
    this.props.getMoveByLocator(locator).then(({ response: { body: move } }) => {
      // eslint-disable-next-line react/prop-types,react/destructuring-assignment
      this.props.loadOrders(move.orders_id);
    });
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { orders } = this.props;
    return (
      <div className="grid-container-desktop-lg" data-cy="too-move-details">
        <h1>Move details</h1>
        <div className="container">
          <OrdersTable
            ordersInfo={{
              // eslint-disable-next-line react/prop-types
              newDutyStation: get(orders.new_duty_station, 'name'),
              // eslint-disable-next-line react/prop-types
              issuedDate: orders.issue_date,
              // eslint-disable-next-line react/prop-types
              reportByDate: orders.report_by_date,
              // eslint-disable-next-line react/prop-types
              departmentIndicator: orders.department_indicator,
              // eslint-disable-next-line react/prop-types
              ordersNumber: orders.orders_number,
              // eslint-disable-next-line react/prop-types
              ordersType: orders.orders_type,
              // eslint-disable-next-line react/prop-types
              ordersTypeDetail: orders.orders_type_detail,
              // eslint-disable-next-line react/prop-types
              tacMDC: orders.tac,
              // eslint-disable-next-line react/prop-types
              sacSDN: orders.sacSDN,
            }}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { locator } = ownProps.match.params;
  const move = selectMoveByLocator(state, locator);
  let moveId;
  if (move) {
    moveId = move.id;
  }
  return {
    move,
    orders: selectOrdersForMove(state, moveId),
  };
};

const mapDispatchToProps = {
  getMoveByLocator,
  loadOrders,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
