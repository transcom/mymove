import React, { Component } from 'react';
import { arrayOf, shape, string } from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { getAllCustomerMoves } from 'shared/Entities/modules/moveTaskOrders';

class TOO extends Component {
  componentDidMount() {
    this.props.getAllCustomerMoves();
  }

  handleCustomerInfoClick = customerId => {
    this.props.history.push(`/too/customer/${customerId}/details`);
  };

  render() {
    const { customerMoves } = this.props;
    return (
      <div>
        <h2>All Customer Moves</h2>
        <table>
          <thead>
            <tr>
              <th>Customer Name</th>
              <th>Confirmation #</th>
              <th>Branch of Service</th>
              <th>Origin Duty Station</th>
            </tr>
          </thead>
          <tbody>
            {customerMoves.map(
              ({
                id,
                customer_name,
                customer_id,
                confirmation_number,
                branch_of_service,
                origin_duty_station_name,
              }) => (
                <tr data-cy="too-row" onClick={() => this.handleCustomerInfoClick(customer_id)} key={id}>
                  <td>{customer_name}</td>
                  <td>{confirmation_number}</td>
                  <td>{branch_of_service}</td>
                  <td>{origin_duty_station_name}</td>
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
  customer_name: string.isRequired,
  confirmation_number: string.isRequired,
  branch_of_service: string.isRequired,
  origin_duty_station_name: string.isRequired,
});

TOO.propTypes = {
  customerMoves: arrayOf(customerMoveProps),
};

const mapStateToProps = state => {
  return {
    customerMoves: Object.values(get(state, 'entities.customerMoveItem', {})),
  };
};
const mapDispatchToProps = {
  getAllCustomerMoves,
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(TOO),
);
