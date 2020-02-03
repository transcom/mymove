import React, { Component } from 'react';
import { arrayOf, shape, string } from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { getAllMoveOrders } from 'shared/Entities/modules/moveOrders';

class TOO extends Component {
  componentDidMount() {
    this.props.getAllMoveOrders();
  }

  handleCustomerInfoClick = (moveOrderId, customerId) => {
    this.props.history.push(`/too/customer-moves/${moveOrderId}/customer/${customerId}`);
  };

  render() {
    const { moveOrders } = this.props;
    return (
      <div>
        <h2>All Customer Moves</h2>
        <table>
          <thead>
            <tr>
              <th>Customer Name</th>
              <th>Confirmation #</th>
              <th>Agency</th>
              <th>Origin Duty Station</th>
            </tr>
          </thead>
          <tbody>
            {moveOrders.map(
              ({
                id: moveOrderId,
                first_name,
                last_name,
                confirmation_number,
                agency,
                originDutyStation,
                customerID,
              }) => (
                <tr
                  data-cy="too-row"
                  onClick={() => this.handleCustomerInfoClick(moveOrderId, customerID)}
                  key={moveOrderId}
                >
                  <td>{`${last_name}, ${first_name}`}</td>
                  <td>{confirmation_number}</td>
                  <td>{agency}</td>
                  <td>{originDutyStation && originDutyStation.name}</td>
                </tr>
              ),
            )}
          </tbody>
        </table>
      </div>
    );
  }
}

const customerMoveProps = shape({
  id: string.isRequired,
  first_name: string.isRequired,
  last_name: string.isRequired,
  confirmation_number: string.isRequired,
  branch_of_service: string,
  originDutyStation: shape({
    name: string.isRequired,
  }).isRequired,
});

TOO.propTypes = {
  moveOrders: arrayOf(customerMoveProps),
};

const mapStateToProps = state => {
  return {
    moveOrders: Object.values(get(state, 'entities.moveOrder', {})),
  };
};
const mapDispatchToProps = {
  getAllMoveOrders,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(TOO));
