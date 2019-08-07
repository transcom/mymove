import React from 'react';
import { string } from 'prop-types';

const ExpenseTicketListItem = ({ amount, type, paymentMethod, showDelete }) => (
  <div className="ticket-item">
    <div className="expense-li-item-container">
      <h4>
        {type} - ${amount}
      </h4>
    </div>
    <div>
      {type} ({paymentMethod})
    </div>
  </div>
);

ExpenseTicketListItem.propTypes = {
  amount: string.isRequired,
  type: string.isRequired,
  paymentMethod: string.isRequired,
};

ExpenseTicketListItem.defaultProps = {
  showDelete: false,
};

export default ExpenseTicketListItem;
