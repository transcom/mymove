import React from 'react';
import { string } from 'prop-types';
import deleteButtonImg from 'shared/images/delete-doc-button.png';

const ExpenseTicketListItem = ({ amount, type, paymentMethod, showDelete }) => (
  <div className="ticket-item">
    <div className="expense-li-item-container">
      <h4>
        {type} - ${amount}
      </h4>
      {showDelete && <img alt="delete document button" onClick={() => console.log('lol')} src={deleteButtonImg} />}
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
