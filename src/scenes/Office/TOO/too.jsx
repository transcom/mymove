import React, { Component } from 'react';
import { arrayOf, shape, string } from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { getAllMoveOrders, selectMoveOrderList } from 'shared/Entities/modules/moveOrders';

class TOO extends Component {
  componentDidMount() {
    this.props.getAllMoveOrders();
  }

  handleCustomerInfoClick = (moveOrderId) => {
    this.props.history.push(`/moves/${moveOrderId}/details`);
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
              <th>MoveOrderID</th>
              <th>Actions</th>
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
                <tr data-testid="too-row" key={moveOrderId}>
                  <td onClick={() => this.handleCustomerInfoClick(moveOrderId)}>{`${last_name}, ${first_name}`}</td>
                  <td onClick={() => this.handleCustomerInfoClick(moveOrderId)}>{confirmation_number}</td>
                  <td onClick={() => this.handleCustomerInfoClick(moveOrderId)}>{agency}</td>
                  <td onClick={() => this.handleCustomerInfoClick(moveOrderId)}>
                    {originDutyStation && originDutyStation.name}
                  </td>
                  <td onClick={() => this.handleCustomerInfoClick(moveOrderId)}>{moveOrderId}</td>
                  <td>
                    <a href={`/too/${moveOrderId}/customer/${customerID}`}>Customer Details Page Skeleton</a>
                  </td>
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
  confirmation_number: string,
  branch_of_service: string,
  originDutyStation: shape({
    name: string.isRequired,
  }).isRequired,
});

TOO.propTypes = {
  moveOrders: arrayOf(customerMoveProps),
};

const mapStateToProps = (state) => {
  return {
    moveOrders: selectMoveOrderList(state),
  };
};
const mapDispatchToProps = {
  getAllMoveOrders,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(TOO));
